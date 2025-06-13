# Interactive context
if status is-interactive
    # Set Abbreviations
    abbr lg lazygit
    abbr vim nvim
    abbr zj zellij

    # Set Editor Variable
    set -gx EDITOR "zed --wait"
    set -gx VISUAL "zed --wait"

    # Initialize tools
    starship init fish | source
    atuin init fish | source
    zoxide init --cmd cd fish | source
    mise activate fish | source

    function atuin_sync_on_exit --on-event fish_exit
        if test -x (command -v atuin)
            atuin sync >/dev/null 2>&1
        end
    end
end

# Add shims for non interactive contexts
mise activate fish --shims | source
