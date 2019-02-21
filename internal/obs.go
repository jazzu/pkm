package internal

import (
	"encoding/json"
	"github.com/jmoiron/jsonq"

	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
)

type (
	Config struct {
		TeamAFile *string
		TeamBFile *string
		TestOnly  *bool
	}

	obsServer struct {
		address    string
		port       string
		connection *websocket.Conn
	}

	// OBS:lle lähetettävä komento
	SetSceneItemProperties struct {
		RequestType string `json:"request-type"`
		MessageId   string `json:"message-id"`
		Item        string `json:"item"`
		Visible     bool   `json:"visible"`
		SceneName   string `json:"scene-name"`
	}

	Player struct {
		PlayerName string `json:"player_name"`
		Camera     string `json:"camera"` //todo: siivoa json-formaatti
		Place      int    `json:"place"`
	}
)

var (
	obsServers        []obsServer
	Players           map[string]interface{}
	Cameras           map[string]interface{}
	previousPlayerSID string
	messageID         int
	testOnly          bool
)

func ConfigureOBS(configuration Config) {
	var err error

	serverSetup()

	testOnly = *configuration.TestOnly

	Cameras = make(map[string]interface{})
	Cameras, err = CQ.Object("cameras")

	teamConfigurations := make(map[string]*jsonq.JsonQuery)
	teamConfigurations["A"] = LoadJsonFile(*configuration.TeamAFile)
	teamConfigurations["B"] = LoadJsonFile(*configuration.TeamBFile)

	log.Println("Load players:")
	//yhdistetään eri tiedostot yhteen
	for teamLetter, confJQ := range teamConfigurations {
		var teamConf map[string]interface{}
		teamConf, err = confJQ.Object("players")

		if err != nil {
			log.Fatalln("Joukkuekonfiguraation lukeminen ei onnistunut: ", err)
		}

		for steamId, iPlayerConf := range teamConf {
			playerConf := iPlayerConf.(Player)
			// Jos pelaajan place-arvoksi on annettu 0, ei tämän videokuvaa näytetä observauksen aikana.
			cameraId := teamLetter + strconv.Itoa(playerConf.Place)
			playerConf.Camera = Cameras[cameraId].(string)
			log.Printf("%s -> %s : %d - %s\n", steamId, playerConf.PlayerName, playerConf.Place, playerConf.Camera)
			Players[steamId] = playerConf
		}
	}

	log.Println(Players)
	log.Println("OBS konfiguraation lataus tehty ja palvelimiin yhdistetty.")
}

// SwitchPlayer käskee tunnettuja palvelimia vaihtamaan inputtia, samat komennot jokaiselle.
// Inputtien nimet pitää olla OBS:ssä uniikkeja jotta vain oikea kone reagoi (muut antavat virheen josta ei välitetä)

//todo: silmukoita palvelimista, vai oma funktio joka lähettää kaikille
func SwitchPlayer(currentPlayerSID string) {
	cp := Players[currentPlayerSID].(Player)
	pp := Players[previousPlayerSID].(Player)

	log.Println("Valittu pelaajakamera: ", cp.Camera)
	if cp.Camera == "" {
		// TODO: formaatti sellaiseksi että voi copypasteta suoraan conffiin
		log.Printf("Pelaajatunnusta %s ei löytynyt. Pelaajakuvan vaihto ei onnistu.\n", currentPlayerSID)
		previousPlayerSID = "0"
		return
	}

	if previousPlayerSID == "" {
		log.Printf("nollataan")
		// Piilotetaan kaikki kamerakuvat, koska muuten saadaan tuplia
		for _, p := range Players {
			player := p.(Player)
			setCameraVisibility(player.Camera, false)
		}
	}

	if currentPlayerSID != previousPlayerSID {
		log.Printf("Observattava pelaaja vaihtui %s -> %s", previousPlayerSID, currentPlayerSID)
		// Uusi pelaaja näkyviin
		setCameraVisibility(cp.Camera, true)
		// Vanha pois. Jos uusi pelaaja on pienemmällä numerolla kuin vanha, näkyvä muutos tapahtuu vasta tässä
		setCameraVisibility(pp.Camera, false)
		previousPlayerSID = currentPlayerSID
	}
}

func serverSetup() {
	servers, err := CQ.ArrayOfObjects("camera_servers")
	if err != nil {
		log.Fatal("OBS-palvelinten konfiguraatioiden luku epäonnistui: %s", err)
	}

	obsServers = make([]obsServer, len(servers))
	for i, v := range servers {
		log.Printf("%d:%s", i)
		obsServers[i].address = v["address"].(string)
		obsServers[i].port = v["port"].(string)
		if err = obsServers[i].Connect(); err != nil {
			log.Fatal("OBS palvelimeen yhdistäminen epäonnistui: ", err)
		}
	}
}

func setCameraVisibility(camera string, visible bool) {
	for _, s := range obsServers {
		s.SetVisibility(camera, visible)
	}
}

func (obs obsServer) Connect() error {
	var err error
	u := url.URL{Scheme: "ws", Host: obs.host(), Path: "/"}
	obs.connection, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("Yhteys OBS-palvelimeen %s epäonnistui: %s", obs.host(), err)
	}
	log.Printf("Yhteys OBS-palvelimeen %s avattu", obs.host)

	for i := 1; i <= 10; i++ {
		obs.SetVisibility("cam"+strconv.Itoa(i), false)
	}

	log.Printf("Kamerakuvat piilotettu")
	return nil
}

func (obs obsServer) SetVisibility(camera string, visible bool) {
	var (
		err           error
		commandToSend *SetSceneItemProperties
		jsonToSend    []byte
	)

	messageID++
	commandToSend = &SetSceneItemProperties{
		RequestType: "SetSceneItemProperties",
		MessageId:   strconv.Itoa(messageID),
		Item:        camera, // cam1..cam10
		Visible:     visible,
		SceneName:   "Scene1"}

	jsonToSend, err = json.Marshal(commandToSend)
	if err != nil {
		log.Fatalf("Lähetettävän JSON-viestin muodostaminen epäonnistui: %s", err)
		return
	}

	//debug ilman servereitä
	if testOnly {
		log.Println("Testimoodi, viestiä ei lähetetä OBS-palvelimelle")
		//log.Println(jsonToSend)
		//log.Println(obsServers.connection)
		return
	}

	err = obs.connection.WriteMessage(websocket.TextMessage, jsonToSend)
	if err != nil {
		log.Printf("Websocket kirjoitus OBS-palvelimelle %s epäonnistui: %s", obs.host(), err)
	}
	return
}

func (obs obsServer) host() string {
	return obs.address + ":" + obs.port
}
