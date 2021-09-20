package turnip

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/gonids"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type DataSources struct {
	DataSources []DataSource `json:"sources"`
}

type DataSource struct {
	Name      string `json:"name"`
	Severity  string `json:"severity"`
	Link      string `json:"link"`
	Type      string `json:"type"`
	IsEnabled bool   `json:"enable"`
}

type AddressData struct {
	Reason string
}

// Bad = Blocked ADdress Map
type BadMap struct {
	Map map[string]AddressData
	Ds  DataSource
}

var isInited bool
var srcData DataSources // contains data sources from json

var badMaps []BadMap // contains blocked address maps

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func httpGet(url string) ([]byte, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return io.ReadAll(rsp.Body)
}

func readFileLink(f string) ([]byte, error) {
	if strings.HasPrefix(f, "http") {
		return httpGet(f)
	}
	return os.ReadFile(f)
}

func getBadMapFromList(buf []byte) BadMap {
	var bm BadMap
	list := strings.Split(string(buf), "\n")
	bm.Map = make(map[string]AddressData)
	for _, ip := range list {
		if strings.HasPrefix(ip, "#") {
			continue
		}
		ip = strings.TrimRight(ip, "\r\n")
		bm.Map[ip] = AddressData{"blocked IP"}
	}
	return bm
}

func getBadMapFromIDSRules(buf []byte) BadMap {
	var bm BadMap
	list := strings.Split(string(buf), "\n")
	bm.Map = make(map[string]AddressData)
	for _, rule := range list {
		if strings.HasPrefix(rule, "#") {
			continue
		}
		r, err := gonids.ParseRule(rule)
		if err != nil {
			log.Error().Msgf("failed parsing ids rule: [%v]", rule)
			continue
		}
		/* log.Info().Msgf("ids rule: action=%v prot=%v dst=%v desc=%v",
		r.Action, r.Protocol, r.Destination.Nets, r.Description) */
		for _, ip := range r.Destination.Nets {
			bm.Map[ip] = AddressData{r.Description}
		}
	}
	return bm
}

func handleSource(src DataSource) (BadMap, error) {
	bm := BadMap{nil, src}
	log.Info().Msgf("src=%v enable=%v", src.Name, src.IsEnabled)
	if !src.IsEnabled {
		log.Info().Msgf("data source:%v is not enabled", src.Name)
		return bm, nil
	}
	log.Info().Msgf("reading from [%v]", src.Link)
	buf, err := readFileLink(src.Link)
	if err != nil {
		log.Error().Msgf("read file failed for %v", src.Link)
		return bm, err
	}
	if src.Type == "list" {
		bm = getBadMapFromList(buf)
	} else if src.Type == "snort" {
		bm = getBadMapFromIDSRules(buf)
	} else {
		return bm, fmt.Errorf("unknown type: %v", src.Type)
	}
	return bm, nil
}

func Setup(dataSrc string) error {
	jsonFile, err := os.Open(dataSrc)
	if err != nil {
		log.Error().Msgf("failed opening json file=%v err=%v", dataSrc, err.Error())
		return err
	}
	defer jsonFile.Close()
	jsonData, _ := io.ReadAll(jsonFile)

	json.Unmarshal([]byte(jsonData), &srcData)

	for _, src := range srcData.DataSources {
		bm, err := handleSource(src)
		if err != nil {
			log.Error().Msgf("handle source failed: %v", err.Error())
			continue
		}
		if bm.Map == nil { // could happen since source's IsEnabled is false
			continue
		}
		bm.Ds = src
		badMaps = append(badMaps, bm)
	}
	return nil
}

func AddressIsBlocked(addr string) (*DataSource, string) {
	for _, bm := range badMaps {
		ad, ok := bm.Map[addr] // ad = address data
		if ok {
			return &bm.Ds, ad.Reason
		}
	}
	return nil, ""
}
