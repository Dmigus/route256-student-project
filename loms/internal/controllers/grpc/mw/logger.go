package mw

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
)

func LogReqAndResp(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	raw, _ := protojson.Marshal((req).(proto.Message))
	log.Printf("request: method: %v, req: %v\n", info.FullMethod, string(raw))

	if resp, err = handler(ctx, req); err != nil {
		log.Printf("response: method: %v, err: %v\n", info.FullMethod, err)
		return
	}

	rawResp, _ := protojson.Marshal((resp).(proto.Message))
	log.Printf("response: method: %v, resp: %v\n", info.FullMethod, string(rawResp))

	return
}
