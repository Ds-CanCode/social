package funcs

type Group struct {
	Id           int
	CreatorId    int    `json:"CreatorId"`
	Title        string `json:"title"`
	Descreption  string `json:"Descreption"`
	NbrMembers   int    `json:"nbr"`
	Path         string `json:"image"`
	CreatedAt    string `json:"createAt"`
	MemberStatus STATUS `json:"memberstatus"`
}

type STATUS struct {
	Status string `json:"status"`
	Since  string `json:"since"`
	Sender SUser  `json:"sender"`
}

type SGroup struct {
	Id          int
	Title       string
	Path        string
	MemberCount int
}

type SUser struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Nickname  string `json:"nickname"`
	Path      string `json:"path"`
}

type Event struct {
    Id          int    `json:"id"`
    Title       string `json:"title"`
    GroupId     int    `json:"groupid"`
    CreatorId   int    `json:"creatorId"`
    User        SUser   `json:"user"`
    Descreption string `json:"Descreption"`
    Time        string `json:"eventtime"`
    Status      int    `json:"status"`
	CountAttends int    `json:"countattends"`
}

var (
	StatusNotMember     string = "notmember"
	StatusMember        string = "member"
	StatusSendRequest   string = "sendReaquest"
	StatusReciverequest string = "Invited"
	StatusCreator       string = "creator"
	StatusErr           string = "can't find status"
)
