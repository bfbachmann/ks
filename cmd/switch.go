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
	Use:     "switch",
	Aliases: []string{"s"},
	Short:   "Switch to a different context",
	Long: `Switch to a different context and/or namespace in one of the kubeconfig files under KSPATH.
`,
	Example: strings.TrimLeft(example, "\n"),
	Run: func(cmd *cobra.Command, args []string) {
		ns, err := cmd.Flags().GetString("namespace")
		handleFatal(err, "Error reading namespace flag: %v", err)

		// Load kubeconfig from file
		confPath := os.Getenv("KUBECONFIG")
		conf, err := loadKubeconfig([]string{confPath})
		handleFatal(err, "Error loading config from %s: %v", confPath, err)

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
		if ns != "" {
			ctx.Namespace = ns
		}

		// Write updated config to file
		err = writeKubeconfig(masterConfigPath, conf)
		handleFatal(err, "Error writing config: %v", err)
		infof(`Switched to context "%s" (namespace: "%s")`, ctxName, ctx.Namespace)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().StringP("namespace", "n", "", "The namespace to use in the context")
}
