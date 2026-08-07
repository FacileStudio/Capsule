<script lang="ts">
	import { Alert, Button, Card, Checkbox, Field, Input, SecretField, Select, Spinner, Textarea, icons } from '@facile/muse';
	import { generateKey, exportKey, encrypt, wrapContentKey, buildPasswordFragment } from '$lib/crypto';
	import { backend } from '$lib/backend';
	import { page } from '$app/state';

	type State = 'idle' | 'sealing' | 'sealed';

	let phase: State = $state('idle');
	let content = $state('');
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

	const expiryOptions = [
		{ value: '1h', label: '1 hour' },
		{ value: '24h', label: '24 hours' },
		{ value: '7d', label: '7 days' },
		{ value: '30d', label: '30 days' },
	];

	const syntaxOptions = [
		'plaintext', 'javascript', 'typescript', 'python', 'go', 'rust',
		'json', 'yaml', 'sql', 'bash', 'html', 'css', 'markdown',
	];

	async function seal() {
		if (!content.trim()) return;
		if (usePassword && !password) return;
		error = '';
		phase = 'sealing';

		try {
			const key = await generateKey();
			const ciphertext = await encrypt(content, key);

			let fragment: string;
			const hasPassword = usePassword && password.length > 0;

			if (hasPassword) {
				const wrapped = await wrapContentKey(key, password);
				fragment = buildPasswordFragment(wrapped.encryptedKey, wrapped.salt, wrapped.iv);
			} else {
				fragment = await exportKey(key);
			}

			const result = await backend.createPaste({
				content: ciphertext,
				burn_after_read: burnAfterRead,
				expires_in: expiresIn,
				has_password: hasPassword,
				syntax: syntax !== 'plaintext' ? syntax : undefined,
			});

			const origin = page.url.origin;
			capsuleUrl = `${origin}/${result.id}#${fragment}`;
			capsuleId = result.id;
			deleteToken = result.delete_token;
			phase = 'sealed';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to seal capsule';
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
			error = err instanceof Error ? err.message : 'Failed to revoke capsule';
		} finally {
			revoking = false;
		}
	}

	function reset() {
		phase = 'idle';
		content = '';
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
	<title>Capsule — Ephemeral Encrypted Sharing</title>
	<meta name="description" content="Share secrets that self-destruct. End-to-end encrypted, zero knowledge." />
</svelte:head>

{#if phase === 'idle' || phase === 'sealing'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">Seal a capsule</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">
				Your content is encrypted in the browser. The server never sees the plaintext.
			</p>
		</div>

		<Textarea
			bind:value={content}
			placeholder="Paste your secret here..."
			rows={8}
			disabled={phase === 'sealing'}
			class="font-fc-mono text-fc-sm"
		/>

		<div class="grid gap-4 sm:grid-cols-2">
			<Field label="Expires in">
				<Select bind:value={expiresIn} disabled={phase === 'sealing'}>
					{#each expiryOptions as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</Select>
			</Field>

			<Field label="Syntax">
				<Select bind:value={syntax} disabled={phase === 'sealing'}>
					{#each syntaxOptions as s}
						<option value={s}>{s}</option>
					{/each}
				</Select>
			</Field>
		</div>

		<div class="flex flex-col gap-3">
			<!-- The 16px box is by design; the <label> Checkbox renders carries the hit target. -->
			<Checkbox
				label="Burn after opening"
				bind:checked={burnAfterRead}
				disabled={phase === 'sealing'}
				class="min-h-11"
			/>
			<Checkbox
				label="Password protect"
				bind:checked={usePassword}
				disabled={phase === 'sealing'}
				class="min-h-11"
			/>
		</div>

		{#if usePassword}
			<Input
				type="password"
				bind:value={password}
				placeholder="Enter a password..."
				disabled={phase === 'sealing'}
			/>
		{/if}

		{#if error}
			<Alert tone="danger">{error}</Alert>
		{/if}

		<div class="flex flex-col gap-3 sm:flex-row">
			<Button
				size="lg"
				icon={phase === 'sealing' ? undefined : 'solar:pill-bold-duotone'}
				onclick={seal}
				disabled={!content.trim() || (usePassword && !password) || phase === 'sealing'}
				class="flex-1"
			>
				{#if phase === 'sealing'}
					<Spinner size="sm" label="Sealing" />
					Sealing...
				{:else}
					Seal capsule
				{/if}
			</Button>
			<Button href="/revoke" variant="outline" size="lg" icon={icons.revoke}>Revoke</Button>
		</div>
	</div>
{:else if phase === 'sealed'}
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-1">
			<iconify-icon
				icon={icons.check}
				width="28"
				height="28"
				class="mb-3 block text-fc-success"
			></iconify-icon>
			<h1 class="text-fc-2xl font-semibold tracking-tight sm:text-fc-3xl">Capsule sealed</h1>
			<p class="text-fc-sm text-fc-fg-muted sm:text-fc-md">
				Share this link{usePassword ? ' and the password' : ''} — {burnAfterRead ? 'the capsule can only be opened once.' : `it expires in ${expiryOptions.find(o => o.value === expiresIn)?.label}.`}
			</p>
		</div>

		<div class="flex flex-col gap-2 sm:flex-row">
			<Input readonly value={capsuleUrl} class="font-fc-mono text-fc-sm sm:flex-1" />
			<Button size="lg" icon={copied ? icons.check : icons.copy} onclick={copyLink}>
				{copied ? 'Copied' : 'Copy link'}
			</Button>
		</div>

		<Card class="flex flex-col gap-4">
			{#if revoked}
				<p class="flex items-center gap-2 text-fc-sm text-fc-danger">
					<iconify-icon icon={icons.revoke} width="16" height="16" class="block"></iconify-icon>
					Capsule revoked
				</p>
			{:else}
				<!-- The token is shown once and is the only way to burn this capsule early, so it
				     gets a copy button. It used to be `truncate`d text: unreadable and unselectable. -->
				<SecretField
					bind:value={deleteToken}
					label="Delete token"
					helper="Copy it now — it is not stored anywhere you can read it back."
				/>
				<Button
					variant="danger"
					icon={icons.remove}
					onclick={revoke}
					disabled={revoking}
					class="w-fit"
				>
					{revoking ? 'Revoking...' : 'Revoke'}
				</Button>
			{/if}
		</Card>

		{#if error}
			<Alert tone="danger">{error}</Alert>
		{/if}

		<Button
			variant="outline"
			size="lg"
			icon="solar:pill-bold-duotone"
			onclick={reset}
			class="w-fit"
		>
			Seal another capsule
		</Button>
	</div>
{/if}
