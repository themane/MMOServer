package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

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

func (p *Profile) Init(userData models.UserData, clanData *models.ClanData, experienceConstants constants.ExperienceConstants) {
	p.Username = userData.Profile.Username
	p.Experience.Init(userData.Profile, experienceConstants)
	if len(userData.Profile.ClanId) > 0 && clanData != nil {
		p.Clan = &Clan{}
		p.Clan.Init(userData, *clanData)
	}
}

func (e *Experience) Init(profileUser models.ProfileUser, experienceConstants constants.ExperienceConstants) {
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

func (c *Clan) Init(userData models.UserData, clan models.ClanData) {
	c.Name = clan.Name
	for _, clanMember := range clan.Members {
		if clanMember.Id == userData.Id {
			c.Role = clanMember.Role
			break
		}
	}
}
