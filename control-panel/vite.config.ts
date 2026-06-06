import { defineConfig } from "vite-plus";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
	plugins: [react()],
	server: {
		open: true,
		proxy: {
			"/api": {
				target: "http://127.0.0.1:" + (process.env.PORT || "3999"),
				changeOrigin: true,
				secure: false,
			},
		},
	},
	build: {
		outDir: "build",
		sourcemap: true,
	},
});
