package funcs

import (
	dataBase "funcs/internal/database"
)

func ReactPost(rectPost reactionPost) error {
	db := dataBase.GetDb()
	var currentRect string
	var err error
	query := "SELECT reaction_type FROM reactionsPosts WHERE post_id = ? AND user_id = ?"
	db.QueryRow(query, rectPost.PostID, rectPost.UserID).Scan(&currentRect)

	// Update the reaction
	if currentRect == "" {
		_, err = db.Exec(
			"INSERT INTO reactionsPosts (post_id, user_id, reaction_type) VALUES (?, ?, ?)",
			rectPost.PostID, rectPost.UserID, rectPost.ReactionType,
		)
	} else if currentRect != rectPost.ReactionType {
		_, err = db.Exec(
			"UPDATE reactionsPosts SET reaction_type = ? WHERE post_id = ? AND user_id = ?",
			rectPost.ReactionType, rectPost.PostID, rectPost.UserID,
		)
	} else if currentRect == rectPost.ReactionType {
		_, err = db.Exec(
			"DELETE FROM reactionsPosts WHERE post_id = ? AND user_id = ?", rectPost.PostID, rectPost.UserID,
		)
	}
	if err != nil {
		return err
	}
	return nil
}

func ReactComment(rectComment reactionComment) error {
	db := dataBase.GetDb()
	var currentRect string
	var err error
	query := "SELECT reaction_type FROM reactionsComments WHERE comment_id = ? AND user_id = ?"
	db.QueryRow(query, rectComment.CommentID, rectComment.UserID).Scan(&currentRect)

	// Update the reaction
	if currentRect == "" {
		_, err = db.Exec(
			"INSERT INTO reactionsComments (comment_id, user_id, reaction_type) VALUES (?, ?, ?)",
			rectComment.CommentID, rectComment.UserID, rectComment.ReactionType,
		)
	} else if currentRect != rectComment.ReactionType {
		_, err = db.Exec(
			"UPDATE reactionsComments SET reaction_type = ? WHERE comment_id = ? AND user_id = ?",
			rectComment.ReactionType, rectComment.CommentID, rectComment.UserID,
		)
	} else if currentRect == rectComment.ReactionType {
		_, err = db.Exec(
			"DELETE FROM reactionsComments WHERE comment_id = ? AND user_id = ?", rectComment.CommentID, rectComment.UserID,
		)
	}
	if err != nil {
		return err
	}
	return nil
}
