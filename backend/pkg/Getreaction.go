package funcs

import (

	dataBase "funcs/internal/database"
)

func CountRect(id int) (reactionCount, error) {
	query := `SELECT 
    COALESCE(SUM(CASE WHEN reaction_type = 'LIKE' THEN 1 ELSE 0 END), 0) AS like_count,
    COALESCE(SUM(CASE WHEN reaction_type = 'DISLIKE' THEN 1 ELSE 0 END), 0) AS dislike_count
FROM reactionsPosts
WHERE post_id = ?;


`
	db := dataBase.GetDb()
	var ReactionCount reactionCount
	err := db.QueryRow(query, id).Scan(&ReactionCount.LikeCount, &ReactionCount.DislikeCount)
	if err != nil {
		return ReactionCount, err
	}

	return ReactionCount, nil
}

func CountCommentRect(commentId int) (reactionCount, error) {
	query := `SELECT 
    COALESCE(SUM(CASE WHEN reaction_type = 'LIKE' THEN 1 ELSE 0 END), 0) AS like_count,
    COALESCE(SUM(CASE WHEN reaction_type = 'DISLIKE' THEN 1 ELSE 0 END), 0) AS dislike_count
FROM reactionsComments
WHERE comment_id = ?;

`

	db := dataBase.GetDb()
	var ReactionCount reactionCount
	err := db.QueryRow(query, commentId).Scan(&ReactionCount.LikeCount, &ReactionCount.DislikeCount)
	if err != nil {
		return ReactionCount, err
	}
	
	return ReactionCount, nil
}
