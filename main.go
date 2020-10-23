package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func urlParse(passedurl string) *url.URL {
	if passedurl[:4] != "http" {
		passedurl = "https://" + passedurl
	}
	u, err := url.Parse(passedurl)
	if err != nil {
		fmt.Println("Please enter a valid URL - ", err.Error())
		os.Exit(0)
	}
	if u.Path == "" {
		u.Path = "/"
	}
	return u
}

func getResponse(u *url.URL) (string, int) {

	// conn, err := net.Dial("tcp", u.Host+":80")
	timeout, _ := time.ParseDuration("10s")

	// Establishing secure TLS connection to the url
	conn, err := tls.DialWithDialer(&(net.Dialer{Timeout: timeout}), "tcp", u.Hostname()+":https", nil)
	if err != nil {
		fmt.Println("TCP Dialup failed - ", err.Error())
		os.Exit(0)
	}
	rt := fmt.Sprintf("GET %v HTTP/1.0\r\n", u.Path)
	rt += fmt.Sprintf("Host: %v\r\n", u.Hostname())
	// rt += fmt.Sprintf("Connection: close\r\n")
	rt += fmt.Sprintf("\r\n")

	defer conn.Close()

	_, err = conn.Write([]byte(rt))
	if err != nil {
		fmt.Println("Connection failed - ", err.Error())
		os.Exit(0)
	}

	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println("Failed to get response - ", err.Error())
		os.Exit(0)
	}
	respstring := string(resp)
	status, _ := strconv.Atoi(strings.Split(respstring, " ")[1])

	return respstring, status

}

func helpMe() {
	println("Usage:\n")
	println("\tgo run . [arguments]\n")
	println("The arguments are:\n")
	println("\t--url <URL>\t\t\t\tMake an HTTP request to the specified URL")
	println("\t--profile <URL> <Number of requests>\tProfile specified number of requests to the specified URL")
	println("\t--help\t\t\t\t\tPrints this help page")
}

func beginProfile(u *url.URL, count int) {
	var totalTimeArr []int
	var errors []int
	maxSize := float64(0)
	minSize := math.MaxFloat64
	sumTotalTime := 0
	for i := 1; i <= count; i++ {
		startTime := time.Now()
		resp, status := getResponse(u)
		resSize := float64(len(resp))
		totalTime := int(time.Since(startTime).Milliseconds())
		sumTotalTime += totalTime
		totalTimeArr = append(totalTimeArr, totalTime)
		if status != 200 {
			errors = append(errors, status)
		}
		maxSize = math.Max(maxSize, resSize)
		minSize = math.Min(minSize, resSize)
		// println(resSize)
		// println(totalTime)
		// // fmt.Println(resp)
	}
	sort.Ints(totalTimeArr)
	println("Profile for ", strings.ToLower(u.Hostname()), ":\n")
	println("Number of requests: ", count)
	println("The fastest time: ", totalTimeArr[count-1], "ms")
	println("The slowest time: ", totalTimeArr[0], "ms")
	println("The mean time: ", int(float64(sumTotalTime)/float64(count)), "ms")
	println("The median time: ", totalTimeArr[count/2], "ms")
	println("The percentage of requests that succeeded: ", int(((float64(count)-float64(len(errors)))*100.0)/float64(count)), "%")
	if len(errors) > 0 {
		fmt.Printf("The error codes returned that weren't a success: ")
		for _, value := range errors {
			fmt.Printf("%d ", value)
		}
		println("")
	}
	println("The size in bytes of the smallest response: ", int(minSize))
	println("The size in bytes of the largest response: ", int(maxSize))
}

func executeCLI() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "--help" {
		helpMe()
	} else if args[0] == "--url" {
		u := urlParse(args[1])
		resp, _ := getResponse(u)
		fmt.Println(resp)
	} else if args[0] == "--profile" && len(args) == 3 {
		u := urlParse(args[1])
		count, err := strconv.Atoi(args[2])
		if err != nil || count <= 0 {
			fmt.Println("Please enter a valid number of requests")
			os.Exit(0)
		}
		beginProfile(u, count)

	} else {
		println("Invalid Usage, please refer usage below.")
		helpMe()
		os.Exit(0)
	}
}

func main() {
	// fmt.Println("Hello World")
	executeCLI()
}
