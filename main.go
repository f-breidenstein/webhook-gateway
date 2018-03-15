package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type IRCMessage struct {
	Messages []string
	Channels []string
}

var messageChannel = make(chan IRCMessage)

func main() {

	// Load configuration frm file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	var moduleList = viper.Sub("modules")

	// Status module
	if moduleList.GetBool("status.enabled") {
		log.Println("Status module is active")
		http.HandleFunc("/status", statusHandler)
	} else {
		log.Println("Status module disabled of not configured")
	}

	// Prometheus module
	if moduleList.GetBool("prometheus.enabled") {
		log.Println("Prometheus module is active")
		http.HandleFunc("/prometheus", prometheusHandler)
	}

	// Start IRC worker process
	go ircChannelWorker()

	// Start HTTP server
	log.Fatal(http.ListenAndServe(viper.GetString("http.listen"), nil))

}
