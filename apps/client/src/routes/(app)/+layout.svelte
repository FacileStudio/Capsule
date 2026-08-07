<script lang="ts">
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';

	let { children } = $props();
</script>

<div class="flex min-h-[100dvh] flex-col">
	<header class="border-b border-fc-border">
		<div class="mx-auto flex max-w-2xl items-center justify-between px-4 py-4 sm:px-6">
			<a href="/" class="flex items-center gap-3 transition-opacity hover:opacity-80">
				<!-- bold-duotone is reserved for brand marks; everything else on the page is linear. -->
				<iconify-icon icon="solar:pill-bold-duotone" width="24" height="24" class="block"></iconify-icon>
				<span class="text-fc-xl font-semibold tracking-tight">Capsule</span>
			</a>
			<ThemeToggle />
		</div>
	</header>

	<main class="mx-auto flex w-full max-w-2xl flex-1 flex-col px-4 py-8 sm:px-6 sm:py-12">
		{@render children()}
	</main>

	<footer class="border-t border-fc-border text-center">
		<div class="mx-auto max-w-2xl px-4 py-6 text-fc-sm text-fc-fg-muted sm:px-6">
			&copy; {new Date().getFullYear()} Capsule by <a href="https://facile.studio" class="font-medium underline">Facile.</a>
		</div>
	</footer>
</div>
