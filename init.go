package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aelpxy/xoniaapp/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dataSources struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	S3Session   *session.Session
}

func initDS() (*dataSources, error) {
	log.Printf("Starting Server...\n")
	dbUrl := os.Getenv("DATABASE_URL")

	log.Printf("Connecting to PostgreSQL... \n")
	db, err := gorm.Open(postgres.Open(dbUrl))

	if err != nil {
		return nil, fmt.Errorf("error opening Database: %w", err)
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
		return nil, fmt.Errorf("error migrating models: %w", err)
	}

	if err := db.SetupJoinTable(&model.Guild{}, "Members", &model.Member{}); err != nil {
		return nil, fmt.Errorf("error creating join table: %w", err)
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
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_S3_REGION")
	// endpointURL := os.Getenv("AWS_S3_URL")

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
			Region:           aws.String(region),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error with S3 session: %w", err)
	}

	return &dataSources{
		DB:          db,
		RedisClient: rdb,
		S3Session:   sess,
	}, nil
}

func (d *dataSources) close() error {
	if err := d.RedisClient.Close(); err != nil {
		return fmt.Errorf("error connecting with Redis Client: %w", err)
	}

	return nil
}
