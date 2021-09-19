package turnip

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

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
type BadMap map[string]AddressData

var srcData DataSources // contains data sources from json
var badMaps []BadMap    // contains blocked address maps

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
	bm = make(BadMap)
	for _, ip := range list {
		if strings.HasPrefix(ip, "#") {
			continue
		}
		ip = strings.TrimRight(ip, "\r\n")
		bm[ip] = AddressData{"unk"}
	}
	return bm
}

func handleSource(src DataSource) (BadMap, error) {
	log.Info().Msgf("src=%v enable=%v", src.Name, src.IsEnabled)
	if !src.IsEnabled {
		log.Info().Msgf("data source:%v is not enabled", src.Name)
		return nil, nil
	}
	log.Info().Msgf("reading from [%v]", src.Link)
	buf, err := readFileLink(src.Link)
	if err != nil {
		log.Error().Msgf("read file failed for %v", src.Link)
		return nil, err
	}
	log.Info().Msg(string(buf))
	bm := getBadMapFromList(buf)
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
		if err != nil || bm == nil {
			continue
		}
		badMaps = append(badMaps, bm)
		log.Info().Msgf("bad:%v err:%v", len(badMaps), err)
	}
	return nil
}

func AddressIsBlocked(addr string) (bool, string) {
	for _, bm := range badMaps {
		ad, ok := bm[addr] // ad = address data
		if ok {
			return true, ad.Reason
		}
		log.Info().Msgf("addr=%v, ad=%v, reason=%v len=%v", addr, ad, ad.Reason, len(bm))
	}
	return false, ""
}
