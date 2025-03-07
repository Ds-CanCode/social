package funcs

import "time"

type POST struct {
	Creator     User      `json:"creator"`
	Id          int       `json:"id"`
	GroupId     int      `json:"group_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Path        string    `json:"image"`
	Type        string    `json:"privacy"`
	CreateAt    time.Time `json:"created_at"`
	Likes       int       `json:"likes"`
	DisLikes    int       `json:"dislikes"`
	IsLiked     int       `json:"isliked"`
	AllowedUser []User    `json:"alloweduser"`
}
type User struct {
	Id        int     `json:"id"`
	Avatar    *string `json:"avatar"`
	Nickname  *string `json:"nickname"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}
