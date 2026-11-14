<script lang="ts">
	import { Badge } from '@facile/muse';
	import { Stage, type Act } from '@facile/mise-en-scene';
	import { t } from '$lib/i18n/index.svelte';
	import './tiles.css';

	const rows = $derived([
		{ initial: 'L', who: t('landing.tile.share.lea'), what: t('landing.tile.share.leaWhat'), link: 'capsule.facile.studio/cap_7Hq2#k3Y9' },
		{ initial: 'W', who: t('landing.tile.share.guest'), what: t('landing.tile.share.guestWhat'), link: 'capsule.facile.studio/cap_Vb81#p2Lm' }
	]);

	let shown = $state(1);
	let sent = $state(0);

	function reset() {
		shown = 1;
		sent = 0;
	}

	function poster() {
		shown = 2;
		sent = 2;
	}

	async function script(act: Act) {
		await act.wait(900);
		sent = 1;
		await act.wait(1200);
		shown = 2;
		await act.wait(900);
		sent = 2;
		await act.wait(2400);
		reset();
		await act.wait(600);
	}
</script>

<Stage width={400} height={260} label={t('landing.tile.share.label')} {script} {reset} {poster}>
	{#snippet children()}
		<div class="td-panel">
			{#each rows as row, i (row.initial)}
				<div class="td-row" class:td-on={i < shown}>
					<span class="td-avatar">{row.initial}</span>
					<div class="min-w-0 flex-1">
						<div class="flex items-center justify-between gap-2">
							<span class="td-name truncate">{row.who}</span>
							<Badge tone="success" class="shrink-0 whitespace-nowrap transition-opacity duration-300 {i < sent ? 'opacity-100' : 'opacity-0'}">
								{t('sealed.copied')}
							</Badge>
						</div>
						<span class="td-muted block truncate">{row.what}</span>
						<span class="td-link">{row.link}</span>
					</div>
				</div>
			{/each}
		</div>
	{/snippet}
</Stage>
