# Set up Homebrew's zsh completion functions
fpath=("$(brew --prefix)/share/zsh/site-functions" $fpath)
fpath=(/Users/luca.fondo/.docker/completions $fpath)

source ~/.zsh_aliases

eval "$(mise activate zsh)"

autoload -Uz compinit
compinit
