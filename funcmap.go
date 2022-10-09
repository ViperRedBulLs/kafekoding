package main

import (
	"bytes"
	"kafekoding/models"
	"log"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func markdownRender(text string) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(text), &buf); err != nil {
		log.Println(err)
		return ""
	}
	return buf.String()
}

func getUserByID(id uint) models.User {
	user, err := models.NewUserModel().GetUserByID(id)
	if err != nil {
		return models.User{}
	}

	return user
}

func isJoinedRoom(username string, members []*models.User) bool {
	var joined bool = false
	for _, member := range members {
		if username == member.Username {
			joined = true
		}
	}

	return joined
}

func isUserByID(userID, targetID uint) bool {
	var joined bool = false
	if userID == targetID {
		joined = true
	} else {
		joined = false
	}
	return joined
}

func countUserMembers(v []*models.User) int {
	var count int
	for range v {
		count += 1
	}
	return count
}

func getUserFullNameByID(id uint) string {
	user, _ := models.NewUserModel().GetUserByID(id)
	return user.FirstName + " " + user.LastName
}

func timeKitchen(v time.Time) string {
	return v.Format(time.Kitchen)
}
