package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/hhdang6637/small_go_apps/vtv_link_go/avaliablelinks"
	"github.com/hhdang6637/small_go_apps/vtv_link_go/htmlHelper"
	"github.com/hhdang6637/small_go_apps/vtv_link_go/vtvutil"
)

var (
	vtvChannel = map[string]string{
		"vtv1":     "https://vtvgo.vn/xem-truc-tuyen-kenh-vtv1-1.html",
		"vtv2":     "https://vtvgo.vn/xem-truc-tuyen-kenh-vtv2-2.html",
		"vtv3":     "https://vtvgo.vn/xem-truc-tuyen-kenh-vtv3-3.html",
		"vtv4":     "https://vtvgo.vn/xem-truc-tuyen-kenh-vtv4-4.html",
		"vtv5":     "https://vtvgo.vn/xem-truc-tuyen-kenh-vtv5-5.html",
		"vtv6":     "https://vtvgo.vn/xem-truc-tuyen-kenh-vtv6-6.html",
		"vtv7":     "https://vtvgo.vn/xem-truc-tuyen-kenh-vtv7-27.html",
		"vtv8":     "https://vtvgo.vn/xem-truc-tuyen-kenh-vtv8-36.html",
		"vtv9":     "https://vtvgo.vn/xem-truc-tuyen-kenh-vtv9-39.html",
		"vtv5_tnb": "https://vtvgo.vn/xem-truc-tuyen-kenh-vtv5-tây-nam-bộ-7.html",
		"vtv5_tn":  "https://vtvgo.vn/xem-truc-tuyen-kenh-kênh-vtv5-tây-nguyên-163.html",
		// "Yeah1TV":  "https://vtvgo.vn/xem-truc-tuyen-kenh-yeah1tv-195.html",
		"HungYen": "https://vtvgo.vn/xem-truc-tuyen-kenh-th-hưng-yên-185.html",
	}

	vtvs = make([]string, 0)

	vtvM3u8Links = map[string][]string{}

	otherAvaliableGroups = []avaliablelinks.Group{}

	logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
)

func testVtvM3u8Link() {

	m3u8Links := vtvGetM3u8Link("vtv1")

	tsLinks, err := vtvutil.M3u8GetTSLinks(m3u8Links[len(m3u8Links)-1])
	if err != nil {
		vtvM3u8Links["vtv1"] = vtvGetM3u8Link("vtv1")
		m3u8Links = vtvM3u8Links["vtv1"]
		tsLinks, err = vtvutil.M3u8GetTSLinks(m3u8Links[len(m3u8Links)-1])
		if err != nil {
			logger.Panic(err)
		}
	}

	if len(m3u8Links) > 0 {
		for _, link := range tsLinks {
			fmt.Println(link)
		}
	} else {
		logger.Panic("m3u8 link is not found")
	}
}

func vtvGetM3u8Link(vtv string) []string {
	logger.Printf("vtvGetM3u8Link('%s') start", vtv)
	defer logger.Printf("vtvGetM3u8Link('%s') end", vtv)
	return vtvutil.M3u8Index2Mono(vtvutil.GetVtvGoM3u8Link(vtvChannel[vtv]))
}

