# ssh-askpass-touchid

`ssh-askpass-touchid` asks for Touch ID each time ssh-agent is about to use a key added with `ssh-add -c`. It plugs into OpenSSH's standard confirmation flow as `SSH_ASKPASS`, so ssh-agent itself is the stock `/usr/bin/ssh-agent`.

## Requirements

- macOS 13 or later with Touch ID
- OpenSSH 8.4 or later for `SSH_ASKPASS_REQUIRE`; macOS ships a newer one

## Install

Download a prebuilt binary from [GitHub Releases](https://github.com/yokonao/ssh-askpass-touchid/releases).
Every release ships with a build provenance attestation.

See [docs/install.md](docs/install.md) for other install methods and how to verify a release.

## Setup

Run a dedicated ssh-agent from launchd with the helper as its `SSH_ASKPASS`. The macOS built-in agent cannot take extra environment variables, and its socket path changes on every boot.

Set `askpass` to the absolute path of the installed binary; see [docs/install.md](docs/install.md#stable-path) for a path that survives upgrades.

```sh
askpass="$HOME/.local/bin/ssh-askpass-touchid"
socket="$HOME/.ssh/agent-touchid.sock"
label=com.github.yokonao.ssh-askpass-touchid
plist="$HOME/Library/LaunchAgents/$label.plist"

mkdir -p "$HOME/Library/LaunchAgents"
cat >"$plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>$label</string>
  <key>LimitLoadToSessionType</key>
  <string>Aqua</string>
  <key>ProgramArguments</key>
  <array>
    <string>/usr/bin/ssh-agent</string>
    <string>-D</string>
    <string>-a</string>
    <string>$socket</string>
  </array>
  <key>EnvironmentVariables</key>
  <dict>
    <key>SSH_ASKPASS</key>
    <string>$askpass</string>
    <key>SSH_ASKPASS_REQUIRE</key>
    <string>force</string>
  </dict>
  <key>KeepAlive</key>
  <true/>
  <key>RunAtLoad</key>
  <true/>
</dict>
</plist>
EOF
launchctl bootstrap "gui/$(id -u)" "$plist"
```

Point your shell at the new agent, for example in `~/.zshenv`:

```sh
export SSH_AUTH_SOCK="$HOME/.ssh/agent-touchid.sock"
```

Then add keys with confirmation:

```sh
ssh-add -c ~/.ssh/id_ed25519
```

The agent forgets its keys when it restarts. To have `ssh` add a key with confirmation on first use instead, set `AddKeysToAgent confirm` in `~/.ssh/config`.

## Behavior

ssh-agent runs the helper with its prompt as the only argument and `SSH_ASKPASS_PROMPT=confirm`. The helper shows the prompt in a Touch ID dialog and exits 0 only if Touch ID succeeds. Every use needs a fresh Touch ID; a recent success is not reused.

The helper exits 1 without a dialog for anything else ssh-agent or ssh asks through `SSH_ASKPASS`, such as passphrases. There is no password fallback: when Touch ID is unavailable, for example with the lid closed and no Touch ID keyboard, key use is denied.

## Security model

Confirmation protects keys from processes that can reach the agent socket without your knowledge, such as a remote host through a forwarded agent. Each signature needs a Touch ID on this Mac.

This helper does not protect against:

- keys added without `-c`;
- a compromised user session, kernel, or helper binary, or anyone who can change the LaunchAgent;
- a process that asks for a signature at the moment you expect your own, since the dialog cannot tell the two apart;
- physical access after successful biometric authentication.

Keep the helper and the LaunchAgent plist writable only by the current user.

## Development

```sh
go test ./...
golangci-lint run
golangci-lint fmt
```

## License

MIT.
