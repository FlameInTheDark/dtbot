package bot

type BotMessages struct {
	Messages map[string][]string
}

// NewMessagesMap creates map of bot's messages
func NewMessagesMap() *BotMessages {
	return &BotMessages{Messages: make(map[string][]string)}
}

func (m *BotMessages) Add(ctx *Context, messageID string) {
	m.Messages[ctx.Message.ChannelID] = append(m.Messages[ctx.Message.ChannelID], messageID)
	if len(m.Messages[ctx.Message.ChannelID]) > ctx.Conf.General.MessagePool {
		m.Messages[ctx.Message.ChannelID] = m.Messages[ctx.Message.ChannelID][1:]
	}
}

func (m *BotMessages) Clear(ctx *Context, from int) {
	channelID := ctx.Message.ChannelID
	if len(m.Messages[channelID][:(len(m.Messages[channelID])-1)-from]) > 0 {
		ctx.Discord.ChannelMessagesBulkDelete(ctx.Message.ChannelID, m.Messages[channelID][:(len(m.Messages[channelID])-1)-from])
		m.Messages[channelID] = m.Messages[channelID][(len(m.Messages[channelID])-1)-from:]
	}
}
