package lib

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/charmbracelet/log"
	pb "github.com/meshnet-gophers/meshtastic-go/meshtastic"
	"github.com/meshnet-gophers/meshtastic-go/mqtt"
	"github.com/meshnet-gophers/meshtastic-go/radio"
	"google.golang.org/protobuf/proto"
	"strings"
)

func ChannelHandler(channel string) mqtt.HandlerFunc {
	return func(m mqtt.Message) {
		var env pb.ServiceEnvelope
		err := proto.Unmarshal(m.Payload, &env)
		if err != nil {
			log.Fatal("failed unmarshalling to service envelope", "err", err, "payload", hex.EncodeToString(m.Payload))
			return
		}

		key, err := GenerateKey("1PG7OiApB1nwvP+rz05pAQ==")
		if err != nil {
			log.Fatal(err)
		}

		decodedMessage, err := radio.XOR(env.Packet.GetEncrypted(), key, env.Packet.Id, env.Packet.From)
		if err != nil {
			log.Error(err)
		}
		var message pb.Data
		err = proto.Unmarshal(decodedMessage, &message)

		log.Info(ProcessMessage(&message), "topic", m.Topic, "channel", channel, "portnum", message.Portnum.String())
	}
}

func GenerateKey(key string) ([]byte, error) {
	// Pad the key with '=' characters to ensure it's a valid base64 string
	padding := (4 - len(key)%4) % 4
	paddedKey := key + strings.Repeat("=", padding)

	// Replace '-' with '+' and '_' with '/'
	replacedKey := strings.ReplaceAll(paddedKey, "-", "+")
	replacedKey = strings.ReplaceAll(replacedKey, "_", "/")

	// Decode the base64-encoded key
	return base64.StdEncoding.DecodeString(replacedKey)
}
