//                  GNU GENERAL PUBLIC LICENSE
//                      Version 3, 29 June 2007
//
// Copyright (C) 2007 Free Software Foundation, Inc. <https://fsf.org/>
// Everyone is permitted to copy and distribute verbatim copies
// of this license document, but changing it is not allowed.
//

package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-redis/redis/v8"
	"github.com/aelpxy/xoniaapp/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type dataSources struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	S3Session   *session.Session
}

func initDS() (*dataSources, error) {
	log.Printf("Starting Server...\n")
	dbUrl := os.Getenv("DATABASE_URL")

	log.Printf("Connecting to PostgreSQL\n")
	db, err := gorm.Open(postgres.Open(dbUrl))

	if err != nil {
		return nil, fmt.Errorf("Error opening Database: %w", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Guild{},
		&model.Member{},
		&model.Channel{},
		&model.DMMember{},
		&model.Message{},
		&model.Attachment{},
	); err != nil {
		return nil, fmt.Errorf("Error migrating models: %w", err)
	}

	if err := db.SetupJoinTable(&model.Guild{}, "Members", &model.Member{}); err != nil {
		return nil, fmt.Errorf("Error creating join table: %w", err)
	}

	redisURL := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	log.Printf("Connecting to Redis\n")
	rdb := redis.NewClient(opt)

	_, err = rdb.Ping(context.Background()).Result()

	if err != nil {
		return nil, fmt.Errorf("Failed to connect with Redis: %w", err)
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_S3_REGION")
	// TODO: Make the AWS Endpoint a env var
	sess, err := session.NewSession(
		&aws.Config{
			Credentials: credentials.NewStaticCredentials(
				accessKey,
				secretKey,
				"",
			),
			Endpoint:         aws.String("https://spaces-s3.services.xoniaapp.com/"),
			DisableSSL:       aws.Bool(false),
			S3ForcePathStyle: aws.Bool(true),
			Region: aws.String(region),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Error making S3 session: %w", err)
	}

	return &dataSources{
		DB:          db,
		RedisClient: rdb,
		S3Session:   sess,
	}, nil
}

func (d *dataSources) close() error {
	if err := d.RedisClient.Close(); err != nil {
		return fmt.Errorf("Error with Redis Client: %w", err)
	}

	return nil
}
// data
