package config

import (
    "fmt"
    "os"
    "encoding/json"
    "io/ioutil"
)

type LocalesMap map[string]map[string]string

func (l LocalesMap) Get(key string) string {
    return l[Language][key]
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
    
    if _, ok := Locales[Language]; ok {
        return
    } else {
        fmt.Printf("Locale file not contain language \"%v\"\n", Language)
        os.Exit(1)
    }
}