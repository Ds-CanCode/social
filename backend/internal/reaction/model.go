package funcs

type reactionPost struct {
	PostID       int `json:"post_id"`
	UserID		 int `json:"user_id"`
	ReactionType string `json:"reaction_type"`
}

type reactionComment struct {
	CommentID       int `json:"comment_id"`
	UserID		 int `json:"user_id"`
	ReactionType string `json:"reaction_type"`
}


