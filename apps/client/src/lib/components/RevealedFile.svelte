<script lang="ts">
	import { onDestroy } from 'svelte';
	import { Button, Card, Icon, icons } from '@facile/muse';
	import { fileToBlob, formatBytes, type CapsuleFile } from '$lib/envelope';
	import { t } from '$lib/i18n/index.svelte';

	let { file }: { file: CapsuleFile } = $props();

	let objectUrl = $state('');

	$effect(() => {
		const blob = fileToBlob(file);
		const url = URL.createObjectURL(blob);
		objectUrl = url;
		return () => URL.revokeObjectURL(url);
	});

	onDestroy(() => {
		if (objectUrl) URL.revokeObjectURL(objectUrl);
	});

	const isImage = $derived(file.type.startsWith('image/'));
	const isAudio = $derived(file.type.startsWith('audio/'));
	const isVideo = $derived(file.type.startsWith('video/'));

	function download() {
		if (!objectUrl) return;
		const a = document.createElement('a');
		a.href = objectUrl;
		a.download = file.name || 'capsule-file';
		document.body.appendChild(a);
		a.click();
		a.remove();
	}
</script>

<Card class="flex flex-col gap-4">
	{#if objectUrl && isImage}
		<div class="flex justify-center overflow-hidden rounded-fc-md bg-fc-component p-2">
			<img src={objectUrl} alt={file.name} class="max-h-96 max-w-full rounded-fc-sm object-contain" />
		</div>
	{:else if objectUrl && isAudio}
		<audio controls src={objectUrl} class="w-full"></audio>
	{:else if objectUrl && isVideo}
		<video controls src={objectUrl} class="max-h-96 w-full rounded-fc-md bg-fc-component">
			<track kind="captions" />
		</video>
	{/if}

	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex min-w-0 items-center gap-3">
			<span class="flex size-10 shrink-0 items-center justify-center rounded-fc-sm bg-fc-component text-fc-fg-muted">
				<Icon icon={icons.folder} size={20} />
			</span>
			<div class="min-w-0">
				<p class="truncate text-fc-sm font-medium text-fc-fg">{file.name}</p>
				<p class="text-fc-xs text-fc-fg-muted">{formatBytes(file.size)}</p>
			</div>
		</div>

		<Button size="sm" icon={icons.download} onclick={download}>
			{t('reveal.download')}
		</Button>
	</div>
</Card>
