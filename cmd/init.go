package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const shellInitScript = `
# Added by ks init. Do not edit.
[[ -f ${HOME}/.ks/init.sh ]] && source ${HOME}/.ks/init.sh`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initialize ks",
	Long: `This command will create the ${HOME}/.ks directory and add an initialization script to the user's 
shell initialization script (.bashrc, .zshrc, config.fish).
`,
	Run: func(cmd *cobra.Command, args []string) {
		force, err := cmd.Flags().GetBool("force")
		handleFatal(err, "Error getting --force flag: %v\n, err")

		if !force {
			// Abort if the .ks directory already exists (i.e. if we're already initialized)
			info, err := os.Stat(ksHomeDir)
			if err == nil && info.IsDir() {
				fatalf("Already initialized. Use --force flag to for reinitialization.")
			}
		}

		// Make sure the .ks directory exists
		err = os.MkdirAll(ksHomeDir, 0755)
		handleFatal(err, "Error creating %s: %v", ksHomeDir, err)

		// Get the user's default shell
		shell := os.Getenv("SHELL")
		shellBaseName := strings.TrimPrefix(filepath.Base(shell), "-")

		// Determine the init file path for some common shells
		shellInitFilePath := ""
		switch shellBaseName {
		case "bash":
			shellInitFilePath = homeDir + "/.bashrc"
		case "zsh":
			shellInitFilePath = homeDir + "/.zshrc"
		case "fish":
			shellInitFilePath = homeDir + ".config/fish/config.fish"
		default:
			fatalf("unknown shell: %s", shellBaseName)
		}

		// Open the shell init file
		initFile, err := os.OpenFile(shellInitFilePath, os.O_APPEND|os.O_WRONLY, 0644)
		handleFatal(err, "Error updating %s: %v", shellInitFilePath, err)
		defer initFile.Close()

		// Write some data to the file
		_, err = fmt.Fprintln(initFile, shellInitScript)
		handleFatal(err, "Error writing file %s: %v", shellInitFilePath, err)

		infof(`Initialized. Use "ks activate" to use kubeconfig generated from KSPATH for new shell sessions.
`)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("force", "f", false, "Force reinitialization")
}
