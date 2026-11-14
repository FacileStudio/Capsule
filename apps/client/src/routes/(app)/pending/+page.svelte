<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { Alert, Button, EmptyState, Spinner, icons } from '@facile/muse';
	import { backend, type PendingAction } from '$lib/backend';
	import PendingCard from '$lib/components/PendingCard.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import { refusesOwnSession } from '$lib/session-refusal';
	import { SIGN_IN_HREF, rememberReturn } from '$lib/signin';

	let actions = $state<PendingAction[]>([]);
	let phase = $state<'loading' | 'signedOut' | 'ready'>('loading');
	let busy = $state('');
	let notice = $state<{ tone: 'success' | 'danger'; text: string } | null>(null);
	let notOwnSession = $state(false);

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

	/** Reports whether err is the refusal of a session no browser sign-in opened. */
	function refusedSession(err: unknown) {
		notOwnSession = refusesOwnSession(err);
		if (notOwnSession) notice = null;
		return notOwnSession;
	}

	/** The server's own words for any other 403, else the fallback text. */
	function reason(err: unknown, fallback: string): string {
		const refusal = err as { status?: unknown } | null;
		return refusal?.status === 403 && err instanceof Error ? err.message : fallback;
	}

	/** Loads the open actions. */
	async function load() {
		try {
			actions = (await backend.listPending()).actions;
		} catch (err) {
			if (!refusedSession(err)) notice = { tone: 'danger', text: reason(err, t('pending.loadFailed')) };
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
			if (!refusedSession(err)) {
				notice = { tone: 'danger', text: err instanceof Error ? err.message : t('pending.failed') };
				await load();
			}
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
			if (!refusedSession(err)) notice = { tone: 'danger', text: err instanceof Error ? err.message : t('pending.failed') };
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

	{#if notOwnSession}
		<Alert tone="warning" title={t('pending.signInAgain')}>
			<p>{t('pending.signInAgainBody')}</p>
			<Button size="sm" class="mt-3" onclick={signIn}>{t('pending.signInAgainButton')}</Button>
		</Alert>
	{:else if phase === 'loading'}
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
