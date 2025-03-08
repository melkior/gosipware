package main

import (
	"os"
	"fmt"
	"encoding/json"
	// "github.com/melkior/sipware-go"
	"github.com/melkior/sipware-go/api"
	"github.com/melkior/sipware-go/message"
	"github.com/melkior/sipware-go/ua/tcp"
)

func readConfig(file string) api.Config {
	data, err1 := os.ReadFile(file)

	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	}

	config := api.Config{}
	err2 := json.Unmarshal(data, &config)

	if err2 != nil {
		fmt.Println(err2)
		os.Exit(1)
	}

	fmt.Println("Config", config)
	return config
}

func readContact(file string) *string {
	fmt.Println("Read contact", file)

	data, err := os.ReadFile(file)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("Contact", data)
	contact := string(data)
	return &contact
}

func writeContact(file string, contact string) error {
	fmt.Println("Write contact", file, contact)

	err := os.WriteFile(file, []byte(contact), 0600)

	if err != nil {
		fmt.Println(err)
		// os.Exit(1)
		return err
	}

	return nil
}

func exitHandler(ua api.Ua) {
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

	contactFile := config.Cache.Contact.File
	contact := readContact(contactFile)
	fmt.Println("Ua contact", contact)

	ua := tcpua.New("Tcp ua")
	ua.SetExitHandler()

	go exitHandler(ua)

	ua.Open(config.Open)
	ua.Start()

	if contact == nil {
		ua.Register(config.Register, func(msg message.Msg) {
			if(msg.Code == 200) {
				fmt.Println("Register response", msg)
				to := msg.Get("To") [0]
				writeContact(contactFile, to)
				return
			}
			os.Exit(1)
		})
	} else {
		ua.Connect(*contact, func(msg message.Msg) {
			fmt.Println("Connect response", msg)

			if(msg.Code == 404) {
				os.Remove(contactFile)
				os.Exit(1)
			}
		})
	}

	ua.Wait()
	ua.Destroy(true)
}
