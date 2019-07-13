package main

import (
	"fmt"
	"os"

	"github.com/hhdang6637/small_go_apps/wireguard_config_go/models"
	"github.com/hhdang6637/small_go_apps/wireguard_config_go/models/globalconf"
)

const configFileName string = "wg.json"
const globalFileName string = "global.json"

func help() {
	fmt.Println("expected [showconf,add,del,genconf] <PeerName> commands, with PeerName isn't empty")
	os.Exit(1)
}

func main() {
	// Global info
	globalconf.LoadFromFile(globalFileName)

	// Peer info
	var peerListPtr *models.PeerList
	peerListPtr = models.NewPeerList()
	peerListPtr.LoadFromFile(configFileName)

	if len(os.Args) == 2 {
		if os.Args[1] == "showconf" {
			fmt.Println("Global config:")
			fmt.Println(globalconf.ToJSON())
			fmt.Println("\nPeer config:")
			fmt.Println(peerListPtr.ToJSON())
			os.Exit(0)
		}
		help()
	}

	if len(os.Args) != 3 || len(os.Args[2]) == 0 {
		help()
	}

	switch os.Args[1] {
	case "add":
		fmt.Println(peerListPtr.AddPeer(os.Args[2]))
	case "del":
		fmt.Println(peerListPtr.DelPeer(os.Args[2]))
	case "genconf":
		fmt.Println(peerListPtr.GenConf(os.Args[2]))
	default:
		fmt.Println("expected add or del commands")
		os.Exit(1)
	}

	peerListPtr.SaveToFile(configFileName)
}
