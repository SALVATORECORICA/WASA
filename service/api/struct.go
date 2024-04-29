package api

// only the nickname of the user
type Data struct {
	Nickname string `json:"nickname"`
}

// only the id of the user
type DataId struct {
	Id int `json:"id"`
}

// th complete Profile of the user
type User struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
}
