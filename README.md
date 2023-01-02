# Simple MQTT Bridge

Bit of a misnomer as it's actually just forwarding messages from one MQTT server to another.

This could be used for forwarding traffic from a public MQTT server to one that is not accessible from the internet.

My personal itch that made me create this was that I had a esp-home sensor that had to work over the internet. The easiest
way I found to fix this was to send the data via a publicly accessible MQTT server and forward the traffic to the Mosquitto
addon's local server.

# Building

```
git clone https://github.com/tobias-/simple_mqtt_bridge
go build -ldflags="-X main.Commit=$(git rev-parse HEAD)"
```

# Running

```
./simple_mqtt_bridge --verbose --source-url ssl://foo:bar@instance.rmq.cloudamqp.com --dest-url tcp://bar:foo@hassio.local
```

