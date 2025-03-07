package funcs

type Followers struct {
	Id        int
	FirstName string
	LastName  string
	Avatar    string
}
type Profile struct {
	Id        int     `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	NickName  *string `json:"nickName"`
	Datebirth string  `json:"datebirth"`
	Avatar    *string `json:"avatar"`
	Aboutme   *string `json:"aboutme"`
	CreatedAt string  `json:"createdAt"`
	Type      *bool   `json:"type"`
	Followers int     `json:"followers"`
	NbrPosts  int     `json:"nbrPosts"`
	IsFollow  *bool     `json:"isfollowing"`
}
