package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xoniaapp/server/model/apperrors"
)

func Timeout(timeout time.Duration, errTimeout *apperrors.Error) gin.HandlerFunc {
	return func(c *gin.Context) {
		tw := &timeoutWriter{ResponseWriter: c.Writer, h: make(http.Header)}
		c.Writer = tw

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-panicChan:
			e := apperrors.NewInternal()
			tw.ResponseWriter.WriteHeader(e.Status())
			eResp, _ := json.Marshal(gin.H{
				"error": e,
			})
			_, _ = tw.ResponseWriter.Write(eResp)
		case <-finished:
			tw.mu.Lock()
			defer tw.mu.Unlock()
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			_, _ = tw.ResponseWriter.Write(tw.wbuf.Bytes())
		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()
			tw.ResponseWriter.Header().Set("Content-Type", "application/json")
			tw.ResponseWriter.WriteHeader(errTimeout.Status())
			eResp, _ := json.Marshal(gin.H{
				"error": errTimeout,
			})
			_, _ = tw.ResponseWriter.Write(eResp)
			c.Abort()
			tw.SetTimedOut()
		}
	}
}

type timeoutWriter struct {
	gin.ResponseWriter
	h    http.Header
	wbuf bytes.Buffer

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, nil
	}

	return tw.wbuf.Write(b)
}
func (tw *timeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}
func (tw *timeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}
func (tw *timeoutWriter) Header() http.Header {
	return tw.h
}
func (tw *timeoutWriter) SetTimedOut() {
	tw.timedOut = true
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}

// LICENSE@GNUv3
