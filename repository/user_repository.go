package repository

import (
	"database/sql"
	"errors"
	"log"
	"regexp"

	"github.com/aelpxy/xoniaapp/model"
	"github.com/aelpxy/xoniaapp/model/apperrors"
	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	user := &model.User{}

	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("uid", id)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

func (r *userRepository) Create(user *model.User) (*model.User, error) {
	if result := r.DB.Create(&user); result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return nil, apperrors.NewBadRequest(apperrors.DuplicateEmail)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", user.Email, result.Error)
		return nil, apperrors.NewInternal()
	}

	return user, nil
}


func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}

	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("email", email)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}


func (r *userRepository) Update(user *model.User) error {
	return r.DB.Save(&user).Error
}


func (r *userRepository) GetFriendAndGuildIds(userId string) (*[]string, error) {
	var ids []string
	result := r.DB.Raw(`
          SELECT g.id
          FROM guilds g
          JOIN members m on m.guild_id = g."id"
          where m.user_id = @userId
          UNION
          SELECT "User__friends"."id"
          FROM "users" "User" LEFT JOIN "friends" "User_User__friends" ON "User_User__friends"."user_id"="User"."id" LEFT
              JOIN "users" "User__friends" ON "User__friends"."id"="User_User__friends"."friend_id"
          WHERE ( "User"."id" = @userId )
	`, sql.Named("userId", userId)).Find(&ids)

	return &ids, result.Error
}

func (r *userRepository) GetRequestCount(userId string) (*int64, error) {
	var count int64
	err := r.DB.
		Table("users").
		Joins("JOIN friend_requests fr ON users.id = fr.sender_id").
		Where("fr.receiver_id = ?", userId).
		Count(&count).
		Error

	return &count, err
}

func isDuplicateKeyError(err error) bool {
	duplicate := regexp.MustCompile(`\(SQLSTATE 23505\)$`)
	return duplicate.MatchString(err.Error())
}
