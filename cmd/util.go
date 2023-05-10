package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd/api/latest"
	"sigs.k8s.io/yaml"
)

// handleFatal prints the given message and exits with code 1 if the error is not nil. Otherwise, it does nothing.
func handleFatal(err error, format string, a ...any) {
	if err != nil {
		fatalf(format, a...)
	}
}

// fatalf prints an error message and immediately exits with code 1.
func fatalf(format string, a ...any) {
	infof(format, a...)
	os.Exit(1)
}

// infof prints an info message.
func infof(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

// printCtx prints information about a context.
func printCtx(name string, ctx *api.Context, verbose bool) {
	if verbose {
		infof("%s\n  Cluster: %s\n  Namespace: %s\n", name, ctx.Cluster, ctx.Namespace)
	} else {
		infof("%s", name)
	}
}

// loadKubeconfig loads all kubeconfig files at the given paths (can be files or dirs).
func loadKubeconfig(paths []string) (*api.Config, error) {
	// Search for kubeconfig files in each path
	loadingRules := clientcmd.ClientConfigLoadingRules{
		Precedence: make([]string, 0),
	}
	for _, path := range paths {
		// Replace "~" with home directory in path
		if strings.HasPrefix(path, "~/") {
			path = filepath.Join(homeDir, path[2:])
		}

		// Find all kubeconfig files in path
		if err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				// Check if this is a valid kubeconfig file by loading it
				_, err = clientcmd.LoadFromFile(currentPath)
				if err != nil {
					return nil
				}

				loadingRules.Precedence = append(loadingRules.Precedence, currentPath)
			}

			return nil
		}); err != nil {
			return nil, err
		}
	}

	// Load all the located kubeconfig files
	return loadingRules.Load()
}

// writeKubeconfig writes the given kubeconfig to a file at the given path.
func writeKubeconfig(path string, conf *api.Config) error {
	// Encode and write to file
	jsonBytes, err := runtime.Encode(latest.Codec, conf)
	if err != nil {
		return fmt.Errorf("error encoding merged kubeconfig as JSON: %v", err)
	}

	output, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		return fmt.Errorf("error converting merged JSON kubeconfig to YAML: %v", err)
	}

	err = os.WriteFile(path, output, 0644)
	if err != nil {
		return fmt.Errorf("error writing merged kubeconfig: %v", err)
	}

	return nil
}
