package funcs

import (
	"fmt"
	"net/http"
	"strconv"

	pkg "funcs/pkg"
)

func ReadComment(w http.ResponseWriter, r *http.Request) (Comment, int, error) {
	var comment Comment
	var err error
	err = r.ParseMultipartForm(2 << 20) // 2 MB limit for uploaded files
	if err != nil {
		fmt.Println(err)
		return comment, http.StatusBadRequest, pkg.ErrMaxSizeImage
	}
	comment.PostID, _ = strconv.Atoi(r.FormValue("post_id"))
	if comment.PostID == 0 {
		fmt.Println("2")
		return comment, http.StatusBadRequest, pkg.ErrInvalidNamber
	}
	comment.Content = r.FormValue("content")
	if comment.Content == "" || len(comment.Content) > 500 {
		fmt.Println("2")
		return comment, http.StatusBadRequest, pkg.ErrContentLenght
	}

	// Handle the image file upload if available
	imageFile, head, err := r.FormFile("image")
	if err == nil || err.Error() != "http: no such file" {
		if err == nil {
			name, err := pkg.SaveImage(imageFile, head)
			comment.PathImg = &name
			if err != nil {
				fmt.Println("3")
				return comment, http.StatusInternalServerError, pkg.ErrProccessingFile
			}
		} else {
			fmt.Println("4")
			return comment, http.StatusInternalServerError, pkg.ErrInvalidFile
		}
	}
	return comment, 0, nil
}
