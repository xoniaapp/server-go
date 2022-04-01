package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/xoniaapp/server/model"
	"github.com/xoniaapp/server/model/apperrors"
)

/*
 * AuthHandler contains all routes related to account actions (/api/account) editing it might break it
 */

type registerReq struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r registerReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.Username, validation.Required, validation.Length(3, 30)),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
	)
}

func (r *registerReq) sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func (h *Handler) Register(c *gin.Context) {
	var req registerReq

	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	initial := &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	user, err := h.userService.Register(initial)

	if err != nil {
		if err.Error() == apperrors.NewBadRequest(apperrors.DuplicateEmail).Error() {
			toFieldErrorResponse(c, "Email", apperrors.DuplicateEmail)
			return
		}
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	setUserSession(c, user.ID)

	c.JSON(http.StatusCreated, user)
}

type loginReq struct {
	// Must be unique
	Email string `json:"email"`
	// Min 6, max 150 characters.
	Password string `json:"password"`
} //@name LoginRequest

func (r loginReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
	)
}

func (r *loginReq) sanitize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq

	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	user, err := h.userService.Login(req.Email, req.Password)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	setUserSession(c, user.ID)

	c.JSON(http.StatusOK, user)
}

func (h *Handler) Logout(c *gin.Context) {
	c.Set("user", nil)

	session := sessions.Default(c)
	session.Set("userId", "")
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	err := session.Save()

	if err != nil {
		log.Printf("error clearing session: %v\n", err.Error())
	}

	c.JSON(http.StatusOK, true)
}

type forgotRequest struct {
	Email string `json:"email"`
} //@ForgotPasswordRequest

func (r forgotRequest) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
	)
}

func (r *forgotRequest) sanitize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
}

// ForgotPassword sends a password reset email to the requested email but it doesn't work right now for some reason. gotta debug that
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req forgotRequest
	if valid := bindData(c, &req); !valid {
		return
	}

	req.sanitize()

	user, err := h.userService.GetByEmail(req.Email)

	if err != nil {
		// No user with the email found
		if err.Error() == apperrors.NewNotFound("email", req.Email).Error() {
			c.JSON(http.StatusOK, true)
			return
		}

		e := apperrors.NewInternal()
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	ctx := c.Request.Context()
	err = h.userService.ForgotPassword(ctx, user)

	if err != nil {
		e := apperrors.NewInternal()
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}

type resetRequest struct {
	Token           string `json:"token"`
	Password        string `json:"newPassword"`
	ConfirmPassword string `json:"confirmNewPassword"`
} //@name ResetPasswordRequest

func (r resetRequest) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Token, validation.Required),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
		validation.Field(&r.ConfirmPassword, validation.Required, validation.Length(6, 150)),
	)
}

func (r *resetRequest) sanitize() {
	r.Token = strings.TrimSpace(r.Token)
	r.Password = strings.TrimSpace(r.Password)
	r.ConfirmPassword = strings.TrimSpace(r.ConfirmPassword)
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req resetRequest

	if valid := bindData(c, &req); !valid {
		return
	}

	req.sanitize()

	// Check if passwords match
	if req.Password != req.ConfirmPassword {
		toFieldErrorResponse(c, "Password", apperrors.PasswordsDoNotMatch)
		return
	}

	ctx := c.Request.Context()
	user, err := h.userService.ResetPassword(ctx, req.Password, req.Token)

	if err != nil {
		if err.Error() == apperrors.NewBadRequest(apperrors.InvalidResetToken).Error() {
			toFieldErrorResponse(c, "Token", apperrors.InvalidResetToken)
			return
		}
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	setUserSession(c, user.ID)

	c.JSON(http.StatusOK, user)
}
