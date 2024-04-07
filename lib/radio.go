package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"github.com/charmbracelet/log"
	pb "github.com/meshnet-gophers/meshtastic-go/meshtastic"
	"github.com/meshnet-gophers/meshtastic-go/transport"
	"github.com/meshnet-gophers/meshtastic-go/transport/serial"
	"google.golang.org/protobuf/proto"
	"strings"
)

func ConnectSerial(comport string, errorOnNoHandler bool) (*transport.Client, error) {
	if comport == "" {
		potentialPorts := serial.GetPorts()
		if potentialPorts == nil {
			return &transport.Client{}, errors.New("no usb serial devices detected")
		}
		if len(potentialPorts) > 1 {
			return &transport.Client{}, errors.New("multiple ports detected")
		}

		comport = potentialPorts[0]
	}
	log.Info("Connecting to", "serial port", comport)

	serialConn, err := serial.Connect(comport)
	if err != nil {
		return &transport.Client{}, err
	}

	streamConn, err := transport.NewClientStreamConn(serialConn)
	if err != nil {
		return &transport.Client{}, err
	}
	client := transport.NewClient(streamConn, errorOnNoHandler)
	return client, nil

}

func generateKey(key string) ([]byte, error) {
	// Pad the key with '=' characters to ensure it's a valid base64 string
	padding := (4 - len(key)%4) % 4
	paddedKey := key + strings.Repeat("=", padding)

	// Replace '-' with '+' and '_' with '/'
	replacedKey := strings.ReplaceAll(paddedKey, "-", "+")
	replacedKey = strings.ReplaceAll(replacedKey, "_", "/")

	// Decode the base64-encoded key
	return base64.StdEncoding.DecodeString(replacedKey)
}

func generateNonce(packetId uint32, node uint32) []byte {
	packetNonce := make([]byte, 8)
	nodeNonce := make([]byte, 8)

	binary.LittleEndian.PutUint32(packetNonce, packetId)
	binary.LittleEndian.PutUint32(nodeNonce, node)

	return append(packetNonce, nodeNonce...)
}

func decode(encryptionKey []byte, encryptedData []byte, nonce []byte) (pb.Data, error) {
	var message pb.Data

	ciphertext := encryptedData

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return message, err
	}
	stream := cipher.NewCTR(block, nonce)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	err = proto.Unmarshal(plaintext, &message)
	return message, err
}
