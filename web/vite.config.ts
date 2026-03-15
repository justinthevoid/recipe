import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { svelteTesting } from "@testing-library/svelte/vite";
import { defineConfig } from "vite";

export default defineConfig({
	plugins: [svelte(), tailwindcss(), svelteTesting()],
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
