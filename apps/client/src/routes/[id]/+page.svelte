<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { importKey, decrypt } from '$lib/crypto';
	import { backend, type PasteMeta } from '$lib/backend';

	type State = 'loading' | 'ready' | 'revealing' | 'revealed' | 'empty' | 'error';

	let state: State = $state('loading');
	let meta: PasteMeta | null = $state(null);
	let plaintext = $state('');
	let error = $state('');
	let copied = $state(false);
	let keyFragment = $state('');

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
			state = 'ready';
		} catch {
			state = 'empty';
		}
	});

	async function reveal() {
		if (!meta) return;
		state = 'revealing';

		try {
			const { content } = await backend.getPasteContent(meta.id);
			const key = await importKey(keyFragment);
			plaintext = await decrypt(content, key);
			state = 'revealed';
		} catch {
			error = 'Failed to decrypt. The key may be wrong or the data may be corrupted.';
			state = 'error';
		}
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
			class="inline-flex h-12 items-center justify-center rounded-md bg-primary px-6 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed"
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
			<pre class="w-full overflow-x-auto rounded-lg border border-input bg-card p-3 pr-16 font-mono text-sm text-card-foreground whitespace-pre-wrap break-words sm:p-4 sm:pr-20">{plaintext}</pre>
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
