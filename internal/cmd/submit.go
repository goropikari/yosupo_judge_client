package cmd

import (
	"fmt"
	"net/url"
	"os"

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
			fmt.Println("invalid argument")
			cmd.Help()
			os.Exit(1)
		}

		langID := "cpp"
		if len(args) == 3 {
			langID = args[2]
		}

		rawurl := args[0]
		cl, err := NewClientFromProblemURL(rawurl)
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

		u, err := url.Parse(rawurl)
		if err != nil {
			handleError(err)
		}

		fmt.Printf("%s://%s/submission/%d", u.Scheme, u.Host, res.GetId())
	},
}
