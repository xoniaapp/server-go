package handler

import (
	"log"
	"net/http"

	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xoniaapp/server/handler/middleware"
	"github.com/xoniaapp/server/model"
	"github.com/xoniaapp/server/model/apperrors"
)

type Handler struct {
	userService    model.UserService
	friendService  model.FriendService
	guildService   model.GuildService
	channelService model.ChannelService
	messageService model.MessageService
	socketService  model.SocketService
	MaxBodyBytes   int64
}

type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	FriendService   model.FriendService
	GuildService    model.GuildService
	ChannelService  model.ChannelService
	MessageService  model.MessageService
	SocketService   model.SocketService
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

func NewHandler(c *Config) {
	h := &Handler{
		userService:    c.UserService,
		friendService:  c.FriendService,
		guildService:   c.GuildService,
		channelService: c.ChannelService,
		messageService: c.MessageService,
		socketService:  c.SocketService,
		MaxBodyBytes:   c.MaxBodyBytes,
	}

	c.R.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No route found.",
		})
	})

	if gin.Mode() != gin.TestMode {
		c.R.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
	}

	ag := c.R.Group("api/account")
	ag.POST("/register", h.Register)
	ag.POST("/login", h.Login)
	ag.POST("/logout", h.Logout)
	ag.POST("/forgot-password", h.ForgotPassword)
	ag.POST("/reset-password", h.ResetPassword)
	ag.Use(middleware.AuthUser())
	ag.GET("", h.GetCurrent)
	ag.PUT("", h.Edit)
	ag.PUT("/change-password", h.ChangePassword)
	ag.GET("/me/friends", h.GetUserFriends)
	ag.GET("/me/pending", h.GetUserRequests)
	ag.POST("/:memberId/friend", h.SendFriendRequest)
	ag.DELETE("/:memberId/friend", h.RemoveFriend)
	ag.POST("/:memberId/friend/accept", h.AcceptFriendRequest)
	ag.POST("/:memberId/friend/cancel", h.CancelFriendRequest)
	gg := c.R.Group("api/guilds")
	gg.Use(middleware.AuthUser())
	gg.GET("/:guildId/members", h.GetGuildMembers)
	gg.GET("", h.GetUserGuilds)
	gg.POST("/create", h.CreateGuild)
	gg.GET("/:guildId/invite", h.GetInvite)
	gg.DELETE("/:guildId/invite", h.DeleteGuildInvites)
	gg.POST("/join", h.JoinGuild)
	gg.GET("/:guildId/member", h.GetMemberSettings)
	gg.PUT("/:guildId/member", h.EditMemberSettings)
	gg.DELETE("/:guildId", h.LeaveGuild)
	gg.PUT("/:guildId", h.EditGuild)
	gg.DELETE("/:guildId/delete", h.DeleteGuild)
	gg.GET("/:guildId/bans", h.GetBanList)
	gg.POST("/:guildId/bans", h.BanMember)
	gg.DELETE("/:guildId/bans", h.UnbanMember)
	gg.POST("/:guildId/kick", h.KickMember)
	cg := c.R.Group("api/channels")
	cg.Use(middleware.AuthUser())
	cg.GET("/:id", h.GuildChannels)
	cg.POST("/:id", h.CreateChannel)
	cg.GET("/:id/members", h.PrivateChannelMembers)
	cg.POST("/:id/dm", h.GetOrCreateDM)
	cg.GET("/me/dm", h.DirectMessages)
	cg.PUT("/:id", h.EditChannel)
	cg.DELETE("/:id", h.DeleteChannel)
	cg.DELETE("/:id/dm", h.CloseDM)
	mg := c.R.Group("api/messages")
	mg.Use(middleware.AuthUser())
	mg.GET("/:channelId", h.GetMessages)
	mg.POST("/:channelId", h.CreateMessage)
	mg.PUT("/:messageId", h.EditMessage)
	mg.DELETE("/:messageId", h.DeleteMessage)
}

func setUserSession(c *gin.Context, id string) {
	session := sessions.Default(c)
	session.Set("userId", id)
	if err := session.Save(); err != nil {
		log.Printf("Error setting the session: %v\n", err.Error())
	}
}

func toFieldErrorResponse(c *gin.Context, field, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errors": []model.FieldError{
			{
				Field:   field,
				Message: message,
			},
		},
	})
}
