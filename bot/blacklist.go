package bot

// BlackList contains ignored guilds and users
type BlackListStruct struct {
	Guilds []string
	Users []string
}

// CheckGuild returns true if guild in blacklist
func (b BlackListStruct) CheckGuild(id string) bool {
	for _,g := range b.Guilds {
		if g == id {
			return true
		}
	}
	return false
}

// CheckUser returns true if user in blacklist
func (b BlackListStruct) CheckUser(id string) bool {
	for _,g := range b.Users {
		if g == id {
			return true
		}
	}
	return false
}

// BlacklistAddGuild adds guild in blacklist
func (ctx *Context) BlacklistAddGuild(id string) {
	ctx.BlackList.Guilds = append(ctx.BlackList.Guilds, id)
	ctx.DB.AddBlacklistGuild(id)
}

// BlacklistAddUser adds user in blacklist
func (ctx *Context) BlacklistAddUser(id string) {
	ctx.BlackList.Users = append(ctx.BlackList.Users, id)
	ctx.DB.AddBlacklistUser(id)
}

// BlacklistRemoveGuild removes guild from blacklist
func (ctx *Context) BlacklistRemoveGuild(id string) {
	var newArray []string
	for _, g := range ctx.BlackList.Guilds {
		if g != id {
			newArray = append(newArray, g)
		}
	}
	ctx.BlackList.Guilds = newArray
	ctx.DB.RemoveBlacklistGuild(id)
}

// BlacklistRemoveUser removes user from blacklist
func (ctx *Context) BlacklistRemoveUser(id string) {
	var newArray []string
	for _, u := range ctx.BlackList.Users {
		if u != id {
			newArray = append(newArray, u)
		}
	}
	ctx.BlackList.Users = newArray
	ctx.DB.RemoveBlacklistUser(id)
}