package model

import (
	"context"
	"mime/multipart"
)

type FileRepository interface {
	UploadAvatar(header *multipart.FileHeader, directory string) (string, error)
	UploadFile(header *multipart.FileHeader, directory, filename, mimetype string) (string, error)
	DeleteImage(key string) error
}

type MailRepository interface {
	SendResetMail(email string, html string) error
}

type RedisRepository interface {
	SetResetToken(ctx context.Context, id string) (string, error)
	GetIdFromToken(ctx context.Context, token string) (string, error)
	SaveInvite(ctx context.Context, guildId string, id string, isPermanent bool) error
	GetInvite(ctx context.Context, token string) (string, error)
	InvalidateInvites(ctx context.Context, guild *Guild)
}
