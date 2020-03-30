package handlers

import (
	"context"
	"fmt"
	"github.com/pytorchtw/go-db-services-starter/proto"
)

type GrpcHandler struct {
}

func (h *GrpcHandler) CreatePage(ctx context.Context, in *proto.CreatePageRequest) (*proto.CreatePageResponse, error) {
	return nil, nil
}

func (h *GrpcHandler) GetPage(ctx context.Context, in *proto.GetPageRequest) (*proto.GetPageResponse, error) {
	return nil, nil
}

func (h *GrpcHandler) SayHello(ctx context.Context, in *proto.SimpleRequest) (*proto.SimpleResponse, error) {
	return &proto.SimpleResponse{Message: fmt.Sprintf("%s", in.Message)}, nil
}
