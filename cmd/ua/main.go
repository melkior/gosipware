package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"sipware/api"
	"sipware/ua/tcp"
	"sipware/message"
)

func readConfig(file string) tcpua.Config {
	var config tcpua.Config

	jsonFile, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer jsonFile.Close()

	data, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(data, &config)

	fmt.Printf("Config %+v\n", err, config)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return config
}

func exitHandler(ua sipware.Ua) {
	<-ua.Exit()
	println("Exit!")
	ua.Done()
	os.Exit(2)
}

func main() {
	args := os.Args[1:]

	if(len(args) != 1) {
		panic("Usage: bin/ua confile")
	}

	config := readConfig(args[0])
	fmt.Println("Ua config", config)
	ua := tcpua.NewUa("Tcp ua")
	ua.SetExitHandler()

	go exitHandler(ua)

	ua.Open(config.Open)
	ua.Start()
	ua.Register(config.Register)

	ua.Wait()
	ua.Destroy(true)
}
