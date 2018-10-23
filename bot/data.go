package bot

import "errors"

type DataType struct {
	Polls map[string]*PollType
}

type PollType struct {
	// Fields array of field names
	Fields []string
	// Votes map of votes Key: UserID Value: FieldKey
	Votes map[string]int
}

// NewDataType creates data type
func NewDataType() *DataType {
	var newData = new(DataType)
	newData.Polls = make(map[string]*PollType)
	return newData
}

// CreatePoll creates new poll in guild
func (data *DataType) CreatePoll(ctx *Context, fields []string) error {
	if _, ok := data.Polls[ctx.Guild.ID]; ok {
		return errors.New(ctx.Loc("polls_already_exists"))
	}

	newPoll := new(PollType)
	newPoll.Fields = fields
	newPoll.Votes = make(map[string]int)
	data.Polls[ctx.Guild.ID] = newPoll
	return nil
}

// AddPollVote votes in poll
func (data *DataType) AddPollVote(ctx *Context, vote int) error {
	if _, ok := data.Polls[ctx.Guild.ID]; !ok {
		return errors.New(ctx.Loc("polls_not_exists"))
	}
	if _, ok := data.Polls[ctx.Guild.ID].Votes[ctx.User.ID]; ok {
		return errors.New(ctx.Loc("polls_already_voted"))
	}
	if (vote - 1) >= len(data.Polls[ctx.Guild.ID].Fields) {
		return errors.New(ctx.Loc("polls_wrong_field"))
	}
	data.Polls[ctx.Guild.ID].Votes[ctx.User.ID] = vote - 1
	return nil
}

// EndPoll ends poll and returns results
func (data *DataType) EndPoll(ctx *Context) (results map[string]int, err error) {
	if _, ok := data.Polls[ctx.Guild.ID]; !ok {
		return nil, errors.New(ctx.Loc("polls_not_exists"))
	}
	var newResults = make(map[string]int)
	for _, field := range data.Polls[ctx.Guild.ID].Votes {
		newResults[data.Polls[ctx.Guild.ID].Fields[field]]++
	}
	delete(data.Polls, ctx.Guild.ID)
	return newResults, nil
}
