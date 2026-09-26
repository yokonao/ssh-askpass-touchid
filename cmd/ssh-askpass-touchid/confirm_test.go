package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfirm(t *testing.T) {
	errDenied := errors.New("denied")

	tests := []struct {
		name       string
		promptType string
		args       []string
		authErr    error
		wantReason string
		wantErr    error
	}{
		{name: "passes the ssh-agent prompt as the reason", promptType: "confirm", args: []string{"Allow use of key id_ed25519?"}, wantReason: "Allow use of key id_ed25519?"},
		{name: "falls back to the default reason without a prompt", promptType: "confirm", wantReason: defaultReason},
		{name: "falls back to the default reason for an empty prompt", promptType: "confirm", args: []string{""}, wantReason: defaultReason},
		{name: "returns the Touch ID failure", promptType: "confirm", authErr: errDenied, wantReason: defaultReason, wantErr: errDenied},
		{name: "refuses a passphrase prompt", promptType: "", args: []string{"Enter passphrase"}, wantErr: errNotConfirm},
		{name: "refuses a notification", promptType: "none", wantErr: errNotConfirm},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotReason string
			called := false
			err := confirm(tt.promptType, tt.args, func(reason string) error {
				called = true
				gotReason = reason
				return tt.authErr
			})

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.wantReason != "", called)
			assert.Equal(t, tt.wantReason, gotReason)
		})
	}
}
