package cmd

import (
	"fmt"
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-spec-parse/pkg/parse"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "carapce-spec-parse",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO other parsers
		c := parse.Bazel(cmd.Flag("name").Value.String(), cmd.Flag("description").Value.String(), os.Stdin)
		m, err := yaml.Marshal(c)
		if err != nil {
			return err
		}
		fmt.Println(string(m))
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
func init() {
	rootCmd.Flags().StringP("name", "n", "", "name of the command")
	rootCmd.Flags().StringP("description", "d", "", "description of the command")
	rootCmd.Flags().StringP("parent", "p", "", "parent of the command")
	rootCmd.Flags().StringP("input", "i", "", "input format")
	rootCmd.Flags().StringP("output", "o", "", "output format")

	carapace.Gen(rootCmd).FlagCompletion(carapace.ActionMap{
		"input":  carapace.ActionValues("bazel"),
		"output": carapace.ActionValues("spec"), // TODO code
	})
}
