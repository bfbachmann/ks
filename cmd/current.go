package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:     "current",
	Aliases: []string{"c"},
	Short:   "Show the current context",
	Long: `Print information about the current context
`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, err := cmd.Flags().GetBool("verbose")
		handleFatal(err, "Error getting verbose flag: %v\n", err)

		// Load kubeconfig from file
		confPath := os.Getenv("KUBECONFIG")
		conf, err := loadKubeconfig([]string{confPath})
		handleFatal(err, "Error loading config from %s: %v\n", confPath, err)

		if len(conf.Contexts) == 0 {
			infof("No contexts found. Please make sure your KSPATH is set correctly.\n")
		}

		// Print current context
		printCtx(conf.CurrentContext, conf.Contexts[conf.CurrentContext], verbose)
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
	currentCmd.Flags().BoolP("verbose", "v", false, "Print all context info")
}
