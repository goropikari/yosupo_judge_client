package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Exp struct {
	Name string `json:"name"`
}

var problemInfoCmd = &cobra.Command{
	Use:   `probinfo <problem url>`,
	Short: "problem info",
	Long:  `get problem info`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("invalid argument")
			cmd.Help()
			os.Exit(1)
		}

		cl, err := NewClientFromProblemURL(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		res, err := cl.ProblemInfo()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		j, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		fmt.Print(string(j))
	},
}
