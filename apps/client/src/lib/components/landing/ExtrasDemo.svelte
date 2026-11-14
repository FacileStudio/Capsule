<script lang="ts">
	import { Icon, icons } from '@facile/muse';
	import { Stage, type Act } from '@facile/mise-en-scene';
	import { formatBytes } from '$lib/envelope';
	import { t } from '$lib/i18n/index.svelte';
	import './tiles.css';

	const DOTS = '••••••••••';

	let ticked = $state(false);
	let dots = $state('');
	let attached = $state(false);

	function reset() {
		ticked = false;
		dots = '';
		attached = false;
	}

	function poster() {
		ticked = true;
		dots = DOTS;
		attached = true;
	}

	async function script(act: Act) {
		await act.wait(800);
		ticked = true;
		await act.wait(500);
		await act.type(DOTS, (v) => (dots = v), 90);
		await act.wait(700);
		attached = true;
		await act.wait(2600);
		reset();
		await act.wait(600);
	}
</script>

<Stage width={400} height={260} label={t('landing.tile.extras.label')} {script} {reset} {poster}>
	{#snippet children()}
		<div class="td-panel">
			<span class="td-check">
				<span class="td-box" class:td-ticked={ticked}>
					{#if ticked}✓{/if}
				</span>
				{t('create.password')}
			</span>
			<div class="td-field" class:td-muted={!dots}>{dots || ' '}</div>
			<div class="td-file" class:td-on={attached}>
				<span class="td-avatar"><Icon icon={icons.folder} size={16} /></span>
				<div class="min-w-0 flex-1">
					<span class="td-name block truncate">contrat.pdf</span>
					<span class="td-muted">{formatBytes(245760)}</span>
				</div>
			</div>
		</div>
	{/snippet}
</Stage>
