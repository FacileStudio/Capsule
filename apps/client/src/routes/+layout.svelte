<script lang="ts">
	import '../app.css';
	import { browser } from '$app/environment';
	import { setTheme, theme } from '$lib/theme.svelte';

	let { children } = $props();

	/*
	 * muse renders <iconify-icon> elements and the consumer has to register the custom
	 * element. Imported dynamically because the package calls customElements.define at
	 * module scope. setTheme re-applies what the pre-paint script in app.html already wrote,
	 * so the store and the <html> classes cannot disagree.
	 */
	if (browser) {
		void import('iconify-icon');
		setTheme(theme.mode);
	}
</script>

{@render children()}
