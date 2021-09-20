package main

import (
	"flag"
	"os"

	turnip "github.com/nyrahul/turnip/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var GitCommit string
var GitBranch string
var BuildDate string

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func printBuildDetails() {
	log.Info().Msgf("commit:%v, branch:%v, date:%v", GitCommit, GitBranch, BuildDate)
}

func tryAddress(ip string) {
	src, reason := turnip.AddressIsBlocked(ip)
	if src != nil {
		log.Info().Msgf("ip=%v\nsrc=%v\nlink=%v\nseverity=%v\nreason=%v",
			ip, src.Name, src.Link, src.Severity, reason)
	} else {
		log.Info().Msgf("Not in any blocked list")
	}
}

func main() {
	source := flag.String("source", "data-sources.json", "Data source to use")
	flag.Parse()
	printBuildDetails()
	err := turnip.Setup(*source)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	tryAddress("97.107.134.115")
	tryAddress("103.248.217.234")
	tryAddress("192.168.10.10")
}
