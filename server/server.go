package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sytem/pkm/obs"
	"github.com/sytem/pkm/tools"
	"github.com/sytem/pkm/vmix"
)

var (
	ActiveInput   int64
	PreviousInput int64
)

func Run() {
	obs.Configure()

	listenAddress := tools.GetEnvParam("PKM_LISTEN_ADDRESS", "127.0.0.1:1999")
	log.Print("PKM palvelin käynnistyy osoitteessa: " + listenAddress)

	router := mux.NewRouter()

	router.HandleFunc("/", ReceiveGameStatus)
	router.HandleFunc("/active_input/{input:[0-9]+}", ReceiveActiveInput)
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
		obs.SwitchPlayer(ActiveInput, data.PlayerID.SteamID)
		log.Print("\"" + data.PlayerID.SteamID + "\": {\"player_name\": \"" + data.PlayerID.Name + "\", \"place\": 0},")

	}
	w.WriteHeader(http.StatusOK)
}

// ReceiveActiveInput käsittelee gorutiinin vMixiltä pollaamaan input-tiedon
func ReceiveActiveInput(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)

	PreviousInput = ActiveInput
	ActiveInput, err = strconv.ParseInt(vars["input"], 10, 64)

	if err != nil {
		log.Fatal("Virheellinen GET-parametri: ", err)
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
	//log.Print("Tulkattu data:")
	var checkData []byte
	var err error
	checkData, err = json.Marshal(data)
	if err != nil {
		log.Fatal("Vertailumerkkijonon muodostaminen epäonnistui: ", err)
		log.Print(string(checkData))
	}
}

// vMixPoller tarkistaa videomikserin tilan silmukassa ohjelman loppuun asti. Gorutiinina ajettava funktio kirjoittaa
// vMixiltä lukemansa input valinnan PKM:n active_input endpointtiin. Tällä vältetään channeleiden kanssa pelaaminen.
func vMixPoller(listenAddress string) {
	for {
		time.Sleep(time.Millisecond * 100)
		input, err := pkm.CheckVmixStatus()
		if err != nil {
			log.Print("VMixin tilan tarkistus epäonnistui: ", err)
		}
		targetUrl := fmt.Sprintf("http://%s/active_input/%d", listenAddress, input)
		resp, err := http.Get(targetUrl)
		if err != nil {
			log.Printf("VMixin tilatiedon kirjoitus osoitteeseen %s epäonnistui: %s", targetUrl, err)
		}
		resp.Body.Close()
	}
}
