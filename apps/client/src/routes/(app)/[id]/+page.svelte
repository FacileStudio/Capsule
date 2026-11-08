<script lang="ts">
	import { onMount } from 'svelte';
	import { Alert, Badge, Button, Icon, Input, Spinner, icons } from '@facile/muse';
	import { page } from '$app/state';
	import { importKey, decrypt, unwrapContentKey } from '$lib/crypto';
	import { parsePayload, type CapsuleFile } from '$lib/envelope';
	import RevealedFile from '$lib/components/RevealedFile.svelte';
	import { backend, type PasteMeta } from '$lib/backend';
	import { highlight } from '$lib/highlight';
	import { locale, t } from '$lib/i18n/index.svelte';
	import { CAPSULE_ICON } from '$lib/icons';
	import { timeAgo } from '$lib/i18n/format';

	type State = 'loading' | 'ready' | 'password' | 'revealing' | 'revealed' | 'empty' | 'error';

	let phase: State = $state('loading');
	let meta: PasteMeta | null = $state(null);
	let plaintext = $state('');
	let revealedFile: CapsuleFile | null = $state(null);
	let error = $state('');
	let copied = $state(false);
	let keyFragment = $state('');
	let password = $state('');
	let passwordError = $state('');
	let highlightedHtml = $state('');

	onMount(async () => {
		const id = page.params.id;
		keyFragment = window.location.hash.slice(1);

		if (!id) {
			error = t('reveal.noId');
			phase = 'error';
			return;
		}

		if (!keyFragment) {
			error = t('reveal.noKey');
			phase = 'error';
			return;
		}

		try {
			const result = await backend.getPasteMeta(id);
			if (!result.exists) {
				phase = 'empty';
				return;
			}
			meta = result;
			phase = result.has_password ? 'password' : 'ready';
		} catch {
			phase = 'empty';
		}
	});

	async function reveal() {
		if (!meta) return;
		phase = 'revealing';

		try {
			let key: CryptoKey;

			if (meta.has_password) {
				if (keyFragment.split('.').length !== 3) throw new Error(t('reveal.badFragment'));
				try {
					key = await unwrapContentKey(keyFragment, password);
				} catch {
					passwordError = t('reveal.wrongPassword');
					phase = 'password';
					return;
				}
			} else {
				key = await importKey(keyFragment);
			}

			const { content } = await backend.getPasteContent(meta.id);
			const decrypted = await decrypt(content, key);
			const payload = parsePayload(decrypted);
			plaintext = payload.text;
			revealedFile = payload.file;
			phase = 'revealed';

			if (plaintext && meta?.syntax && meta.syntax !== 'plaintext') {
				highlight(plaintext, meta.syntax).then((html) => {
					if (html) highlightedHtml = html;
				});
			}
		} catch {
			error = t('reveal.decryptFailed');
			phase = 'error';
		}
	}

	function submitPassword() {
		if (!password) return;
		passwordError = '';
		reveal();
	}

	async function copyContent() {
		await navigator.clipboard.writeText(plaintext);
		copied = true;
		setTimeout(() => (copied = false), 2000);
	}
</script>

<svelte:head>
	<title>{t('reveal.title')}</title>
	<meta name="description" content={t('reveal.description')} />
</svelte:head>

