<script lang="ts">
	import { Button, DownloadSection, Icon, icons } from '@facile/muse';
	import { CHECKSUMS_HREF, CLI_DOWNLOADS, RELEASES_HREF, cliCommands } from '$lib/downloads';
	import { t } from '$lib/i18n/index.svelte';
	import { CAPSULE_ICON } from '$lib/icons';
	import SealDemo from './SealDemo.svelte';

	const features = $derived([
		{ icon: icons.shield, title: t('landing.encrypted.title'), body: t('landing.encrypted.body') },
		{ icon: icons.eye, title: t('landing.once.title'), body: t('landing.once.body') },
		{ icon: icons.key, title: t('landing.extras.title'), body: t('landing.extras.body') }
	]);

	const commands = $derived(cliCommands(t('landing.script')));
</script>

<section class="mx-auto w-full max-w-5xl px-4 py-12 sm:px-6 sm:py-16" aria-labelledby="demo-title">
	<h2 id="demo-title" class="mb-8 text-center font-fc-title text-fc-2xl font-bold tracking-tight">
		{t('landing.demoTitle')}
	</h2>
	<SealDemo />
</section>

<section class="mx-auto w-full max-w-5xl px-4 py-12 sm:px-6 sm:py-16">
	<ul class="grid gap-4 md:grid-cols-3">
		{#each features as feature (feature.icon)}
			<li class="rounded-fc-lg bg-fc-surface p-6">
				<span class="mb-3 flex size-10 items-center justify-center rounded-fc-md bg-fc-component">
					<Icon icon={feature.icon} size={20} />
				</span>
				<h3 class="text-fc-lg font-bold">{feature.title}</h3>
				<p class="mt-1.5 text-fc-sm text-fc-fg-muted">{feature.body}</p>
			</li>
		{/each}
	</ul>
</section>

<section id="download" class="mx-auto w-full max-w-3xl scroll-mt-20 px-4 py-12 sm:px-6 sm:py-16" aria-labelledby="download-title">
	<div class="mb-6 text-center">
		<h2 id="download-title" class="font-fc-title text-fc-2xl font-bold tracking-tight">{t('landing.downloadTitle')}</h2>
		<p class="mt-2 text-fc-sm text-fc-fg-muted sm:text-fc-md">{t('landing.downloadBody')}</p>
	</div>
	<DownloadSection
		mode="device"
		allHref="/downloads"
		downloads={CLI_DOWNLOADS}
		{commands}
		headingLevel={3}
		releasesHref={RELEASES_HREF}
		checksumsHref={CHECKSUMS_HREF}
	/>
</section>

<section class="mx-auto w-full max-w-5xl px-4 py-16 text-center sm:px-6 sm:py-20">
	<h2 class="font-fc-title text-fc-2xl font-bold tracking-tight">{t('landing.ctaTitle')}</h2>
	<p class="mt-2 text-fc-sm text-fc-fg-muted">{t('landing.ctaBody')}</p>
	<div class="mt-6 flex justify-center">
		<Button href="#seal" size="lg" icon={CAPSULE_ICON}>{t('create.heading')}</Button>
	</div>
</section>
