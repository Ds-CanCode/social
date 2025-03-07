package funcs

import (
	"net/http"
	"strconv"

	users "funcs/internal/user"
	pkg "funcs/pkg"
)

func GetPostForProfile(userId, profileId int) ([]POST, int, error) {
	var posts []POST
	profileType, err := users.GetProfileType(profileId)
	if err != nil {
		return nil, 500, err
	}
	if profileType && userId != profileId {
		is, err := users.IsFriends(userId, profileId)
		if err != nil {
			return posts, 500, err
		}
		if is == 0 {
			return posts, http.StatusForbidden, pkg.ErrNotFriends
		}
	}
	posts, err = GetUserPostDB(userId, profileId)
	if err != nil {
		return posts, 500, err
	}
	for p := range posts {
		posts[p].Creator, err = GetUserInfo(profileId)
		if err != nil {
			return posts, 500, err
		}
	}

	return posts, 0, nil
}

func ValidatePost(post POST) error {
	if len(post.Title) > 100 || post.Title == "" {
		return pkg.ErrTitlePost
	}
	if len(post.Title) > 1000 || post.Title == "" {
		return pkg.ErrContentPost
	}
	if post.Type != "private" && post.Type != "public" && post.Type != "semi-private" && post.Type != "group" {
		return pkg.ErrTyoePost
	}
	// if post.IsEvent && post.EventDate == nil {
	// 	return pkg.ErrEventTime
	// }

	return nil
}

func ReadPost(w http.ResponseWriter, r *http.Request, id int) (POST, int, error) {
	var post POST
	var err error
	post.Creator.Id = id

	// Parse multipart form data (allows us to handle both files and form fields)
	err = r.ParseMultipartForm(2 << 20) // 2 MB limit for uploaded files
	if err != nil {
		return post, http.StatusBadRequest, pkg.ErrMaxSizeImage
	}

	// Retrieve form values
	post.Title = r.FormValue("title")
	post.Content = r.FormValue("content")
	post.Type = r.FormValue("privacy")
	friends := r.Form["targetedFriends"]
	if post.Type == "semi-private" {
		if len(friends) == 0 {
			return post, http.StatusBadRequest, pkg.ErrBadRequest
		}

		for _, f := range friends {
			var user User
			id, _ := strconv.Atoi(f)
			user.Id = id
			post.AllowedUser = append(post.AllowedUser, user)
		}

		post.AllowedUser = append(post.AllowedUser, post.Creator)

	}

	if post.Type == "group" {
		post.GroupId, _ = strconv.Atoi(r.FormValue("groupId"))
		if post.GroupId == 0 {
			return post, http.StatusBadRequest, pkg.ErrInvalidNamber
		}
		/************************/
		/***check if following***/
		/************************/
	}

	// Handle the image file upload if available
	imageFile, head, err := r.FormFile("image")
	if err == nil || err.Error() != "http: no such file" {
		if err == nil {
			post.Path, err = pkg.SaveImage(imageFile, head)
			if err != nil {
				return post, http.StatusInternalServerError, pkg.ErrProccessingFile
			}
		} else {
			return post, http.StatusInternalServerError, pkg.ErrInvalidFile
		}
	}

	err = ValidatePost(post)
	if err != nil {
		return post, http.StatusBadRequest, err
	}
	return post, 0o0, nil
}
