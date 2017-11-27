package caspar

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	. "assembly.org/tools"
	// "github.com/jmoiron/jsonq"
	"github.com/jmoiron/jsonq"
	"strconv"
	"time"
)

type (
	casparConfig struct {
		address  string
		amcpPort string
		oscPort  string
		out      string
		amcpConn net.Conn
	}

	Player struct {
		Server  int    `json:"server"`
		Channel string `json:"channel"`
	}

	Server struct {
		VMixInput int64 `json:"vmix-input"`
		Observer  int64 `json:"observer"`
	}

	ConfigFile struct {
		Servers  []Server          `json:"servers"`
		Players  map[string]Player `json:"players"`
		Commands map[string]string `json:"commands"`
	}
)

var (
	cs             []casparConfig
	commands       map[string]string
	Players        map[string]Player
	Servers        map[int64]int64 // vMix -> Caspar mapping
	previousPlayer string
	previousInput  int
)

func Configure() {
	cs = make([]casparConfig, 2)

	cs[0].address = GetEnvParam("CASPAR_0_ADDRESS", "127.0.0.1")
	cs[0].amcpPort = GetEnvParam("CASPAR_0_AMCP_PORT", "5250")
	cs[0].oscPort = GetEnvParam("CASPAR_0_OSC_PORT", "5253")
	cs[0].out = "1-1"

	cs[1].address = GetEnvParam("CASPAR_1_ADDRESS", "127.0.0.1")
	cs[1].amcpPort = GetEnvParam("CASPAR_1_AMCP_PORT", "5251")
	cs[1].oscPort = GetEnvParam("CASPAR_1_OSC_PORT", "5254")
	cs[1].out = "1-1"

	conffile := ConfigFile{}
	confFilename := GetEnvParam("CASPAR_CONFIG", "pkm.json")
	readConfig(&conffile, confFilename)

	Servers = make(map[int64]int64)
	for _, v := range conffile.Servers {
		Servers[v.VMixInput] = v.Observer
	}

	Players = make(map[string]Player)
	Players = conffile.Players

	commands = make(map[string]string)
	commands = conffile.Commands

	connectCaspar(cs[0].address, cs[0].amcpPort, 0)
	connectCaspar(cs[1].address, cs[1].amcpPort, 1)
}

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
		plr.Server = Players[strconv.Itoa(n)].Server
		plr.Channel = Players[strconv.Itoa(n)].Channel
		plrs[k] = plr
		n++
	}
}

// SwitchPlayer käskee konfiguraation mukaista CasparCG palvelinta tarvittaessa vaihtamaan pelaajatunnuksen mukaiseen
// pelaajakameraan. Lisäksi input-parametrilla menee tällä hetkellä tunnettu vMixin valittuna oleva kanava, joka saattaa
// olla CasparCG-kuva. VMix input<->CasparCG palvelin mäppäys pitää hoitaa konfiguraatiossa.
func SwitchPlayer(input int64, currentPlayer string) {
	if Servers[input] == 0 {
		return
	}

	log.Printf("Observattava pelaaja vaihtui %d -> %d", previousPlayer, currentPlayer)
	if Players[currentPlayer].Channel == "" {
		log.Printf("Pelaajatunnusta %s ei löytynyt. Pelaajakuvan vaihto ei onnistu.", currentPlayer)
		return
	}
	prevSrv := Players[previousPlayer].Server
	currSrv := Players[currentPlayer].Server
	if prevSrv != currSrv {
		sendCommand(fmt.Sprintf(commands["clear_layer"], cs[prevSrv].out), prevSrv)
	}
	cmd := fmt.Sprintf(commands["select_player"], cs[currSrv].out, Players[currentPlayer].Channel)
	sendCommand(cmd, currSrv)
	previousPlayer = currentPlayer
}

func PlayClip(filename string) {
	cmd := fmt.Sprintf(commands["play_clip"], cs[0].out, filename)
	sendCommand(cmd, 0)
}

func ClearOut() {
	cmd := fmt.Sprintf(commands["clear_layer"], cs[0].out)
	sendCommand(cmd, 0)
}

func sendCommand(cmd string, server int) {
	var crnl = "\x0D\x0A"
	var code int
	var err error

	code, err = cs[server].amcpConn.Write([]byte(strings.Join([]string{cmd, crnl}, "")))
	if err != nil {
		log.Print("Caspar yhteys hukattu: ", code, err)
		connectCaspar(cs[server].address, cs[server].amcpPort, server)
	}
}

func connectCaspar(address string, port string, server int) {
	var err error
	cs[server].amcpConn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%s", address, port), time.Second*2)
	if err != nil {
		log.Printf("Yhteys Caspar-palvelimeen %s:%s epäonnistui: %s", address, port, err)
	}
	log.Printf("Yhteys Caspar-palvelimeen %s:%s avattu", address, port)
}

func readConfig(conf *ConfigFile, filename string) {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(conf)
	if err != nil {
		log.Fatal("Konfiguraatiotiedoston lukuvirhe: ", err)
	}
}
