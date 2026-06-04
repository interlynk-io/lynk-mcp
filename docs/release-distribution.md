# Interlynk OSS Release Distribution

This document defines the shared release model for Interlynk OSS tools:

- `sbomqs`
- `sbomasm`
- `sbommv`
- `sbomgr`
- `bomtique`
- `lynk-mcp`
- `pylynk`

The goal is consistent installation and upgrade behavior across macOS, Windows, and Linux while keeping package-manager updates auditable through signed pull requests.

## Shared Repositories

| Channel | Repository | Usage |
|---------|------------|-------|
| Homebrew | `interlynk-io/homebrew-interlynk` | Homebrew casks for macOS CLI installs |
| Scoop | `interlynk-io/homebrew-interlynk` | Scoop manifests for Windows installs |
| winget | `interlynk-io/winget-pkgs` | Interlynk fork used to open PRs to `microsoft/winget-pkgs` |
| Linux packages | Source repository releases | `.deb`, `.rpm`, and `.apk` release assets |

Homebrew and Scoop intentionally share `interlynk-io/homebrew-interlynk` so Interlynk has one public package repository for direct package-manager installs. winget must use the `interlynk-io/winget-pkgs` fork because official winget packages are accepted through PRs to `microsoft/winget-pkgs`.

## Shared Secrets

Create these as organization-level GitHub Actions secrets and scope them to the OSS source repositories that publish releases.

| Secret | Purpose |
|--------|---------|
| `INTERLYNK_RELEASE_GITHUB_TOKEN` | Fine-grained PAT for opening package-manager PRs |
| `INTERLYNK_RELEASE_SSH_KEY` | Private SSH key used by GoReleaser to push signed PR branches |
| `INTERLYNK_RELEASE_GPG_PRIVATE_KEY` | ASCII-armored private GPG key used to sign tap commits |
| `INTERLYNK_RELEASE_GPG_PASSPHRASE` | Passphrase for the release GPG key |

The release workflows still use the built-in `GITHUB_TOKEN` for creating releases in the source repository.

## Release Bot

Use `interlynk-support-bot` as the GitHub machine user for all OSS release automation.

The machine user should have:

- 2FA enabled.
- Write access to `interlynk-io/homebrew-interlynk`.
- Write access to `interlynk-io/winget-pkgs`.
- A fine-grained PAT stored as `INTERLYNK_RELEASE_GITHUB_TOKEN`.
- An SSH public key matching `INTERLYNK_RELEASE_SSH_KEY`.
- The GPG public key matching `INTERLYNK_RELEASE_GPG_PRIVATE_KEY`.

The GoReleaser configs should use this commit identity:

```yaml
commit_author:
  name: interlynk-support-bot
  email: support_eng@interlynk.io
```

GitHub marks signed commits as verified only when the public GPG key is uploaded to the GitHub account that owns the commit email. The `support_eng@interlynk.io` email must be verified on the `interlynk-support-bot` GitHub account.

## Token Permissions

`INTERLYNK_RELEASE_GITHUB_TOKEN` should be a fine-grained PAT from the release bot account.

Repository access:

- `interlynk-io/homebrew-interlynk`
- `interlynk-io/winget-pkgs`

Permissions:

- Contents: read and write.
- Pull requests: read and write.
- Metadata: read.

If the organization requires approval for fine-grained PAT access, approve the token for the selected repositories before running releases.

## SSH Key

Generate one SSH key for the release bot:

```bash
ssh-keygen -t ed25519 -C "interlynk-support-bot" -f ./interlynk-support-bot
```

Add `./interlynk-support-bot.pub` to the release bot GitHub account as an SSH key, or add it as a write-enabled deploy key on `interlynk-io/homebrew-interlynk`.

Store the private key:

```bash
gh secret set INTERLYNK_RELEASE_SSH_KEY --org interlynk-io --visibility selected < ./interlynk-support-bot
```

Then select the OSS repositories that should be allowed to use it.

## GPG Key

Create or reuse a release signing key owned by the release bot identity:

```bash
gpg --full-generate-key
gpg --list-secret-keys --keyid-format=long
gpg --armor --export-secret-keys <KEY_ID> > ./interlynk-release-signing-key.asc
gpg --armor --export <KEY_ID> > ./interlynk-release-signing-key.pub
```

Upload `interlynk-release-signing-key.pub` to the release bot GitHub account.

Store the private key and passphrase:

```bash
gh secret set INTERLYNK_RELEASE_GPG_PRIVATE_KEY --org interlynk-io --visibility selected < ./interlynk-release-signing-key.asc
gh secret set INTERLYNK_RELEASE_GPG_PASSPHRASE --org interlynk-io --visibility selected
```

## GoReleaser Pattern

Each source repository owns its release artifacts and package metadata. Keep the GoReleaser structure consistent and change only the tool-specific values.

Required sections:

- `builds` or equivalent artifact build configuration.
- `archives` for `.tar.gz` and `.zip` artifacts.
- `nfpms` for `.deb`, `.rpm`, and `.apk`.
- `checksum`.
- `homebrew_casks`.
- `scoops`.
- `winget`.

The package-manager sections should use shared credentials:

