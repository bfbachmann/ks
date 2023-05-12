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
	Long: `This command will create the ${HOME}/.ks directory and add initialization code to the user's rc file (.bashrc, .zshrc, config.fish).
`,
	Run: func(cmd *cobra.Command, args []string) {
		if !getBoolFlag(cmd, "force") {
			// Abort if the .ks directory already exists (i.e. if we're already initialized)
			info, err := os.Stat(ksHomeDir)
			if err == nil && info.IsDir() {
				fatalf("Already initialized. Use --force flag to force reinitialization.")
			}
		}

		// Make sure the .ks directory exists
		err := os.MkdirAll(ksHomeDir, 0755)
		handleFatalf(err, "Error creating %s: %v", ksHomeDir, err)

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
			infof(`Unknown shell "%s". Changes would not be made automatically.`, shellBaseName)
			infof(`To initialize manually, place the following code at the bottom of your shell's equivalent of .bashrc.`)
			fatalf(shellInitScript)
		}

		// Open the shell init file
		initFile, err := os.OpenFile(shellInitFilePath, os.O_APPEND|os.O_WRONLY, 0644)
		handleFatalf(err, "Error updating %s: %v", shellInitFilePath, err)
		defer initFile.Close()

		// Write some data to the file
		_, err = fmt.Fprintln(initFile, shellInitScript)
		handleFatalf(err, "Error writing file %s: %v", shellInitFilePath, err)

		infof(
			`Initialized. Changes were made to %s. Use "ks activate" to use kubeconfig generated from KSPATH for new shell sessions.`,
			shellInitFilePath,
		)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("force", "f", false, "Force reinitialization")
}
