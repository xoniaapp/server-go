package service

import (
	"log"

	"github.com/aelpxy/xoniaapp/model/apperrors"
	"github.com/bwmarrin/snowflake"
)

func GenerateId() (string, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Printf("Failed to genenerate an snowflake id: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}
	id := node.Generate()
	return id.String(), nil
}
