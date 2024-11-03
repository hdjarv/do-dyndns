package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/digitalocean/godo"
)

func checkErr(err error, msg string, code int) {
	if err != nil {
		errExit(msg, code)
	}
}

func errExit(msg string, code int) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(code)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if defaultValue == "" {
		errExit(fmt.Sprintf("Error: environment variable '%s' not set", key), 1)
	}
	return defaultValue
}

func getExtIP(url string, ipRe *regexp.Regexp) string {
	resp, err := http.Get(url)
	checkErr(err, fmt.Sprintf("Error: can't get external IP-address: %s", err), 2)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	checkErr(err, fmt.Sprintf("Error: reading response body: %s", err), 2)
	ip := ipRe.FindString(string(body))
	if ip == "" {
		errExit(fmt.Sprintf("Error: IP-address not found in data: %s", string(body)), 2)
	}
	return ip
}

func main() {
	isDryRun := getEnv("DO_DYN_DRY_RUN", "0") == "1"
	url := getEnv("DO_DYN_EXT_IP_URL", "")
	re := getEnv("DO_DYN_IP_REGEX", "\\b(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})\\b")
	ipRe, err := regexp.Compile(re)
	checkErr(err, fmt.Sprintf("Error compiling regex: %s", err), 1)

	apiToken := getEnv("DO_DYN_API_TOKEN", "")
	domain := getEnv("DO_DYN_DO_DOMAIN", "")
	record := getEnv("DO_DYN_RECORD_NAME", "")
	ttl, err := strconv.Atoi(getEnv("DO_DYN_DNS_TTL", "60"))
	checkErr(err, fmt.Sprintf("Error: non-integer value for 'DO_DYN_DNS_TTL': %s", err), 1)
	if ttl < 1 {
		errExit(fmt.Sprintf("Error: invalid value for 'DO_DYN_DNS_TTL': %d", ttl), 1)
	}
	fullDomain := record + "." + domain

	fmt.Println("Starting DO DynDNS")
	extIp := getExtIP(url, ipRe)
	fmt.Printf("External IP address is: %s\n", extIp)

	ctx := context.Background()
	doClient := godo.NewFromToken(apiToken)

	doRecords, _, err := doClient.Domains.RecordsByTypeAndName(ctx, domain, "A", fullDomain, &godo.ListOptions{})
	checkErr(err, fmt.Sprintf("Error: from Digital Ocean API: %s", err), 3)
	if len(doRecords) == 0 {
		errExit(fmt.Sprintf("Error: domain '%s' not found in Digital Ocean API", fullDomain), 3)
	}
	doRecord := doRecords[0]
	fmt.Printf("Current domain record for %s: IP=%s TTL=%d\n", fullDomain, doRecord.Data, doRecord.TTL)
	if doRecord.Data == extIp && doRecord.TTL == ttl {
		fmt.Println("Domain record is up to date")
	} else {
		logMsg := "Updated domain record:"
		if doRecord.Data != extIp {
			logMsg += " IP=" + extIp
		}
		if doRecord.TTL != ttl {
			logMsg += " TTL=" + strconv.Itoa(ttl)
		}
		if isDryRun {
			logMsg = "[DRY-RUN MODE] " + logMsg
		} else {
			dr, _, err := doClient.Domains.EditRecord(ctx, domain, doRecord.ID, &godo.DomainRecordEditRequest{Data: extIp, TTL: ttl})
			checkErr(err, fmt.Sprintf("Error: from Digital Ocean API: %s", err), 5)
			fmt.Println(dr)
		}
		fmt.Println(logMsg)
	}
	fmt.Println("Finished DO DynDNS")
}
