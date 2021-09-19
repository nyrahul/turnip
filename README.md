# Turnip
Scourge the given IP address in the published blocklists. Attackers connect to their C&C (command and control) servers and there are services that publish a list of IP addresses of such C&C servers. Blocklists are published as:
* [List of IP Addresses](https://feodotracker.abuse.ch/downloads/ipblocklist_recommended.txt)
* [Set of snort rules](https://feodotracker.abuse.ch/downloads/feodotracker_aggressive.rules)

This project provides a golang API to verify if the given address is in any of the published blocklist.

