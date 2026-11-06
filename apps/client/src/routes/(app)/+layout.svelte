<script lang="ts">
	import { LandingFooter, LandingHeader, LocaleSwitch } from '@facile/muse';
	import { page } from '$app/state';
	import { CAPSULE_ICON } from '$lib/icons';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';
	import { rememberLocale, t } from '$lib/i18n/index.svelte';

	let { children } = $props();

	const path = $derived(page.url.pathname);
	const links = $derived([
		{ href: '/', label: t('shell.seal'), active: path === '/' },
		{ href: '/revoke', label: t('shell.revoke'), active: path === '/revoke' },
		{ href: '/docs', label: t('shell.api'), external: true }
	]);
	const columns = $derived([
		{ title: 'Capsule', links: links.map(({ href, label, external }) => ({ href, label, external })) },
		{ title: 'Facile', links: [{ href: 'https://facile.studio', label: t('shell.facile'), external: true }] }
	]);
</script>

<div class="flex min-h-[100dvh] flex-col">
	<LandingHeader appName="Capsule" icon={CAPSULE_ICON} {links} skipHref="#main">
		{#snippet actions()}
			<LocaleSwitch onChange={rememberLocale} />
			<ThemeToggle />
		{/snippet}
	</LandingHeader>

	<main id="main" class="mx-auto flex w-full max-w-2xl flex-1 flex-col px-4 py-8 sm:px-6 sm:py-12">
		{@render children()}
	</main>

	<LandingFooter
		appName="Capsule"
		icon={CAPSULE_ICON}
		tagline={t('shell.tagline')}
		{columns}
		legal={{ hideLocaleSwitch: true }}
		year={new Date().getFullYear()}
	/>
</div>
