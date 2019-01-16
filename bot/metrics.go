package bot

import (
	"bytes"
	"fmt"
	"net/http"
)

func (ctx *Context) MetricsCommand(command string) {
	query := []byte(fmt.Sprintf("commands,server=%v,user=%v command=\"%v\"", ctx.Guild.ID, ctx.Message.Author.ID, command))
	addr := fmt.Sprintf("%v/write?db=%v", ctx.Conf.Metrics.Address, ctx.Conf.Metrics.Database)
	r := bytes.NewReader(query)
	_, _ = http.Post(addr, "", r)
}

func (ctx *Context) MetricsLog(module string) {
	query := []byte(fmt.Sprintf("logs,server=%v module=\"%v\"", ctx.Guild.ID, module))
	addr := fmt.Sprintf("%v/write?db=%v", ctx.Conf.Metrics.Address, ctx.Conf.Metrics.Database)
	r := bytes.NewReader(query)
	_, _ = http.Post(addr, "", r)
}
