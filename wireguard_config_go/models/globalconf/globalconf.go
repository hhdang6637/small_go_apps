package globalconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hhdang6637/small_go_apps/wireguard_config_go/util"
)

type globalConf struct {
	Subnet           string
	SubnetMask       int
	ServerDomainName string
	WgPort           int
}

var global = globalConf{
	Subnet:           "10.0.0.0",
	SubnetMask:       24,
	ServerDomainName: "example.wireguard.com",
	WgPort:           51820,
}

// ToJSON export all node to json format which can be used for storegae
func ToJSON() string {

	var (
		err   error
		Bytes []byte
	)
	Bytes, err = json.MarshalIndent(global, "", "    ")
	util.CheckErr(err)

	return string(Bytes)
}

// FromJSON load from json string which was created by ToJOSN
func FromJSON(str string) {
	err := json.Unmarshal([]byte(str), &global)
	util.CheckErr(err)
}

// LoadFromFile load global config from json file
func LoadFromFile(fileName string) {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s, create default config\n", err.Error())
		SaveToFile(fileName)
		return
	}

	FromJSON(string(dat))
}

// SaveToFile save config as json to text file
func SaveToFile(fileName string) {
	file, err := os.Create(fileName)
	util.CheckErr(err)
	defer file.Close()

	file.WriteString(ToJSON())
}

// GetGlobalInfo return global config
func GetGlobalInfo() (string, int, string, int) {
	return global.Subnet, global.SubnetMask, global.ServerDomainName, global.WgPort
}

// GetFisrtIPFromSubnet return the first ip in subnet
func GetFisrtIPFromSubnet() string {
	strs := strings.Split(global.Subnet, ".")
	if len(strs) != 4 {
		panic(global.Subnet + " is not valid ip, something wrong happend")
	}

	return fmt.Sprintf("%s.%s.%s.1", strs[0], strs[1], strs[2])
}
