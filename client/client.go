package client

import (
	"context"

	nwPB "github.com/SoumyadipPayra/NightsWatchProtobufs/gogenproto/nightswatch"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type NightsWatchInstallationClient interface {
	Register(ctx context.Context, req *nwPB.RegisterRequest) error
	Close() error
}

type NightsWatchInitClient interface {
	Login(ctx context.Context, req *nwPB.LoginRequest) error
	SendDeviceData(ctx context.Context, username string, req *nwPB.DeviceDataRequest) error
	Close() error
}

const (
	defaultTarget = "localhost:50051"
)

var token string

type nightsWatchClientImpl struct {
	conn   *grpc.ClientConn
	client nwPB.NightsWatchServiceClient
}

func NewNightsWatchInstallationClient() (NightsWatchInstallationClient, error) {
	conn, err := grpc.NewClient(defaultTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &nightsWatchClientImpl{
		conn:   conn,
		client: nwPB.NewNightsWatchServiceClient(conn),
	}, nil
}

func NewNightsWatchInitClient() (NightsWatchInitClient, error) {
	conn, err := grpc.NewClient(defaultTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &nightsWatchClientImpl{
		conn:   conn,
		client: nwPB.NewNightsWatchServiceClient(conn),
	}, nil
}

func (c *nightsWatchClientImpl) Login(ctx context.Context, req *nwPB.LoginRequest) error {
	resp, err := c.client.Login(ctx, req)
	if err != nil {
		return err
	}
	token = resp.Token
	return nil
}

func (c *nightsWatchClientImpl) Register(ctx context.Context, req *nwPB.RegisterRequest) error {
	_, err := c.client.Register(ctx, req)
	return err
}

func (c *nightsWatchClientImpl) SendDeviceData(ctx context.Context, username string, req *nwPB.DeviceDataRequest) error {
	outgoingContext := metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"jwt_token": token,
		"user_name": username,
	}))
	_, err := c.client.SendDeviceData(outgoingContext, req)
	return err
}

func (c *nightsWatchClientImpl) Close() error {
	return c.conn.Close()
}
