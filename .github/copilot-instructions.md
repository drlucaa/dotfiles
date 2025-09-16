# Personal Dotfiles Repository

This is a personal dotfiles repository using GNU Stow for managing configuration files via symlinks. The repository contains shell configurations (zsh, fish), editor configs (neovim with LazyVim, helix), and various development tool configurations.

**Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

## Working Effectively

### Prerequisites and Dependencies
- Install GNU Stow: `sudo apt-get update && sudo apt-get install -y stow` 
- **NEVER CANCEL**: Package installation takes 2-5 minutes depending on system updates. Set timeout to 10+ minutes for safety.
- **Validated timing**: Stow installation completes in under 5 minutes on standard Ubuntu systems
- Most tools referenced in configs (mise, starship, helix, etc.) are not required for basic dotfiles installation
- The repository works primarily through symlink management, not compilation
- **No build dependencies required** - GNU Stow is the only external dependency

### Installation and Setup Commands
- Clone the repository: `git clone https://github.com/drlucaa/dotfiles.git`
- Navigate to directory: `cd dotfiles`
- **Primary installation method**: `stow .` 
  - Creates symlinks from home directory to dotfiles
  - **Timing**: Completes instantly (< 0.1 seconds)
  - **WARNING**: Will fail with detailed error message if target files already exist
- **Preview changes first**: `stow --simulate --verbose .` (shows what will be linked without making changes)
- **If conflicts exist**: `stow --adopt .`
  - **CRITICAL**: This overwrites files in the dotfiles directory with existing home directory files
  - Use with extreme caution as it modifies the repository content
  - **Validated behavior**: Existing files are moved into dotfiles directory and symlinks created
- **Removal**: `stow -D .` (removes all symlinks created by stow)
  - **Timing**: Completes instantly (< 0.1 seconds)

### Validation and Testing
- **ALWAYS validate symlinks after stow**: `ls -la ~ | grep "dotfiles"`
- **Expected symlinks created**:
  - `.zshrc -> dotfiles/.zshrc`
  - `.config -> dotfiles/.config`
  - `.ideavimrc -> dotfiles/.ideavimrc`
  - `.zprofile -> dotfiles/.zprofile`
  - `.zsh_functions -> dotfiles/.zsh_functions`
  - `Brewfile -> dotfiles/Brewfile`
- **Verify file accessibility**: `cat ~/.zshrc | head -5` (should show zsh configuration)
- **Test shell configuration**: Source the file `source ~/.zshrc` or open new terminal
  - **Expected behavior**: May show warnings about missing tools (mise, starship, etc.) but should not fail
- **Test editor configs**: Check if config files exist: `ls ~/.config/nvim/` 
- **No build process required** - changes take effect immediately upon symlink creation
- **Conflict testing**: Create dummy file, run stow, observe clear error message

## Repository Structure and Key Components

### Shell Configurations
- **zsh**: `.zshrc` with zinit plugin manager, includes fzf, zoxide, atuin integration
- **fish**: `.config/fish/config.fish` with abbreviations and tool integrations
- **Shell functions**: `.zsh_functions` for zsh-specific functions

### Editor Configurations
- **Neovim**: `.config/nvim/` - LazyVim starter configuration
- **Helix**: `.config/helix/` - Modern terminal editor config
- **IdeaVim**: `.ideavimrc` - Vim emulation for JetBrains IDEs

### Development Tool Configs
- **mise**: `.config/mise/config.toml` - Runtime version manager
- **jj**: `.config/jj/config.toml` - Jujutsu VCS configuration  
- **atuin**: `.config/atuin/config.toml` - Shell history sync
- **starship**: `.config/starship.toml` - Cross-shell prompt
- **zed**: `.config/zed/` - Modern code editor settings

### Package Management
- **Brewfile**: macOS Homebrew bundle file listing all packages and apps
- Contains CLI tools (git, neovim, helix, etc.) and GUI applications

## Common Tasks and Workflows

### Making Configuration Changes
- Edit files directly in the dotfiles directory
- Changes take effect immediately (configs are symlinked)
- **Always test changes**: Open new shell/editor session to verify
- **For shell configs**: Source the file or open new terminal
- **For editor configs**: Restart the editor to load changes

### Adding New Configurations
- Add new config files/directories to the dotfiles repository
- Run `stow .` again to create new symlinks
- **Check for conflicts**: Use `stow --simulate .` first

