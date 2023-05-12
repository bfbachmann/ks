package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const example = `
  ks switch my-context -n my-namespace  # switch to context "my-context" with namespace "my-namespace"
  ks switch my-other-context	        # switch to context "my-other-context" with default/existing namespace
  ks switch -n my-namespace             # use "my-namespace" in the current context
`

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:     "switch [name]",
	Aliases: []string{"s"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Switch to a different context",
	Long: `Switch to a different context and/or namespace in one of the kubeconfig files under KSPATH.
`,
	Example: strings.TrimLeft(example, "\n"),
	Run: func(cmd *cobra.Command, args []string) {
		flagNamespace := getStringFlag(cmd, "namespace")

		// Load kubeconfig from file
		confPath := os.Getenv("KUBECONFIG")
		conf, err := loadKubeconfig([]string{confPath})
		handleFatalf(err, "Error loading config from %s: %v", confPath, err)

		// Get context to switch to, defaulting to current context if one was not specified
		ctxName := conf.CurrentContext
		if len(args) > 0 {
			ctxName = args[0]
		}

		ctx, ok := conf.Contexts[ctxName]
		if !ok {
			fatalf("No such context: %s", ctxName)
		}

		// Set current context and namespace, if specified
		conf.CurrentContext = ctxName
		if flagNamespace != "" {
			ctx.Namespace = flagNamespace
		}

		// Write updated config to file
		err = writeKubeconfig(masterConfigPath, conf)
		handleFatalf(err, "Error writing config: %v", err)
		infof(`Switched to context "%s" (namespace: "%s")`, ctxName, ctx.Namespace)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().StringP("namespace", "n", "", "The namespace to use in the context")
}
