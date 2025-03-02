package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ktr0731/grpc-web-go-client/grpcweb"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/goropikari/yosupo_judge_client/proto/librarychecker"
)

type Client struct {
	client *grpcweb.ClientConn
}

func NewClient() (*Client, error) {
	host := os.Getenv("YOSUPO_JUDGE_HOST")
	port := os.Getenv("YOSUPO_JUDGE_PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "12380"
	}

	addr := host + ":" + port
	client, err := grpcweb.DialContext(addr, grpcweb.WithDefaultCallOptions())
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

func (c *Client) ProblemInfo(req ProblemInfoRequest) (string, error) {
	in, out := new(pb.ProblemInfoRequest), new(pb.ProblemInfoResponse)
	in.Name = req.Name

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()
	if err := c.client.Invoke(
		ctx,
		"/librarychecker.LibraryCheckerService/ProblemInfo",
		in,
		out,
	); err != nil {
		return "", err
	}

	bs, err := protojson.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}
	return string(bs), nil
}

func (c *Client) Submit(req SubmitRequest) (string, error) {
	in, out := new(pb.SubmitRequest), new(pb.SubmitResponse)
	in.Problem = req.Problem
	in.Source = req.Source
	in.Lang = req.Lang

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()
	if err := c.client.Invoke(
		ctx,
		"/librarychecker.LibraryCheckerService/Submit",
		in,
		out,
	); err != nil {
		return "", err
	}

	bs, err := protojson.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}
	return string(bs), nil
}
