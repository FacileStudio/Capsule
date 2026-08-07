<script lang="ts">
	import { IconButton, icons } from '@facile/muse';
	import { setTheme, theme } from '$lib/theme.svelte';

	let systemDark = $state(false);
	let dark = $derived(theme.mode === 'dark' || (theme.mode === 'system' && systemDark));

	$effect(() => {
		const mql = window.matchMedia('(prefers-color-scheme: dark)');
		systemDark = mql.matches;
		function onChange(e: MediaQueryListEvent) {
			systemDark = e.matches;
		}
		mql.addEventListener('change', onChange);
		return () => mql.removeEventListener('change', onChange);
	});
</script>

<IconButton
	variant="ghost"
	aria-label="Toggle theme"
	onclick={() => setTheme(dark ? 'light' : 'dark')}
>
	<iconify-icon
		icon={dark ? icons.sun : icons.moon}
		width="18"
		height="18"
		class="block"
	></iconify-icon>
</IconButton>
