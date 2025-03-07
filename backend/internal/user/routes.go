package funcs

import "net/http"

func Rootes(mux *http.ServeMux) {
	mux.HandleFunc("/api/users/info", GetUserInfo)
	mux.HandleFunc("/api/users/profile", getProfile)
	mux.HandleFunc("/api/users/request", request)
	mux.HandleFunc("/api/users/handlerequest",handlerequest)
	mux.HandleFunc("/api/users/followers", GetFollowers)	
	mux.HandleFunc("/api/users/changeprofiletype",ChangeProfileType)
	mux.HandleFunc("/api/users/userfollowers", GetUserFollowers)
	mux.HandleFunc("/api/users/userfollowing", GetUserFollowing)
	mux.HandleFunc("/api/users/nowRequest", NowRequest)
}
