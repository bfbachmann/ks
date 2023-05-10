package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:     "deactivate",
	Aliases: []string{"d"},
	Short:   "Return to regular KUBECONFIG",
	Long: `This command will make it so ks has no effect on your KUBECONFIG.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := os.RemoveAll(initPath)
		handleFatal(err, "Error removing %s: %v", initPath, err)

		infof("Deactivated. Your KUBECONFIG will take its normal value for future shell sessions.")
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)
}
