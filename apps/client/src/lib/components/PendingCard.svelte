<script lang="ts">
	import { Badge, Button, Card, Icon, icons } from '@facile/muse';
	import type { PendingAction } from '$lib/backend';
	import { formatDateTime } from '$lib/i18n/format';
	import { locale, t } from '$lib/i18n/index.svelte';
	import { CAPSULE_ICON } from '$lib/icons';

	let {
		action,
		busy,
		onconfirm,
		ondismiss
	}: {
		action: PendingAction;
		busy: boolean;
		onconfirm: (action: PendingAction) => void;
		ondismiss: (action: PendingAction) => void;
	} = $props();

	const capsule = $derived(action.capsule);

	/** Names the client that asked: the host of a client URL, else its id. */
	function clientName(clientId: string): string {
		if (clientId === 'bulle') return 'Bulle';
		try {
			return new URL(clientId).host;
		} catch {
			return clientId || '?';
		}
	}
</script>

<Card class="flex flex-col gap-4">
	<div class="flex flex-wrap items-center justify-between gap-2">
		<Badge tone="danger" class="shrink-0 whitespace-nowrap">{t('pending.kind.revoke')}</Badge>
		<span class="text-fc-xs text-fc-fg-muted">
			{t('pending.askedBy', { client: clientName(action.client_id) })} · {formatDateTime(action.created_at, locale())}
		</span>
	</div>

	{#if capsule}
		<div class="flex items-center gap-3 rounded-fc-md bg-fc-surface p-4">
			<Icon icon={CAPSULE_ICON} size={24} class="shrink-0 text-fc-fg-muted" />
			<div class="flex min-w-0 flex-1 flex-col gap-1">
				<p class="font-fc-mono text-fc-sm font-medium break-all text-fc-fg">{capsule.id}</p>
				<p class="text-fc-sm text-fc-fg-muted">
					{t('pending.sealed', { date: formatDateTime(capsule.created_at, locale()) })}
					{#if capsule.expires_at}
						· {t('pending.expires', { date: formatDateTime(capsule.expires_at, locale()) })}
					{/if}
				</p>
				<div class="flex flex-wrap gap-2">
					{#if capsule.burn_after_read}
						<Badge class="shrink-0 whitespace-nowrap">{t('pending.burns')}</Badge>
					{/if}
					{#if capsule.has_password}
						<Badge class="shrink-0 whitespace-nowrap">{t('pending.password')}</Badge>
					{/if}
				</div>
			</div>
		</div>
		<p class="text-fc-sm text-fc-fg-muted">{t('pending.revokeNote')}</p>
	{:else}
		<p class="text-fc-sm break-words text-fc-fg-muted">{action.summary}</p>
	{/if}

	{#if action.state === 'changed'}
		<p class="text-fc-sm text-fc-warning">{t('pending.changed')}</p>
	{:else if action.state === 'missing'}
		<p class="text-fc-sm text-fc-warning">{t('pending.missing')}</p>
	{/if}

	<div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
		<Button variant="ghost" class="w-full sm:w-auto [@media(pointer:coarse)]:h-11" disabled={busy} onclick={() => ondismiss(action)}>
			{t('pending.dismiss')}
		</Button>
		{#if action.state === 'ready'}
			<Button
				variant="danger"
				class="w-full sm:w-auto [@media(pointer:coarse)]:h-11"
				disabled={busy}
				icon={icons.remove}
				onclick={() => onconfirm(action)}
			>
				{t('pending.confirmRevoke')}
			</Button>
		{/if}
	</div>
</Card>
