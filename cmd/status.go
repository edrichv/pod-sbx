package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the current status of the sandbox environment.",
	Run: func(cmd *cobra.Command, args []string) {
		checkStatus(sandboxName)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func checkStatus(name string) {
	cmd := exec.Command("podman", "pod", "inspect", name, "--format", "{{.State}}")
	out, err := cmd.Output()

	if err == nil {
		state := strings.TrimSpace(string(out))
		fmt.Printf("Sandbox '%s' Status: %s\n", name, state)
	} else {
		fmt.Printf("Sandbox '%s' Status: Not running / Does not exist\n", name)
	}
}
