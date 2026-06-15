// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'Nipo Tunnel',
			logo: {
				src: './src/assets/logo.png',
			},
			favicon: '/favicon.png',
			defaultLocale: 'root',
			locales: {
				root: {
					label: 'English',
					lang: 'en',
				},
				vi: {
					label: 'Tiếng Việt',
					lang: 'vi',
				},
			},
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/ngtuonghy/nipo-tunnel' }],
			sidebar: [
				{
					label: 'Guides',
					translations: {
						vi: 'Hướng dẫn',
					},
					items: [{ autogenerate: { directory: 'guides' } }],
				},
			],
		}),
	],
});
