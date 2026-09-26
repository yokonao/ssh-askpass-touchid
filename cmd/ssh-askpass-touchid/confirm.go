package main

import "errors"

const defaultReason = "Allow SSH key use"

var errNotConfirm = errors.New("only key use confirmations are supported; SSH_ASKPASS_PROMPT is not confirm")

// confirm approves a key use only when ssh-agent asks for a confirmation and Touch ID succeeds.
func confirm(promptType string, args []string, authenticate func(reason string) error) error {
	if promptType != "confirm" {
		return errNotConfirm
	}
	reason := defaultReason
	if len(args) == 1 && args[0] != "" {
		reason = args[0]
	}
	return authenticate(reason)
}
