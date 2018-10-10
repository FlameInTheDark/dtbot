package currency

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../bot"
)

// Data currency structure
type Data struct {
	Date         string            `json:"Date"`
	PreviousDate string            `json:"PreviousDate"`
	PreviousURL  string            `json:"PreviousURL"`
	Timestamp    string            `json:"Timestamp"`
	Valutes      map[string]Valute `json:"Valute"`
}

// Valute valute structure
type Valute struct {
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
	fmt.Println("Current val:")
	resp, err := http.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		response = fmt.Sprintf("API error: %v", err)
		return
	}

	bbytes, berr := ioutil.ReadAll(resp.Body)
	if berr != nil {
		response = fmt.Sprintf("Response read error: %v", berr)
		return
	}

	jerr := json.Unmarshal(bbytes, &newData)
	if jerr != nil {
		response = fmt.Sprintf("Response parse error: %v", jerr)
		return
	}
	var arrow string
	response = ""
	for _, arg := range args {
		if newData.Valutes[arg].Value > 0 {
			if newData.Valutes[arg].Value > newData.Valutes[arg].Previous {
				arrow = "▲"
			} else {
				arrow = "▼"
			}
			response = fmt.Sprintf("%v%v\n`%v %v  %0.2v`\n", response, newData.Valutes[arg].Name, newData.Valutes[arg].Value, arrow, newData.Valutes[arg].Value-newData.Valutes[arg].Previous)
		}
	}
	return
}
