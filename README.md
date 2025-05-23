# My dotfiles

This directory containes the dotfiles for my system

## Requirements

Ensure you have the following installed on your system

### Git

```
brew install git
```

### Stow
```
brew install stow
```

## Installation

First, check out the dotfiles repo in your $HOME directory using git

```
$ git clone git@github.com:drLucaa/dotfiles.git
$ cd dotfiles
```

then use Stow to create symlinks

```
$ stow .
```

when you get errors that the files already exists, use this:

```
$ stow --adopt .
```

but be aware that this can ovveride the fiels in the `$HOME/dotfiles` directory.
