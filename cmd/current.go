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
		flagVerbose := getBoolFlag(cmd, "verbose")

		// Load kubeconfig from file
		confPath := os.Getenv("KUBECONFIG")
		if confPath == "" {
			fatalf("KUBECONFIG is not set. Please make sure KUBECONFIG is set correctly.")
		}

		conf, err := loadKubeconfig([]string{confPath})
		handleFatalf(
			err,
			"Error loading config from %s: %v. Please make sure KUBECONFIG is set correctly.",
			confPath,
			err,
		)

		if len(conf.Contexts) == 0 {
			infof("No contexts found. Please make sure your KSPATH is set correctly.")
		}

		// Print current context
		printCtx(conf.CurrentContext, conf.Contexts[conf.CurrentContext], flagVerbose)
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
	currentCmd.Flags().BoolP("verbose", "v", false, "Print all context info")
}
