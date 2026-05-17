package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Tear down the sandbox pod and firewall container.",
	Run: func(cmd *cobra.Command, args []string) {
		teardown(sandboxName)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}

func teardown(name string) {
	fmt.Printf("Tearing down existing '%s' environment\n", name)

	firewallName := fmt.Sprintf("init-firewall-%s", name)

	rmFirewall := exec.Command("podman", "rm", "-f", firewallName)
	rmFirewall.Run()

	rmPod := exec.Command("podman", "pod", "rm", "-f", name)
	rmPod.Run()

	fmt.Println("Teardown complete.")
}
