# gitctx: use multiple GitHub's SSH accounts without headaches

This repository contains `gitctx` tool. 

## What's `gitctx`?

**gitctx** is a tool that allows you to switch from a GitHub account to another in a fast way.
Currently the only way to support multiple accounts on Github is to add to the ~/.ssh/config file multiple Hosts (check [here](https://gist.github.com/rahularity/86da20fe3858e6b311de068201d279e3) or use the `.gitconfig` file with the **includeIf** directive [here](https://blog.gitguardian.com/8-easy-steps-to-set-up-multiple-git-accounts/)).

Using this tool you can switch from an account (referred as context) to another just like in the example:

![demo](img/demo.gif)

## Help

```
Usage:
  gitctx [OPTIONS]

Application Options:
      --add        Create a new host in the selected config file.
      --rm=        Remove an existing host in the selected config file.
      --edit=      Edit an existing host in the selected config file.
  -s, --sshconfig= Set the path of the config (default: ~/.ssh/config).
      --hostname=  Set the hostname to use for context change (default: github.com).
      --config=    Set the path of the gitctx folder (default: ~/.gitctx).

Help Options:
  -h, --help       Show this help message
```

## Installation

> $ go install github.com/Matrix86/gitctx/cmd/gitctx@latest

### From sources

> $ git clone git@github.com:Matrix86/gitctx.git

> $ cd gitctx

> $ make install

## Completion

To enable the shell completion, you need to move add the following line to the end of the `~/.bashrc` file:

> . $HOME/.gitctx/gitctx.bash

## Git configuration

It is possible to define some GIT configs that should be called when you switch from a context to another. 
In order to define them you can manually add them in the configuration file (stored on `~/.gitctx/config.yml` as default) as in the following example:

```
git_settings:
    username: <-- it has to be the same used in the `hosts` field
        user.name: MyName
        user.email: something@me.me
        user.signingkey: /home/ubuntu/.ssh/id_rsa.pub
        commit.gpgsign: true
        gpg.format: ssh
```