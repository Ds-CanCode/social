package funcs

import "time"

type User struct {
	Id        int     `json:"id"`
	Avatar    *string `json:"avatar"`
	Nickname  *string `json:"nickname"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type Comment struct {
    User     User      `json:"user"`
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    PathImg   *string    `json:"pathimg"`
    CreatedAt time.Time `json:"created_at"`
    Likes       int       `json:"likes"`
	DisLikes    int       `json:"dislikes"`
}