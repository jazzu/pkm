package pkm

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type Result struct {
	XMLName xml.Name `xml:"vmix"`
	Active  int      `xml:"active"`
}

// CheckVmixStatus kertoo vMixissä aktiivisena olevan inputin, jolloin CasparCG-koneen ollessa aktiivisena
// pystytään lähettämään oikeanlaiset käskyt oikealle koneelle.
func CheckVmixStatus() (int, error) {
	// TODO: Parametrisoi vMixin pollaus-URL
	vmixUrl := "http://192.168.133.100:8088/api/"
	statusXml, err := getVmixStatus(vmixUrl)
	if err != nil {
		// Logita lisää? Periaatteessa tieto on kerrottu jo tuolla toisessa funktiossa
		// ¯\_(ツ)_/¯
		return 0, err
	}
	return statusXml.Active, nil
}

func getVmixStatus(vmixUrl string) (Result, error) {
	v := Result{}
	resp, err := http.Get(vmixUrl)
	if err != nil {
		log.Printf("vMixin %s tilaa ei voitu lukea: %s", vmixUrl, err)
		return v, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("vMix tila-XML:n lukuvirhe: %v", err)
		return v, err
	}

	err = xml.Unmarshal([]byte(data), &v)
	if err != nil {
		log.Printf("XML parsintavirhe: %v", err)
		return v, err
	}
	return v, nil
}
