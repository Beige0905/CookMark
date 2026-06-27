import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

const apiTarget = process.env.API_URL ?? 'http://localhost:8080';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	ssr: {
		noExternal: ['lucide-svelte']
	},
	server: {
		host: true,
		hmr: { host: 'localhost' },
		proxy: {
			'/api': {
				target: apiTarget,
				changeOrigin: true
			}
		}
	}
});
