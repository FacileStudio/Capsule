<script lang="ts">
	import { Alert, Button, Field, Input, Spinner, icons } from '@facile/muse';
	import { backend } from '$lib/backend';

	type State = 'idle' | 'revoking' | 'revoked' | 'error';

	let phase: State = $state('idle');
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

		phase = 'revoking';
		error = '';

		try {
			await backend.deletePaste(id, deleteToken.trim());
			phase = 'revoked';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to revoke capsule';
			phase = 'error';
		}
	}

	function reset() {
		phase = 'idle';
		capsuleUrl = '';
		deleteToken = '';
		error = '';
	}
</script>

<svelte:head>
	<title>Revoke — Capsule</title>
</svelte:head>

{#if phase === 'revoked'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<iconify-icon icon={icons.remove} width="28" height="28" class="mb-3 block text-fc-danger"></iconify-icon>
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">Capsule revoked</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">
				The capsule has been burned. Its content is permanently gone.
			</p>
		</div>

		<div class="flex flex-col gap-3 sm:flex-row">
			<Button variant="outline" size="lg" onclick={reset} class="sm:w-fit">Revoke another</Button>
			<Button href="/" size="lg" icon="solar:pill-bold-duotone" class="sm:w-fit">Seal a capsule</Button>
		</div>
	</div>
{:else}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">Revoke a capsule</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">
				Burn a capsule before it's opened. You'll need the capsule URL or ID and the delete token you received when sealing.
			</p>
		</div>

		<div class="flex flex-col gap-4">
			<Field label="Capsule URL or ID">
				<Input
					bind:value={capsuleUrl}
					placeholder="https://capsule.facile.studio/cap_... or cap_..."
					disabled={phase === 'revoking'}
					class="font-fc-mono text-fc-sm"
				/>
			</Field>

			<Field label="Delete token">
				<Input
					bind:value={deleteToken}
					placeholder="Paste your delete token here..."
					disabled={phase === 'revoking'}
					class="font-fc-mono text-fc-sm"
				/>
			</Field>
		</div>

		{#if error}
			<Alert tone="danger">{error}</Alert>
		{/if}

		<div class="flex flex-col gap-3 sm:flex-row">
			<Button
				variant="danger"
				size="lg"
				icon={phase === 'revoking' ? undefined : icons.remove}
				onclick={revoke}
				disabled={!capsuleUrl.trim() || !deleteToken.trim() || phase === 'revoking'}
				class="sm:w-fit"
			>
				{#if phase === 'revoking'}
					<Spinner size="sm" label="Revoking" />
					Revoking...
				{:else}
					Revoke capsule
				{/if}
			</Button>
			<Button href="/" variant="outline" size="lg" icon="solar:pill-bold-duotone" class="sm:w-fit">
				Seal a capsule
			</Button>
		</div>
	</div>
{/if}