{#snippet meterBar()}
	<div class="flex flex-wrap gap-2">
		{#if meta?.syntax}
			<Badge>{meta.syntax}</Badge>
		{/if}
		{#if meta?.created_at}
			<Badge>{t('reveal.sealedAgo', { when: timeAgo(meta.created_at, locale()) })}</Badge>
		{/if}
		{#if meta?.has_password}
			<Badge tone="warning">{t('reveal.badgePassword')}</Badge>
		{/if}
		{#if meta?.burn_after_read}
			<Badge tone="danger">{t('reveal.badgeBurn')}</Badge>
		{/if}
	</div>
{/snippet}

{#if phase === 'loading'}
	<div class="flex flex-1 items-center justify-center">
		<Spinner />
	</div>
{:else if phase === 'empty'}
	<!-- Hand-built rather than `EmptyState`: that component is a Card for a list with nothing
	     in it, and this is the whole page — the title has to be the page's own <h1>. -->
	<div class="flex flex-col items-center gap-4 py-12 text-center sm:py-20">
		<Icon icon={CAPSULE_ICON} size={28} class="text-fc-fg-muted" />
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-xl font-semibold tracking-tight sm:text-fc-2xl">{t('reveal.emptyHeading')}</h1>
			<p class="mx-auto max-w-sm text-fc-sm text-fc-fg-muted">
				{t('reveal.emptyText')}
			</p>
		</div>
		<Button href="/" size="lg" icon={CAPSULE_ICON}>{t('reveal.sealNew')}</Button>
	</div>
{:else if phase === 'error'}
	<div class="flex flex-col items-center gap-4 py-12 text-center sm:py-20">
		<Icon icon={icons.warning} size={28} class="text-fc-danger" />
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-xl font-semibold tracking-tight sm:text-fc-2xl">{t('reveal.errorHeading')}</h1>
			<p class="mx-auto max-w-sm text-fc-sm text-fc-fg-muted">{error}</p>
		</div>
		<Button href="/" variant="secondary" size="lg">{t('reveal.back')}</Button>
	</div>
{:else if phase === 'password'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">{t('reveal.passwordHeading')}</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">
				{t('reveal.passwordText')}
			</p>
		</div>

		{@render meterBar()}

		<form
			onsubmit={(e: Event) => { e.preventDefault(); submitPassword(); }}
			class="flex flex-col gap-4"
		>
			<Input type="password" bind:value={password} placeholder={t('reveal.passwordPlaceholder')} />

			{#if passwordError}
				<Alert tone="danger">{passwordError}</Alert>
			{/if}

			<Button type="submit" size="lg" icon={icons.key} disabled={!password} class="w-fit">
				{t('reveal.unlock')}
			</Button>
		</form>
	</div>
{:else if phase === 'ready' || phase === 'revealing'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">{t('reveal.readyHeading')}</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">
				{#if meta?.burn_after_read}
					{t('reveal.readyBurn')}
				{:else}
					{t('reveal.description')}
				{/if}
			</p>
		</div>

		{@render meterBar()}

		<Button
			size="lg"
			icon={phase === 'revealing' ? undefined : icons.key}
			onclick={reveal}
			disabled={phase === 'revealing'}
			class="w-fit"
		>
			{#if phase === 'revealing'}
				<Spinner size="sm" label={t('reveal.openingLabel')} />
				{t('reveal.opening')}
			{:else}
				{t('reveal.open')}
			{/if}
		</Button>
	</div>
{:else if phase === 'revealed'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">{t('reveal.openedHeading')}</h1>
			{#if meta?.burn_after_read}
				<p class="text-fc-sm text-fc-danger">
					{t('reveal.destroyed')}
				</p>
			{/if}
		</div>

		{#if revealedFile}
			<RevealedFile file={revealedFile} />
		{/if}

		{#if plaintext || !revealedFile}
			<div class="relative">
				{#if highlightedHtml}
					<div class="w-full overflow-x-auto rounded-fc-md bg-fc-component p-4 pr-16 font-fc-mono text-fc-sm sm:pr-24 [&_pre]:!bg-transparent [&_pre]:!m-0 [&_pre]:!p-0 [&_code]:!bg-transparent">
						{@html highlightedHtml}
					</div>
				{:else}
					<pre class="w-full overflow-x-auto rounded-fc-md bg-fc-component p-4 pr-16 font-fc-mono text-fc-sm text-fc-fg whitespace-pre-wrap break-words sm:pr-24">{plaintext}</pre>
				{/if}
				<div class="absolute right-2 top-2 sm:right-3 sm:top-3">
					<Button variant="secondary" size="sm" icon={copied ? icons.check : icons.copy} onclick={copyContent}>
						{copied ? t('reveal.copied') : t('reveal.copy')}
					</Button>
				</div>
			</div>
		{/if}

		<Button href="/" variant="secondary" size="lg" icon={CAPSULE_ICON} class="w-fit">
			{t('reveal.sealNew')}
		</Button>
	</div>
{/if}
