<script lang="ts">
	import { onMount } from 'svelte';
	import { Alert, Badge, Button, Input, Spinner, icons } from '@facile/muse';
	import { page } from '$app/state';
	import { importKey, decrypt, unwrapContentKey, parsePasswordFragment } from '$lib/crypto';
	import { backend, type PasteMeta } from '$lib/backend';
	import { highlight } from '$lib/highlight';

	type State = 'loading' | 'ready' | 'password' | 'revealing' | 'revealed' | 'empty' | 'error';

	let phase: State = $state('loading');
	let meta: PasteMeta | null = $state(null);
	let plaintext = $state('');
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
			error = 'No capsule id in the URL.';
			phase = 'error';
			return;
		}

		if (!keyFragment) {
			error = 'No decryption key found in URL. The link may be incomplete.';
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
				const parsed = parsePasswordFragment(keyFragment);
				if (!parsed) throw new Error('Invalid URL fragment');
				try {
					key = await unwrapContentKey(parsed.encryptedKey, parsed.salt, parsed.iv, password);
				} catch {
					passwordError = 'Wrong password or corrupted link.';
					phase = 'password';
					return;
				}
			} else {
				key = await importKey(keyFragment);
			}

			const { content } = await backend.getPasteContent(meta.id);
			plaintext = await decrypt(content, key);
			phase = 'revealed';

			if (meta?.syntax && meta.syntax !== 'plaintext') {
				highlight(plaintext, meta.syntax).then((html) => {
					if (html) highlightedHtml = html;
				});
			}
		} catch {
			error = 'Failed to decrypt. The key may be wrong or the data may be corrupted.';
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

	function timeAgo(dateStr: string): string {
		const seconds = Math.floor((Date.now() - new Date(dateStr).getTime()) / 1000);
		if (seconds < 60) return 'just now';
		const minutes = Math.floor(seconds / 60);
		if (minutes < 60) return `${minutes}m ago`;
		const hours = Math.floor(minutes / 60);
		if (hours < 24) return `${hours}h ago`;
		const days = Math.floor(hours / 24);
		return `${days}d ago`;
	}
</script>

<svelte:head>
	<title>Capsule — Open</title>
	<meta name="description" content="Someone sent you an encrypted capsule." />
</svelte:head>

{#snippet meterBar()}
	<div class="flex flex-wrap gap-2">
		{#if meta?.syntax}
			<Badge>{meta.syntax}</Badge>
		{/if}
		{#if meta?.created_at}
			<Badge>Sealed {timeAgo(meta.created_at)}</Badge>
		{/if}
		{#if meta?.has_password}
			<Badge tone="warning">Password protected</Badge>
		{/if}
		{#if meta?.burn_after_read}
			<Badge tone="danger">Burns after opening</Badge>
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
		<iconify-icon icon="solar:pill-bold-duotone" width="28" height="28" class="block text-fc-fg-muted"></iconify-icon>
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-xl font-semibold tracking-tight sm:text-fc-2xl">This capsule is empty</h1>
			<p class="mx-auto max-w-sm text-fc-sm text-fc-fg-muted">
				It was already opened, expired, or never existed.
			</p>
		</div>
		<Button href="/" size="lg" icon="solar:pill-bold-duotone">Seal a new capsule</Button>
	</div>
{:else if phase === 'error'}
	<div class="flex flex-col items-center gap-4 py-12 text-center sm:py-20">
		<iconify-icon icon={icons.warning} width="28" height="28" class="block text-fc-danger"></iconify-icon>
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-xl font-semibold tracking-tight sm:text-fc-2xl">Something went wrong</h1>
			<p class="mx-auto max-w-sm text-fc-sm text-fc-fg-muted">{error}</p>
		</div>
		<Button href="/" variant="outline" size="lg">Back to Capsule</Button>
	</div>
{:else if phase === 'password'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">Password required</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">
				This capsule is password-protected. Enter the password to decrypt it.
			</p>
		</div>

		{@render meterBar()}

		<form
			onsubmit={(e: Event) => { e.preventDefault(); submitPassword(); }}
			class="flex flex-col gap-4"
		>
			<Input type="password" bind:value={password} placeholder="Enter password..." />

			{#if passwordError}
				<Alert tone="danger">{passwordError}</Alert>
			{/if}

			<Button type="submit" size="lg" icon={icons.key} disabled={!password} class="w-fit">
				Unlock capsule
			</Button>
		</form>
	</div>
{:else if phase === 'ready' || phase === 'revealing'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">A capsule for you</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">
				{#if meta?.burn_after_read}
					This capsule will self-destruct after you open it.
				{:else}
					Someone sent you an encrypted capsule.
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
				<Spinner size="sm" label="Opening" />
				Opening...
			{:else}
				Open capsule
			{/if}
		</Button>
	</div>
{:else if phase === 'revealed'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">Capsule opened</h1>
			{#if meta?.burn_after_read}
				<p class="text-fc-sm text-fc-danger">
					Content destroyed from server. This is the only time you'll see it.
				</p>
			{/if}
		</div>

		<div class="relative">
			{#if highlightedHtml}
				<div class="w-full overflow-x-auto rounded-fc-md bg-fc-component p-4 pr-16 font-fc-mono text-fc-sm sm:pr-24 [&_pre]:!bg-transparent [&_pre]:!m-0 [&_pre]:!p-0 [&_code]:!bg-transparent">
					{@html highlightedHtml}
				</div>
			{:else}
				<pre class="w-full overflow-x-auto rounded-fc-md bg-fc-component p-4 pr-16 font-fc-mono text-fc-sm text-fc-fg whitespace-pre-wrap break-words sm:pr-24">{plaintext}</pre>
			{/if}
			<div class="absolute right-2 top-2 sm:right-3 sm:top-3">
				<Button variant="outline" size="sm" icon={copied ? icons.check : icons.copy} onclick={copyContent}>
					{copied ? 'Copied' : 'Copy'}
				</Button>
			</div>
		</div>

		<Button href="/" variant="outline" size="lg" icon="solar:pill-bold-duotone" class="w-fit">
			Seal a new capsule
		</Button>
	</div>
{/if}
