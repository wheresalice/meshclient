package lib

import (
	pb "github.com/meshnet-gophers/meshtastic-go/meshtastic"
	"google.golang.org/protobuf/proto"
)

func ProcessMessage(message *pb.Data) string {
	if message.Portnum == pb.PortNum_NODEINFO_APP {
		var user = pb.User{}
		proto.Unmarshal(message.Payload, &user)
		return user.String()
	}
	if message.Portnum == pb.PortNum_POSITION_APP {
		var pos = pb.Position{}
		proto.Unmarshal(message.Payload, &pos)
		return pos.String()
	}
	if message.Portnum == pb.PortNum_TELEMETRY_APP {
		var t = pb.Telemetry{}
		proto.Unmarshal(message.Payload, &t)
		return t.String()
	}
	if message.Portnum == pb.PortNum_NEIGHBORINFO_APP {
		var n = pb.NeighborInfo{}
		proto.Unmarshal(message.Payload, &n)
		return n.String()
	}
	if message.Portnum == pb.PortNum_STORE_FORWARD_APP {
		var s = pb.StoreAndForward{}
		proto.Unmarshal(message.Payload, &s)
		return s.String()
	}
	if message.Portnum == pb.PortNum_TEXT_MESSAGE_APP {
		return string(message.Payload)
	}

	return "unknown message type"
}
