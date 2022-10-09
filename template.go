package main

import (
	"html/template"

	"github.com/gin-contrib/multitemplate"
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	r.AddFromFilesFuncs(
		"login",
		template.FuncMap{},
		"views/users/base.html", "views/users/login.html")

	r.AddFromFilesFuncs("register",
		template.FuncMap{},
		"views/users/base.html", "views/users/register.html")

	r.AddFromFilesFuncs("index",
		template.FuncMap{},
		"views/rooms/base.html", "views/rooms/index.html")

	r.AddFromFilesFuncs("room_detail",
		template.FuncMap{
			"markdown":     markdownRender,
			"getUserByID":  getUserByID,
			"isJoinedRoom": isJoinedRoom,
		},
		"views/rooms/base.html", "views/rooms/detail.html")

	// chats
	r.AddFromFilesFuncs("chat_index",
		template.FuncMap{
			"countUserMembers": countUserMembers,
			"isJoinedRoom":     isJoinedRoom,
		},
		"views/chats/base.html", "views/chats/index.html")
	r.AddFromFilesFuncs("chat_detail",
		template.FuncMap{
			"isUserByID":          isUserByID,
			"getUserByID":         getUserByID,
			"getUserFullNameByID": getUserFullNameByID,
			"timeKitchen":         timeKitchen,
			"countUserMembers":    countUserMembers,
		},
		"views/chats/base.html", "views/chats/detail.html")

	// Admin
	r.AddFromFilesFuncs("admin_login", template.FuncMap{}, "views/admin/login.html")
	r.AddFromFilesFuncs("admin_index", template.FuncMap{}, "views/admin/base.html", "views/admin/index.html")
	r.AddFromFilesFuncs("admin_detail", template.FuncMap{}, "views/admin/base.html", "views/admin/detail.html")
	r.AddFromFilesFuncs("admin_db_users", template.FuncMap{}, "views/admin/base.html", "views/admin/users.html")
	r.AddFromFilesFuncs("admin_db_rooms", template.FuncMap{}, "views/admin/base.html", "views/admin/rooms.html")

	return r
}
