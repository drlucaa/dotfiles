return {
  "folke/snacks.nvim",
  ---@type snacks.Config
  opts = {
    explorer = {},
    picker = {
      ---@class snacks.picker.formatters.Config
      formatters = {
        file = {
          filename_first = true,
        },
      },
      win = {
        input = {
          keys = {
            ["<S-h>"] = { "toggle_hidden", mode = { "n" } },
            ["<S-i>"] = { "toggle_ignored", mode = { "n" } },
          },
        },
      },
    },
  },
}
