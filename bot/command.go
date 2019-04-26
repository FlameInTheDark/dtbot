package bot

type (
	// Command : Executable command function
	Command func(Context)

	// CmdMap : Map with executable functions
	CmdMap map[string]Command

	// CommandHandler : Command handler struct
	CommandHandler struct {
		cmds CmdMap
	}
)

// NewCommandHandler creates command handler
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{make(CmdMap)}
}

// GetCmds returns handler commands
func (handler CommandHandler) GetCmds() CmdMap {
	return handler.cmds
}

// Get returns command by command name
func (handler CommandHandler) Get(name string) (*Command, bool) {
	cmd, found := handler.cmds[name]
	return &cmd, found
}

// Register adds new command in handler
func (handler CommandHandler) Register(name string, command Command) {
	handler.cmds[name] = command
	if len(name) > 1 {
		handler.cmds[name[:1]] = command
	}
}


// New command handler
type NodeTree struct {
	Elements map[string]*NodeElement
}

type NodeElement struct {
	Current string
	Elements map[string]*NodeElement
	Workers []NodeWorker
}

// NodeWorker contains workers of command
type NodeWorker struct {
	CommandWorker func(Context)
	Middlewares []func(*Context) bool
}

// CommandSignature contains data of command
type CommandSignature struct {
	Path []string
	Command func(Context)
	Middlewares []func(*Context) bool
}

// NodeTree creates new node tree
func NewTree() NodeTree {
	return NodeTree{Elements:make(map[string]*NodeElement)}
}

func (n *NodeElement) AddElement(element string) {
	if n.Elements == nil {
		n.Elements = make(map[string]*NodeElement)
	}
	n.Elements[element] = &NodeElement{Current:element, Elements:make(map[string]*NodeElement)}
}

func (n *NodeElement) GetElement(element string) *NodeElement {
	if n.Elements == nil {
		return nil
	}
	return n.Elements[element]
}

func (t *NodeTree) Execute(ctx Context) {

}

func (t *NodeTree) GetElement(element string) *NodeElement {
	if t.Elements == nil {
		return nil
	}
	return t.Elements[element]
}

func (t *NodeTree) AddCommand(command *CommandSignature, ) {
	if len(command.Path) == 0 {
		return
	}
	var element = t.GetElement(command.Path[0])
	if len(command.Path) > 1 {
		for _, c := range command.Path[1:] {
			element = element.GetElement(c)
		}
	}
	element.Workers = append(element.Workers, NodeWorker{CommandWorker:command.Command, Middlewares:command.Middlewares})
}

func (w *NodeWorker) CheckMiddlewares(ctx *Context) bool {
	for _, m := range w.Middlewares {
		if m(ctx) == false {
			return false
		}
	}
	return true
}
