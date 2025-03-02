package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type ProblemInfoRequest struct {
	Name string `json:"name"`
}

var problemInfoCmd = &cobra.Command{
	Use:   `probinfo '{"name":"<problem name>"}'`,
	Short: "problem info",
	Long:  `get problem info`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("invalid argument")
			cmd.Help()
			os.Exit(1)
			return
		}
		client, err := NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		req := ProblemInfoRequest{}
		if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		res, err := client.ProblemInfo(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		fmt.Println(res)
	},
}
