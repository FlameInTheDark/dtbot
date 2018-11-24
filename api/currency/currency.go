package currency

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/FlameInTheDark/dtbot/bot"
)

// Data currency structure
type Data struct {
	Date         string              `json:"Date"`
	PreviousDate string              `json:"PreviousDate"`
	PreviousURL  string              `json:"PreviousURL"`
	Timestamp    string              `json:"Timestamp"`
	Currencies   map[string]Currency `json:"Valute"`
}

// Currency structure
type Currency struct {
	ID       string
	NumCode  string
	CharCode string
	Nominal  int
	Name     string
	Value    float32
	Previous float32
}

// GetCurrency returns string of parsed currency data
func GetCurrency(ctx *bot.Context) (response string) {
	var (
		newData Data
		args    = ctx.Conf.Currency.Default
	)

	if len(ctx.Args) > 0 {
		args = ctx.Args
	}

	resp, err := http.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		response = fmt.Sprintf("API error: %v", err)
		return
	}

	bbytes, berr := ioutil.ReadAll(resp.Body)
	if berr != nil {
		response = fmt.Sprintf("%v: %v", ctx.Loc("curr_resp_read_err"), berr)
		return
	}

	jerr := json.Unmarshal(bbytes, &newData)
	if jerr != nil {
		response = fmt.Sprintf("%v: %v", ctx.Loc("curr_resp_parse_err"), jerr)
		return
	}

	response = ""

	// List of currencies
	if args[0] == "list" {
		response = fmt.Sprintf("%v: ", ctx.Loc("available_currencies"))
		for key := range newData.Currencies {
			response = fmt.Sprintf("%v %v", response, key)
		}
		return
	}

	var arrow string
	// Current currency
	for _, arg := range args {
		if newData.Currencies[arg].Value > 0 {
			if newData.Currencies[arg].Value > newData.Currencies[arg].Previous {
				arrow = "▲"
			} else {
				arrow = "▼"
			}
			response = fmt.Sprintf("%v%v\n`%v %v = %v RUB %v  %0.2v`\n", response, newData.Currencies[arg].Name, newData.Currencies[arg].Nominal, newData.Currencies[arg].CharCode, newData.Currencies[arg].Value, arrow, newData.Currencies[arg].Value-newData.Currencies[arg].Previous)
		}
	}
	return
}
