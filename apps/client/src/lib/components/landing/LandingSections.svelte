<script lang="ts">
	import {
		Button,
		FeatureSplit,
		ForDevelopers,
		LANDING_HEADING,
		PricingPlans,
		SignUpWithFacile,
		WorksWith,
		icons,
		type DeveloperRow,
		type FeatureItem,
		type WorksWithApp
	} from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';
	import { BULLE_ICON, CAPSULE_ICON, PLACE_ICON, PRISME_ICON } from '$lib/icons';
	import DemoTiles from './DemoTiles.svelte';
	import SealDemo from './SealDemo.svelte';

	const block = 'mx-auto w-full max-w-5xl px-4 sm:px-6';

	const features: FeatureItem[] = $derived([
		{ icon: icons.shield, title: t('blocks.features.device.title'), text: t('blocks.features.device.text') },
		{ icon: icons.clock, title: t('blocks.features.expiry.title'), text: t('blocks.features.expiry.text') },
		{ icon: icons.revoke, title: t('blocks.features.revoke.title'), text: t('blocks.features.revoke.text') },
		{ icon: icons.userCircle, title: t('blocks.features.account.title'), text: t('blocks.features.account.text') }
	]);

	const apps: WorksWithApp[] = $derived([
		{ name: 'Prisme', icon: PRISME_ICON, line: t('blocks.with.prisme'), href: 'https://prisme.facile.studio' },
		{ name: 'Place', icon: PLACE_ICON, line: t('blocks.with.place'), href: 'https://place.facile.studio' },
		{ name: 'Bulle', icon: BULLE_ICON, line: t('blocks.with.bulle'), href: 'https://bulle.facile.studio' }
	]);

	const rows: DeveloperRow[] = $derived([
		{ kind: 'cli', text: t('blocks.dev.cli'), href: '/downloads', command: 'brew install facilestudio/tap/capsule' },
		{ kind: 'api', text: t('blocks.dev.api'), href: '/docs' },
		{ kind: 'mcp', text: t('blocks.dev.mcp') },
		{ kind: 'source', text: t('blocks.dev.source'), href: 'https://github.com/FacileStudio/Capsule' }
	]);

	const notes = $derived({ solo: t('blocks.pricing.free'), '*': t('blocks.pricing.paid') });
</script>

<div class="flex flex-col gap-20 pt-12 pb-20 sm:gap-28 sm:pt-16 sm:pb-28">
	<section class={block} aria-labelledby="demo-title">
		<h2 id="demo-title" class="mb-8 text-center {LANDING_HEADING}">{t('landing.demoTitle')}</h2>
		<SealDemo />
	</section>

	<DemoTiles />

	<FeatureSplit class={block} title={t('blocks.features.title')} description={t('blocks.features.line')} {features} />

	<WorksWith class={block} {apps} description={t('blocks.with.line')} />

	<PricingPlans
		id="pricing"
		class={block}
		title={t('blocks.pricing.title')}
		headingClass={LANDING_HEADING}
		highlight={['ai.access']}
		highlightOnly
		{notes}
		recommended="pro"
		freeLabel={t('blocks.cta.seal')}
		freeHref="/#seal"
		suiteHref="https://suite.facile.studio/#tarifs"
	/>

	<div class={block}>
		<ForDevelopers {rows} selfHostHref={t('blocks.dev.selfHostHref')} description={t('blocks.dev.line')} />
	</div>

	<section class="{block} flex flex-col items-center gap-4 text-center">
		<h2 class={LANDING_HEADING}>{t('blocks.cta.title')}</h2>
		<p class="max-w-xl text-fc-sm text-fc-fg-muted sm:text-fc-md">{t('blocks.cta.body')}</p>
		<div class="mt-2 flex w-full flex-col gap-3 sm:w-auto sm:flex-row">
			<Button href="#seal" size="lg" icon={CAPSULE_ICON}>{t('blocks.cta.seal')}</Button>
			<SignUpWithFacile size="lg" class="bg-fc-component text-fc-fg" />
		</div>
	</section>
</div>
