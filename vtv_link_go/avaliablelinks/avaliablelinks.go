package avaliablelinks

import (
	"bufio"
	"encoding/json"
	"net/http"
)

// JSONDataLink point to online json file
const JSONDataLink = "https://raw.githubusercontent.com/hhdang6637/small_go_apps/master/vtv_link_go/playlist.w3u.txt"

// Station provides a convenient data struct of station object json
type Station struct {
	Image     string `json:"image"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	UserAgent string `json:"userAgent"`
	IsHost    bool   `json:"isHost"`
}

// Group provides a convenient data struct of group object json
type Group struct {
	Name     string    `json:"name"`
	Image    string    `json:"image"`
	Stations []Station `json:"stations"`
}

type onlineTv struct {
	Name   string  `json:"name"`
	Author string  `json:"author"`
	Image  string  `json:"image"`
	Info   string  `json:"info"`
	URL    string  `json:"url"`
	Groups []Group `json:"groups"`
}

// GetAvaliableGroups return all m3u8 link via Station structs
func GetAvaliableGroups() []Group {
	var (
		resp *http.Response
		err  error
	)
	resp, err = http.Get(JSONDataLink)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var httpBytes []byte
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		httpBytes = append(httpBytes, scanner.Bytes()...)
	}

	// fmt.Println(string(httpData))
	var onlineTvObj onlineTv
	err = json.Unmarshal(httpBytes, &onlineTvObj)
	if err != nil {
		panic(err)
	}

	// fmt.Println(onlineTvObj.Groups)
	return onlineTvObj.Groups
}

// GroupHaveM3u8Links return true if any stations have m3u8 link otherwise return false
func GroupHaveM3u8Links(groups Group) bool {
	for _, s := range groups.Stations {
		if s.IsHost == false {
			return true
		}
	}
	return false
}

// filter links die
// func main() {
// 	http.DefaultClient.Timeout = time.Second * 5
// 	online := getAvaliableGroups()
// 	for gIndex, g := range online.Groups {
// 		for sIndex, s := range g.Stations {
// 			if s.IsHost == true {
// 				continue
// 			}
// 			fmt.Fprintf(os.Stderr, "Verfiy link: %s - %s\n", s.Name, s.URL)
// 			resp, err := http.Get(s.URL)
// 			if err != nil || resp.StatusCode >= 300 {
// 				fmt.Fprintf(os.Stderr, "Link die: remove later %s\n", s.Name)
// 				online.Groups[gIndex].Stations[sIndex].IsHost = true
// 				continue
// 			}
// 			defer resp.Body.Close()
// 		}
// 	}

// 	for gIndex, g := range online.Groups {
// 		newStations := []Station{}
// 		for _, s := range g.Stations {
// 			if s.IsHost == true {
// 				continue
// 			}
// 			newStations = append(newStations, s)
// 		}
// 		online.Groups[gIndex].Stations = newStations
// 	}

// 	newGs := []Group{}
// 	for _, g := range online.Groups {
// 		if len(g.Stations) == 0 {
// 			continue
// 		}
// 		newGs = append(newGs, g)
// 	}
// 	online.Groups = newGs

// 	data, err := json.MarshalIndent(online, "", "    ")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(data))
// }
