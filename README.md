# Turnip
Scourge the given IP address in the published blocklists. Attackers connect to their C&C (command and control) servers and there are services that publish a list of IP addresses of such C&C servers. Blocklists are published as:
* [List of IP Addresses](https://feodotracker.abuse.ch/downloads/ipblocklist_recommended.txt)
* [Set of snort rules](https://feodotracker.abuse.ch/downloads/feodotracker_aggressive.rules)

This project provides a golang API to verify if the given address is in any of the published blocklist.

### Example Usage

```go
import turnip "github.com/nyrahul/turnip/api"
:::
	err := turnip.Setup("data-sources.json")	// sample data-sources.json provided
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	isBlocked, reason, src := turnip.AddressIsBlocked("97.107.134.115")
	log.Info().Msgf("IsBlocked:%v, reason:%v, src=%v", isBlocked, reason, src)

	/* Sample Output
	9:06PM INF IsBlocked:true, reason:Feodo Tracker: potential TrickBot CnC Traffic detected, src=feodo-snort
	*/
```
