package main

import (
	"context"
	"github.com/charmbracelet/log"
	pb "github.com/meshnet-gophers/meshtastic-go/meshtastic"
	"github.com/meshnet-gophers/meshtastic-go/mqtt"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/proto"
	"meshclient/lib"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "meshclient",
		Description: "A CLI to receive Meshtastic messages from MQTT or a radio connected to a serial port",
		Commands: []*cli.Command{
			{
				Name:    "mqtt",
				Aliases: []string{"m"},
				Usage:   "receive from mqtt",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "url",
						Value: "tcp://mqtt.meshtastic.org:1883",
						Usage: "mqtt broker url",
					},
					&cli.StringFlag{
						Name:  "username",
						Value: "meshdev",
						Usage: "mqtt username",
					},
					&cli.StringFlag{
						Name:  "password",
						Value: "large4cats",
						Usage: "mqtt password",
					},
					&cli.StringFlag{
						Name:  "topic",
						Value: "msh/EU_868",
						Usage: "mqtt topic",
					},
					&cli.StringFlag{
						Name:  "channel",
						Value: "LongFast",
						Usage: "meshtastic channel",
					},
				},
				Action: func(c *cli.Context) error {
					client := mqtt.NewClient(c.String("url"), c.String("username"), c.String("password"), c.String("topic"))
					err := client.Connect()
					if err != nil {
						log.Fatal(err)
					}
					client.Handle("LongFast", lib.ChannelHandler(c.String("channel")))
					log.Info("Started")
					select {}
				},
			},
			{
				Name:    "radio",
				Aliases: []string{"r"},
				Usage:   "receive from radio",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "port",
						Value: "/dev/ttyUSB0",
						Usage: "serial port to connect to",
					},
				},
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					client, err := lib.ConnectSerial(c.String("port"), false)
					if err != nil {
						log.Fatal("error connecting to serial node", "err", err)
					}
					client.Handle(new(pb.MeshPacket), func(msg proto.Message) {
						pkt := msg.(*pb.MeshPacket)
						//log.Info(pkt.String())
						data := pkt.GetDecoded()
						//log.Info(processMessage(*data))
						log.Info(lib.ProcessMessage(data), "from", pkt.From, "portnum", data.Portnum.String())

					})

					err = client.Connect(ctx)
					if err != nil {
						log.Fatal(err)
					}
					log.Info("Started")
					select {}
				},
			},
		}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
