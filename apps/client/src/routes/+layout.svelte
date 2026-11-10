<script lang="ts">
	import '../app.css';
	import { browser } from '$app/environment';
	import { registerAppIcons } from '$lib/icons';
	import { provideLocale } from '@facile/muse';
	import { theme } from '$lib/theme.svelte';
	import { bindLocale, localeSource } from '$lib/i18n/index.svelte';

	let { children } = $props();

	const i18n = provideLocale(localeSource);
	bindLocale(() => i18n.locale);

	registerAppIcons();

	/*
	 * restore re-applies what the pre-paint script in app.html already wrote, so the store
	 * and the <html> classes cannot disagree.
	 */
	if (browser) theme.restore();
</script>

{@render children()}
