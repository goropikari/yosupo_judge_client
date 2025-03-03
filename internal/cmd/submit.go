package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type SubmitRequest struct {
	Problem string `json:"problem"`
	Source  string `json:"source"`
	Lang    string `json:"lang"`
}

var submitCmd = &cobra.Command{
	Use:   "submit <problem url> <file path> [lang id]",
	Short: "Submit your code",
	Long:  "Submit your code. language id is optional. default is cpp.\nref: https://github.com/yosupo06/library-checker-judge/blob/6efba7d2120e85dee1f96507c2c535fa5f303b50/langs/langs.toml",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			return
		}

		langID := "cpp"
		if len(args) == 3 {
			langID = args[2]
		}

		url := args[0]
		cl, err := NewClientFromProblemURL(url)
		if err != nil {
			handleError(err)
			return
		}

		path := args[1]

		res, err := cl.Submit(path, langID)
		if err != nil {
			handleError(err)
			return
		}

		fmt.Print(res.GetId())
	},
}
