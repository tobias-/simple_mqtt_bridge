package main

import (
	log "github.com/amoghe/distillog"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jessevdk/go-flags"
	"os"
	"time"
)

var Options struct {
	SourceMqttUrl      string `short:"s" long:"source-url" description:"ssl://username:password@host/virtual_host use tcp instead of ssl for cleartext and use virtual_host only for rabbitmq'"`
	DestinationMqttUrl string `short:"d" long:"dest-url" description:"ssl://username:password@host/virtual_host use tcp instead of ssl for cleartext and use virtual_host only for rabbitmq"`
	Version            bool   `short:"V" long:"version" description:"Print version"`
	Verbose            bool   `short:"v" long:"verbose" description:"Verbose"`
}

var Commit string = "not set"

func main() {
	_, err := flags.ParseArgs(&Options, os.Args)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	if Options.Version {
		log.Errorf("Compiled from %s", Commit)
		os.Exit(0)
	}

	sourceMqttClientOptions, err := ParseMqttArgs(Options.SourceMqttUrl, "source")
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	destinationMqttClientOptions, err := ParseMqttArgs(Options.DestinationMqttUrl, "destination")
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
	sourceClient := mqtt.NewClient(sourceMqttClientOptions)
	if !sourceClient.Connect().WaitTimeout(10 * time.Second) {
		log.Errorf("Error connecting to source")
		os.Exit(1)
	}
	destinationClient := mqtt.NewClient(destinationMqttClientOptions)
	if !destinationClient.Connect().WaitTimeout(10 * time.Second) {
		log.Errorf("Error connecting to destination")
		os.Exit(1)
	}
	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		if Options.Verbose {
			log.Infof("%s -> %s", msg.Payload(), msg.Topic())
		}
		destinationClient.Publish(msg.Topic(), 0, true, msg.Payload())
	}
	sourceClient.Subscribe("#", 0, messageHandler)
	for true {
		time.Sleep(1 * time.Second)
	}
}
