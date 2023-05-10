# ks

`ks` (short for "kube switch") is a simple command line utility for managing and switching between Kubernetes contexts
and namespaces. It is primarily designed to handle cases where the user has multiple and frequently-changing kubeconfig
files.

## How It Works

The main function of `ks` is to find all kubeconfig files listed in the environment variable `KSPATH` (defaults to
`${HOME}/.kube`) and merge them into one file at `${HOME}/.ks/config`. `ks`, when active, will ensure the merged 
kubeconfig file stays up to date and will change the `KUBECONFIG` environment variable to point to it. It then offers
a simple command, `ks switch`, to switch between available contexts and namespaces.

This has a few benefits:
1. Original kubeconfig files are never changed, moved, or deleted.
2. No need to constantly change the `KUBECONFIG` manually.
3. Any newly-added, edited, or deleted kubeconfig under `KSPATH` will automatically be available on execution of any 
`ks` command.

## Basic Usage

Here's an example of how you might set up and use `ks` if you use Bash (Zsh is similar).

```shell
# Initialize ks
ks init

# Optionally set KSPATH (defaults to ~/.kube)
echo "export KSPATH=~/project/k8s:~/cluster/config.yaml" >> .bashrc

# Activate ks so all new shell sessions use ks-managed config
ks activate

# Reload shell for changes to take effect in current session
source ~/.bashrc

# List all available contexts
ks list -v

# Switch context (and namespace)
ks switch new-context -n new-namespace
```
