<script lang="ts">
	import { Alert, Button, Card, Checkbox, Field, Icon, Input, SecretField, Select, Spinner, Textarea, icons } from '@facile/muse';
	import { generateKey, exportKey, encrypt, wrapContentKey } from '$lib/crypto';
	import { encodeFile, serializePayload } from '$lib/envelope';
	import FileAttachment from '$lib/components/FileAttachment.svelte';
	import { backend } from '$lib/backend';
	import { page } from '$app/state';
	import { t } from '$lib/i18n/index.svelte';
	import { CAPSULE_ICON } from '$lib/icons';

	type State = 'idle' | 'sealing' | 'sealed';

	let phase: State = $state('idle');
	let content = $state('');
	let attachedFile: File | null = $state(null);
	let burnAfterRead = $state(true);
	let expiresIn = $state('24h');
	let syntax = $state('plaintext');
	let usePassword = $state(false);
	let password = $state('');
	let error = $state('');
	let capsuleUrl = $state('');
	let capsuleId = $state('');
	let deleteToken = $state('');
	let copied = $state(false);
	let revoking = $state(false);
	let revoked = $state(false);

	const expiryValues = ['1h', '24h', '7d', '30d'] as const;
	const expiryOptions = $derived(expiryValues.map((value) => ({ value, label: t(`expiry.${value}`) })));
	const duration = $derived(expiryOptions.find((o) => o.value === expiresIn)?.label ?? '');
	const sealedHelp = $derived(
		burnAfterRead
			? t(usePassword ? 'sealed.burnOncePassword' : 'sealed.burnOnce')
			: t(usePassword ? 'sealed.expiresPassword' : 'sealed.expires', { duration })
	);

	const syntaxOptions = [
		'plaintext', 'javascript', 'typescript', 'python', 'go', 'rust',
		'json', 'yaml', 'sql', 'bash', 'html', 'css', 'markdown',
	];

	async function seal() {
		if (!content.trim() && !attachedFile) return;
		if (usePassword && !password) return;
		error = '';
		phase = 'sealing';
		try {
			const encoded = attachedFile ? await encodeFile(attachedFile) : null;
			const payload = serializePayload(content, encoded);
			const key = await generateKey();
			const ciphertext = await encrypt(payload, key);
			const hasPassword = usePassword && password.length > 0;
			const fragment = hasPassword
				? (await wrapContentKey(key, password)).fragment
				: await exportKey(key);
			const result = await backend.createPaste({
				content: ciphertext,
				burn_after_read: burnAfterRead,
				expires_in: expiresIn,
				has_password: hasPassword,
				syntax: syntax !== 'plaintext' ? syntax : undefined,
			});
			capsuleUrl = `${page.url.origin}/${result.id}#${fragment}`;
			capsuleId = result.id;
			deleteToken = result.delete_token;
			phase = 'sealed';
		} catch (err) {
			error = err instanceof Error ? err.message : t('create.sealFailed');
			phase = 'idle';
		}
	}

	async function copyLink() {
		await navigator.clipboard.writeText(capsuleUrl);
		copied = true;
		setTimeout(() => (copied = false), 2000);
	}

	async function revoke() {
		if (revoking || revoked) return;
		revoking = true;
		try {
			await backend.deletePaste(capsuleId, deleteToken);
			revoked = true;
		} catch (err) {
			error = err instanceof Error ? err.message : t('revoke.failed');
		} finally {
			revoking = false;
		}
	}

	function reset() {
		phase = 'idle';
		content = '';
		attachedFile = null;
		usePassword = false;
		password = '';
		capsuleUrl = '';
		capsuleId = '';
		deleteToken = '';
		error = '';
		copied = false;
		revoking = false;
		revoked = false;
	}
</script>

<svelte:head>
	<title>{t('create.title')}</title>
	<meta name="description" content={t('create.description')} />
</svelte:head>

