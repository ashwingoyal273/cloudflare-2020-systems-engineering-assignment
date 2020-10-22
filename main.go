package main

import (
	"fmt"
	"os"
)

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
	} else {
		println("HALLOWEEN")
	}
}

func main() {
	fmt.Println("Hello World")
	parseArgs()
}
