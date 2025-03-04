package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
)

var downloadTestCasesCmd = &cobra.Command{
	Use:   `download-test <problem url> <output dir>`,
	Short: "download test cases",
	Long:  `doanload test cases`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("invalid argument")
			cmd.Help()
			os.Exit(1)
		}

		url := args[0]
		outdir := args[1]

		cl, err := NewClientFromProblemURL(url)
		handleError(err)

		res, err := cl.ProblemInfo()
		handleError(err)

		var testcases []testcase
		if strings.HasPrefix(url, "https://judge.yosupo.jp") {
			testcases, err = downloadTestcasesFromGoogleStorage(cl.ProblemName(), res.GetVersion(), res.GetTestcasesVersion())
			handleError(err)
		} else {
			testcases, err = downloadTestcasesFromLocal(cl.ProblemName(), res.GetTestcasesVersion())
			handleError(err)
		}

		if err := os.MkdirAll(outdir, 0755); err != nil {
			handleError(err)
		}

		for _, testcase := range testcases {
			fmt.Println(testcase.name)
			filename := fmt.Sprintf("%s/%s", outdir, testcase.name)
			if err := os.WriteFile(filename, testcase.data, 0644); err != nil {
				handleError(err)
			}
		}
	},
}

type TestCaseInfo struct {
	Name   string `toml:"name"`
	Number int    `toml:"number"`
}

type InfoToml struct {
	Tests []TestCaseInfo `toml:"tests"`
}

type testcase struct {
	name string
	data []byte
}

func downloadTestcasesFromGoogleStorage(problemName string, version string, testcasesVersion string) ([]testcase, error) {
	cl := http.DefaultClient
	cl.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName: "storage.googleapis.com",
		},
	}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://storage.googleapis.com/v2-prod-library-checker-data-public/v3/%s/files/%s/info.toml",
			problemName,
			version,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	info := &InfoToml{}
	if _, err := toml.Decode(string(bs), info); err != nil {
		return nil, err
	}

	testcnt := 0
	for _, test := range info.Tests {
		if strings.HasPrefix(test.Name, "example") {
			testcnt = test.Number
			break
		}
	}

	testcases := make([]testcase, 0)
	for _, suffix := range []string{"in", "out"} {
		for i := 0; i < testcnt; i++ {
			filename := fmt.Sprintf("example_%02d.%s", i, suffix)
			req, err := http.NewRequest(
				"GET",
				fmt.Sprintf(
					"https://storage.googleapis.com/v2-prod-library-checker-data-public/v3/%s/testcase/%s/%s/%s",
					problemName,
					testcasesVersion,
					suffix,
					filename,
				),
				nil,
			)
			if err != nil {
				return nil, err
			}

			res, err := cl.Do(req)
			if err != nil {
				return nil, err
			}

			bs, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}

			testcases = append(testcases, testcase{name: filename, data: bs})
		}
	}

	return testcases, nil
}

func downloadTestcasesFromLocal(problemName string, testcasesVersion string) ([]testcase, error) {
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	if s3Endpoint == "" {
		s3Endpoint = "http://127.0.0.1:9000"
	}

	const bucket = "testcase-public"

	svc := s3.New(s3.Options{
		Region:       "ap-northeast-1",
		BaseEndpoint: aws.String(s3Endpoint),
	})
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(fmt.Sprintf("v3/%s/testcase/%s", problemName, testcasesVersion)),
	}

	testcases := make([]testcase, 0)
	objectPaginator := s3.NewListObjectsV2Paginator(svc, input)
	for objectPaginator.HasMorePages() {
		output, err := objectPaginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		} else {
			for _, object := range output.Contents {
				keys := strings.Split(*object.Key, "/")
				name := keys[len(keys)-1]

				result, err := svc.GetObject(context.Background(), &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    object.Key,
				})
				if err != nil {
					return nil, err
				}

				bs, err := io.ReadAll(result.Body)
				if err != nil {
					return nil, err
				}

				testcases = append(testcases, testcase{name: name, data: bs})
			}
		}
	}

	return testcases, nil
}
