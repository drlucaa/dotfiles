fpath=("$(brew --prefix)/share/zsh/site-functions" $fpath)
fpath=("$HOME/.docker/completions" $fpath)

source ~/.zsh_aliases
source ~/.config/op/plugins.sh

eval "$(mise activate zsh)"
eval "$(zoxide init --cmd cd zsh)"

autoload -Uz compinit
compinit
