<script lang="ts">
	import { Icon, IconButton, icons } from '@facile/muse';
	import { theme } from '$lib/theme.svelte';
	import { t } from '$lib/i18n/index.svelte';

	let dark = $derived(theme.resolved === 'dark');
</script>

<IconButton
	variant="ghost"
	aria-label={t('shell.theme')}
	onclick={() => theme.set(dark ? 'light' : 'dark')}
>
	<Icon icon={dark ? icons.sun : icons.moon} size={18} />
</IconButton>
