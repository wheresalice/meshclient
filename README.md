# Meshclient

A CLI to receive [Meshtastic](https://meshtastic.org/) messages from MQTT or a radio connected to a serial port, written in Go

## Usage

### MQTT

```shell
meshclient m  # Use the defaults to connect to LongFast on EU_868
meshclient mqtt --url "tcp://mqtt.meshtastic.org:1883" --username "meshdev" --password "large4cats" --topic "msh/EU_868" --channel "LongFast"  # The long form of the above
```

### Radio

```shell
meshclient r  # Use the defaults to connect to a radio on /dev/ttyUSB0
meshclient radio --port "/dev/ttyUSB0"
```
