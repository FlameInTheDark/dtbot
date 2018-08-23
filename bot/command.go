package bot

type (
    
    // Command: Executable command function
	Command func(Context)
    
    // CmdMap: Map with executable functions
	CmdMap map[string]Command

    // CommandHandler: Command handler struct
	CommandHandler struct {
		cmds CmdMap
	}
)

// Create command handler
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{make(CmdMap)}
}

// Return handler commands
func (handler CommandHandler) GetCmds() CmdMap {
	return handler.cmds
}

// Return handler command by command name
func (handler CommandHandler) Get(name string) (*Command, bool) {
	cmd, found := handler.cmds[name]
	return &cmd, found
}

// Register new command in handler
func (handler CommandHandler) Register(name string, command Command) {
	handler.cmds[name] = command
	if len(name) > 1 {
		handler.cmds[name[:1]] = command
	}
}
