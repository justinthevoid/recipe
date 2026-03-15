import path from "node:path";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { svelteTesting } from "@testing-library/svelte/vite";
import { defineConfig } from "vitest/config";

export default defineConfig({
	plugins: [svelte(), svelteTesting(), tailwindcss()],
	resolve: {
		alias: {
			$lib: path.resolve(__dirname, "src/lib"),
		},
		conditions: ["browser"],
	},
	test: {
		environment: "jsdom",
		setupFiles: ["./vitest-setup.ts"],
		server: {
			deps: {
				inline: [/svelte/],
			},
		},
	},
	build: {
		outDir: "../extension/dist/webview",
		emptyOutDir: true,
		rollupOptions: {
			output: {
				entryFileNames: "webview.js",
				assetFileNames: (info) =>
					info.name?.endsWith(".css") ? "webview.css" : "assets/[name]-[hash][extname]",
				// NO code splitting — single file for CSP
				manualChunks: undefined,
			},
		},
	},
});
