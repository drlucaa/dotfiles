fpath=("$(brew --prefix)/share/zsh/site-functions" $fpath)
fpath=("$HOME/.docker/completions" $fpath)

source ~/.zsh_aliases

eval "$(mise activate zsh)"

autoload -Uz compinit
compinit
