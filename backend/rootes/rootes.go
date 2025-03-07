package funcs

import (
	"net/http"

	group "funcs/internal/Groups"
	auth "funcs/internal/auth"
	comment "funcs/internal/comment"
	imageserv "funcs/internal/fileserver"
	notif "funcs/internal/notification"
	Post "funcs/internal/post"
	reaction "funcs/internal/reaction"
	search "funcs/internal/search"
	users "funcs/internal/user"
	ws "funcs/internal/chat"
)

func Rootes() *http.ServeMux {
	mux := http.NewServeMux()
	comment.Rootes(mux)
	Post.Rootes(mux)
	group.Rootes(mux)
	users.Rootes(mux)
	auth.Rootes(mux)
	reaction.Rootes(mux)
	imageserv.Rootes(mux)
	search.Rootes(mux)
	notif.Rootes(mux)
	ws.Rootes(mux)

	return mux
}
