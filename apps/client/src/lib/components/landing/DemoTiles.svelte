<script lang="ts">
	import { LANDING_HEADING } from '@facile/muse';
	import { Tile } from '@facile/mise-en-scene';
	import { t } from '$lib/i18n/index.svelte';
	import ExtrasDemo from './ExtrasDemo.svelte';
	import OnceDemo from './OnceDemo.svelte';
	import ShareDemo from './ShareDemo.svelte';
</script>

<section id="features" class="mx-auto w-full max-w-5xl scroll-mt-20 px-4 sm:px-6">
	<h2 class={LANDING_HEADING}>{t('landing.tilesTitle')}</h2>
	<div class="mt-8 grid gap-4 lg:grid-cols-3">
		<Tile title={t('landing.tile.share.title')} body={t('landing.tile.share.body')}>
			<div class="mx-auto w-full max-w-md"><ShareDemo /></div>
		</Tile>
		<Tile title={t('landing.tile.once.title')} body={t('landing.tile.once.body')}>
			<div class="mx-auto w-full max-w-md"><OnceDemo /></div>
		</Tile>
		<Tile title={t('landing.tile.extras.title')} body={t('landing.tile.extras.body')}>
			<div class="mx-auto w-full max-w-md"><ExtrasDemo /></div>
		</Tile>
	</div>
</section>
