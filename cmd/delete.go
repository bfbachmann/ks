package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete <name>",
	Aliases: []string{"rm", "remove"},
	Args:    cobra.ExactArgs(1),
	Short:   "Delete a context",
	Long: `This command deletes the given context from the kubeconfig pointed to by the KUBECONFIG env var. If this 
context is the current context, the current context will be set to empty.
`,
	Run: func(cmd *cobra.Command, args []string) {
		argName := args[0]

		// Load kubeconfig from file
		confPath := os.Getenv("KUBECONFIG")
		if confPath == "" {
			fatalf("KUBECONFIG is not set. Please make sure KUBECONFIG is set correctly.")
		}

		conf, err := loadKubeconfig([]string{confPath})
		handleFatal(
			err,
			"Error loading config from %s: %v. Please make sure KUBECONFIG is set correctly.",
			confPath,
			err,
		)

		// Delete the context
		delete(conf.Contexts, argName)

		// Update current context if necessary
		if conf.CurrentContext == argName {
			conf.CurrentContext = ""
		}

		// Write config to file
		err = writeKubeconfig(confPath, conf)
		handleFatal(err, "Error writing config to %s: %v", confPath, err)

		infof("Context %s deleted.", argName)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
