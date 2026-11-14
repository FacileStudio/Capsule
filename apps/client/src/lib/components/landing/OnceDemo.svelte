<script lang="ts">
	import { Stage, type Act } from '@facile/mise-en-scene';
	import { t } from '$lib/i18n/index.svelte';
	import './tiles.css';

	type View = 'ready' | 'opened' | 'empty';

	let view = $state<View>('ready');

	function reset() {
		view = 'ready';
	}

	function poster() {
		view = 'empty';
	}

	async function script(act: Act) {
		await act.wait(1200);
		view = 'opened';
		await act.wait(2000);
		view = 'empty';
		await act.wait(2400);
		reset();
		await act.wait(600);
	}
</script>

<Stage width={400} height={260} label={t('landing.tile.once.label')} {script} {reset} {poster}>
	{#snippet children()}
		<div class="td-panel">
			<div class="td-stack">
				<div class="td-view" class:td-off={view !== 'ready'}>
					<p class="td-heading">{t('reveal.readyHeading')}</p>
					<p class="td-muted m-0">{t('reveal.readyBurn')}</p>
					<span class="td-button">{t('reveal.open')}</span>
				</div>
				<div class="td-view" class:td-off={view !== 'opened'}>
					<p class="td-heading">{t('reveal.openedHeading')}</p>
					<p class="td-danger">{t('reveal.destroyed')}</p>
					<pre class="td-secret">{t('landing.demoSecret')}</pre>
				</div>
				<div class="td-view" class:td-off={view !== 'empty'}>
					<p class="td-heading">{t('reveal.emptyHeading')}</p>
					<p class="td-muted m-0">{t('reveal.emptyText')}</p>
					<span class="td-button td-quiet">{t('reveal.sealNew')}</span>
				</div>
			</div>
		</div>
	{/snippet}
</Stage>
