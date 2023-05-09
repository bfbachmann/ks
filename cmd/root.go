package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

var (
	homeDir          string
	ksHomeDir        string
	masterConfigPath string
	initPath         string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ks",
	Short: "Quickly navigate kubectl config files, contexts, and namespaces.",
	Long: `Use the KSPATH environment variable to list files and directories in which to search for kubeconfig files. 
By default, KSPATH will be set to ${HOME}/.kube.

Example: KSPATH="~/.kube:~/code/my-project/conf:~/clusters/local.yaml"

This program, when run, will also find all valid kubeconfig files in the paths specified in KSPATH, merge them, 
and write them to ${HOME}/.ks/config.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize package-level vars
	homeDir = homedir.HomeDir()
	ksHomeDir = homeDir + "/.ks"
	masterConfigPath = ksHomeDir + "/config"
	initPath = ksHomeDir + "/init.sh"

	// Abort if we're not initialized (i.e. the .ks director doesn't exist)
	info, err := os.Stat(ksHomeDir)
	if err != nil {
		if !os.IsNotExist(err) {
			fatalf("Error checking %s: %v\n", ksHomeDir, err)
		}

		// Not initialized, abort
		return
	} else if !info.IsDir() {
		// Not initialized, abort
		return
	}

	// Load path from env with default
	ksPath := os.Getenv("KSPATH")
	if ksPath == "" {
		ksPath = "~/.kube"
	}

	// Parse paths from envvar appending the master config file if it exists
	paths := strings.Split(ksPath, ":")
	_, err = os.Stat(masterConfigPath)
	if err != nil && !os.IsNotExist(err) {
		fatalf("Error checking file %s: %v\n", masterConfigPath, err)
	} else if err == nil {
		paths = append(paths, masterConfigPath)
	}

	// Load kubeconfig
	conf, err := loadKubeconfig(paths)
	handleFatal(err, "Error loading config: %v\n", err)

	// Encode and write to file
	err = writeKubeconfig(masterConfigPath, conf)
	handleFatal(err, "Error writing config: %v\n", err)
}
