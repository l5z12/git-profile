# Git Profile switcher

[![build](https://github.com/dotzero/git-profile/actions/workflows/ci.yml/badge.svg)](https://github.com/dotzero/git-profile/actions/workflows/ci.yml)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dotzero/git-profile/blob/master/LICENSE)

Git Profile allows you to switch between multiple user profiles in git repositories

## Installation

If you are a macOS user, you can use [Homebrew](http://brew.sh/):

```bash
brew install dotzero/tap/git-profile
```

### Prebuilt binaries

Download the binary from the [releases](https://github.com/dotzero/git-profile/releases) page and place it under `$PATH` directory.

### Building from source

If your operating system does not have a binary release, but does run Go, you can build it from the source.

```bash
go get -u github.com/dotzero/git-profile
```

The binary will then be installed to `$GOPATH/bin` (or your `$GOBIN`).

## Usage

Adds an entry to a profile or updates an existing profile

```bash
git profile add home user.name dotzero
git profile add home user.email "mail@mail.com"
git profile add home user.signingkey AAAAAAAA
```

If no arguments are provided, `add` starts an interactive mode:

```bash
git profile add
```

It will ask for the profile name and values for `user.name`, `user.email`, and
`user.signingkey`. Leaving a field empty skips writing that field for a new
profile. For an existing profile, the current values are prefilled and pressing
Enter keeps them unchanged. Press `Esc` or `Ctrl+C` to cancel interactive mode.

Displays a list of available profiles

```bash
git profile list
```

Delete a profile or a specific key from a profile

```bash
git profile del
git profile del home
git profile del home user.email
```

When called without arguments, `del` opens an interactive profile selector and
removes the selected profile entirely.

Applies the selected profile entries to the current git repository

```bash
git profile use home

# Under the hood it runs following commands:
# git config --local user.name dotzero
# git config --local user.email "me@dotzero.ru"
# git config --local user.signingkey AAAAAAAA
```

If no profile name is provided, an interactive selector will appear to choose a profile

```bash
git profile use
```

Export a profile in JSON format

```bash
git profile export home > home.json
```

Import profile from JSON format

```bash
cat home.json | xargs -0 git profile import home
```

## Shell Completion

Shell completion is available for `git-profile` command.

**Note:** Use `git-profile TAB` instead of `git profile TAB` for completion.

### Bash

Requires [bash-completion](https://github.com/scop/bash-completion).

```bash
# Temporary (current session only)
source <(git-profile completion bash)

# Permanent (Linux)
git-profile completion bash | sudo tee /etc/bash_completion.d/git-profile

# Permanent (macOS)
git-profile completion bash > $(brew --prefix)/etc/bash_completion.d/git-profile
```

### Zsh

```bash
# Create completion directory if it doesn't exist
mkdir -p ~/.zsh/completions

# Generate completion file
git-profile completion zsh > ~/.zsh/completions/_git-profile

# Add to ~/.zshrc if not already present
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo "autoload -U compinit && compinit" >> ~/.zshrc

# Reload shell
source ~/.zshrc
```

### Fish

```bash
# Temporary (current session only)
git-profile completion fish | source

# Permanent
git-profile completion fish > ~/.config/fish/completions/git-profile.fish
```

## License

http://www.opensource.org/licenses/mit-license.php
