/// <reference types="vitest" />
import { defineConfig } from "vite-plus";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
	plugins: [react()],
	server: {
		// Wails provides the desktop window during `wails dev`, so don't open a browser.
		open: false,
	},
	build: {
		outDir: "build",
		sourcemap: true,
	},
	test: {
		environment: "jsdom",
		globals: true,
		setupFiles: ["./src/test-setup.ts"],
	},
});
