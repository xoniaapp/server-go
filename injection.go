//                  GNU GENERAL PUBLIC LICENSE
//                      Version 3, 29 June 2007
//
// Copyright (C) 2007 Free Software Foundation, Inc. <https://fsf.org/>
// Everyone is permitted to copy and distribute verbatim copies
// of this license document, but changing it is not allowed.
//
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/aelpxy/xoniaapp/handler"
	"github.com/aelpxy/xoniaapp/handler/middleware"
	"github.com/aelpxy/xoniaapp/model"
	"github.com/aelpxy/xoniaapp/repository"
	"github.com/aelpxy/xoniaapp/service"
	"github.com/aelpxy/xoniaapp/ws"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	userRepository := repository.NewUserRepository(d.DB)
	friendRepository := repository.NewFriendRepository(d.DB)
	guildRepository := repository.NewGuildRepository(d.DB)
	channelRepository := repository.NewChannelRepository(d.DB)
	messageRepository := repository.NewMessageRepository(d.DB)

	bucketName := os.Getenv("AWS_STORAGE_BUCKET_NAME")
	fileRepository := repository.NewFileRepository(d.S3Session, bucketName)
	redisRepository := repository.NewRedisRepository(d.RedisClient)

	mailUser := os.Getenv("MAIL_USER")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	origin := os.Getenv("CORS_ORIGIN")
	mailRepository := repository.NewMailRepository(mailUser, mailPassword, origin)

	userService := service.NewUserService(&service.USConfig{
		UserRepository:  userRepository,
		FileRepository:  fileRepository,
		RedisRepository: redisRepository,
		MailRepository:  mailRepository,
	})

	friendService := service.NewFriendService(&service.FSConfig{
		UserRepository:   userRepository,
		FriendRepository: friendRepository,
	})

	guildService := service.NewGuildService(&service.GSConfig{
		UserRepository:    userRepository,
		FileRepository:    fileRepository,
		RedisRepository:   redisRepository,
		GuildRepository:   guildRepository,
		ChannelRepository: channelRepository,
	})

	channelService := service.NewChannelService(&service.CSConfig{
		ChannelRepository: channelRepository,
		GuildRepository:   guildRepository,
	})

	messageService := service.NewMessageService(&service.MSConfig{
		MessageRepository: messageRepository,
		FileRepository:    fileRepository,
	})

	router := gin.Default()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{origin},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})
	router.Use(c)

	redisURL := d.RedisClient.Options().Addr
	password := d.RedisClient.Options().Password

	secret := os.Getenv("SECRET")
	store, _ := redis.NewStore(10, "tcp", redisURL, password, []byte(secret))

	domain := os.Getenv("DOMAIN")

	store.Options(sessions.Options{
		Domain:   domain,
		MaxAge:   60 * 60 * 24 * 7,
		Secure:   gin.Mode() == gin.ReleaseMode,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	router.Use(sessions.Sessions(model.CookieName, store))

	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	maxBodyBytes := os.Getenv("MAX_BODY_BYTES")
	mbb, err := strconv.ParseInt(maxBodyBytes, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse MAX_BODY_BYTES as int: %w", err)
	}

	rate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  1000,
	}

	limitStore, _ := sredis.NewStore(d.RedisClient)

	rateLimiter := mgin.NewMiddleware(limiter.New(limitStore, rate))
	router.Use(rateLimiter)

	hub := ws.NewWebsocketHub(&ws.Config{
		UserService:    userService,
		GuildService:   guildService,
		ChannelService: channelService,
		Redis:          d.RedisClient,
	})
	go hub.Run()

	router.GET("/ws", middleware.AuthUser(), func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	socketService := service.NewSocketService(&service.SSConfig{
		Hub:               *hub,
		GuildRepository:   guildRepository,
		ChannelRepository: channelRepository,
	})

	handler.NewHandler(&handler.Config{
		R:               router,
		UserService:     userService,
		FriendService:   friendService,
		GuildService:    guildService,
		ChannelService:  channelService,
		MessageService:  messageService,
		SocketService:   socketService,
		TimeoutDuration: time.Duration(ht) * time.Second,
		MaxBodyBytes:    mbb,
	})

	return router, nil
}
