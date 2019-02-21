package internal

import (
	"encoding/json"
	"github.com/jmoiron/jsonq"
	"io"
	"log"
	"os"
)

var (
	CQ *jsonq.JsonQuery
)

func ConfigurePKM(filename string) {
	CQ = LoadJsonFile(filename)
}

func LoadJsonFile(filename string) *jsonq.JsonQuery {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("JSON-tiedoston %s lataaminen epäonnistui: %s", filename, err)
		return nil
	}
	return DecodeJsonToJsonQ(file)
}

func DecodeJsonToJsonQ(reader io.Reader) *jsonq.JsonQuery {
	var err error
	decoder := json.NewDecoder(reader)
	configuration := map[string]interface{}{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatalf("Konfiguraation lukuvirhe: %s", err)
		return nil
	}
	return jsonq.NewQuery(configuration)
}
