package cmd

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

		if strings.HasPrefix(url, "https://judge.yosupo.jp") {
			downloadTestcasesFromGoogleStorage(outdir, cl.ProblemName(), res.GetVersion(), res.GetTestcasesVersion())
		} else {
			downloadTestcasesFromLocal(outdir, cl.ProblemName(), res.GetTestcasesVersion())
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

func downloadTestcasesFromGoogleStorage(outdir string, problemName string, version string, testcasesVersion string) {
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
	handleError(err)

	res, err := cl.Do(req)
	handleError(err)

	bs, err := io.ReadAll(res.Body)
	handleError(err)

	info := &InfoToml{}
	_, err = toml.Decode(string(bs), info)
	handleError(err)

	testcnt := 0
	for _, test := range info.Tests {
		if strings.HasPrefix(test.Name, "example") {
			testcnt = test.Number
			break
		}
	}

	for _, suffix := range []string{"in", "out"} {
		for i := 0; i < testcnt; i++ {
			req, err := http.NewRequest(
				"GET",
				fmt.Sprintf(
					"https://storage.googleapis.com/v2-prod-library-checker-data-public/v3/%s/testcase/%s/%s/example_%02d.%s",
					problemName,
					testcasesVersion,
					suffix,
					i,
					suffix,
				),
				nil,
			)
			handleError(err)

			res, err := cl.Do(req)
			handleError(err)

			bs, err := io.ReadAll(res.Body)
			handleError(err)

			os.Mkdir(outdir, 0755)
			filename := fmt.Sprintf("%s/example_%02d.%s", outdir, i, suffix)
			fmt.Println(filename)
			if err := os.WriteFile(filename, bs, 0644); err != nil {
				handleError(err)
			}
		}
	}
}

func downloadTestcasesFromLocal(outdir string, problemName string, testcasesVersion string) {
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

	objectPaginator := s3.NewListObjectsV2Paginator(svc, input)
	for objectPaginator.HasMorePages() {
		output, err := objectPaginator.NextPage(context.Background())
		if err != nil {
			var noBucket *types.NoSuchBucket
			if errors.As(err, &noBucket) {
				log.Printf("Bucket %s does not exist.\n", "hoge")
				err = noBucket
			}
			break
		} else {
			for _, object := range output.Contents {
				keys := strings.Split(*object.Key, "/")
				name := keys[len(keys)-1]
				fmt.Println(name)

				result, err := svc.GetObject(context.Background(), &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    object.Key,
				})
				handleError(err)

				bs, err := io.ReadAll(result.Body)
				handleError(err)

				os.Mkdir(outdir, 0755)
				if err := os.WriteFile(fmt.Sprintf("%s/%s", outdir, name), bs, 0644); err != nil {
					handleError(err)
				}
			}
		}
	}
}
