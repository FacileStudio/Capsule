<script lang="ts">
	import { Icon, IconButton, icons } from '@facile/muse';
	import { setTheme, theme } from '$lib/theme.svelte';
	import { t } from '$lib/i18n/index.svelte';

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
	aria-label={t('shell.theme')}
	onclick={() => setTheme(dark ? 'light' : 'dark')}
>
	<Icon icon={dark ? icons.sun : icons.moon} size={18} />
</IconButton>
