package main

import "fmt"

func main() {
	done := make(chan struct{})
	if err := launchDefaultService(done); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	<-done
}
