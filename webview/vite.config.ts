import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

export default defineConfig({
	plugins: [svelte(), tailwindcss()],
	resolve: {
		alias: {
			$lib: path.resolve(__dirname, "src/lib"),
		},
	},
	build: {
		outDir: "../extension/dist/webview",
		emptyOutDir: true,
		rollupOptions: {
			output: {
				entryFileNames: "webview.js",
				assetFileNames: "webview.css",
				// NO code splitting — single file for CSP
				manualChunks: undefined,
			},
		},
	},
});
