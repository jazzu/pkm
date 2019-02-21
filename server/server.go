package server

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sytem/pkm/obs"
	"github.com/sytem/pkm/tools"
)

func Run() {
	setup()

	listenAddress := listenAddress()
	log.Print("PKM palvelin käynnistyy osoitteessa: " + listenAddress)

	router := mux.NewRouter()

	router.HandleFunc("/", ReceiveGameStatus)
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

// ReceiveGameStatus käsittelee CS:GO observerin lähettämän pelidatapaketin
func ReceiveGameStatus(w http.ResponseWriter, r *http.Request) {
	rawPost := getRawPost(r)

	var data GameData
	err := json.Unmarshal(rawPost, &data)
	if err != nil {
		log.Fatal(err)
	}
	logComparisonJson(data)

	// Varmista että JSON:issa tuli mukana pelaajatieto ja yritä vaihtaa kuvaa ainoastaan jos se löytyy
	if data.PlayerID != nil {
		obs.SwitchPlayer(data.PlayerID.SteamID)
		log.Print("\"" + data.PlayerID.SteamID + "\": {\"player_name\": \"" + data.PlayerID.Name + "\", \"place\": 0},")
	}
	w.WriteHeader(http.StatusOK)
}

func getRawPost(r *http.Request) (body []byte) {
	var err error
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func logComparisonJson(data GameData) {
	var checkData []byte
	var err error
	checkData, err = json.Marshal(data)
	if err != nil {
		log.Fatal("Vertailumerkkijonon muodostaminen epäonnistui: ", err)
		log.Print(string(checkData))
	}
}

func setup() {
	pConfFilename := flag.String("conf", "pkm.json", "JSON konfiguraatiotiedosto yleisille asetuksille")

	obsConfig := obs.Config{}
	obsConfig.TeamAFile = flag.String("A", "team1.json", "JSON konfiguraatiotiedosto A-tiimille")
	obsConfig.TeamBFile = flag.String("B", "team2.json", "JSON konfiguraatiotiedosto B-tiimille")
	obsConfig.TestOnly = flag.Bool("test", false, "testaa palvelinsovellusta paikallisesti lähettämättä ohjauskomentoja")
	flag.Parse()

	tools.Configure(*pConfFilename)
	obs.Configure(obsConfig)
}

func listenAddress() string {
	var address, port string
	var err error

	address, err = tools.CQ.String("pkm", "address")
	if err != nil {
		log.Fatal("Puuttuva tai virheellinen PKM osoitekonfiguraatio: ", err)
		os.Exit(1)
	}

	port, err = tools.CQ.String("pkm", "port")
	if err != nil {
		log.Fatal("Puuttuva tai virheellinen PKM porttikonfiguraatio: ", err)
		os.Exit(1)
	}

	return address + ":" + port
}
