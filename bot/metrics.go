package bot

import (
	"bytes"
	"fmt"
	"net/http"
)

// MetricsCommand sends command metrics
func (ctx *Context) MetricsCommand(command string) {
	query := []byte(fmt.Sprintf("commands,server=%v,user=%v command=\"%v\"", ctx.Guild.ID, ctx.Message.Author.ID, command))
	addr := fmt.Sprintf("%v/write?db=%v&u=%v&p=%v",
		ctx.Conf.Metrics.Address, ctx.Conf.Metrics.Database, ctx.Conf.Metrics.User, ctx.Conf.Metrics.Password)
	r := bytes.NewReader(query)
	_, _ = http.Post(addr, "", r)
}

// MetricsLog sends log metrics
func (ctx *Context) MetricsLog(module string) {
	query := []byte(fmt.Sprintf("logs,server=%v module=\"%v\"", ctx.Guild.ID, module))
	addr := fmt.Sprintf("%v/write?db=%v&u=%v&p=%v",
		ctx.Conf.Metrics.Address, ctx.Conf.Metrics.Database, ctx.Conf.Metrics.User, ctx.Conf.Metrics.Password)
	r := bytes.NewReader(query)
	_, _ = http.Post(addr, "", r)
}

// MetricsMessage sends message metrics
func (ctx *Context) MetricsMessage() {
	query := []byte(fmt.Sprintf("messages,server=%v user=\"%v\"", ctx.Guild.ID, ctx.Message.Author.ID))
	addr := fmt.Sprintf("%v/write?db=%v&u=%v&p=%v",
		ctx.Conf.Metrics.Address, ctx.Conf.Metrics.Database, ctx.Conf.Metrics.User, ctx.Conf.Metrics.Password)
	r := bytes.NewReader(query)
	_, _ = http.Post(addr, "", r)
}
