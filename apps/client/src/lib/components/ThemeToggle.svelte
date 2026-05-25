<script lang="ts">
	type Theme = 'light' | 'dark' | 'system';

	let theme = $state<Theme>('system');
	let systemDark = $state(false);
	let isDark = $derived(theme === 'dark' || (theme === 'system' && systemDark));

	function toggle() {
		theme = isDark ? 'light' : 'dark';
		localStorage.setItem('capsule-theme', theme);
	}

	$effect(() => {
		const stored = localStorage.getItem('capsule-theme') as Theme | null;
		if (stored === 'dark' || stored === 'light') {
			theme = stored;
		} else {
			theme = 'system';
		}

		const mql = window.matchMedia('(prefers-color-scheme: dark)');
		systemDark = mql.matches;
		function onChange(e: MediaQueryListEvent) {
			systemDark = e.matches;
		}
		mql.addEventListener('change', onChange);
		return () => mql.removeEventListener('change', onChange);
	});

	$effect(() => {
		document.documentElement.classList.toggle('dark', isDark);
	});
</script>

<button
	onclick={toggle}
	aria-label="Toggle theme"
	class="inline-flex h-11 w-11 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground sm:h-10 sm:w-10"
>
	{#if isDark}
		<iconify-icon icon="solar:sun-bold" width="20"></iconify-icon>
	{:else}
		<iconify-icon icon="solar:moon-bold" width="20"></iconify-icon>
	{/if}
</button>
