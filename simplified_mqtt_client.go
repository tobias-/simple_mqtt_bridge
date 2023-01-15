package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	log "github.com/amoghe/distillog"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func parseUrl(mqttUrl string, name string) (parsedUrl string, username string, password string, err error) {
	parsedMqttUrl, err := url.Parse(mqttUrl)
	if err != nil {
		return
	}
	protocol := parsedMqttUrl.Scheme
	if protocol != "ssl" && protocol != "tcp" && protocol != "ws" && protocol != "wss" {
		err = errors.New(name + ": unknown protocol/scheme '" + parsedMqttUrl.Scheme + "'")
		return
	}

	var port int
	if len(parsedMqttUrl.Port()) > 0 {
		port, err = strconv.Atoi(parsedMqttUrl.Port())
		if err != nil {
			err = errors.New(name + ": invalid port in mqtt-url " + mqttUrl)
			return
		}
	} else {
		if protocol == "ssl" || protocol == "wss" {
			port = 8883
		} else {
			port = 1883
		}
	}
	virtualHost := strings.Trim(parsedMqttUrl.Path, "/")
	if len(virtualHost) != 0 {
		username = parsedMqttUrl.User.Username() + ":" + virtualHost
	} else {
		username = parsedMqttUrl.User.Username()
	}

	password, hasPassword := parsedMqttUrl.User.Password()
	if hasPassword && strings.Contains(password, ":") {
		err = errors.New(name + ": invalid username/password. Only exactly one colon is allowed")
		return
	}
	host := parsedMqttUrl.Hostname()
	if len(host) == 0 {
		err = errors.New(name + ": invalid mqtt-url")
		return
	}
	parsedUrl = fmt.Sprintf("%s://%s:%d", protocol, host, port)
	return
}

func ParseMqttArgs(mqttUrl string, name string) (mqttOptions *mqtt.ClientOptions, err error) {
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Infof("Connected to %s mqtt server", name)
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Warningf("Connect lost to %s mqtt server: %v", name, err)
	}

	var connectAttemptHandler mqtt.ConnectionAttemptHandler = func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		log.Debugf("Connecting to %s mqtt server: %s", name, broker)
		return tlsCfg
	}

	parsedUrl, username, password, err := parseUrl(mqttUrl, name)

	mqttOptions = mqtt.NewClientOptions()
	mqttOptions.AddBroker(parsedUrl)
	if len(username) != 0 {
		mqttOptions.SetUsername(username)
	}
	if len(password) != 0 {
		mqttOptions.SetPassword(password)
	}
	mqttOptions.SetConnectTimeout(5 * time.Second)
	mqttOptions.SetConnectionAttemptHandler(connectAttemptHandler)
	mqttOptions.SetConnectRetry(true)
	mqttOptions.SetConnectRetryInterval(10 * time.Second)
	mqttOptions.SetAutoReconnect(true)
	mqttOptions.SetMaxReconnectInterval(10 * time.Second)
	mqttOptions.OnConnect = connectHandler
	mqttOptions.OnConnectionLost = connectLostHandler
	return
}