func vtvTableLinks(w http.ResponseWriter) {
	for _, k := range vtvs {

		if vtvM3u8Links[k] == nil || len(vtvM3u8Links[k]) == 0 {
			vtvM3u8Links[k] = vtvGetM3u8Link(k)
		}

		if len(vtvM3u8Links[k]) == 0 {
			logger.Panicf(`Fail to get m3u8 link %s from vtv.go`, k)
		}

		fmt.Fprintf(w, `
		<tr>
	  		<td><h2><a href="/%s.m3u8" >%s</a></h2></td>
	  		<td>%s</td>
		</tr>
`,
			k, k, vtvM3u8Links[k][len(vtvM3u8Links[k])-1])
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	logger.Printf("%s: %s", r.RemoteAddr, r.RequestURI)

	if r.RequestURI != "/" {
		w.WriteHeader(404)
		fmt.Fprintln(w, "NOT FOUND URL REQUEST")
		return
	}

	htmlHelper.HTMLDocOpen(w)
	// vtvgo links
	fmt.Fprintln(w, "<h2>VTV Link Table</h2>")
	htmlHelper.TableDocOpen(w, []string{"Channel", "M3U8 link"})
	vtvTableLinks(w)
	htmlHelper.TableDocClose(w)

	// other links
	for _, g := range otherAvaliableGroups {
		if avaliablelinks.GroupHaveM3u8Links(g) == false {
			continue
		}
		fmt.Fprintf(w, "<h2>%s</h2>", g.Name)
		htmlHelper.TableDocOpen(w, []string{"Channel", "M3U8 link"})
		for _, s := range g.Stations {
			if s.IsHost != true {
				fmt.Fprintf(w, `
		<tr>
	  		<td><h2><a href="%s" >%s</a></h2></td>
	  		<td>%s</td>
		</tr>
`,
					s.URL, s.Name, s.URL)
			}
		}
		htmlHelper.TableDocClose(w)
	}

	htmlHelper.TableDocClose(w)
	htmlHelper.HTMLDocClose(w)
}

func vtvHandler(w http.ResponseWriter, r *http.Request) {

	logger.Printf("%s: %s", r.RemoteAddr, r.RequestURI)

	vtvC := "vtv1"
	switch r.RequestURI {
	case "/vtv1.m3u8":
		vtvC = "vtv1"
	case "/vtv2.m3u8":
		vtvC = "vtv2"
	case "/vtv3.m3u8":
		vtvC = "vtv3"
	case "/vtv4.m3u8":
		vtvC = "vtv4"
	case "/vtv5.m3u8":
		vtvC = "vtv5"
	case "/vtv6.m3u8":
		vtvC = "vtv6"
	case "/vtv7.m3u8":
		vtvC = "vtv7"
	case "/vtv8.m3u8":
		vtvC = "vtv8"
	case "/vtv9.m3u8":
		vtvC = "vtv9"
	}

	if vtvM3u8Links[vtvC] == nil || len(vtvM3u8Links[vtvC]) == 0 {
		vtvM3u8Links[vtvC] = vtvGetM3u8Link(vtvC)
	}

	if len(vtvM3u8Links[vtvC]) == 0 {
		logger.Printf(`Fail to get m3u8 link %s from vtv.go`, vtvC)
		return
	}

	w.Header().Add("Content-Type", "application/vnd.apple.mpegurl")

	links := vtvM3u8Links[vtvC]
	tsLinks, err := vtvutil.M3u8GetTSLinks(links[len(links)-1])
	if err != nil {
		logger.Panicf("Fail to get m3u8 from %s, try to update m3u8 link", links[len(links)-1])

		vtvM3u8Links[vtvC] = vtvGetM3u8Link(vtvC)
		links = vtvM3u8Links[vtvC]

		logger.Panicf("New m3u8 link: %s", links[len(links)-1])

		tsLinks, err = vtvutil.M3u8GetTSLinks(links[len(links)-1])
		if err != nil {
			logger.Panic(err)
		}
	}

	for _, link := range tsLinks {
		fmt.Fprintf(w, "%s\n", link)
	}
}

func httpServerLoop(port int) {

	http.HandleFunc("/", rootHandler)
	for k := range vtvChannel {
		http.HandleFunc("/"+k+".m3u8", vtvHandler)
	}

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		logger.Panic(err)
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "vtv_link_go <port_number>\n")
		fmt.Fprintf(os.Stderr, "you must provide port number for web server\n")
		os.Exit(1)
	}

	logger.Printf("Collect link started")

	for k := range vtvChannel {

		vtvs = append(vtvs, k)

		if vtvM3u8Links[k] == nil || len(vtvM3u8Links[k]) == 0 {
			vtvM3u8Links[k] = vtvGetM3u8Link(k)
		}

		if len(vtvM3u8Links[k]) == 0 {
			logger.Panicf(`Fail to get m3u8 link %s from vtv.go`, k)
		}
	}
	sort.Strings(vtvs)

	otherAvaliableGroups = avaliablelinks.GetAvaliableGroups()

	logger.Printf("Collect link done")

	portNumber, err := strconv.Atoi(os.Args[1])
	if err != nil {
		logger.Panic(err)
	}

	httpServerLoop(portNumber)
}
