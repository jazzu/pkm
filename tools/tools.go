package tools

import (
	"encoding/json"
	"github.com/jmoiron/jsonq"
	"log"
	"os"
)

var (
	CQ *jsonq.JsonQuery
)

func Configure(filename string) {
	CQ = LoadJsonFile(filename)
}

func LoadJsonFile(filename string) *jsonq.JsonQuery {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("JSON-tiedoston %s lataaminen epäonnistui: %s", filename, err)
		return nil
	}
	decoder := json.NewDecoder(file)
	configuration := map[string]interface{}{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatalf("Konfiguraatiotiedoston %s lukuvirhe: %s", filename, err)
		return nil
	}
	return jsonq.NewQuery(configuration)
}
