package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func validURL(urlStr string) {
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

func vtvRequestM3u8(urlStr string, cookie string, postData string) string {

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

	html := ""
	postScanner := bufio.NewScanner(postResp.Body)
	for postScanner.Scan() {
		html += postScanner.Text() + "\n"
	}

	dst := bytes.Buffer{}
	json.Indent(&dst, []byte(html), "", "    ")
	return dst.String()
}

func exampleHTTPNewRequest(urlStr string) string {

	validURL(urlStr)

	cookie, postData := vtvGoExtractCookieAndPostData(urlStr)

	return vtvRequestM3u8(urlStr, cookie, postData)
}

func main() {
	fmt.Println(exampleHTTPNewRequest(os.Args[1]))
}
