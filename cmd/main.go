package main

import "fmt"

func main() {
	if err := launchDefaultService(); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}
