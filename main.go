package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"
)

func urlParse(passedurl string) *url.URL {
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

func getResponse(u *url.URL) string {
	conn, err := net.Dial("tcp", u.Host+":80")
	if err != nil {
		fmt.Println("TCP Dialup failed - ", err.Error())
		os.Exit(0)
	}
	rt := fmt.Sprintf("GET %v HTTP/1.1\r\n", u.Path)
	rt += fmt.Sprintf("Host: %v\r\n", u.Host)
	rt += fmt.Sprintf("Connection: close\r\n")
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
	return string(resp)

}

func helpMe() {
	println("Usage:\n")
	println("\tgo run . [arguments]\n")
	println("The arguments are:\n")
	println("\t--url <URL>\t\t\t\tMake an HTTP request to the specified URL")
	println("\t--profile <URL> <Number of requests>\tProfile specified number of requests to the specified URL")
	println("\t--help\t\t\t\t\tPrints this help page")
}

func parseArgs() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "--help" {
		helpMe()
	} else if args[0] == "--url" {
		u := urlParse(args[1])
		resp := getResponse(u)
		fmt.Println(resp)
	} else if args[0] == "--profile" && len(args) == 3 {
		u := urlParse(args[1])
		count, err := strconv.Atoi(args[2])
		if err != nil || count <= 0 {
			fmt.Println("Please enter a valid number of requests")
			os.Exit(0)
		}
		for i := 1; i <= count; i++ {
			startTime := time.Now()
			resp := getResponse(u)
			resSize := len(resp)
			totalTime := time.Since(startTime)
			println(totalTime)
			println(resSize)
			// fmt.Println(resp)
		}
	} else {
		println("Invalid Usage, please refer usage below.")
		helpMe()
		os.Exit(0)
	}
}

func main() {
	fmt.Println("Hello World")
	parseArgs()
}
