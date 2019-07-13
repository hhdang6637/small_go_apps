package vtvUtil

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
)

func validVtvURL(urlStr string) {
	var (
		urlRe *regexp.Regexp
		err   error
	)
	urlRe, err = regexp.Compile("https://vtvgo\\.vn/xem-truc-tuyen-kenh-")
	if err != nil {
		panic(err)
	}
	if urlRe.MatchString(urlStr) == false {
		panic("Cannot understand URL " + urlStr)
	}
}

func vtvGoExtractCookieAndPostData(urlStr string) (string, string) {
	var (
		paramsRe *regexp.Regexp
		resp     *http.Response
		err      error
		html     string
		postData string
		params   []string
	)

	resp, err = http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cookieStr := ""
	for _, s := range resp.Header["Set-Cookie"] {
		cookieStr += s + ";"
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		html += scanner.Text() + "\n"
	}

	// fmt.Println(html)
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	paramsRe, err = regexp.Compile("\\s(?P<key>(?:type_)?id|time|token)\\s*=\\s*[\"']?(?P<value>[^\"']+)[\"']?;")
	params = paramsRe.FindAllString(html, -1)

	if len(params) != 4 {
		panic("POST data is not expected")
	}

	// fmt.Println(params)
	for _, p := range params {
		strs := strings.Split(p, " = ")
		if len(strs) == 2 {

			// append & if we have more than one item
			if len(postData) != 0 {
				postData += "&"
			}

			// trim some specical charecters " ;'"
			postData += strings.Trim(strs[0], " ") + "=" + strings.Trim(strs[1], ";'")
		} else {
			panic("wrong param value " + p)
		}
	}

	return cookieStr, postData
}

// Example return value of vtvRequestM3u8
// {
//     "stream_url": [
//         "https:\/\/1414383384.rsc.cdn77.org\/aT1MGbeKvgA_jkLlxCFn4w==,1562917179\/ls-46961-2\/index.m3u8"
//     ],
//     "ads_tags": "https:\/\/pubads.g.doubleclick.net\/gampad\/live\/ads?iu=\/276136803\/vtvgo.desktop.vtv2.video\u0026description_url=http%3A%2F%2Fvtvgo.vn\u0026env=vp\u0026impl=s\u0026correlator=\u0026tfcd=0\u0026npa=0\u0026gdfp_req=1\u0026output=vast\u0026sz=640x480\u0026unviewed_position_start=1",
//     "chromecast_url": "https:\/\/1414383384.rsc.cdn77.org\/aT1MGbeKvgA_jkLlxCFn4w==,1562917179\/ls-46961-2\/index.m3u8",
//     "remoteip": "174.138.20.203",
//     "content_id": 2,
//     "stream_info": [
//         {
//             "bandwidth": 528000,
//             "resolution": "360"
//         },
//         {
//             "bandwidth": 928000,
//             "resolution": "480"
//         },
//         {
//             "bandwidth": 1728000,
//             "resolution": "720"
//         }
//     ],
//     "date": "",
//     "player_type": "native",
//     "channel_name": "vtv2",
//     "geoname_id": 1880251,
//     "ads_time": "1"
// }

func vtvRequestM3u8(urlStr string, cookie string, postData string) map[string]interface{} {

	const urlPostRequest string = "https://vtvgo.vn/ajax-get-stream"
	const userAgentFirefox string = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36"

	var (
		req      *http.Request
		postResp *http.Response
		err      error
	)

	httpClient := &http.Client{}

	// fmt.Println(postData)
	req, err = http.NewRequest("POST", urlPostRequest, bytes.NewBuffer([]byte(postData)))
	if err != nil {
		panic(err.Error())
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://vtvgo.vn")
	req.Header.Set("Referer", urlStr)
	req.Header.Set("User-Agent", userAgentFirefox)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Cookie", cookie)

	postResp, err = httpClient.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer postResp.Body.Close()

	jsonData := ""
	postScanner := bufio.NewScanner(postResp.Body)
	for postScanner.Scan() {
		jsonData += postScanner.Text() + "\n"
	}

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(jsonData), &jsonMap)
	return jsonMap
}

// GetVtvGoM3u8Link accept vtv link as input and return chromecast_url
//
// Example:
// 		input is https://vtvgo.vn/xem-truc-tuyen-kenh-vtv2-2.html
//
// 		output is https://1414383384.rsc.cdn77.org/MOsdz-XmM1ZnjC421_e0WA==,1562917765/ls-46961-2/index.m3u8
func GetVtvGoM3u8Link(urlStr string) string {

	validVtvURL(urlStr)

	cookie, postData := vtvGoExtractCookieAndPostData(urlStr)

	mapJSON := vtvRequestM3u8(urlStr, cookie, postData)

	if mapJSON["chromecast_url"] != nil {
		switch mapJSON["chromecast_url"].(type) {
		case string:
			return mapJSON["chromecast_url"].(string)
		}
	}

	panic("chromecast_url is not found")
}

// M3u8Index2Mono return all m3u8 sublinks
func M3u8Index2Mono(link string) []string {
	var (
		linkTypes []string
		resp      *http.Response
		err       error
	)

	link = strings.Replace(link, "https://", "http://", -1)
	baseURL := link[:strings.LastIndex(link, "/")+1]

	resp, err = http.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		str := scanner.Text()
		if len(str) > 0 && str[0] != '#' {
			linkTypes = append(linkTypes, baseURL+str)
		}
	}

	return linkTypes
}

// M3u8GetTSLinks return all ts links in m3u8 link
func M3u8GetTSLinks(m3u8URL string) ([]string, error) {
	var (
		tsLinks []string
		resp    *http.Response
		err     error
	)

	baseURL := m3u8URL[:strings.LastIndex(m3u8URL, "/")+1]

	resp, err = http.Get(m3u8URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode%200 > 99 {
		return nil, errors.New("Fail to request ts links")
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		str := scanner.Text()
		// fmt.Println(str)
		if len(str) > 0 && str[0] != '#' {
			tsLinks = append(tsLinks, baseURL+str)
		} else {
			tsLinks = append(tsLinks, str)
		}
	}

	return tsLinks, nil
}
