eval "$(/opt/homebrew/bin/brew shellenv)"

# Added by Toolbox App
export PATH="$PATH:$HOME/Library/Application Support/JetBrains/Toolbox/scripts"

eval "$(mise activate bash --shims)"

# Zed
export EDITOR="zed --wait"
export VISUAL="zed --wait"
