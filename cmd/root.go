package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Global variable shared across all subcommands
var sandboxName string

var rootCmd = &cobra.Command{
	Use:   "pod-sbx",
	Short: "Manages a Podman sandbox environment",
	Long: `Manages a Podman sandbox environment with a restricted egress network.
It creates a pod and injects an alpine-based init container running nftables 
to restrict outbound traffic to a specific IP/CIDR allowlist.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flag available to all subcommands
	rootCmd.PersistentFlags().StringVarP(&sandboxName, "name", "n", "sandbox", "The name of the sandbox pod")
}
