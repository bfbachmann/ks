package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete <context> [context...]",
	Aliases: []string{"rm", "remove"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Delete contexts",
	Long: `This command deletes the given contexts from the kubeconfig pointed to by the KUBECONFIG env var. If any of 
the contexts being deleted are the current context, the current context will be set to empty.
`,
	Run: func(cmd *cobra.Command, args []string) {
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

		for _, name := range args {
			// Delete the context
			delete(conf.Contexts, name)

			// Update current context if necessary
			if conf.CurrentContext == name {
				conf.CurrentContext = ""
				warnf("Current context was deleted and has been set to empty.")
			}
		}

		// Write config to file
		err = writeKubeconfig(confPath, conf)
		handleFatalf(err, "Error writing config to %s: %v", confPath, err)

		infof("Deleted contexts %v.", args)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
