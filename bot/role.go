package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// UserRoles struct with array of user roles in guild
type UserRoles struct {
	Roles []*discordgo.Role
}

// IsAdmin returns true if user is admin
func (ctx *Context) IsAdmin() bool {
	return ctx.User.ID == ctx.Conf.General.AdminID
}

// GetRole returns UserRoles struct pointer
func (ctx *Context) GetRoles() *UserRoles {
	var userRoles = new(UserRoles)
	memb, err := ctx.Discord.GuildMember(ctx.Guild.ID, ctx.User.ID)
	if err != nil {
		fmt.Println("Getting member error: " + err.Error())
	}
	for _, grole := range ctx.Guild.Roles {
		for _, urole := range memb.Roles {
			if grole.ID == urole {
				userRoles.Roles = append(userRoles.Roles, grole)
			}
		}
	}
	return userRoles
}

// ExistsName checks if user role nema exists on user
func (r *UserRoles) ExistsName(name string) bool {
	for _, val := range r.Roles {
		if val.Name == name {
			return true
		}
	}
	return false
}
