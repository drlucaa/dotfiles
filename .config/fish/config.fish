# Interactive context
if status is-interactive
    # Set Abbreviations
    abbr gg lazygit
    abbr vim nvim
    abbr zj zellij
    abbr ll ls -alF
    abbr ju jjui

    # Set Editor Variable
    set -gx EDITOR hx
    set -gx VISUAL hx
    set -gx LANG en_US.UTF-8
    set -gx SSH_AUTH_SOCK "$HOME/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock"

    # Initialize tools
    starship init fish | source
    atuin init fish | source
    zoxide init --cmd cd fish | source
    mise activate fish | source

    function cx --description "Crete and move into new dir"
        mkdir -p -- $argv; and cd -- $argv
    end

    function jgf --description "JJ git fetch"
        jj git fetch
    end

    function jbm --description "JJ create bookmark at current revision"
        jj bookmark create $argv --revision @
    end

    function atuin_sync_on_exit --on-event fish_exit
        if test -x (command -v atuin)
            atuin sync >/dev/null 2>&1
        end
    end
end

if not status is-interactive
    # Add shims for non interactive contexts
    mise activate fish --shims | source
end
