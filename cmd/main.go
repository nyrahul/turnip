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

func main() {
	source := flag.String("source", "data-sources.json", "Data source to use")
	flag.Parse()
	printBuildDetails()
	err := turnip.Setup(*source)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	isBlocked, reason, src := turnip.AddressIsBlocked("97.107.134.115")
	log.Info().Msgf("IsBlocked:%v, reason:%v, src=%v", isBlocked, reason, src)

	isBlocked, reason, src = turnip.AddressIsBlocked("103.248.217.234")
	log.Info().Msgf("IsBlocked:%v, reason:%v, src=%v", isBlocked, reason, src)
}
