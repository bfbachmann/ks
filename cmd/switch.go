package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:     "switch",
	Aliases: []string{"s"},
	Short:   "Switch to a different context",
	Long: `Switch to a different context in one of the kubeconfig files under KSPATH.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fatalf(cmd.Long)
		}

		name := args[0]

		ns, err := cmd.Flags().GetString("namespace")
		handleFatal(err, "Error reading namespace flag: %v\n", err)

		// Load kubeconfig from file
		confPath := os.Getenv("KUBECONFIG")
		conf, err := loadKubeconfig([]string{confPath})
		handleFatal(err, "Error loading config from %s: %v\n", confPath, err)

		ctx, ok := conf.Contexts[name]
		if !ok {
			fatalf("No such context: %s\n", name)
		}

		conf.CurrentContext = name
		if ns != "" {
			ctx.Namespace = ns
		}

		err = writeKubeconfig(masterConfigPath, conf)
		handleFatal(err, "Error writing config: %v\n", err)

		infof("Switched to context %s (namespace: %s)\n", name, ctx.Namespace)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().StringP("namespace", "n", "", "The namespace to use in the context")
}
