import { defineConfig } from "vite-plus";

const projectSpecificIgnores = ["*.json", "pear-api-queue.mts", "pear-api-queue-search-songs.mts"];

export default defineConfig({
	staged: {
		"*": "vp check --fix",
	},
	fmt: {
		ignorePatterns: projectSpecificIgnores,
		useTabs: true,
	},
	lint: {
		jsPlugins: [{ name: "vite-plus", specifier: "vite-plus/oxlint-plugin" }],
		rules: { "vite-plus/prefer-vite-plus-imports": "error" },
		options: { typeAware: true, typeCheck: true },
		ignorePatterns: projectSpecificIgnores,
	},
});
