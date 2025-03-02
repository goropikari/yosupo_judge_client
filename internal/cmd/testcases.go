package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/cobra"
)

type DownloadProblemTestCases struct {
	Name string `json:"name"`
}

var downloadTestCasesCmd = &cobra.Command{
	Use:   `download-test`,
	Short: "download test cases",
	Long:  `doanload test cases`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("invalid argument")
			cmd.Help()
			os.Exit(1)
			return
		}
		s3Endpoint := os.Getenv("S3_ENDPOINT")
		if s3Endpoint == "" {
			s3Endpoint = "http://127.0.0.1:9000"
		}

		client, err := NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		req := &struct {
			Name   string `json:"name"`
			Outdir string `json:"outdir"`
		}{}
		if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		res, err := client.ProblemInfo(ProblemInfoRequest{Name: req.Name})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		problemInfo := struct {
			Title            string `json:"title"`
			SourceUrl        string `json:"sourceUrl"`
			TimeLimit        int    `json:"timeLimit"`
			Version          string `json:"version"`
			TestcasesVersion string `json:"testcasesVersion"`
		}{}
		json.Unmarshal([]byte(res), &problemInfo)

		const bucket = "testcase-public"

		svc := s3.New(s3.Options{
			Region:       "ap-northeast-1",
			BaseEndpoint: aws.String(s3Endpoint),
		})
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(fmt.Sprintf("v3/%s/testcase/%s", req.Name, problemInfo.TestcasesVersion)),
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
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
						return
					}
					bs, err := io.ReadAll(result.Body)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
						return
					}
					os.Mkdir(req.Outdir, 0755)
					if err := os.WriteFile(fmt.Sprintf("%s/%s", req.Outdir, name), bs, 0644); err != nil {
						fmt.Println(err)
						os.Exit(1)
						return
					}
				}
			}
		}
	},
}
