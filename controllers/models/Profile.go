package models

import (
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type Profile struct {
	Username   string     `json:"username" example:"devashish"`
	Experience Experience `json:"experience"`
	Species    string     `json:"species" example:"KLAYANS"`
	Clan       *Clan      `json:"clan,omitempty"`
}

type Clan struct {
	Name string `json:"name,omitempty" example:"Mind Krackers"`
	Role string `json:"role,omitempty" example:"MEMBER"`
}

type Experience struct {
	Level    int `json:"level" example:"4"`
	Current  int `json:"current" example:"185"`
	Required int `json:"required" example:"368"`
}

func (p *Profile) Init(userData models.UserData, clanData *models.ClanData, experienceConstants constants.ExperienceConstants) {
	p.Username = userData.Profile.Username
	p.Species = userData.Profile.Species
	p.Experience.Init(userData.Profile, experienceConstants)
	if len(userData.Profile.ClanId) > 0 && clanData != nil {
		p.Clan = &Clan{}
		p.Clan.Init(userData, *clanData)
	}
}

func (e *Experience) Init(profileUser models.ProfileUser, userExperienceConstants constants.ExperienceConstants) {
	var nextLevelString string
	for key, experienceRequired := range userExperienceConstants.ExperiencesRequired {
		if experienceRequired.ExperienceRequired > profileUser.Experience {
			nextLevelString = key
		}
	}
	nextLevel, _ := strconv.Atoi(nextLevelString)
	e.Level = nextLevel - 1
	e.Current = profileUser.Experience
	e.Required = userExperienceConstants.ExperiencesRequired[nextLevelString].ExperienceRequired
}

func (c *Clan) Init(userData models.UserData, clan models.ClanData) {
	c.Name = clan.Name
	for _, clanMember := range clan.Members {
		if clanMember.MemberId == userData.Id {
			c.Role = clanMember.Role
			break
		}
	}
}
