<script lang="ts">
	import { generateKey, exportKey, encrypt, wrapContentKey, buildPasswordFragment } from '$lib/crypto';
	import { backend } from '$lib/backend';
	import { page } from '$app/state';

	type State = 'idle' | 'sealing' | 'sealed';

	let state: State = $state('idle');
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
		state = 'sealing';

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
			state = 'sealed';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to seal capsule';
			state = 'idle';
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
		state = 'idle';
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

{#if state === 'idle' || state === 'sealing'}
	<div class="flex flex-col gap-5 sm:gap-6">
		<div>
			<h1 class="text-2xl font-bold font-heading tracking-tight sm:text-3xl">Seal a capsule</h1>
			<p class="mt-2 text-sm text-muted-foreground sm:text-base">
				Your content is encrypted in the browser. The server never sees the plaintext.
			</p>
		</div>

		<textarea
			bind:value={content}
			placeholder="Paste your secret here..."
			rows="8"
			disabled={state === 'sealing'}
			class="w-full resize-y rounded-lg border border-input bg-card p-3 font-mono text-sm text-card-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring sm:p-4"
		></textarea>

		<div class="flex flex-col gap-4 sm:flex-row sm:flex-wrap sm:items-center sm:gap-6">
			<label class="flex min-h-[44px] items-center gap-2 text-sm">
				<input
					type="checkbox"
					bind:checked={burnAfterRead}
					disabled={state === 'sealing'}
					class="h-5 w-5 rounded border-input accent-primary sm:h-4 sm:w-4"
				/>
				Burn after opening
			</label>

			<label class="flex min-h-[44px] items-center gap-2 text-sm">
				<span class="text-muted-foreground">Expires in</span>
				<select
					bind:value={expiresIn}
					disabled={state === 'sealing'}
					class="min-h-[44px] rounded-md border border-input bg-card px-3 py-2 text-sm sm:min-h-0 sm:px-2 sm:py-1"
				>
					{#each expiryOptions as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</select>
			</label>

			<label class="flex min-h-[44px] items-center gap-2 text-sm">
				<span class="text-muted-foreground">Syntax</span>
				<select
					bind:value={syntax}
					disabled={state === 'sealing'}
					class="min-h-[44px] rounded-md border border-input bg-card px-3 py-2 text-sm sm:min-h-0 sm:px-2 sm:py-1"
				>
					{#each syntaxOptions as s}
						<option value={s}>{s}</option>
					{/each}
				</select>
			</label>

			<label class="flex min-h-[44px] items-center gap-2 text-sm">
				<input
					type="checkbox"
					bind:checked={usePassword}
					disabled={state === 'sealing'}
					class="h-5 w-5 rounded border-input accent-primary sm:h-4 sm:w-4"
				/>
				Password protect
			</label>
		</div>

		{#if usePassword}
			<input
				type="password"
				bind:value={password}
				placeholder="Enter a password..."
				disabled={state === 'sealing'}
				class="w-full rounded-lg border border-input bg-card px-4 py-3 text-sm text-card-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
			/>
		{/if}

		{#if error}
			<p class="text-sm text-red-400">{error}</p>
		{/if}

		<div class="flex flex-col gap-3 sm:flex-row">
			<button
				onclick={seal}
				disabled={!content.trim() || (usePassword && !password) || state === 'sealing'}
				class="inline-flex h-12 flex-1 items-center justify-center rounded-md bg-primary px-6 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed sm:h-11"
			>
				{#if state === 'sealing'}
					<iconify-icon icon="solar:refresh-bold" width="16" class="mr-2 animate-spin"></iconify-icon>
					Sealing...
				{:else}
					<iconify-icon icon="solar:pill-bold-duotone" width="16" class="mr-2"></iconify-icon>
					Seal capsule
				{/if}
			</button>
			<a
				href="/revoke"
				class="inline-flex h-12 items-center justify-center rounded-md border border-border bg-background px-5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground sm:h-11"
			>
				<iconify-icon icon="solar:fire-bold" width="16" class="mr-2"></iconify-icon>
				Revoke
			</a>
		</div>
	</div>
{:else if state === 'sealed'}
	<div class="flex flex-col gap-5 sm:gap-6">
		<div>
			<div class="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-green-500/10">
				<iconify-icon icon="solar:check-circle-bold" width="28" class="text-green-400"></iconify-icon>
			</div>
			<h1 class="text-2xl font-bold font-heading tracking-tight sm:text-3xl">Capsule sealed</h1>
			<p class="mt-2 text-sm text-muted-foreground sm:text-base">
				Share this link{usePassword ? ' and the password' : ''} — {burnAfterRead ? 'the capsule can only be opened once.' : `it expires in ${expiryOptions.find(o => o.value === expiresIn)?.label}.`}
			</p>
		</div>

		<div class="flex flex-col gap-2 sm:flex-row sm:items-stretch">
			<input
				type="text"
				readonly
				value={capsuleUrl}
				class="w-full rounded-lg border border-input bg-card px-3 py-3 font-mono text-sm text-card-foreground sm:flex-1 sm:px-4"
			/>
			<button
				onclick={copyLink}
				class="inline-flex h-12 items-center justify-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 sm:h-auto"
			>
				{#if copied}
					<iconify-icon icon="solar:check-circle-bold" width="16" class="mr-2"></iconify-icon>
					Copied
				{:else}
					<iconify-icon icon="solar:copy-bold" width="16" class="mr-2"></iconify-icon>
					Copy link
				{/if}
			</button>
		</div>

		<div class="rounded-lg border border-border bg-card p-4">
			{#if revoked}
				<div class="flex items-center gap-2 text-sm text-red-400">
					<iconify-icon icon="solar:fire-bold" width="16"></iconify-icon>
					Capsule revoked
				</div>
			{:else}
				<div class="flex items-center justify-between gap-3">
					<div class="min-w-0 flex-1">
						<p class="text-xs text-muted-foreground">Delete token</p>
						<code class="mt-1 block truncate text-xs text-muted-foreground/70">{deleteToken}</code>
					</div>
					<button
						onclick={revoke}
						disabled={revoking}
						class="inline-flex h-9 shrink-0 items-center justify-center rounded-md border border-red-500/30 bg-red-500/10 px-3 text-xs font-medium text-red-400 transition-colors hover:bg-red-500/20 disabled:opacity-50"
					>
						{#if revoking}
							Revoking...
						{:else}
							<iconify-icon icon="solar:fire-bold" width="14" class="mr-1.5"></iconify-icon>
							Revoke
						{/if}
					</button>
				</div>
			{/if}
		</div>

		{#if error}
			<p class="text-sm text-red-400">{error}</p>
		{/if}

		<button
			onclick={reset}
			class="inline-flex h-12 items-center justify-center rounded-md border border-border bg-background px-6 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground sm:h-11"
		>
			<iconify-icon icon="solar:pill-bold-duotone" width="16" class="mr-2"></iconify-icon>
			Seal another capsule
		</button>
	</div>
{/if}
