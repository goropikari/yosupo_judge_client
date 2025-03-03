package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	pb "github.com/goropikari/yosupo_judge_client/proto/librarychecker"
	"google.golang.org/protobuf/proto"
)

type ProtoMessage interface {
	ProtoMessage()
}

func NewClientFromProblemURL(rawurl string) (*client, error) {
	problemURL, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	paths := strings.Split(problemURL.Path, "/")
	problemName := paths[len(paths)-1]

	if problemURL.Host == "judge.yosupo.jp" {
		return newClient(clientOptions{
			tls:         true,
			hostname:    "v2.api.judge.yosupo.jp",
			port:        "443",
			problemName: problemName,
		}), nil
	}

	opts := clientOptions{
		tls:         false,
		hostname:    "127.0.0.1",
		port:        "12380",
		problemName: problemName,
	}

	return newClient(opts), nil
}

type client struct {
	cl          *http.Client
	tls         bool
	host        string
	port        string
	problemName string
}

type clientOptions struct {
	tls         bool
	hostname    string
	port        string
	problemName string
}

func newClient(opts clientOptions) *client {
	cl := http.DefaultClient
	if opts.tls {
		cl.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: opts.hostname,
			},
		}
	}

	return &client{
		cl:          http.DefaultClient,
		tls:         opts.tls,
		host:        opts.hostname,
		port:        opts.port,
		problemName: opts.problemName,
	}
}

// proto message を marshal したものを返却する
func (c *client) post(path string, msg proto.Message) ([]byte, error) {
	url, err := url.JoinPath(c.baseURL(), path)
	if err != nil {
		return nil, err
	}

	pbbuf, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 5)
	i := uint32(len(pbbuf))
	binary.BigEndian.PutUint32(buf[1:], i)
	buf = append(buf, pbbuf...)

	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/grpc-web+proto")
	req.Header.Set("x-grpc-web", "1")

	res, err := c.cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	pbres, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(pbres[1:5])
	return pbres[5 : 5+length], nil
}

func (c *client) baseURL() string {
	scheme := "http"
	if c.tls {
		scheme = "https"
	}
	return scheme + "://" + c.host + ":" + c.port
}

func (c *client) ProblemName() string {
	return c.problemName
}

func (c *client) ProblemInfo() (*pb.ProblemInfoResponse, error) {
	res, err := c.post(
		"/librarychecker.LibraryCheckerService/ProblemInfo",
		&pb.ProblemInfoRequest{Name: c.problemName},
	)
	if err != nil {
		return nil, err
	}

	pbres := &pb.ProblemInfoResponse{}
	if err := proto.Unmarshal(res, pbres); err != nil {
		return nil, err
	}

	return pbres, nil
}

func (c *client) Submit(filepath string, langID string) (*pb.SubmitResponse, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	req := &pb.SubmitRequest{
		Problem: c.problemName,
		Source:  string(data),
		Lang:    langID,
	}

	res, err := c.post(
		"/librarychecker.LibraryCheckerService/Submit",
		req,
	)
	if err != nil {
		return nil, err
	}

	pbres := &pb.SubmitResponse{}
	if err := proto.Unmarshal(res, pbres); err != nil {
		return nil, err
	}

	return pbres, nil
}
