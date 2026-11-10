<script lang="ts">
	import { LandingFooter, LandingHeader, LocaleSwitch, UserMenu } from '@facile/muse';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { CAPSULE_ICON } from '$lib/icons';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';
	import { backend, type Invitations, type SessionUser } from '$lib/backend';
	import { rememberLocale, t } from '$lib/i18n/index.svelte';
	import { setTheme, theme, type ThemeMode } from '$lib/theme.svelte';

	let { children } = $props();

	let user = $state<SessionUser | null>(null);
	let invitations = $state<Invitations | null>(null);

	onMount(async () => {
		user = await backend.sessionUser().catch(() => null);
		if (user) invitations = await backend.invitations().catch(() => null);
	});

	async function signOut() {
		await backend.logout().catch(() => undefined);
		user = null;
		invitations = null;
	}

	const menu = $derived({
		subscription: { href: 'https://sso.facile.studio/spaces' },
		helpHref: 'https://github.com/FacileStudio/Capsule/tree/main/docs',
		invitations,
		onSignOut: signOut,
		onLocaleChange: rememberLocale,
		theme: theme.mode,
		onThemeChange: (next: string) => setTheme(next as ThemeMode)
	});

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
			{#if user}
				<UserMenu
					trigger="avatar"
					user={{ name: user.name?.trim() || user.email, email: user.email, avatar: user.avatar_url }}
					{...menu}
				/>
			{:else}
				<LocaleSwitch onChange={rememberLocale} />
				<ThemeToggle />
			{/if}
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
