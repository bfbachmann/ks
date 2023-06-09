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
	kubeconfigPaths  []string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ks",
	Short: "Quickly navigate kubectl config files, contexts, and namespaces.",
	Long: `Use the KSPATH environment variable to list files and directories in which to search for kubeconfig files. 
By default, KSPATH will be set to ${HOME}/.kube.

Example: KSPATH="~/.kube:~/code/my-project/conf:~/clusters/local.yaml"

This program, when run, will find all valid kubeconfig files in the paths specified in KSPATH, merge them, and write 
them to ${HOME}/.ks/config. Higher precedence is given to files or directories that appear closer to the beginning of 
KSPATH. Existing config at ${HOME}/.ks/config will always get lowest precedence, but the current context and namespace
listed there will persist unless changed manually.
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

	// Abort if we're not initialized (i.e. the .ks directory doesn't exist)
	info, err := os.Stat(ksHomeDir)
	if err != nil {
		if !os.IsNotExist(err) {
			fatalf("Error checking %s: %v", ksHomeDir, err)
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

	// First we need to get the current context and namespace, so we can make sure not to overwrite it later
	var (
		currentCtxName string
		currentNs      string
	)
	if existingConfPath := os.Getenv("KUBECONFIG"); existingConfPath != "" {
		// Make sure the file still exists before trying to load it. If it doesn't we'll just skip this step since
		// there is no current context in this case.
		_, err = os.Stat(existingConfPath)
		if err != nil && !os.IsNotExist(err) {
			// Some unexpected error occurred
			fatalf("Error checking for config at %s: %v.", existingConfPath, err)
		} else if err == nil {
			// The file exists
			conf, err := loadKubeconfig([]string{existingConfPath})
			handleFatalf(err, "Error loading existing config from %s: %v", existingConfPath, err)

			currentCtxName = conf.CurrentContext
			if ctx, ok := conf.Contexts[currentCtxName]; ok {
				currentNs = ctx.Namespace
			}
		}
	}

	// Parse paths from envvar, appending the master config file if it exists
	kubeconfigPaths = strings.Split(ksPath, ":")
	_, err = os.Stat(masterConfigPath)
	if err != nil && !os.IsNotExist(err) {
		fatalf("Error checking file %s: %v", masterConfigPath, err)
	} else if err == nil {
		kubeconfigPaths = append(kubeconfigPaths, masterConfigPath)
	}

	// Load kubeconfig
	conf, err := loadKubeconfig(kubeconfigPaths)
	handleFatalf(err, "Error loading config: %v", err)

	// Make sure we restore the current context and namespace, if the context still exists. Otherwise, print a warning
	// message to let the user know that their current context has changed.
	if ctx, exists := conf.Contexts[currentCtxName]; exists {
		conf.CurrentContext = currentCtxName
		ctx.Namespace = currentNs
	} else {
		var newNs string
		if newCtx, ok := conf.Contexts[conf.CurrentContext]; ok {
			newNs = newCtx.Namespace
		}
		warnf(
			`Context "%s" no longer exists. Current context is now "%s" (namespace: "%s").`,
			currentCtxName,
			conf.CurrentContext,
			newNs,
		)
	}

	// Encode and write to file
	err = writeKubeconfig(masterConfigPath, conf)
	handleFatalf(err, "Error writing config: %v", err)
}
