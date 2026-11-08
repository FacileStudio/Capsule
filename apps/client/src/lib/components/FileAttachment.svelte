<script lang="ts">
	import { Dropzone, Icon, IconButton, icons } from '@facile/muse';
	import { formatBytes, MAX_FILE_SIZE } from '$lib/envelope';
	import { t } from '$lib/i18n/index.svelte';

	let {
		file = $bindable<File | null>(null),
		disabled = false,
		onError,
	}: {
		file?: File | null;
		disabled?: boolean;
		onError?: (message: string) => void;
	} = $props();

	function handleFiles(files: File[]) {
		const picked = files[0];
		if (!picked) return;
		file = picked;
		onError?.('');
	}

	function handleReject() {
		onError?.(t('create.fileTooLarge'));
	}

	function removeFile() {
		file = null;
	}
</script>

{#if file}
	<div class="flex items-center justify-between gap-3 rounded-fc-md bg-fc-surface px-4 py-3">
		<div class="flex min-w-0 items-center gap-3">
			<span class="flex size-9 shrink-0 items-center justify-center rounded-fc-sm bg-fc-component text-fc-fg-muted">
				<Icon icon={icons.folder} size={18} />
			</span>
			<div class="min-w-0">
				<p class="truncate text-fc-sm font-medium text-fc-fg">{file.name}</p>
				<p class="text-fc-xs text-fc-fg-muted">{formatBytes(file.size)}</p>
			</div>
		</div>
		<IconButton
			type="button"
			aria-label={t('create.fileRemove')}
			onclick={removeFile}
			{disabled}
			class="text-fc-fg-muted hover:text-fc-danger"
		>
			<Icon icon={icons.remove} size={18} />
		</IconButton>
	</div>
{:else}
	<Dropzone
		multiple={false}
		maxSize={MAX_FILE_SIZE}
		{disabled}
		label={t('create.fileDrop')}
		hint={t('create.fileHint')}
		onFiles={handleFiles}
		onReject={handleReject}
	/>
{/if}
