<script lang="ts">
	import { tick } from 'svelte';
	import { Cursor, Stage, Window, cursorAt, type Act } from '@facile/mise-en-scene';
	import { t } from '$lib/i18n/index.svelte';
	import './seal-demo.css';

	const SECRET = 'PGPASSWORD=horse-battery-staple';
	const HOST = 'capsule.facile.studio';

	let draft = $state('');
	let sealed = $state(false);
	let copied = $state(false);
	let arrived = $state(false);
	let opened = $state(false);
	let pressed = $state<'seal' | 'copy' | 'open' | null>(null);
	let cursor = $state(cursorAt());
	let sealButton = $state<HTMLElement>();
	let copyButton = $state<HTMLElement>();
	let openButton = $state<HTMLElement>();

	function reset() {
		draft = '';
		sealed = false;
		copied = false;
		arrived = false;
		opened = false;
		pressed = null;
		cursor.visible = false;
	}

	function poster() {
		draft = SECRET;
		sealed = true;
		copied = true;
		arrived = true;
		opened = true;
	}

	async function press(act: Act, target: HTMLElement | undefined, name: 'seal' | 'copy' | 'open') {
		await act.move(cursor, act.point(target, 0.8, 0.7));
		await act.click(cursor, () => (pressed = name));
		pressed = null;
	}

	async function script(act: Act) {
		await act.wait(700);
		await act.type(SECRET, (v) => (draft = v), 55);
		await act.wait(300);
		cursor.visible = true;
		await press(act, sealButton, 'seal');
		sealed = true;
		await tick();
		await act.wait(900);
		await press(act, copyButton, 'copy');
		copied = true;
		await act.wait(500);
		arrived = true;
		await tick();
		await act.wait(1000);
		await press(act, openButton, 'open');
		opened = true;
		await act.move(cursor, act.point(openButton, 3.2, 3.4), 500);
		await act.wait(2000);
		cursor.visible = false;
		arrived = false;
		sealed = false;
		draft = '';
		await act.wait(800);
	}
</script>

<Stage
	width={960}
	height={520}
	compact={{ below: 640, width: 358, height: 520 }}
	label={t('landing.demoLabel')}
	{script}
	{reset}
	{poster}
>
	{#snippet children({ compact })}
		<div class="sd" class:sd-compact={compact}>
			<Window title={HOST} class="sd-seal {arrived ? '' : 'sd-solo'}">
				<div class="sd-body sd-stack">
					<div class="sd-view" class:sd-off={sealed}>
						<p class="sd-heading">{t('create.heading')}</p>
						<div class="sd-area" class:sd-placeholder={!draft}>{draft || t('create.placeholder')}</div>
						<div class="sd-selects">
							<span class="sd-select">{t('expiry.24h')}</span>
							<span class="sd-select">plaintext</span>
						</div>
						<span class="sd-check"><span class="sd-box"></span>{t('create.burn')}</span>
						<span class="sd-button" class:sd-pressed={pressed === 'seal'} bind:this={sealButton}>{t('create.seal')}</span>
					</div>
					<div class="sd-view" class:sd-off={!sealed}>
						<p class="sd-heading">{t('sealed.heading')}</p>
						<p class="sd-muted">{t('sealed.burnOnce')}</p>
						<div class="sd-link">
							<span class="sd-url"><span class="sd-dim">{HOST}/cap_7Hq2</span><span class="sd-key">#k3Y9pQ</span></span>
							<span class="sd-button sd-small" class:sd-pressed={pressed === 'copy'} bind:this={copyButton}>
								{copied ? t('sealed.copied') : t('sealed.copy')}
							</span>
						</div>
					</div>
				</div>
			</Window>
			<Window title={HOST} class="sd-reveal {arrived ? '' : 'sd-away'}">
				<div class="sd-body sd-stack">
					<div class="sd-view" class:sd-off={opened}>
						<p class="sd-heading">{t('reveal.readyHeading')}</p>
						<p class="sd-muted">{t('reveal.readyBurn')}</p>
						<span class="sd-button" class:sd-pressed={pressed === 'open'} bind:this={openButton}>{t('reveal.open')}</span>
					</div>
					<div class="sd-view" class:sd-off={!opened}>
						<p class="sd-heading">{t('reveal.openedHeading')}</p>
						<p class="sd-danger">{t('reveal.destroyed')}</p>
						<pre class="sd-secret">{SECRET}</pre>
					</div>
				</div>
			</Window>
		</div>
		<Cursor {cursor} />
	{/snippet}
</Stage>
