package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// whenceCmd represents the whence command
var whenceCmd = &cobra.Command{
	Use:     "whence",
	Args:    cobra.ExactArgs(0),
	Aliases: []string{"w"},
	Short:   "List all available kubeconfig files",
	Long: `This command reveals the locations of all kubeconfig files available under KSPATH in order of loading 
precedence.
`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, path := range kubeconfigPaths {
			// Replace "~" with home directory in path
			if strings.HasPrefix(path, "~/") {
				path = filepath.Join(homeDir, path[2:])
			}

			// Find all kubeconfig files in path
			err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() {
					// Check if this is a valid kubeconfig file by loading it
					_, err = clientcmd.LoadFromFile(currentPath)
					if err != nil {
						return nil
					}

					infof(currentPath)
				}

				return nil
			})
			handleFatal(err, "Error loading config: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(whenceCmd)
}
