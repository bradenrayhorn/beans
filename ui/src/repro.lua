-- DO NOT change the paths and don't remove the colorscheme
local root = vim.fn.fnamemodify("./.repro", ":p")

-- set stdpaths to use .repro
for _, name in ipairs({ "config", "data", "state", "cache" }) do
	vim.env[("XDG_%s_HOME"):format(name:upper())] = root .. "/" .. name
end

-- bootstrap lazy
local lazypath = root .. "/plugins/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
	vim.fn.system({
		"git",
		"clone",
		"--filter=blob:none",
		"--single-branch",
		"https://github.com/folke/lazy.nvim.git",
		lazypath,
	})
end
vim.opt.runtimepath:prepend(lazypath)

vim.filetype.add({
	extension = {
		svelte = "svelte",
	},
})

-- install plugins
local plugins = {
	"folke/tokyonight.nvim",
	{
		"stevearc/conform.nvim",
		config = function()
			require("conform").setup({
				log_level = vim.log.levels.DEBUG,
				formatters_by_ft = {
					svelte = { "prettierd" },
				},
				format_on_save = {
					timeout_ms = 500,
					lsp_fallback = true,
				},
			})
		end,
	},
	-- add any other plugins here
}
require("lazy").setup(plugins, {
	root = root .. "/plugins",
})

vim.cmd.colorscheme("tokyonight")
-- add anything else here
--
vim.opt.statusline = "%f %h%w%m%r%=%-14.(%l,%c%V%) %y %P"
