package models

type MembershipRole string

const (
	LEADER MembershipRole = "LEADER"
	SUB_LEADER
	MEMBER
)

type Profile struct {
	Username   string     `json:"username"`
	Experience Experience `json:"experience"`
	Clan       Clan       `json:"clan"`
}

type Clan struct {
	Name string         `json:"name"`
	Role MembershipRole `json:"role"`
}

type Experience struct {
	Level    int `json:"level"`
	Current  int `json:"current"`
	Required int `json:"required"`
}