{#if phase === 'idle' || phase === 'sealing'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">{t('create.heading')}</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">{t('create.intro')}</p>
		</div>

		<Textarea
			bind:value={content}
			placeholder={t('create.placeholder')}
			rows={8}
			disabled={phase === 'sealing'}
			class="font-fc-mono text-fc-sm"
		/>

		<FileAttachment bind:file={attachedFile} disabled={phase === 'sealing'} onError={(msg) => (error = msg)} />

		<div class="grid gap-4 sm:grid-cols-2">
			<Field label={t('create.expires')}>
				<Select bind:value={expiresIn} disabled={phase === 'sealing'}>
					{#each expiryOptions as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</Select>
			</Field>
			<Field label={t('create.syntax')}>
				<Select bind:value={syntax} disabled={phase === 'sealing'}>
					{#each syntaxOptions as s}
						<option value={s}>{s}</option>
					{/each}
				</Select>
			</Field>
		</div>

		<div class="flex flex-col gap-3">
			<!-- The 16px box is by design; the <label> Checkbox renders carries the hit target. -->
			<Checkbox label={t('create.burn')} bind:checked={burnAfterRead} disabled={phase === 'sealing'} class="min-h-11" />
			<Checkbox label={t('create.password')} bind:checked={usePassword} disabled={phase === 'sealing'} class="min-h-11" />
		</div>

		{#if usePassword}
			<Input type="password" bind:value={password} placeholder={t('create.passwordPlaceholder')} disabled={phase === 'sealing'} />
		{/if}

		{#if error}
			<Alert tone="danger">{error}</Alert>
		{/if}

		<div class="flex flex-col gap-3 sm:flex-row">
			<Button
				size="lg"
				icon={phase === 'sealing' ? undefined : CAPSULE_ICON}
				onclick={seal}
				disabled={(!content.trim() && !attachedFile) || (usePassword && !password) || phase === 'sealing'}
				class="flex-1"
			>
				{#if phase === 'sealing'}
					<Spinner size="sm" label={t('create.sealingLabel')} />
					{t('create.sealing')}
				{:else}
					{t('create.seal')}
				{/if}
			</Button>
			<Button href="/revoke" variant="secondary" size="lg" icon={icons.revoke}>{t('shell.revoke')}</Button>
		</div>
	</div>
{:else if phase === 'sealed'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<Icon icon={icons.check} size={28} class="mb-3 text-fc-success" />
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">{t('sealed.heading')}</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">{sealedHelp}</p>
		</div>

		<div class="flex flex-col gap-2 sm:flex-row">
			<Input readonly value={capsuleUrl} class="font-fc-mono text-fc-sm sm:flex-1" />
			<Button size="lg" icon={copied ? icons.check : icons.copy} onclick={copyLink}>
				{copied ? t('sealed.copied') : t('sealed.copy')}
			</Button>
		</div>

		<Card class="flex flex-col gap-4">
			{#if revoked}
				<p class="flex items-center gap-2 text-fc-sm text-fc-danger">
					<Icon icon={icons.revoke} />
					{t('sealed.revoked')}
				</p>
			{:else}
				<!-- The token is shown once and is the only way to burn this capsule early, so it
				     gets a copy button. It used to be `truncate`d text: unreadable and unselectable. -->
				<SecretField bind:value={deleteToken} label={t('sealed.token')} helper={t('sealed.tokenHelp')} />
				<Button variant="danger" icon={icons.remove} onclick={revoke} disabled={revoking} class="w-fit">
					{revoking ? t('sealed.revoking') : t('sealed.revoke')}
				</Button>
			{/if}
		</Card>

		{#if error}
			<Alert tone="danger">{error}</Alert>
		{/if}

		<Button variant="secondary" size="lg" icon={CAPSULE_ICON} onclick={reset} class="w-fit">
			{t('sealed.another')}
		</Button>
	</div>
{/if}
