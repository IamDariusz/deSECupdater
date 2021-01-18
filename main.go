package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"
)

func getIP() string {
	r, err := http.Get("https://api.ipify.org/?format=json") // get json from API
	if err != nil {                                          // handle error if any
		log.Fatalf("Error on getting WAN IP: %v\n", err)
	}
	defer r.Body.Close() // close body

	buf := new(bytes.Buffer)  // create a new byte slice to read to (to save UTF-8 code points - avoids encoding errors)
	buf.ReadFrom(r.Body)      // read body to byte slice
	ipifyData := buf.String() // byte slice to string

	// Ipify : Struct for Current WAN IP
	type Ipify struct {
		IP string
	}

	// Get Struct and unmarshal
	var ipify Ipify
	err = json.Unmarshal([]byte(ipifyData), &ipify)
	if err != nil { // handle error if any
		log.Fatalf("Error on unmarshal for WAN IP: %v\n", err)
	}

	return fmt.Sprintf("%v", ipify.IP)
}

func setDNS(token, domain, ipwan string) (int, string) {
	jsonStr := []byte(fmt.Sprintf(`{"records": ["%v"]}`, ipwan))                      // create body
	deSecURL := fmt.Sprintf("https://desec.io/api/v1/domains/%v/rrsets/@/A/", domain) // create domain
	req, err := http.NewRequest("PATCH", deSecURL, bytes.NewBuffer(jsonStr))          // new patch request
	if err != nil {
		log.Fatalf("Error on creating new request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %v", token))

	httpClient := &http.Client{Timeout: 10 * time.Second}
	r, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Error on performing DESEC request: %v\n", err)
	}
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	return r.StatusCode, string(body)
}

func paramCheck(token, domain, ipwan string) (bool, string) {
	if len(token) != 28 {
		return false, "Token is not 28 characters long"
	}

	if domainCheck, _ := regexp.MatchString("\\.", domain); domainCheck == false {
		return false, "Domain does not contain a '.'"
	}

	if net.ParseIP(ipwan) == nil {
		return false, "IP seems to be invalid"
	}
	return true, ""
}

func main() {
	desecToken := flag.String("token", "XXX", "token for DESEC's API")
	desecDomain := flag.String("domain", "domain.com", "domain to change IP for")
	desecIP := flag.String("ip", "", "IPv4 address to use, will automatically use WAN IP if empty")
	desecDbg := flag.Bool("debug", false, "debug desec call's body")
	flag.Parse()

	wanIP := *desecIP
	if *desecIP == "" {
		log.Printf("Trying to get IP address...")
		wanIP = getIP()
		log.Printf("Got IP Address: %v", wanIP)
	}

	// Check Parameters
	log.Printf("Verifying Parameters and IP...")
	paramCheckResult, paramCheckMsg := paramCheck(*desecToken, *desecDomain, wanIP)
	if paramCheckResult == false {
		log.Fatalf("Parameters seem to be incorrect: %v", paramCheckMsg)
	}
	log.Printf("Parameters seem to be valid!")

	log.Printf("Try to set DESEC DNS Settings...")

	setRetry := 1
	for setRetry <= 10 {
		log.Printf("Try number: %v", setRetry)
		desecStatus, desecBody := setDNS(*desecToken, *desecDomain, wanIP)
		if desecStatus != 200 {
			log.Printf("Was not successful. Got status code: %v", desecStatus)
			log.Printf("Will retry in 5s")
			time.Sleep(5 * time.Second)
		} else {
			log.Printf("Was successful. Got status code: %v", desecStatus)
			setRetry = 11
		}
		if *desecDbg {
			log.Printf("DEBUG Got body: %v", desecBody)
		}

		setRetry = setRetry + 1
	}
}
