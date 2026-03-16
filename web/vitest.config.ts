import { svelte } from "@sveltejs/vite-plugin-svelte";
import { svelteTesting } from "@testing-library/svelte/vite";
import { defineConfig } from "vitest/config";

export default defineConfig({
	plugins: [svelte(), svelteTesting()],
	resolve: {
		alias: {
			$lib: new URL("./src/lib", import.meta.url).pathname,
		},
	},
	test: {
		environment: "jsdom",
		setupFiles: ["./vitest-setup.ts"],
		include: ["src/**/*.test.ts"],
		deps: {
			inline: [/svelte/],
		},
	},
});