### Troubleshooting Installation
- **Conflict errors**: File already exists in target location
  - Option 1: Backup and remove existing file, then `stow .`
  - Option 2: Use `stow --adopt .` (modifies repository)
- **Permission errors**: Ensure write access to home directory
- **Broken symlinks**: Remove with `stow -D .` then reinstall with `stow .`

### Working with Different Operating Systems
- **Primary target**: macOS (Brewfile, PATH configurations)
- **Linux compatibility**: Most configs work, some tool paths may differ
- **Tool availability**: Many tools in Brewfile not available on all systems
- **No Windows support**: Unix-style configurations only

## Important Files and Locations

### Frequently Modified Files
```
.zshrc                    # Main zsh configuration
.config/fish/config.fish  # Fish shell configuration  
.config/nvim/             # Neovim configuration directory
.config/helix/config.toml # Helix editor settings
Brewfile                  # Package list for macOS
```

### Configuration Dependencies
- `.zshrc` depends on zinit plugin manager (auto-installed on first run)
- Fish config depends on starship, atuin, zoxide (not required for basic function)
- Neovim config is LazyVim starter (plugins managed automatically)

### Files to Ignore in Changes
- `.config/nvim/lazy-lock.json` - Plugin lockfile (auto-generated)
- `.config/nvim/.neoconf.json` - Neovim project settings (auto-generated)

## Validation Scenarios

## Validation Scenarios

### Complete Setup Validation (ALWAYS RUN THESE STEPS)
1. **Install prerequisite and verify**:
   ```bash
   sudo apt-get update && sudo apt-get install -y stow  # NEVER CANCEL: 2-5 minutes
   which stow && stow --version  # Should show GNU Stow version 2.3.1+
   ```

2. **Preview installation without changes**:
   ```bash
   cd dotfiles
   stow --simulate --verbose .  # Shows planned symlinks, makes no changes
   ```

3. **Install dotfiles**:
   ```bash
   stow .  # Creates symlinks instantly
   ```

4. **Verify symlinks created correctly**:
   ```bash
   cd ..
   ls -la | grep dotfiles  # Should show 7 symlinks
   cat .zshrc | head -3    # Should show "# Brew completions" etc.
   ```

5. **Test basic functionality**:
   ```bash
   # Verify shell config loads (may show tool warnings - this is normal)
   source .zshrc  
   # Check editor config exists
   ls .config/nvim/
   ```

### Change Validation Workflow
- **Before making changes**: Note current working state with `ls -la | grep dotfiles`
- **After changes**: Verify symlinks still point correctly
- **For shell changes**: `source ~/.zshrc` to reload
- **For editor changes**: Restart editor to pick up new configuration
- **Rollback if needed**: `git checkout HEAD -- filename` to revert changes

### Conflict Resolution Testing
1. **Create conflict scenario**:
   ```bash
   echo "test file" > ~/.zshrc  # Create conflicting file
   cd dotfiles && stow .       # Should fail with clear error message
   ```

2. **Resolve with adopt (USE WITH CAUTION)**:
   ```bash
   stow --adopt .  # Moves existing file into dotfiles, creates symlink
   git status      # Shows modified files that were adopted
   ```

### Complete Removal Testing
```bash
cd dotfiles
stow -D .                    # Remove all symlinks
cd .. && ls -la | grep dotfiles  # Should show no symlinks
```

## Timeline Expectations
- **Stow installation**: 2-5 minutes (including system updates) - **NEVER CANCEL**, use 10+ minute timeout
- **Dotfiles simulation**: Instant (< 0.1 seconds)
- **Dotfiles linking**: Instant (< 0.1 seconds) - **Validated timing**: 0.03 seconds average
- **Dotfiles removal**: Instant (< 0.1 seconds)
- **Shell configuration loading**: Instant (may show warnings about missing tools)
- **Editor startup**: 1-3 seconds (first time may be longer for plugin loading)
- **No compilation or build time** - pure configuration management through symlinks

## Limitations and Notes
- **No automated testing**: Validation is manual/functional only
- **Tool dependencies**: Many configs reference tools not available on all systems
- **macOS-centric**: Some configurations assume macOS file paths
- **No backup mechanism**: Stow overwrites existing files with --adopt
- **Version control**: All changes should be committed to maintain history

This repository is designed for personal use and may require adaptation for different environments or tool preferences.