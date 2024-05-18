package structures

import (
	"time"
)

// only the nickname of the user
type UserNickname struct {
	Nickname string `json:"nickname"`
}

// only the id of the user
type UserId struct {
	Id int `json:"id"`
}

// the Profile of the user
type User struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
}

// the comments
type Comment struct {
	Comment_id int    `json:"comment_id"`
	Comment    string `json:"comment"`
	Photo_id   int    `json:"photo_id"`
	User       User   `json:"user"`
}

// the photos

type Photos struct {
	PhotoId  int       `json:"photo_Id"`
	Owner    User      `json:"owner"`
	Date     time.Time `json:"date"`
	Likes    []User    `json:"likes"`
	Comments []Comment `json:"comments"`
}

// the complete Profile of the user
type UserProfile struct {
	Id         int      `json:"id"`
	Nickname   string   `json:"nickname"`
	Followers  []User   `json:"followers"`
	Followings []User   `json:"following"`
	Photos     []Photos `json:"photos"`
}
type Image struct {
	PhotoData string `json:"photo_data"`
}
