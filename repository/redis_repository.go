package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/aelpxy/xoniaapp/model"
	"github.com/aelpxy/xoniaapp/model/apperrors"
	"log"
	"time"
)

type redisRepository struct {
	rds *redis.Client
}

func NewRedisRepository(rds *redis.Client) model.RedisRepository {
	return &redisRepository{
		rds: rds,
	}
}

const (
	InviteLinkPrefix     = "inviteLink"
	ForgotPasswordPrefix = "forgot-password"
)

func (r *redisRepository) SetResetToken(ctx context.Context, id string) (string, error) {
	uid, err := gonanoid.New()

	if err != nil {
		log.Printf("Failed to generate id: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	if err = r.rds.Set(ctx, fmt.Sprintf("%s:%s", ForgotPasswordPrefix, uid), id, 24*time.Hour).Err(); err != nil {
		log.Printf("Failed to set link in redis: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	return uid, nil
}

func (r *redisRepository) GetIdFromToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("%s:%s", ForgotPasswordPrefix, token)
	val, err := r.rds.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", apperrors.NewBadRequest(apperrors.InvalidResetToken)
	}
	if err != nil {
		log.Printf("Failed to get value from redis: %v\n", err)
		return "", apperrors.NewInternal()
	}

	r.rds.Del(ctx, key)

	return val, nil
}

func (r *redisRepository) SaveInvite(ctx context.Context, guildId string, id string, isPermanent bool) error {

	invite := model.Invite{GuildId: guildId, IsPermanent: isPermanent}

	value, err := json.Marshal(invite)

	if err != nil {
		log.Printf("Error marshalling: %v\n", err.Error())
		return apperrors.NewInternal()
	}

	expiration := 24 * time.Hour
	if isPermanent {
		expiration = 0
	}

	if result := r.rds.Set(ctx, fmt.Sprintf("%s:%s", InviteLinkPrefix, id), value, expiration); result.Err() != nil {
		log.Printf("Failed to set invite link in redis: %v\n", err.Error())
		return apperrors.NewInternal()
	}

	return nil
}

func (r *redisRepository) GetInvite(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("%s:%s", InviteLinkPrefix, token)
	val, err := r.rds.Get(ctx, key).Result()

	if err != nil {
		log.Printf("Failed to get invite link from redis: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	var invite model.Invite
	err = json.Unmarshal([]byte(val), &invite)

	if err != nil {
		log.Printf("Error unmarshalling: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	if !invite.IsPermanent {
		r.rds.Del(ctx, key)
	}

	return invite.GuildId, nil
}

func (r *redisRepository) InvalidateInvites(ctx context.Context, guild *model.Guild) {
	for _, v := range guild.InviteLinks {
		key := fmt.Sprintf("%s:%s", InviteLinkPrefix, v)
		r.rds.Del(ctx, key)
	}
}
