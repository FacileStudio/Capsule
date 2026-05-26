<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { importKey, decrypt, unwrapContentKey, parsePasswordFragment } from '$lib/crypto';
	import { backend, type PasteMeta } from '$lib/backend';
	import { highlight } from '$lib/highlight';

	type State = 'loading' | 'ready' | 'password' | 'revealing' | 'revealed' | 'empty' | 'error';

	let state: State = $state('loading');
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

		if (!keyFragment) {
			error = 'No decryption key found in URL. The link may be incomplete.';
			state = 'error';
			return;
		}

		try {
			const result = await backend.getPasteMeta(id);
			if (!result.exists) {
				state = 'empty';
				return;
			}
			meta = result;
			state = result.has_password ? 'password' : 'ready';
		} catch {
			state = 'empty';
		}
	});

	async function reveal() {
		if (!meta) return;
		state = 'revealing';

		try {
			let key: CryptoKey;

			if (meta.has_password) {
				const parsed = parsePasswordFragment(keyFragment);
				if (!parsed) throw new Error('Invalid URL fragment');
				try {
					key = await unwrapContentKey(parsed.encryptedKey, parsed.salt, parsed.iv, password);
				} catch {
					passwordError = 'Wrong password or corrupted link.';
					state = 'password';
					return;
				}
			} else {
				key = await importKey(keyFragment);
			}

			const { content } = await backend.getPasteContent(meta.id);
			plaintext = await decrypt(content, key);
			state = 'revealed';

			if (meta?.syntax && meta.syntax !== 'plaintext') {
				highlight(plaintext, meta.syntax).then((html) => {
					if (html) highlightedHtml = html;
				});
			}
		} catch {
			error = 'Failed to decrypt. The key may be wrong or the data may be corrupted.';
			state = 'error';
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

{#if state === 'loading'}
	<div class="flex flex-1 items-center justify-center">
		<iconify-icon icon="solar:refresh-bold" width="24" class="animate-spin text-muted-foreground"></iconify-icon>
	</div>
{:else if state === 'empty'}
	<div class="flex flex-col items-center gap-4 py-12 text-center sm:py-20">
		<div class="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
			<iconify-icon icon="solar:pill-bold-duotone" width="32" class="text-muted-foreground"></iconify-icon>
		</div>
		<h1 class="text-xl font-bold font-heading sm:text-2xl">This capsule is empty</h1>
		<p class="max-w-sm text-sm text-muted-foreground sm:text-base">
			It was already opened, expired, or never existed.
		</p>
		<a
			href="/"
			class="mt-4 inline-flex h-12 items-center justify-center rounded-md bg-primary px-5 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 sm:h-10"
		>
			Seal a new capsule
		</a>
	</div>
{:else if state === 'error'}
	<div class="flex flex-col items-center gap-4 py-12 text-center sm:py-20">
		<div class="flex h-16 w-16 items-center justify-center rounded-full bg-red-500/10">
			<iconify-icon icon="solar:danger-triangle-bold" width="32" class="text-red-400"></iconify-icon>
		</div>
		<h1 class="text-xl font-bold font-heading sm:text-2xl">Something went wrong</h1>
		<p class="max-w-sm text-sm text-muted-foreground sm:text-base">{error}</p>
		<a
			href="/"
			class="mt-4 inline-flex h-12 items-center justify-center rounded-md border border-border bg-background px-5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground sm:h-10"
		>
			Back to Capsule
		</a>
	</div>
{:else if state === 'password'}
	<div class="flex flex-col gap-5 sm:gap-6">
		<div>
			<h1 class="text-2xl font-bold font-heading tracking-tight sm:text-3xl">Password required</h1>
			<p class="mt-2 text-sm text-muted-foreground sm:text-base">
				This capsule is password-protected. Enter the password to decrypt it.
			</p>
		</div>

		<div class="flex flex-wrap gap-2 text-sm text-muted-foreground sm:gap-3">
			{#if meta?.syntax}
				<span class="rounded-md bg-secondary px-2 py-1">{meta.syntax}</span>
			{/if}
			{#if meta?.created_at}
				<span class="rounded-md bg-secondary px-2 py-1">Sealed {timeAgo(meta.created_at)}</span>
			{/if}
			<span class="rounded-md bg-amber-500/10 px-2 py-1 text-amber-400">Password protected</span>
		</div>

		<form
			onsubmit={(e: Event) => { e.preventDefault(); submitPassword(); }}
			class="flex flex-col gap-4"
		>
			<input
				type="password"
				bind:value={password}
				placeholder="Enter password..."
				class="w-full rounded-lg border border-input bg-card px-4 py-3 text-sm text-card-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
			/>

			{#if passwordError}
				<p class="text-sm text-red-400">{passwordError}</p>
			{/if}

			<button
				type="submit"
				disabled={!password}
				class="inline-flex h-12 items-center justify-center rounded-md bg-primary px-6 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed sm:h-11"
			>
				<iconify-icon icon="solar:lock-keyhole-unlocked-bold" width="16" class="mr-2"></iconify-icon>
				Unlock capsule
			</button>
		</form>
	</div>
{:else if state === 'ready' || state === 'revealing'}
	<div class="flex flex-col gap-5 sm:gap-6">
		<div>
			<h1 class="text-2xl font-bold font-heading tracking-tight sm:text-3xl">A capsule for you</h1>
			<p class="mt-2 text-sm text-muted-foreground sm:text-base">
				{#if meta?.burn_after_read}
					This capsule will self-destruct after you open it.
				{:else}
					Someone sent you an encrypted capsule.
				{/if}
			</p>
		</div>

		<div class="flex flex-wrap gap-2 text-sm text-muted-foreground sm:gap-3">
			{#if meta?.syntax}
				<span class="rounded-md bg-secondary px-2 py-1">{meta.syntax}</span>
			{/if}
			{#if meta?.created_at}
				<span class="rounded-md bg-secondary px-2 py-1">Sealed {timeAgo(meta.created_at)}</span>
			{/if}
			{#if meta?.burn_after_read}
				<span class="rounded-md bg-red-500/10 px-2 py-1 text-red-400">Burns after opening</span>
			{/if}
		</div>

		<button
			onclick={reveal}
			disabled={state === 'revealing'}
			class="inline-flex h-12 items-center justify-center rounded-md bg-primary px-6 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed sm:h-11"
		>
			{#if state === 'revealing'}
				<iconify-icon icon="solar:refresh-bold" width="16" class="mr-2 animate-spin"></iconify-icon>
				Opening...
			{:else}
				<iconify-icon icon="solar:lock-keyhole-unlocked-bold" width="16" class="mr-2"></iconify-icon>
				Open capsule
			{/if}
		</button>
	</div>
{:else if state === 'revealed'}
	<div class="flex flex-col gap-5 sm:gap-6">
		<div>
			<h1 class="text-2xl font-bold font-heading tracking-tight sm:text-3xl">Capsule opened</h1>
			{#if meta?.burn_after_read}
				<p class="mt-2 text-sm text-red-400">
					Content destroyed from server. This is the only time you'll see it.
				</p>
			{/if}
		</div>

		<div class="relative">
			{#if highlightedHtml}
				<div class="w-full overflow-x-auto rounded-lg bg-secondary/50 p-3 pr-16 font-mono text-sm sm:p-4 sm:pr-20 [&_pre]:!bg-transparent [&_pre]:!m-0 [&_pre]:!p-0 [&_code]:!bg-transparent">
					{@html highlightedHtml}
				</div>
			{:else}
				<pre class="w-full overflow-x-auto rounded-lg border border-input bg-card p-3 pr-16 font-mono text-sm text-card-foreground whitespace-pre-wrap break-words sm:p-4 sm:pr-20">{plaintext}</pre>
			{/if}
			<button
				onclick={copyContent}
				class="absolute right-2 top-2 inline-flex min-h-[44px] min-w-[44px] items-center justify-center rounded-md bg-secondary px-3 py-1.5 text-xs font-medium text-secondary-foreground transition-colors hover:bg-secondary/80 sm:right-3 sm:top-3 sm:min-h-0 sm:min-w-0"
			>
				{#if copied}
					Copied
				{:else}
					Copy
				{/if}
			</button>
		</div>

		<a
			href="/"
			class="inline-flex h-12 items-center justify-center rounded-md border border-border bg-background px-5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground sm:h-10"
		>
			<iconify-icon icon="solar:pill-bold-duotone" width="16" class="mr-2"></iconify-icon>
			Seal a new capsule
		</a>
	</div>
{/if}
