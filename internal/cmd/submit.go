package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type SubmitRequest struct {
	Problem string `json:"problem"`
	Source  string `json:"source"`
	Lang    string `json:"lang"`
}

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit your code",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := NewClient()
		if err != nil {
			log.Fatal(err)
		}

		req := SubmitRequest{}
		if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		resp, err := client.Submit(req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(resp)
	},
}
