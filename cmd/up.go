package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	allowlist string
	podArgs   string
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Initialize the sandbox pod and firewall.",
	Run: func(cmd *cobra.Command, args []string) {
		setupSandbox(sandboxName, allowlist, podArgs)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringVarP(&allowlist, "allowlist", "a", "", "A comma-separated list of allowed IPs or CIDRs (e.g., '192.168.1.1,8.8.8.8'). Defaults to none (only loopback allowed).")
	upCmd.Flags().StringVarP(&podArgs, "pod-args", "p", "", "A quoted string of additional arguments to pass to 'podman pod create'.")
}

func setupSandbox(name, allowlist, extraArgs string) {
	// teardown is defined in down.go, but accessible here because they share the cmd package
	teardown(name)

	fmt.Printf("Creating sandbox pod '%s'\n", name)
	podCmdArgs := []string{"pod", "create", "--network=slirp4netns", "--name", name}

	if extraArgs != "" {
		extraArgs = strings.TrimSpace(extraArgs)
		podCmdArgs = append(podCmdArgs, strings.Split(extraArgs, " ")...)
	}

	podCreateCmd := exec.Command("podman", podCmdArgs...)
	if err := podCreateCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating pod: %v\n", err)
		os.Exit(1)
	}

	allowlistRule := ""
	if allowlist != "" {
		allowlistRule = fmt.Sprintf("ip daddr { %s } accept", allowlist)
	}

	nftRules := fmt.Sprintf(`table inet %s {
    chain output {
        type filter hook output priority filter; policy drop;
        oifname "lo" accept
        %s
    }
}`, name, allowlistRule)

	firewallName := fmt.Sprintf("%s-init-firewall", name)
	fmt.Printf("Initializing firewall container '%s'\n", firewallName)

	shScript := fmt.Sprintf("apk add --no-cache nftables && echo '%s' | nft -f -", nftRules)

	fwCreateArgs := []string{
		"create",
		"--pod", name,
		"--init-ctr=always",
		"--name", firewallName,
		"--cap-add=NET_ADMIN",
		"docker.io/alpine",
		"sh", "-c", shScript,
	}

	fwCreateCmd := exec.Command("podman", fwCreateArgs...)
	if err := fwCreateCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating firewall container: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting sandbox pod '%s'\n", name)
	startPodCmd := exec.Command("podman", "pod", "start", name)
	if err := startPodCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting pod: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sandbox '%s' successfully initialized.\n\n", name)
	fmt.Printf("To put a container in the sandbox, add \"--pod %s\" to the argument list.\n", name)
}
