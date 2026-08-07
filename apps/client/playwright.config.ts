import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
	testDir: 'tests',
	fullyParallel: true,
	forbidOnly: !!process.env.CI,
	retries: process.env.CI ? 2 : 0,
	workers: process.env.CI ? 1 : undefined,
	reporter: 'html',
	use: {
		baseURL: process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:5173',
		trace: 'on-first-retry',
	},
	projects: [
		{
			name: 'chromium',
			use: { ...devices['Desktop Chrome'] },
		},
	],
	/*
	 * Only boot a dev server when the suite is pointed at the default local URL. The pages
	 * call `/api` with no dev proxy in front of them, so a local `vite dev` has no backend
	 * and every sealing test fails — the way to actually run this suite today is
	 * `PLAYWRIGHT_BASE_URL=https://capsule.facile.studio bun run test:e2e`, and starting a
	 * server it will not talk to just costs 30s and a misleading failure mode.
	 */
	webServer: process.env.PLAYWRIGHT_BASE_URL
		? undefined
		: {
				command: 'bun run dev',
				url: 'http://localhost:5173',
				reuseExistingServer: !process.env.CI,
				timeout: 30_000,
			},
});
