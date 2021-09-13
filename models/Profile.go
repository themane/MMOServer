package models

import "strconv"

const (
	LEADER     string = "LEADER"
	SUB_LEADER string = "SUB_LEADER"
	MEMBER     string = "MEMBER"
)

type Profile struct {
	Username   string     `json:"username"`
	Experience Experience `json:"experience"`
	Clan       *Clan      `json:"clan,omitempty"`
}

type Clan struct {
	Name string `json:"name,omitempty"`
	Role string `json:"role,omitempty"`
}

type Experience struct {
	Level    int `json:"level"`
	Current  int `json:"current"`
	Required int `json:"required"`
}

func (p *Profile) Init(profileUser ProfileUser, clan ClanData, experienceConstants ExperienceConstants) {
	p.Username = profileUser.Username
	p.Experience.Init(profileUser, experienceConstants)
	p.Clan = &Clan{Name: clan.Name, Role: getRole(profileUser, clan)}
}

func (e *Experience) Init(profileUser ProfileUser, experienceConstants ExperienceConstants) {
	var nextLevelString string
	for key, experienceRequired := range experienceConstants.User.ExperiencesRequired {
		if experienceRequired.ExperienceRequired > profileUser.Experience {
			nextLevelString = key
		}
	}
	nextLevel, _ := strconv.Atoi(nextLevelString)
	e.Level = nextLevel - 1
	e.Current = profileUser.Experience
	e.Required = experienceConstants.User.ExperiencesRequired[nextLevelString].ExperienceRequired
}

func getRole(profileUser ProfileUser, clan ClanData) string {
	if len(profileUser.ClanId) > 0 {
		for _, clanMember := range clan.Members {
			if clanMember.Id == profileUser.Id {
				return clanMember.Role
			}
		}
	}
	return ""
}
