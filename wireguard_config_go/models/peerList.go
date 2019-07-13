package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/hhdang6637/small_go_apps/wireguard_config_go/models/globalconf"
	"github.com/hhdang6637/small_go_apps/wireguard_config_go/util"
)

// PeerInfo store peer config
type PeerInfo struct {
	Name       string
	PrivateKey string
	PublicKey  string
	IPAddr     string
}

// PeerList contain information for all wiregaurd peer nodes include master node
type PeerList struct {
	Peers []PeerInfo
}

// NewPeerList like contructor to create new peerList
func NewPeerList() *PeerList {
	return &PeerList{
		Peers: make([]PeerInfo, 0),
	}
}

// NameValid return false if the name is already existed, otherwire it return true
func (peerL *PeerList) NameValid(name string) bool {
	for _, p := range peerL.Peers {
		if name == p.Name {
			return false
		}
	}

	return true
}

// AddPeer add new peer node to the list, if any problem, call panic with an error message
func (peerL *PeerList) AddPeer(name string) string {
	var err error

	if util.WgToolIsExisted() == false {
		panic("Wg tool isn't existed")
	}

	if len(peerL.Peers) == 0 && name != "master" {
		panic("the first peer's name must be master")
	}

	if peerL.NameValid(name) == false {
		panic(name + " was existed")
	}

	p := PeerInfo{
		Name:       name,
		IPAddr:     peerL.findFreeIP(),
		PrivateKey: "",
		PublicKey:  "",
	}

	if p.IPAddr == "" {
		panic("Cannot found any free IP address to allocate for new Peer")
	}

	p.PrivateKey, err = util.WgGenPrivateKey()
	util.CheckErr(err)

	p.PublicKey, err = util.WgGenPublicKey(p.PrivateKey)
	util.CheckErr(err)

	p.PrivateKey = strings.Trim(p.PrivateKey, "\n")
	p.PublicKey = strings.Trim(p.PublicKey, "\n")

	peerL.Peers = append(peerL.Peers, p)

	Bytes, err := json.MarshalIndent(p, "", "    ")
	util.CheckErr(err)
	return string(Bytes)
}

// DelPeer delete a peer node and return json encodeing of the node to the calller
// if any problem, call panic with an error message
func (peerL *PeerList) DelPeer(name string) string {
	var foundPeer *PeerInfo
	foundIndex := -1

	for index, p := range peerL.Peers {
		if p.Name == name {
			foundIndex = index
			foundPeer = &p
			break
		}
	}

	if foundIndex == -1 {
		return ""
	}

	peerL.Peers = append(peerL.Peers[:foundIndex], peerL.Peers[foundIndex+1:]...)

	Bytes, err := json.MarshalIndent(*foundPeer, "", "    ")
	util.CheckErr(err)

	return string(Bytes)
}

func (peerL *PeerList) genMasterConf() string {

	var (
		buffer    bytes.Buffer
		masterPtr *PeerInfo
	)

	for index := range peerL.Peers {
		if peerL.Peers[index].Name == "master" {
			masterPtr = &peerL.Peers[index]
			break
		}
	}

	if masterPtr == nil {
		panic("master is not found")
	}

	_, subnetmask, _, listenPort := globalconf.GetGlobalInfo()

	buffer.WriteString("# master\n")
	buffer.WriteString("[Interface]\n")

	buffer.WriteString("ListenPort = ")
	buffer.WriteString(strconv.Itoa(listenPort))
	buffer.WriteString("\n")

	buffer.WriteString("PrivateKey = ")
	buffer.WriteString(masterPtr.PrivateKey)
	buffer.WriteString("\n")

	buffer.WriteString("#PublicKey = ")
	buffer.WriteString(masterPtr.PublicKey)
	buffer.WriteString("\n")

	buffer.WriteString("Address = ")
	buffer.WriteString(masterPtr.IPAddr + "/" + strconv.Itoa(subnetmask))
	buffer.WriteString("\n")
	buffer.WriteString("PostUp = iptables -I FORWARD 2 -i wg0 -j ACCEPT;\n")
	buffer.WriteString("PostDown = iptables -D FORWARD 2 -i wg0 -j ACCEPT;\n")

	buffer.WriteString("\n")

	for _, p := range peerL.Peers {
		if p.Name == "master" {
			continue
		}
		buffer.WriteString("\n")

		buffer.WriteString("# " + p.Name + "\n")
		buffer.WriteString("[Peer]\n")

		buffer.WriteString("#PrivateKey = ")
		buffer.WriteString(p.PrivateKey)
		buffer.WriteString("\n")

		buffer.WriteString("PublicKey = ")
		buffer.WriteString(p.PublicKey)
		buffer.WriteString("\n")

		buffer.WriteString("AllowedIPs = ")
		buffer.WriteString(p.IPAddr + "/" + strconv.Itoa( /*subnetmask*/ 32))
		buffer.WriteString("\n")

		buffer.WriteString("PersistentKeepalive = ")
		buffer.WriteString("30")

		buffer.WriteString("\n")
	}

	return buffer.String()
}

