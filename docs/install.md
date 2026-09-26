# Install

ssh-askpass-touchid is a single binary for macOS. Install it in one of these ways:

1. [GitHub Releases](#github-releases)
1. [mise](#mise)
1. [go install](#go-install)

ssh-agent runs the binary from the absolute path in the LaunchAgent, so it does not need to be in `PATH`. Pick a [stable path](#stable-path) so an upgrade does not break the LaunchAgent.

## GitHub Releases

Download the archive for your Mac from [GitHub Releases](https://github.com/yokonao/ssh-askpass-touchid/releases), then install the binary:

```sh
asset=ssh-askpass-touchid_darwin_arm64.tar.gz # ssh-askpass-touchid_darwin_amd64.tar.gz on Intel
gh release download -R yokonao/ssh-askpass-touchid -p "$asset" # the latest release; pass a tag for another
tar -xzf "$asset" ssh-askpass-touchid
mkdir -p ~/.local/bin
install -m 755 ssh-askpass-touchid ~/.local/bin/
```

The binaries are not notarized. `gh` and `curl` do not mark downloads as quarantined, but a browser does, and macOS then refuses to run the binary; remove the mark with `xattr -d com.apple.quarantine ssh-askpass-touchid`.

### Verify the archive

Each archive has a [build provenance attestation](https://docs.github.com/en/actions/security-for-github-actions/using-artifact-attestations) from the release workflow. Verify it with the [GitHub CLI](https://cli.github.com/):

```sh
gh attestation verify "$asset" \
  -R yokonao/ssh-askpass-touchid \
  --signer-workflow yokonao/ssh-askpass-touchid/.github/workflows/release.yaml
```

## mise

With [mise](https://mise.jdx.dev/):

```sh
mise use -g github:yokonao/ssh-askpass-touchid
```

mise verifies the [build provenance attestation](#verify-the-archive) automatically.

## go install

With Go 1.27.1 or later:

```sh
go install github.com/yokonao/ssh-askpass-touchid/cmd/ssh-askpass-touchid@latest
```

The binary lands in `$(go env GOBIN)`, or `~/go/bin` when that is unset.

## Stable path

ssh-agent starts the helper anew for every confirmation, so replacing the binary at the same path takes effect without restarting the agent.

| Install | Path for `SSH_ASKPASS` |
| --- | --- |
| GitHub Releases | `~/.local/bin/ssh-askpass-touchid` |
| mise | `~/.local/share/mise/installs/github-yokonao-ssh-askpass-touchid/latest/ssh-askpass-touchid` |
| go install | `$(go env GOPATH)/bin/ssh-askpass-touchid` |

mise keeps the `latest` symlink pointing at the newest installed version. Avoid the mise shim: it resolves the version from the working directory's config, which launchd does not provide.
