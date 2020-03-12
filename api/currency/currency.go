package currency

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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

// CurrencyCheck returns true if currency is real
func (d *Data) CurrencyCheck(currency string) bool {
	if currency == "RUB" {
		return true
	}

	if d.Currencies[currency].Value > 0 {
		return true
	}

	return false
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

	newData.Currencies["RUB"] = Currency{"R00000", "000", "RUB", 1, "Российский Рубль", 1, 1}

	response = ""

	// List of currencies
	if args[0] == "list" {
		ctx.MetricsCommand("currency", "list")
		response = fmt.Sprintf("%v: ", ctx.Loc("available_currencies"))
		for key := range newData.Currencies {
			response = fmt.Sprintf("%v %v", response, key)
		}
		return
	}

	// TODO: i should complete currency converter
	// Converting currencies
	if len(args) > 3 && args[0] == "conv" {
		ctx.MetricsCommand("currency", "conv")
		count, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			response = fmt.Sprintf("%v: %v", ctx.Loc("error"), ctx.Loc("nan"))
			return
		}
		if c1, ok1 := newData.Currencies[strings.ToUpper(args[1])]; ok1 {
			if c2, ok2 := newData.Currencies[strings.ToUpper(args[2])]; ok2 {
				cur1 := c1.Value / float32(c1.Nominal)
				cur2 := c2.Value / float32(c2.Nominal)
				//cur1Delta := 1 / cur1
				cur2Delta := 1 / cur2
				res := cur2Delta * (cur1 * float32(count))
				response = fmt.Sprintf("`%v %v = %0.2f %v`\n", args[3], args[1], res, args[2])
			}
		}
		return
	}
	ctx.MetricsCommand("currency", "main")
	var arrow string
	// Current currency
	for _, arg := range args {
		if newData.Currencies[arg].Value > 0 {
			if newData.Currencies[arg].Value > newData.Currencies[arg].Previous {
				arrow = "▲"
			} else {
				arrow = "▼"
			}
			response = fmt.Sprintf("%v%v\n`%v %v = %v RUB %v  %0.2v`\n",
				response,
				newData.Currencies[arg].Name,
				newData.Currencies[arg].Nominal,
				newData.Currencies[arg].CharCode,
				newData.Currencies[arg].Value,
				arrow,
				newData.Currencies[arg].Value-newData.Currencies[arg].Previous)
		}
	}
	return
}
