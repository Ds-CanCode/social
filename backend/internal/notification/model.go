package funcs

type SGroup struct {
	Id    int
	Title string
	Path  string
}

type SUser struct {
	Id        int
	FirstName string
	LastName  string
	Nickname  string
	Path      string
}

type NOTIF struct {
	Id        int
	Sender    SUser
	Group     SGroup
	IsRead    int
	Accepted  bool
	Type      string
	CreatedAt string
	Title     string
}
