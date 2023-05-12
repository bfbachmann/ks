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
	Use:     "whence [context]",
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"w"},
	Short:   "List kubeconfig files in which contexts exist",
	Long: `This command prints the locations and contexts of all kubeconfig files found under KSPATH in order of loading 
precedence. If a context argument is provided, only paths in which that context exists will be printed.
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
					conf, err := clientcmd.LoadFromFile(currentPath)
					if err != nil {
						// We'll assume an error means this is not a valid kubeconfig file
						return nil
					}

					var currentMsg string
					if currentPath == os.Getenv("KUBECONFIG") {
						currentMsg = " (current)"
					}

					// Print contexts from the file if no specific context was listed. Otherwise, only print the name
					// of the file if it contains the given context.
					if len(args) == 0 {
						infof(currentPath + currentMsg)
						for ctxName, _ := range conf.Contexts {
							infof("  %s", ctxName)
						}
					} else if conf.Contexts[args[0]] != nil {
						infof(currentPath + currentMsg)
					}
				}

				return nil
			})
			handleFatalf(err, "Error loading config: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(whenceCmd)
}
