package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yokonao/ssh-askpass-touchid/internal/touchid"
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	root := &cobra.Command{
		Use:   "ssh-askpass-touchid [prompt]",
		Short: "Approve ssh-agent key use with Touch ID",
		Long: "Approve ssh-agent key use with Touch ID.\n\n" +
			"ssh-agent runs this as SSH_ASKPASS for keys added with ssh-add -c. " +
			"It exits 0 only after Touch ID succeeds, and refuses passphrase prompts.",
		Version:       version,
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			return confirm(os.Getenv("SSH_ASKPASS_PROMPT"), args, touchid.Authenticate)
		},
	}
	root.CompletionOptions.DisableDefaultCmd = true
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "ssh-askpass-touchid: %v\n", err)
		os.Exit(1)
	}
}
