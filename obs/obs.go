package obs

import (
	"encoding/json"

	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/jmoiron/jsonq"
	"strconv"

	"github.com/gorilla/websocket"
)

type (
	Player struct {
		PlayerName string `json:"player_name"`
		//Server  int    `json:"server"`
		Channel string `json:"channel"` //todo: siivoa json-formaatti
		Place   int    `json:"place"`
	}

	CameraServer struct {
		Ip   string `json:"ip"`
		Port string `json:"port"`
	}

	ConfigFile struct {
		CameraServers []CameraServer    `json:"cameraServers"`
		Players       map[string]Player `json:"players"`
		Cameras       map[string]string `json:"cameras"`
		//todo: pkm listen port
	}

	obsConfig struct {
		address string
		port    string
		conn    *websocket.Conn
	}

	// OBS:lle lähetettävä komento
	SetSceneItemRender struct {
		RequestType string `json:"request-type"`
		MessageId   string `json:"message-id"`
		Source      string `json:"source"`
		Render      bool   `json:"render"`
		SceneName   string `json:"scene-name"`
	}
)

var (
	obs            []obsConfig
	commands       map[string]string
	Players        map[string]Player
	Cameras        map[string]string
	CameraServers  map[string]string
	previousPlayer string
	previousInput  int
	messageID      int
	testOnly       bool
)

func Configure() {

	confFilenamePtr := flag.String("conf", "pkm.json", "json-file for basic configuration")
	teamAPtr := flag.String("A", "team1.json", "json-file for team A")
	teamBPtr := flag.String("B", "team2.json", "json-file for team B")
	testOnlyPtr := flag.Bool("test", false, "test server without sending commands")
	flag.Parse()

	conffile := ConfigFile{}
	readConfig(&conffile, *confFilenamePtr)

	teamAfile := ConfigFile{}
	readConfig(&teamAfile, *teamAPtr)

	teamBfile := ConfigFile{}
	readConfig(&teamBfile, *teamBPtr)

	obs = make([]obsConfig, 2) //todo: tästä kovakoodaus pois
	for i, v := range conffile.CameraServers {
		log.Printf("%d:%s", i, v.Ip)
		obs[i].address = v.Ip
		obs[i].port = v.Port
		connectOBS(obs[i].address, obs[i].port, 0)
	}

	testOnly = *testOnlyPtr

	Players = make(map[string]Player)
	Players = conffile.Players

	Cameras = make(map[string]string)
	Cameras = conffile.Cameras

	fmt.Println("Load players:")
	//yhdistetään eri tiedostot yhteen
	for k, v := range teamAfile.Players {
		fmt.Printf("%s -> %s : %d\n", k, v.PlayerName, v.Place)
		var camera = "A" + strconv.Itoa(v.Place)
		v.Channel = Cameras[camera]
		Players[k] = v
	}

	for k, v := range teamBfile.Players {
		fmt.Printf("%s -> %s : %d\n", k, v.PlayerName, v.Place)
		//todo poikkeus paikka0
		var camera = "B" + strconv.Itoa(v.Place)
		v.Channel = Cameras[camera]
		fmt.Println("set channel to :" + v.Channel)
		Players[k] = v
	}

	messageID = 0
	previousPlayer = ""

	fmt.Println(Players)

	log.Printf("load valmis")
}

// todo: tekeekö tämä mitään?
func PopulatePlayerConf(jsonData string) {
	plrs := make(map[string]Player)

	// Jos player conffi on tyhjä, ota allplayers tieto pelidatasta ja laita niistä SteamID:t talteen
	testing := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(jsonData))
	dec.Decode(&testing)
	jq := jsonq.NewQuery(testing)
	obj, _ := jq.Object("allplayers")

	var n int = 1
	for k := range obj {
		plr := Player{}
		//plr.Server = Players[strconv.Itoa(n)].Server
		plr.Channel = Players[strconv.Itoa(n)].Channel
		plrs[k] = plr
		n++
	}
}

// SwitchPlayer käskee tunnettuja palvelimia vaihtamaan inputtia, samat komennot jokaiselle.
//Inputtien nimet pitää olla OBS:ssä uniikkeja jotta vain oikea kone reagoi (muut antavat virheen josta ei välitetä)

//todo: silmukoita palvelimista, vai oma funktio joka lähettää kaikille
func SwitchPlayer(input int64, currentPlayer string) {

	log.Print(Players[currentPlayer].Channel)
	if Players[currentPlayer].Channel == "" {
		// todo formaatti sellaiseksi että voi copypasteta suoraan conffiin
		log.Printf("Pelaajatunnusta %s ei löytynyt. Pelaajakuvan vaihto ei onnistu.", currentPlayer)
		sendCommand(Players[previousPlayer].Channel, false, 0)
		sendCommand(Players[previousPlayer].Channel, false, 1)
		previousPlayer = "0"
		return
	}

	if previousPlayer == "" {
		log.Printf("nollataan")
		//tähän reset all pimeäksi, koska muuten saadaan tuplia
		for _, player := range Players {
			sendCommand(player.Channel, false, 0)
			sendCommand(player.Channel, false, 1)
		}
	}

	if currentPlayer != previousPlayer {

		log.Printf("Observattava pelaaja vaihtui %d -> %d", previousPlayer, currentPlayer)
		//todo: tässä voisi olla myös for-luuppi käydä kaikki yhdistetyt serverit läpi
		//uusi pelaaja näkyviin
		sendCommand(Players[currentPlayer].Channel, true, 0)
		sendCommand(Players[currentPlayer].Channel, true, 1)

		//vanha pois. jos uusi pelaaja on pienemmällä numerolla kuin vanha, näkyvä muutosta tapahtuu vasta tässä
		sendCommand(Players[previousPlayer].Channel, false, 0)
		sendCommand(Players[previousPlayer].Channel, false, 1)

		previousPlayer = currentPlayer
	}
}

func sendCommand(input string, vis bool, server int) {

	//debug ilman servereitä
	if testOnly {
		log.Println("Not sending command, just test")
		return
	}

	messageID++

	commandToSend := &SetSceneItemRender{
		RequestType: "SetSceneItemRender",
		MessageId:   strconv.Itoa(messageID),
		Source:      input, // cam1..cam10
		Render:      vis,
		SceneName:   "Scene1"}

	jsonToSend, _ := json.Marshal(commandToSend)

	err := obs[server].conn.WriteMessage(websocket.TextMessage, jsonToSend)
	if err != nil {
		log.Println("write:", err)
		return
	}

}

func connectOBS(address string, port string, server int) {
	var err error
	addr := address + ":" + port
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	obs[server].conn = c
	if err != nil {
		log.Printf("Yhteys OBS-palvelimeen %s:%s epäonnistui: %s", address, port, err)
		return
	}
	log.Printf("Yhteys OBS-palvelimeen %s:%s avattu", address, port)

	//todo loop
	sendCommand("cam1", false, server)
	sendCommand("cam2", false, server)
	sendCommand("cam3", false, server)
	sendCommand("cam4", false, server)
	sendCommand("cam5", false, server)
	sendCommand("cam6", false, server)
	sendCommand("cam7", false, server)
	sendCommand("cam8", false, server)
	sendCommand("cam9", false, server)
	sendCommand("cam10", false, server)

	log.Printf("tyhjennetty")
}

func readConfig(conf *ConfigFile, filename string) {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(conf)
	if err != nil {
		log.Fatal("Konfiguraatiotiedoston lukuvirhe: ", err)
	}
}
