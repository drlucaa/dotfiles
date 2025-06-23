# My dotfiles

This directory containes the dotfiles for my system

## Requirements

```
brew
```

## Installation

First, check out the dotfiles repo in your $HOME directory using git:

```bash
$ git clone git@github.com:drlucaa/dotfiles.git
$ cd dotfiles
```

then use Stow to create symlinks:

```bash
$ stow .
```

when you get errors that the files already exists, use this:

```bash
$ stow --adopt .
```

but be aware that this can ovveride the fiels in the `$HOME/dotfiles` directory.

## TouchID for Sudo

Run this to copy the teamplte:

```bash
$ sudo cp /etc/pam.d/sudo_local.template /etc/pam.d/sudo_local
```

this to edit the file:

```bash
$ sudo hx /etc/pam.d/sudo_local
```

then uncomment or add the following line:

```
auth       sufficient     pam_tid.so
```