```yaml
homebrew_casks:
  - name: <tool>
    repository:
      owner: interlynk-io
      name: homebrew-interlynk
      branch: "<tool>-{{ .Version }}"
      token: "{{ .Env.INTERLYNK_RELEASE_GITHUB_TOKEN }}"
      git:
        url: "ssh://git@github.com/interlynk-io/homebrew-interlynk.git"
        private_key: "{{ .Env.INTERLYNK_RELEASE_SSH_KEY }}"
      pull_request:
        enabled: true
        draft: false
        base:
          owner: interlynk-io
          name: homebrew-interlynk
          branch: main
    directory: Casks
    commit_author:
      name: interlynk-support-bot
      email: support_eng@interlynk.io
      signing:
        enabled: true
        key: "{{ .Env.GPG_SIGNING_KEY }}"
        program: gpg
        format: openpgp

scoops:
  - name: <tool>
    repository:
      owner: interlynk-io
      name: homebrew-interlynk
      branch: "<tool>-scoop-{{ .Version }}"
      token: "{{ .Env.INTERLYNK_RELEASE_GITHUB_TOKEN }}"
      git:
        url: "ssh://git@github.com/interlynk-io/homebrew-interlynk.git"
        private_key: "{{ .Env.INTERLYNK_RELEASE_SSH_KEY }}"
      pull_request:
        enabled: true
        draft: false
        base:
          owner: interlynk-io
          name: homebrew-interlynk
          branch: main
    commit_author:
      name: interlynk-support-bot
      email: support_eng@interlynk.io
      signing:
        enabled: true
        key: "{{ .Env.GPG_SIGNING_KEY }}"
        program: gpg
        format: openpgp

winget:
  - name: <tool>
    publisher: Interlynk
    package_identifier: Interlynk.<tool>
    package_name: <tool>
    repository:
      owner: interlynk-io
      name: winget-pkgs
      branch: "<tool>-{{ .Version }}"
      token: "{{ .Env.INTERLYNK_RELEASE_GITHUB_TOKEN }}"
      pull_request:
        enabled: true
        draft: false
        base:
          owner: microsoft
          name: winget-pkgs
          branch: master
```

Add the macOS quarantine hook to each cask unless the binary is Apple-signed and notarized:

```yaml
hooks:
  post:
    install: |
      if OS.mac?
        system_command "/usr/bin/xattr", args: ["-dr", "com.apple.quarantine", "#{staged_path}/<tool>"]
      end
```

## Workflow Pattern

Each source repository should import the release GPG key before running GoReleaser and pass the shared secrets through the GoReleaser environment:

```yaml
- name: Import release signing key
  id: import_gpg
  uses: crazy-max/ghaction-import-gpg@v6
  with:
    gpg_private_key: ${{ secrets.INTERLYNK_RELEASE_GPG_PRIVATE_KEY }}
    passphrase: ${{ secrets.INTERLYNK_RELEASE_GPG_PASSPHRASE }}
    git_user_signingkey: true
    git_commit_gpgsign: true

- name: Run GoReleaser
  run: make release
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    INTERLYNK_RELEASE_GITHUB_TOKEN: ${{ secrets.INTERLYNK_RELEASE_GITHUB_TOKEN }}
    INTERLYNK_RELEASE_SSH_KEY: ${{ secrets.INTERLYNK_RELEASE_SSH_KEY }}
    GPG_SIGNING_KEY: ${{ steps.import_gpg.outputs.fingerprint }}
```

## Per-Tool Values

Use these values when replicating the release config.

| Source repo | Project/package name | Homebrew cask | Scoop manifest | winget identifier |
|-------------|----------------------|---------------|----------------|-------------------|
| `sbomqs` | `sbomqs` | `sbomqs` | `sbomqs` | `Interlynk.sbomqs` |
| `sbomasm` | `sbomasm` | `sbomasm` | `sbomasm` | `Interlynk.sbomasm` |
| `sbommv` | `sbommv` | `sbommv` | `sbommv` | `Interlynk.sbommv` |
| `sbomgr` | `sbomgr` | `sbomgr` | `sbomgr` | `Interlynk.sbomgr` |
| `bomtique` | `bomtique` | `bomtique` | `bomtique` | `Interlynk.bomtique` |
| `lynk-mcp` | `lynk-mcp` | `lynk-mcp` | `lynk-mcp` | `Interlynk.lynk-mcp` |
| `pylynk` | `pylynk` | `pylynk` | `pylynk` | `Interlynk.pylynk` |

For Python projects such as `pylynk`, keep the distribution targets consistent but adapt the artifact build step to the project. If `pylynk` publishes a CLI binary, the package-manager manifests should install that binary. If it is only a Python library, package-manager distribution should be deferred until a CLI artifact exists.

## Release Checklist

For each repo:

1. Ensure `.goreleaser.yaml` uses the shared release sections and `INTERLYNK_RELEASE_*` environment variables.
2. Ensure the release workflow imports the shared GPG key before `make release` or `goreleaser release`.
3. Run `goreleaser check`.
4. Run the repository test suite.
5. Push a test tag in a non-critical repo first.
6. Confirm the release creates:
   - GitHub release assets.
   - `.deb`, `.rpm`, and `.apk` packages.
   - A signed Homebrew cask PR in `interlynk-io/homebrew-interlynk`.
   - A signed Scoop manifest PR in `interlynk-io/homebrew-interlynk`.
   - A winget PR from `interlynk-io/winget-pkgs` to `microsoft/winget-pkgs`.
7. Confirm tap PR commits show as verified in GitHub.

## Linux Package-Manager Updates

The current standard publishes Linux packages as GitHub release assets. Users can install newer `.deb`, `.rpm`, or `.apk` packages manually from each release.

True `apt upgrade`, `dnf upgrade`, or `apk upgrade` support requires a hosted package repository. If that becomes a requirement, add a shared repository service such as Cloudsmith, Packagecloud, or an Interlynk-hosted apt/rpm/apk repository and publish all seven tools to that service from the same release workflows.
