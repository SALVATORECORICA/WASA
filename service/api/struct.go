package api

import (
	"time"
)

// only the nickname of the user
type Data struct {
	Nickname string `json:"nickname"`
}

// only the id of the user
type DataId struct {
	Id int `json:"id"`
}

// the Profile of the user
type User struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
}

// the comments
type Comments struct {
	Comment_id int    `json:"comment_id"`
	Comment    string `json:"comment"`
	photo_id   int    `json:"photo_id"`
	User       User   `json:"user"`
}

// the photos

type Photos struct {
	Photo_id int        `json:"photo_Id"`
	Owner    User       `json:"owner"`
	date     time.Time  `json:"date"`
	likes    []User     `json:"likes"`
	comments []Comments `json:"comments"`
}

// the complete Profile of the user
type User_Profile struct {
	Id         int      `json:"id"`
	Nickname   string   `json:"nickname"`
	Followers  []User   `json:"followers"`
	Followings []User   `json:"following"`
	Photos     []Photos `json:"photos"`
}

type Image struct {
	Photo_data image64 `json:"photo_Data"`
}

// the likes
type Like struct {
	id_photo int `json:"id_Photo"`
	id_user  int `json:"id_User"`
}

// the comment
type Comment struct {
	id_comment int    `json:"id_comment"`
	id_photo   int    `json:"id_photo"`
	id_user    int    `json:"id_user"`
	comment    string `json:"comment"`
}
