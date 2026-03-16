import { defineConfig } from "astro/config";
import svelte from "@astrojs/svelte";
import tailwindcss from "@tailwindcss/vite";
import sitemap from "@astrojs/sitemap";

export default defineConfig({
	site: "https://recipe.shuttercoach.app",
	output: "static",
	integrations: [svelte(), sitemap()],
	vite: {
		plugins: [tailwindcss()],
		resolve: {
			alias: {
				$lib: new URL("./src/lib", import.meta.url).pathname,
			},
		},
	},
});
