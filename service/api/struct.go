package api

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

//the comments
type Comments struct {
	Comment_id int    `json:"comment_id"`
	Comment    string `json:"comment"`
	photo_id   int    `json:"photo_id"`
	User       User   `json:"user"`
}

// the photos

type Photos struct {
	Photo_id int       `json:"photo_Id"`
	Owner    User      `json:"owner"`
	date     Date      `json:"date"`
	likes    []User    `json:"likes"`
	comments []Comment `json:"comments"`
}

// the complete Profile of the user
type User_Profile struct {
	Id         int      `json:"id"`
	Nickname   string   `json:"nickname"`
	Followers  []User   `json:"followers"`
	Followings []User   `json:"following"`
	Photos     []Photos `json:"photos"`
}
