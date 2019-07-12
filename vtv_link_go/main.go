package main

import (
	"fmt"
	"net/http"
	"time"
	"vtv_link_go/vtvUtil"
)

func testVtvM3u8Link() {

	const vtv2URL = "https://vtvgo.vn/xem-truc-tuyen-kenh-vtv2-2.html"

	m3u8Links := vtvUtil.M3u8Index2Mono(vtvUtil.GetVtvGoM3u8Link(vtv2URL))
	if len(m3u8Links) > 0 {
		for _, link := range vtvUtil.M3u8GetTSLinks(m3u8Links[len(m3u8Links)-1]) {
			fmt.Println(link)
		}
	} else {
		panic("m3u8 link is not found")
	}
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func httpServerLoop() {
	http.HandleFunc("/", greet)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}

func main() {

}
