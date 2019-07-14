package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

var (
	logger = log.New(os.Stdout, "", log.Lshortfile|log.Lmicroseconds)

	wgDevLink      string
	wgEndPointConf string
)

func getWgDevLink(file string) {
	filename := path.Base(file)

	wgDevLink = strings.Trim(filename, ".conf")
	logger.Printf("wgDevLink = %s\n", wgDevLink)
}

func getWgEndPointConf(file string) {
	pFile, err := os.Open(file)
	if err != nil {
		logger.Panicln(err)
	}
	defer pFile.Close()

	scanner := bufio.NewScanner(pFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "Endpoint = %s", &wgEndPointConf)
	}

	strs := strings.Split(wgEndPointConf, ":")
	if len(strs) != 2 {
		logger.Panicf("Wrong endpoint (%s) in file %s\n", wgEndPointConf, file)
	}
	wgEndPointConf = strs[0]
	logger.Printf("wgEndPointConf = %s\n", wgEndPointConf)
}

func getWgEndPointRunning(file string) string {
	wgCmd := exec.Command("wg", "show", wgDevLink)

	out, err := wgCmd.StdoutPipe()
	if err != nil {
		logger.Panicln(err)
	}

	wgCmd.Start()

	var wgEndPoint string

	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "  endpoint: %s", &wgEndPoint)
	}
	wgCmd.Wait()

	// 	logger.Printf("wgEndPoint = %s\n", wgEndPoint)
	strs := strings.Split(wgEndPoint, ":")
	if len(strs) != 2 {
		logger.Panicf("getWgEndPointRunning: Wrong endpoint (%s)\n", wgEndPoint)
	}
	logger.Printf("wg show %s -> EndPoint = %s\n", wgDevLink, wgEndPoint)
	return strs[0]
}

func updateWgConfg() {
	logger.Printf("updateWgConfg is called\n")
	wgCmd := exec.Command("wg-quick", "down", wgDevLink)
	_, err := wgCmd.Output()
	if err != nil {
		logger.Panicf("wg-quick down %s is FAIL\n", wgDevLink)
	}

	time.Sleep(time.Second * 2)

	wgCmd = exec.Command("wg-quick", "up", wgDevLink)
	_, err = wgCmd.Output()
	if err != nil {
		logger.Panicf("wg-quick up %s is FAIL\n", wgDevLink)
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("%s <config> [forced]\n", os.Args[0])
		os.Exit(1)
	}

	getWgDevLink(os.Args[1])
	getWgEndPointConf(os.Args[1])

	if len(os.Args) == 3 && os.Args[2] == "forced" {
		updateWgConfg()
		os.Exit(0)
	}

	for {
		ips, err := net.LookupIP(wgEndPointConf)
		if err == nil {

			sameIP := false
			wgRunningIP := getWgEndPointRunning(os.Args[1])

			for _, ip := range ips {
				logger.Printf("wgRunningIP : %s, DNS IP %s\n", wgRunningIP, ip.String())
				if wgRunningIP == ip.String() {
					sameIP = true
				}
			}

			if sameIP == false {
				updateWgConfg()
			}

		} else {
			logger.Printf("ERROR: %s\n", err.Error())
		}

		time.Sleep(time.Second * 60 * 5)
	}
}
