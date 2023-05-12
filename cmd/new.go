package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd/api"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:     "new <context>",
	Aliases: []string{"n", "create"},
	Args:    cobra.ExactArgs(1),
	Short:   "Create a new context",
	Long: `This command creates a new context in the kubeconfig pointed to by the KUBECONFIG env var.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get args and flags
		argName := args[0]
		flagFrom := getStringFlag(cmd, "from")
		flagCluster := getStringFlag(cmd, "cluster")
		flagUser := getStringFlag(cmd, "user")
		flagNamespace := getStringFlag(cmd, "namespace")

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

		// Create the new context
		newCtx := &api.Context{}
		if flagFrom != "" {
			// Copy from existing context
			existingCtx, exists := conf.Contexts[flagFrom]
			if !exists {
				fatalf(`No context exists with name %s`, flagFrom)
			}

			newCtx = existingCtx.DeepCopy()
		}

		// Set fields on new context
		if flagCluster != "" {
			if conf.Clusters[flagCluster] == nil {
				fatalf("No cluster exists with name %s", flagCluster)
			}
			newCtx.Cluster = flagCluster
		}
		if flagUser != "" {
			if conf.AuthInfos[flagUser] == nil {
				fatalf("No user exists with name %s", flagCluster)
			}
			newCtx.AuthInfo = flagUser
		}
		if flagNamespace != "" {
			newCtx.Namespace = flagNamespace
		}

		// Add new context to config
		conf.Contexts[argName] = newCtx

		// Write config to file
		err = writeKubeconfig(confPath, conf)
		handleFatalf(err, "Error writing config to %s: %v", confPath, err)

		infof("Created context %s.", argName)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("from", "f", "", "Copy from existing context by name")
	newCmd.Flags().StringP("cluster", "c", "", "The cluster for the new context")
	newCmd.Flags().StringP("user", "u", "", "The user for the new context")
	newCmd.Flags().StringP("namespace", "n", "", "The namespace for the new context")
}