func (peerL *PeerList) genPeerConf(peerName string) string {
	var (
		masterPtr, peerPtr *PeerInfo
	)

	for index := range peerL.Peers {
		if peerL.Peers[index].Name == "master" {
			masterPtr = &peerL.Peers[index]
		}
		if peerL.Peers[index].Name == peerName {
			peerPtr = &peerL.Peers[index]
		}
	}

	if masterPtr == nil || peerPtr == nil {
		panic("canot found master or " + peerName)
	}

	subnet, subnetmask, hostname, port := globalconf.GetGlobalInfo()

	var buffer bytes.Buffer

	buffer.WriteString("[Interface]\n")

	buffer.WriteString("PrivateKey = ")
	buffer.WriteString(peerPtr.PrivateKey)
	buffer.WriteString("\n")

	buffer.WriteString("#PublicKey = ")
	buffer.WriteString(peerPtr.PublicKey)
	buffer.WriteString("\n")

	buffer.WriteString("Address = ")
	buffer.WriteString(peerPtr.IPAddr + "/" + strconv.Itoa(subnetmask))
	buffer.WriteString("\n")

	buffer.WriteString("\n")

	buffer.WriteString("[Peer]\n")
	buffer.WriteString("PublicKey = ")
	buffer.WriteString(masterPtr.PublicKey)
	buffer.WriteString("\n")

	buffer.WriteString("AllowedIPs = ")
	buffer.WriteString(subnet + "/" + strconv.Itoa(subnetmask))
	buffer.WriteString("\n")

	buffer.WriteString("Endpoint = ")
	buffer.WriteString(hostname + ":" + strconv.Itoa(port))
	buffer.WriteString("\n")

	buffer.WriteString("PersistentKeepalive = ")
	buffer.WriteString("30")
	buffer.WriteString("\n")

	return buffer.String()
}

// GenConf return the content of wg0.conf
func (peerL *PeerList) GenConf(name string) string {

	if name == "master" {
		return peerL.genMasterConf()
	}

	return peerL.genPeerConf(name)
}

// ToJSON export all node to json format which can be used for storegae
func (peerL *PeerList) ToJSON() string {

	var (
		err   error
		Bytes []byte
	)
	Bytes, err = json.MarshalIndent(peerL.Peers, "", "    ")
	util.CheckErr(err)

	return string(Bytes)
}

func (peerL *PeerList) findFreeIP() string {
	var ints = []int{}

	if len(peerL.Peers) == 0 {
		return globalconf.GetFisrtIPFromSubnet()
	}

	str := peerL.Peers[0].IPAddr

	strs := strings.Split(str, ".")
	if len(strs) != 4 {
		panic(str + " is not valid ip, something wrong happend")
	}

	for _, s := range strs {
		i, err := strconv.Atoi(s)
		util.CheckErr(err)
		ints = append(ints, i)
	}

	intArrayToString := func(a []int, delim string) string {
		return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	}

	for n := 1; n <= 254; n++ {
		freeIP := true
		ints[3] = n
		strIP := intArrayToString(ints, ".")
		for _, p := range peerL.Peers {
			if p.IPAddr == strIP {
				freeIP = false
				break
			}
		}
		if freeIP == true {
			return strIP
		}
	}

	return ""
}

// FromJSON load peer node conf from json string which was created by ToJOSN
func (peerL *PeerList) FromJSON(str string) {
	err := json.Unmarshal([]byte(str), &peerL.Peers)
	util.CheckErr(err)
}

// SaveToFile save the list as json to text file
func (peerL *PeerList) SaveToFile(filename string) {
	file, err := os.Create(filename)
	util.CheckErr(err)
	defer file.Close()

	file.WriteString(peerL.ToJSON())
}

// LoadFromFile load the list of peer node from json file
func (peerL *PeerList) LoadFromFile(filename string) {

	dat, err := ioutil.ReadFile(filename)
	util.CheckErr(err)

	peerL.FromJSON(string(dat))
}
