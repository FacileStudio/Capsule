<script lang="ts">
	import { Alert, Button, Field, Icon, Input, Spinner, icons } from '@facile/muse';
	import { backend } from '$lib/backend';
	import { t } from '$lib/i18n/index.svelte';
	import { CAPSULE_ICON } from '$lib/icons';

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
			error = err instanceof Error ? err.message : t('revoke.failed');
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
	<title>{t('revoke.title')}</title>
</svelte:head>

{#if phase === 'revoked'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<Icon icon={icons.remove} size={28} class="mb-3 text-fc-danger" />
			<h1 class="text-fc-2xl font-bold tracking-tight sm:text-fc-3xl">{t('revoke.doneHeading')}</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">
				{t('revoke.doneText')}
			</p>
		</div>

		<div class="flex flex-col gap-3 sm:flex-row">
			<Button variant="secondary" size="lg" onclick={reset} class="sm:w-fit">{t('revoke.another')}</Button>
			<Button href="/" size="lg" icon={CAPSULE_ICON} class="sm:w-fit">{t('revoke.seal')}</Button>
		</div>
	</div>
{:else}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-bold tracking-tight sm:text-fc-3xl">{t('revoke.heading')}</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">
				{t('revoke.intro')}
			</p>
		</div>

		<div class="flex flex-col gap-4">
			<Field label={t('revoke.urlLabel')}>
				<Input
					bind:value={capsuleUrl}
					placeholder="https://capsule.facile.studio/cap_... or cap_..."
					disabled={phase === 'revoking'}
					class="font-fc-mono text-fc-sm"
				/>
			</Field>

			<Field label={t('revoke.tokenLabel')}>
				<Input
					bind:value={deleteToken}
					placeholder={t('revoke.tokenPlaceholder')}
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
					<Spinner size="sm" label={t('revoke.revokingLabel')} />
					{t('revoke.revoking')}
				{:else}
					{t('revoke.submit')}
				{/if}
			</Button>
			<Button href="/" variant="secondary" size="lg" icon={CAPSULE_ICON} class="sm:w-fit">
				{t('revoke.seal')}
			</Button>
		</div>
	</div>
{/if}
