<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { Alert, Button, EmptyState, Spinner, icons } from '@facile/muse';
	import { backend, type PendingAction } from '$lib/backend';
	import PendingCard from '$lib/components/PendingCard.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import { SIGN_IN_HREF, rememberReturn } from '$lib/signin';

	let actions = $state<PendingAction[]>([]);
	let phase = $state<'loading' | 'signedOut' | 'ready'>('loading');
	let busy = $state('');
	let notice = $state<{ tone: 'success' | 'danger'; text: string } | null>(null);

	const focused = $derived(page.url.searchParams.get('action') ?? '');
	const ordered = $derived([
		...actions.filter((action) => action.id === focused),
		...actions.filter((action) => action.id !== focused)
	]);

	onMount(async () => {
		const user = await backend.sessionUser().catch(() => null);
		if (!user) {
			phase = 'signedOut';
			return;
		}
		await load();
	});

	/** Loads the open actions. */
	async function load() {
		try {
			actions = (await backend.listPending()).actions;
		} catch {
			notice = { tone: 'danger', text: t('pending.loadFailed') };
		}
		phase = 'ready';
	}

	/** Sends the person through sign-in, then back to this page. */
	function signIn() {
		rememberReturn(page.url.pathname + page.url.search);
		location.assign(SIGN_IN_HREF);
	}

	/** Runs one action once, then drops it from the list. */
	async function confirm(action: PendingAction) {
		busy = action.id;
		try {
			await backend.confirmPending(action.id);
			notice = { tone: 'success', text: t('pending.revoked') };
			actions = actions.filter((item) => item.id !== action.id);
		} catch (err) {
			notice = { tone: 'danger', text: err instanceof Error ? err.message : t('pending.failed') };
			await load();
		}
		busy = '';
	}

	/** Closes one action without running it. */
	async function dismiss(action: PendingAction) {
		busy = action.id;
		try {
			await backend.dismissPending(action.id);
			notice = { tone: 'success', text: t('pending.dismissed') };
			actions = actions.filter((item) => item.id !== action.id);
		} catch (err) {
			notice = { tone: 'danger', text: err instanceof Error ? err.message : t('pending.failed') };
		}
		busy = '';
	}
</script>

<svelte:head>
	<title>{t('pending.title')}</title>
</svelte:head>

<div class="flex w-full flex-col gap-6">
	<div class="flex flex-col gap-1">
		<h1 class="text-fc-2xl font-bold tracking-tight sm:text-fc-3xl">{t('pending.heading')}</h1>
		<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">{t('pending.intro')}</p>
	</div>

	{#if notice}
		<Alert tone={notice.tone}>{notice.text}</Alert>
	{/if}

	{#if phase === 'loading'}
		<div class="flex h-48 items-center justify-center"><Spinner /></div>
	{:else if phase === 'signedOut'}
		<div class="flex flex-col gap-4">
			<EmptyState icon={icons.key} title={t('pending.signInHeading')} description={t('pending.signInBody')} />
			<Button size="lg" class="w-full sm:w-fit sm:self-center" onclick={signIn}>{t('pending.signIn')}</Button>
		</div>
	{:else if actions.length === 0}
		<EmptyState icon={icons.check} title={t('pending.none')} description={t('pending.noneBody')} />
	{:else}
		<section class="flex flex-col gap-4">
			{#each ordered as action (action.id)}
				<PendingCard {action} busy={busy === action.id} onconfirm={confirm} ondismiss={dismiss} />
			{/each}
		</section>
	{/if}
</div>
