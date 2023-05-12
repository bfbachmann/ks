package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:     "rename <old-name> <new-name>",
	Aliases: []string{"r"},
	Args:    cobra.ExactArgs(2),
	Short:   "Rename an existing context",
	Long: `This command changes the name assigned to an existing context in the kubeconfig pointed to by the KUBECONFIG
env var. If this context is the current context, the current context will be also updated.
`,
	Run: func(cmd *cobra.Command, args []string) {
		argOldName, argNewName := args[0], args[1]

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

		// Make sure a context exists with the old name
		ctx, exists := conf.Contexts[argOldName]
		if !exists {
			fatalf(`No context exists with name %s`, argOldName)
		}

		// Make sure no context exists with the new name
		if _, conflict := conf.Contexts[argNewName]; conflict {
			fatalf(`A context already exists with the name %s`, argNewName)
		}

		// Assign new name to context
		delete(conf.Contexts, argOldName)
		conf.Contexts[argNewName] = ctx

		// Update current context if necessary
		if conf.CurrentContext == argOldName {
			conf.CurrentContext = argNewName
		}

		// Write config to file
		err = writeKubeconfig(confPath, conf)
		handleFatalf(err, "Error writing config to %s: %v", confPath, err)

		infof("Context %s renamed to %s.", argOldName, argNewName)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
