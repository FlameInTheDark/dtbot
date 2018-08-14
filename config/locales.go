package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type LocalesMap map[string]map[string]string

func (l LocalesMap) Get(key string) string {
	return l[General.Language][key]
}

// Loading locales from file
func LoadLocales() {
	file, e := ioutil.ReadFile("./locales.json")
	if e != nil {
		fmt.Printf("Locale file error: %v\n", e)
		os.Exit(1)
	}

	err := json.Unmarshal(file, &Locales)
	if err != nil {
		panic(err)
	}

	if _, ok := Locales[General.Language]; ok {
		return
	} else {
		fmt.Printf("Locale file not contain language \"%v\"\n", General.Language)
		os.Exit(1)
	}
}
