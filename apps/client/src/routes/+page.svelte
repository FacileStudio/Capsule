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
	let deleteToken = $state('');
	let copied = $state(false);

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

	function reset() {
		state = 'idle';
		content = '';
		usePassword = false;
		password = '';
		capsuleUrl = '';
		deleteToken = '';
		error = '';
		copied = false;
	}
</script>

<svelte:head>
	<title>Capsule — Ephemeral Encrypted Sharing</title>
	<meta name="description" content="Share secrets that self-destruct. End-to-end encrypted, zero knowledge." />
</svelte:head>

<div class="flex min-h-screen flex-col bg-background text-foreground">
	<header class="border-b border-border">
		<div class="mx-auto flex max-w-2xl items-center gap-3 px-6 py-4">
			<iconify-icon icon="solar:pill-bold-duotone" width="28" class="text-foreground"></iconify-icon>
			<span class="text-2xl font-bold font-heading tracking-tight">Capsule</span>
		</div>
	</header>

	<main class="mx-auto flex w-full max-w-2xl flex-1 flex-col px-6 py-12">
		{#if state === 'idle' || state === 'sealing'}
			<div class="flex flex-col gap-6">
				<div>
					<h1 class="text-3xl font-bold font-heading tracking-tight">Seal a capsule</h1>
					<p class="mt-2 text-muted-foreground">
						Your content is encrypted in the browser. The server never sees the plaintext.
					</p>
				</div>

				<textarea
					bind:value={content}
					placeholder="Paste your secret here..."
					rows="10"
					disabled={state === 'sealing'}
					class="w-full resize-y rounded-lg border border-input bg-card p-4 font-mono text-sm text-card-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
				></textarea>

				<div class="flex flex-wrap items-center gap-6">
					<label class="flex items-center gap-2 text-sm">
						<input
							type="checkbox"
							bind:checked={burnAfterRead}
							disabled={state === 'sealing'}
							class="h-4 w-4 rounded border-input accent-primary"
						/>
						Burn after opening
					</label>

					<label class="flex items-center gap-2 text-sm">
						<span class="text-muted-foreground">Expires in</span>
						<select
							bind:value={expiresIn}
							disabled={state === 'sealing'}
							class="rounded-md border border-input bg-card px-2 py-1 text-sm"
						>
							{#each expiryOptions as opt}
								<option value={opt.value}>{opt.label}</option>
							{/each}
						</select>
					</label>

					<label class="flex items-center gap-2 text-sm">
						<span class="text-muted-foreground">Syntax</span>
						<select
							bind:value={syntax}
							disabled={state === 'sealing'}
							class="rounded-md border border-input bg-card px-2 py-1 text-sm"
						>
							{#each syntaxOptions as s}
								<option value={s}>{s}</option>
							{/each}
						</select>
					</label>

					<label class="flex items-center gap-2 text-sm">
						<input
							type="checkbox"
							bind:checked={usePassword}
							disabled={state === 'sealing'}
							class="h-4 w-4 rounded border-input accent-primary"
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

				<button
					onclick={seal}
					disabled={!content.trim() || (usePassword && !password) || state === 'sealing'}
					class="inline-flex h-11 items-center justify-center rounded-md bg-primary px-6 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed"
				>
					{#if state === 'sealing'}
						<iconify-icon icon="solar:refresh-bold" width="16" class="mr-2 animate-spin"></iconify-icon>
						Sealing...
					{:else}
						<iconify-icon icon="solar:pill-bold-duotone" width="16" class="mr-2"></iconify-icon>
						Seal capsule
					{/if}
				</button>
			</div>
		{:else if state === 'sealed'}
			<div class="flex flex-col gap-6">
				<div>
					<div class="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-green-500/10">
						<iconify-icon icon="solar:check-circle-bold" width="28" class="text-green-400"></iconify-icon>
					</div>
					<h1 class="text-3xl font-bold font-heading tracking-tight">Capsule sealed</h1>
					<p class="mt-2 text-muted-foreground">
						Share this link{usePassword ? ' and the password' : ''} — {burnAfterRead ? 'the capsule can only be opened once.' : `it expires in ${expiryOptions.find(o => o.value === expiresIn)?.label}.`}
					</p>
				</div>

				<div class="flex items-stretch gap-2">
					<input
						type="text"
						readonly
						value={capsuleUrl}
						class="flex-1 rounded-lg border border-input bg-card px-4 py-3 font-mono text-sm text-card-foreground"
					/>
					<button
						onclick={copyLink}
						class="inline-flex items-center justify-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
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
					<p class="text-xs text-muted-foreground">
						Delete token (save this to revoke the capsule before it's opened):
					</p>
					<code class="mt-1 block break-all text-xs text-muted-foreground/70">{deleteToken}</code>
				</div>

				<button
					onclick={reset}
					class="inline-flex h-11 items-center justify-center rounded-md border border-border bg-background px-6 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
				>
					<iconify-icon icon="solar:pill-bold-duotone" width="16" class="mr-2"></iconify-icon>
					Seal another capsule
				</button>
			</div>
		{/if}
	</main>

	<footer class="border-t border-border text-center">
		<div class="mx-auto max-w-2xl px-6 py-6 text-sm text-muted-foreground">
			&copy; {new Date().getFullYear()} Capsule by <a href="https://facile.studio" class="font-semibold underline hover:cursor-pointer">Facile.</a>
		</div>
	</footer>
</div>
