package funcs

type ERROR struct {
	Err string
	Status int
}


type reactionCount struct {
	LikeCount int `json:"like_count"`
	DislikeCount int `json:"dislike_count"`
}