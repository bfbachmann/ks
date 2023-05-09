package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:     "activate",
	Aliases: []string{"a"},
	Short:   "Use kubeconfig generated using kubeconfig files from KSPATH for new shell sessions.",
	Long: `This command will set KUBECONFIG=${HOME}/.ks/config for the all future shell sessions.

Use "ks deactivate" to undo.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Write the init.sh file that will modify KUBECONFIG
		err := os.WriteFile(initPath, []byte(fmt.Sprintf("export KUBECONFIG=%s\n", masterConfigPath)), 0644)
		handleFatal(err, "Error writing %s: %v\n", initPath, err)

		infof("Activated. KUBECONFIG will be set to %s for future shell sessions.\n", masterConfigPath)
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
