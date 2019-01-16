package bot

import (
	"bytes"
	"fmt"
	"net/http"
)

func (ctx *Context) MetricsCommand(command string) {
	query := []byte(fmt.Sprintf("commands,server=%v command=%v", ctx.Guild.ID, command))
	addr := fmt.Sprintf("%v/write?db=%v", ctx.Conf.Metrics.Address, ctx.Conf.Metrics.Database)
	r := bytes.NewReader(query)
	_, err := http.Post(addr, "", r)
	if err != nil {
		fmt.Println(err.Error())
	}

}
