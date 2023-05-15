# ks

`ks` (short for "kube switch") is a simple command-line tool for managing and switching between Kubernetes contexts
and namespaces. It is primarily designed to handle cases where the user has multiple and frequently-changing kubeconfig
files.

## Installing

Find the appropriate binary for your system under [releases](https://github.com/bfbachmann/ks/releases).

You can also build it yourself using `make install` (`go build && go install`).

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

```
Use the KSPATH environment variable to list files and directories in which to search for kubeconfig files.
By default, KSPATH will be set to ${HOME}/.kube.

Example: KSPATH="~/.kube:~/code/my-project/conf:~/clusters/local.yaml"

This program, when run, will find all valid kubeconfig files in the paths specified in KSPATH, merge them, and write
them to ${HOME}/.ks/config. Higher precedence is given to files or directories that appear closer to the beginning of
KSPATH. Existing config at ${HOME}/.ks/config will always get lowest precedence, but the current context and namespace
listed there will persist unless changed manually.

Usage:
  ks [command]

Available Commands:
  activate    Use kubeconfig generated using kubeconfig files from KSPATH for new shell sessions
  completion  Generate the autocompletion script for the specified shell
  current     Show the current context
  deactivate  Return to regular KUBECONFIG
  delete      Delete contexts
  help        Help about any command
  init        Initialize ks
  list        List available contexts
  new         Create a new context
  rename      Rename an existing context
  switch      Switch to a different context
  whence      List kubeconfig files in which contexts exist

Flags:
  -h, --help   help for ks

Use "ks [command] --help" for more information about a command.
```

## Usage Examples

Here's an example of how you might set up and use `ks` if you use Bash (Zsh is similar).

```shell
# Initialize ks
ks init

# Optionally set KSPATH (defaults to ~/.kube)
echo "export KSPATH=~/project/k8s:~/cluster/config.yaml" >> ~/.bashrc

# Activate ks so all new shell sessions use ks-managed config
ks activate

# Reload shell for changes to take effect in current session
source ~/.bashrc

# List all available contexts
ks list -v

# Switch context (and namespace)
ks switch new-context -n new-namespace

# View current context
ks current -v

# Rename a context
ks rename new-context my-favorite
```
