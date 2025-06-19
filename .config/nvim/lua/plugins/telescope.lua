return {
  {
    "nvim-telescope/telescope.nvim",
    enabled = false,
    config = function()
      local ts = require("telescope")
      local h_pct = 0.90
      local w_pct = 0.80
      local w_limit = 75

      local standard_setup = {
        borderchars = { "─", "│", "─", "│", "┌", "┐", "┘", "└" },
        preview = { hide_on_startup = true },
      }

      ts.setup({
        defaults = vim.tbl_extend("error", standard_setup, {
          sorting_strategy = "descending",
          path_display = { "filename_first" },
          mappings = {
            n = {
              ["o"] = require("telescope.actions.layout").toggle_preview,
              ["<C-c>"] = require("telescope.actions").close,
            },
            i = {
              ["<C-o>"] = require("telescope.actions.layout").toggle_preview,
            },
          },
        }),
        pickers = {
          find_files = {
            find_command = {
              "fd",
              "--type",
              "f",
              "-H",
              "--strip-cwd-prefix",
            },
          },
        },
      })

      ts.load_extension("fzf")
    end,
  },
}
