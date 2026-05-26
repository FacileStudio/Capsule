<script lang="ts">
	import { backend } from '$lib/backend';

	type State = 'idle' | 'revoking' | 'revoked' | 'error';

	let state: State = $state('idle');
	let capsuleUrl = $state('');
	let deleteToken = $state('');
	let error = $state('');

	function extractId(input: string): string {
		const trimmed = input.trim();
		if (trimmed.startsWith('cap_')) return trimmed.split('#')[0];
		try {
			const url = new URL(trimmed);
			const segments = url.pathname.split('/').filter(Boolean);
			return segments[segments.length - 1] || '';
		} catch {
			return trimmed;
		}
	}

	async function revoke() {
		const id = extractId(capsuleUrl);
		if (!id || !deleteToken.trim()) return;

		state = 'revoking';
		error = '';

		try {
			await backend.deletePaste(id, deleteToken.trim());
			state = 'revoked';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to revoke capsule';
			state = 'error';
		}
	}

	function reset() {
		state = 'idle';
		capsuleUrl = '';
		deleteToken = '';
		error = '';
	}
</script>

<svelte:head>
	<title>Revoke — Capsule</title>
</svelte:head>

{#if state === 'revoked'}
	<div class="flex flex-col gap-5 sm:gap-6">
		<div>
			<div class="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-red-500/10">
				<iconify-icon icon="solar:fire-bold" width="28" class="text-red-400"></iconify-icon>
			</div>
			<h1 class="text-2xl font-bold font-heading tracking-tight sm:text-3xl">Capsule revoked</h1>
			<p class="mt-2 text-sm text-muted-foreground sm:text-base">
				The capsule has been burned. Its content is permanently gone.
			</p>
		</div>

		<div class="flex gap-3">
			<button
				onclick={reset}
				class="inline-flex h-12 items-center justify-center rounded-md border border-border bg-background px-5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground sm:h-10"
			>
				Revoke another
			</button>
			<a
				href="/"
				class="inline-flex h-12 items-center justify-center rounded-md bg-primary px-5 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 sm:h-10"
			>
				<iconify-icon icon="solar:pill-bold-duotone" width="16" class="mr-2"></iconify-icon>
				Seal a capsule
			</a>
		</div>
	</div>
{:else}
	<div class="flex flex-col gap-5 sm:gap-6">
		<div>
			<h1 class="text-2xl font-bold font-heading tracking-tight sm:text-3xl">Revoke a capsule</h1>
			<p class="mt-2 text-sm text-muted-foreground sm:text-base">
				Burn a capsule before it's opened. You'll need the capsule URL or ID and the delete token you received when sealing.
			</p>
		</div>

		<div class="flex flex-col gap-4">
			<div class="flex flex-col gap-1.5">
				<label for="capsule-url" class="text-sm font-medium">Capsule URL or ID</label>
				<input
					id="capsule-url"
					type="text"
					bind:value={capsuleUrl}
					placeholder="https://capsule.facile.studio/cap_... or cap_..."
					disabled={state === 'revoking'}
					class="w-full rounded-lg border border-input bg-card px-4 py-3 font-mono text-sm text-card-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
				/>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="delete-token" class="text-sm font-medium">Delete token</label>
				<input
					id="delete-token"
					type="text"
					bind:value={deleteToken}
					placeholder="Paste your delete token here..."
					disabled={state === 'revoking'}
					class="w-full rounded-lg border border-input bg-card px-4 py-3 font-mono text-sm text-card-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
				/>
			</div>
		</div>

		{#if error}
			<p class="text-sm text-red-400">{error}</p>
		{/if}

		<button
			onclick={revoke}
			disabled={!capsuleUrl.trim() || !deleteToken.trim() || state === 'revoking'}
			class="inline-flex h-12 items-center justify-center rounded-md border border-red-500/30 bg-red-500/10 px-6 text-sm font-medium text-red-400 transition-colors hover:bg-red-500/20 disabled:opacity-50 disabled:cursor-not-allowed sm:h-11"
		>
			{#if state === 'revoking'}
				<iconify-icon icon="solar:refresh-bold" width="16" class="mr-2 animate-spin"></iconify-icon>
				Revoking...
			{:else}
				<iconify-icon icon="solar:fire-bold" width="16" class="mr-2"></iconify-icon>
				Revoke capsule
			{/if}
		</button>

		<a
			href="/"
			class="inline-flex h-12 items-center justify-center rounded-md border border-border bg-background px-6 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground sm:h-11"
		>
			<iconify-icon icon="solar:pill-bold-duotone" width="16" class="mr-2"></iconify-icon>
			Seal a capsule
		</a>
	</div>
{/if}
