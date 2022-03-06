package model

import (
	"context"
	"time"

	"github.com/lib/pq"
)

type Guild struct {
	BaseModel
	Name        string `gorm:"not null"`
	OwnerId     string `gorm:"not null"`
	Icon        *string
	InviteLinks pq.StringArray `gorm:"type:text[]"`
	Members     []User         `gorm:"many2many:members;constraint:OnDelete:CASCADE;"`
	Channels    []Channel      `gorm:"constraint:OnDelete:CASCADE;"`
	Bans        []User         `gorm:"many2many:bans;constraint:OnDelete:CASCADE;"`
}
type GuildResponse struct {
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	OwnerId          string    `json:"ownerId"`
	Icon             *string   `json:"icon"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	HasNotification  bool      `json:"hasNotification"`
	DefaultChannelId string    `json:"default_channel_id"`
}

func (g Guild) SerializeGuild(channelId string) GuildResponse {
	return GuildResponse{
		Id:               g.ID,
		Name:             g.Name,
		OwnerId:          g.OwnerId,
		Icon:             g.Icon,
		CreatedAt:        g.CreatedAt,
		UpdatedAt:        g.UpdatedAt,
		HasNotification:  false,
		DefaultChannelId: channelId,
	}
}

type GuildService interface {
	GetUser(uid string) (*User, error)
	GetGuild(id string) (*Guild, error)
	GetUserGuilds(uid string) (*[]GuildResponse, error)
	GetGuildMembers(userId string, guildId string) (*[]MemberResponse, error)
	CreateGuild(guild *Guild) (*Guild, error)
	GenerateInviteLink(ctx context.Context, guildId string, isPermanent bool) (string, error)
	UpdateGuild(guild *Guild) error
	GetGuildIdFromInvite(ctx context.Context, token string) (string, error)
	GetDefaultChannel(guildId string) (*Channel, error)
	InvalidateInvites(ctx context.Context, guild *Guild)
	RemoveMember(userId string, guildId string) error
	UnbanMember(userId string, guildId string) error
	DeleteGuild(guildId string) error
	GetBanList(guildId string) (*[]BanResponse, error)
	GetMemberSettings(userId string, guildId string) (*MemberSettings, error)
	UpdateMemberSettings(settings *MemberSettings, userId string, guildId string) error
	FindUsersByIds(ids []string, guildId string) (*[]User, error)
	UpdateMemberLastSeen(userId, guildId string) error
}
type GuildRepository interface {
	FindUserByID(uid string) (*User, error)
	FindByID(id string) (*Guild, error)
	List(uid string) (*[]GuildResponse, error)
	GuildMembers(userId string, guildId string) (*[]MemberResponse, error)
	Create(guild *Guild) (*Guild, error)
	Save(guild *Guild) error
	RemoveMember(userId string, guildId string) error
	Delete(guildId string) error
	UnbanMember(userId string, guildId string) error
	GetBanList(guildId string) (*[]BanResponse, error)
	GetMemberSettings(userId string, guildId string) (*MemberSettings, error)
	UpdateMemberSettings(settings *MemberSettings, userId string, guildId string) error
	FindUsersByIds(ids []string, guildId string) (*[]User, error)
	GetMember(userId, guildId string) (*User, error)
	UpdateMemberLastSeen(userId, guildId string) error
	GetMemberIds(guildId string) (*[]string, error)
}
