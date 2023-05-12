package cmd

import (
	"os"
	"sort"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Args:    cobra.ExactArgs(0),
	Short:   "List available contexts",
	Long: `List all contexts found in files or directories listed in $ks_PATH.
`,
	Run: func(cmd *cobra.Command, args []string) {
		flagVerbose := getBoolFlag(cmd, "verbose")

		// Load kubeconfig from file
		confPath := os.Getenv("KUBECONFIG")
		conf, err := loadKubeconfig([]string{confPath})
		handleFatalf(err, "Error loading config from %s: %v", confPath, err)

		if len(conf.Contexts) == 0 {
			infof("No contexts found. Please make sure your KSPATH is set correctly.")
		}

		// Print current context first
		printCtx(conf.CurrentContext+" (current)", conf.Contexts[conf.CurrentContext], flagVerbose)

		// Put remaining contexts in alphabetical order
		others := make([]string, 0)
		for name, _ := range conf.Contexts {
			if name != conf.CurrentContext {
				others = append(others, name)
			}
		}
		sort.Strings(others)

		// Print remaining contexts in order
		for _, name := range others {
			printCtx(name, conf.Contexts[name], flagVerbose)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("verbose", "v", false, "Print all context info")
}
