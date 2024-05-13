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
	Photo_id   int    `json:"photo_id"`
	User       User   `json:"id_user"`
}

// the photos

type Photo struct {
	Photo_id int        `json:"photo_Id"`
	Owner    User       `json:"owner"`
	Date     time.Time  `json:"date"`
	Likes    []User     `json:"likes"`
	Comments []Comments `json:"comments"`
	nLikes   int        `json:"nLIkes"`
	Image    Image      `json:"image"`
}

// the complete Profile of the user
type User_Profile struct {
	Id         int     `json:"id"`
	Nickname   string  `json:"nickname"`
	Followers  []User  `json:"followers"`
	Followed   []User  `json:"following"`
	Photos     []Image `json:"photos"`
	NFollowers int     `json:"nFollowers"`
	NFollowed  int     `json:"nFollowing"`
}

type Image struct {
	Photo_data image64 `json:"photo_Data"`
}

// the likes
type Like struct {
	IdPhoto int `json:"id_Photo"`
	IdUser  int `json:"id_User"`
}

type PhotosProfile struct {
	PhotoId   int     `json:"id_Photo"`
	PhotoData image64 `json:"photo_Data"`
}
