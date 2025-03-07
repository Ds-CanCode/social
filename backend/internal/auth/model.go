package funcs

import "time"

type User struct {
    ID          int       `json:"id"`
    Email       string    `json:"email"`
    Password    string    `json:"password"`
    FirstName   string    `json:"firstName"`
    LastName    string    `json:"lastName"`
    DateOfBirth string    `json:"dateOfBirth"`  // ✅ Renamed for consistency
    Avatar      string   `json:"avatar"`
    Nickname    string   `json:"nickname"`
    AboutMe     string   `json:"aboutMe"`
    ProfileType bool     `json:"profileType"`
    CreatedAt   time.Time `json:"createdAt"`
}


