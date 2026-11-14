<script lang="ts">
	import { DownloadSection } from '@facile/muse';
	import { CHECKSUMS_HREF, CLI_DOWNLOADS, RELEASES_HREF, cliCommands } from '$lib/downloads';
	import { t } from '$lib/i18n/index.svelte';

	const commands = $derived(cliCommands(t('landing.script')));
</script>

<svelte:head>
	<title>{t('downloads.title')}</title>
</svelte:head>

<section class="mx-auto w-full max-w-4xl px-4 py-12 sm:px-6 sm:py-16" aria-labelledby="downloads-title">
	<div class="mb-8 text-center">
		<h1 id="downloads-title" class="font-fc-title text-fc-2xl font-bold tracking-tight sm:text-fc-3xl">{t('downloads.heading')}</h1>
		<p class="mt-2 text-fc-sm text-fc-fg-muted sm:text-fc-md">{t('landing.downloadBody')}</p>
	</div>
	<DownloadSection downloads={CLI_DOWNLOADS} {commands} releasesHref={RELEASES_HREF} checksumsHref={CHECKSUMS_HREF} />
</section>
