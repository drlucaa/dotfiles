return {
  "Pocco81/auto-save.nvim",
  -- load immediately so it catches InsertLeave etc.
  lazy = false,
  -- bind a toggle command into LazyVim's "+ui" menu under <leader>uv
  keys = {
    { "<leader>us", "<cmd>ASToggle<CR>", desc = "Toggle AutoSave" },
  },
  config = function()
    require("auto-save").setup({
      -- enable by default
      enabled = true,
      -- events that trigger a save
      trigger_events = { "InsertLeave", "TextChanged" },
      -- wait 300 ms after the last change before saving
      debounce_delay = 300,
      -- remove the “saved” message
      execution_message = {
        message = function()
          return ""
        end,
      },
      -- only save buffers with these filetypes
      condition = function(buf)
        buf = buf or vim.api.nvim_get_current_buf()
        if not vim.api.nvim_buf_is_valid(buf) then
          return false
        end

        local ft = vim.bo[buf].filetype
        local allowed = {
          go = true,
          rust = true,
          java = true,
          python = true,
          lua = true,
          javascript = true,
          typescript = true,
          html = true,
          css = true,
          markdown = true,
        }
        return allowed[ft]
      end,
    })
  end,
}
