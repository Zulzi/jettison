package serverclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/Zulzi/jettison/errors"
	"github.com/Zulzi/jettison/example/examplepb"
	jetgrpc "github.com/Zulzi/jettison/grpc"
	"github.com/Zulzi/jettison/j"
	"github.com/Zulzi/jettison/log"
)

type Client struct {
	conn   *grpc.ClientConn
	client examplepb.HopperClient
}

func NewClient(url string) *Client {
	conn, err := grpc.Dial(url, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(jetgrpc.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(jetgrpc.StreamClientInterceptor))
	if err != nil {
		panic(fmt.Errorf("grpc.Dial error: %s", err.Error()))
	}

	return &Client{
		conn:   conn,
		client: examplepb.NewHopperClient(conn),
	}
}

func (cl *Client) Hop(ctx context.Context, hops int64) error {
	if hops <= 0 {
		return errors.New("no hops")
	}

	log.Info(ctx, "serverclient: Performing a hop to another server",
		j.KV("hops", hops))

	_, err := cl.client.Hop(ctx, &examplepb.HopRequest{Hops: hops - 1})
	return err
}

func (cl *Client) Close() error {
	return cl.conn.Close()
}
