package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:     "deactivate",
	Aliases: []string{"d"},
	Args:    cobra.ExactArgs(0),
	Short:   "Return to regular KUBECONFIG for new shell sessions",
	Long: `This command will make it so ks has no effect on your KUBECONFIG. You may still use ks commands to view and
manage contexts and namespaces for your current KUBECONFIG.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := os.RemoveAll(initPath)
		handleFatalf(err, "Error removing %s: %v", initPath, err)

		infof("Deactivated. Your KUBECONFIG will take its normal value for future shell sessions.")
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)
}
