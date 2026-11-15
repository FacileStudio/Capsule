import { test, expect, type Page } from '@playwright/test';

const PNG = Buffer.from(
	'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR4nGP4z8DwHwAFAAH/q842iQAAAABJRU5ErkJggg==',
	'base64',
);
const PDF = Buffer.from('%PDF-1.1\n1 0 obj<</Type/Catalog>>endobj\ntrailer<</Root 1 0 R>>\n%%EOF\n');

function wav(): Buffer {
	const samples = 800;
	const buf = Buffer.alloc(44 + samples * 2);
	buf.write('RIFF', 0);
	buf.writeUInt32LE(36 + samples * 2, 4);
	buf.write('WAVEfmt ', 8);
	buf.writeUInt32LE(16, 16);
	buf.writeUInt16LE(1, 20);
	buf.writeUInt16LE(1, 22);
	buf.writeUInt32LE(8000, 24);
	buf.writeUInt32LE(16000, 28);
	buf.writeUInt16LE(2, 32);
	buf.writeUInt16LE(16, 34);
	buf.write('data', 36);
	buf.writeUInt32LE(samples * 2, 40);
	return buf;
}

async function watchCsp(page: Page): Promise<string[]> {
	const violations: string[] = [];
	await page.addInitScript(() => {
		document.addEventListener('securitypolicyviolation', (event) => {
			(window as unknown as { __csp: string[] }).__csp ??= [];
			(window as unknown as { __csp: string[] }).__csp.push(`${event.violatedDirective} ${event.blockedURI}`);
		});
	});
	page.on('console', (message) => {
		if (/content security policy/i.test(message.text())) violations.push(message.text());
	});
	return violations;
}

async function collected(page: Page, consoleViolations: string[]): Promise<string[]> {
	const events = await page.evaluate(() => (window as unknown as { __csp?: string[] }).__csp ?? []);
	return [...events, ...consoleViolations];
}

async function seal(page: Page, text: string, file?: { name: string; mimeType: string; buffer: Buffer }) {
	await page.goto('/');
	await page.locator('textarea').fill(text);
	if (file) await page.locator('input[type="file"]').setInputFiles(file);
	await page.getByRole('button', { name: /seal capsule/i }).click();
	await expect(page.locator('h1')).toContainText('Capsule sealed', { timeout: 10_000 });
	return page.locator('input[readonly]').inputValue();
}

async function reveal(page: Page, url: string) {
	await page.goto(url);
	await expect(page.locator('h1')).toContainText('A capsule for you', { timeout: 10_000 });
	await page.getByRole('button', { name: /open capsule/i }).click();
	await expect(page.locator('h1')).toContainText('Capsule opened', { timeout: 10_000 });
}

test.describe('content security policy', () => {
	test('sends the full policy, not a pared-down one', async ({ request }) => {
		const response = await request.get('/');
		const enforced = response.headers()['content-security-policy'] ?? '';
		expect(enforced).toContain("default-src 'self'");
		expect(enforced).toContain("script-src 'self'");
		expect(enforced).not.toContain('vision.facile.studio');
		expect(enforced).not.toContain('unsafe-eval');
	});

	test('the API reference renders under the enforced policy', async ({ page, request }) => {
		const cdn = await request
			.get('https://cdn.jsdelivr.net/npm/@scalar/api-reference@1.64.0', { timeout: 10_000 })
			.catch(() => null);
		test.skip(!cdn?.ok(), 'the Scalar bundle on cdn.jsdelivr.net is unreachable from here');

		const logged = await watchCsp(page);
		await page.goto('/docs');
		await page.getByRole('button', { name: /^pastes Open Group/ }).click({ timeout: 20_000 });
		await expect(page.getByText('Store encrypted content on the server').first()).toBeVisible({ timeout: 20_000 });
		await page.waitForLoadState('networkidle');
		expect(await collected(page, logged)).toEqual([]);
	});

	test('the pages and flows raise no violation', async ({ page }) => {
		const logged = await watchCsp(page);
		const seen: string[] = [];

		await page.goto('/');
		await expect(page.locator('h1')).toContainText('Send a password');
		await page.waitForLoadState('networkidle');
		seen.push(...(await collected(page, logged)));

		const textUrl = await seal(page, 'plain secret');
		await reveal(page, textUrl);
		await expect(page.locator('pre')).toContainText('plain secret');
		seen.push(...(await collected(page, logged)));

		await page.goto('/');
		await page.locator('textarea').fill('const answer = 42;');
		await page.getByRole('button', { name: 'Options' }).click();
		await page.getByRole('combobox', { name: 'Syntax' }).click();
		await page.getByRole('option', { name: 'javascript' }).click();
		await page.getByRole('button', { name: /seal capsule/i }).click();
		await expect(page.locator('h1')).toContainText('Capsule sealed', { timeout: 10_000 });
		await reveal(page, await page.locator('input[readonly]').inputValue());
		await expect(page.locator('pre .shiki, pre span[style]').first()).toBeVisible({ timeout: 10_000 });
		seen.push(...(await collected(page, logged)));

		const files = [
			{ name: 'pixel.png', mimeType: 'image/png', buffer: PNG },
			{ name: 'doc.pdf', mimeType: 'application/pdf', buffer: PDF },
			{ name: 'tone.wav', mimeType: 'audio/wav', buffer: wav() },
		];
		for (const file of files) {
			const url = await seal(page, 'attached', file);
			await reveal(page, url);
			await expect(page.getByText(file.name, { exact: true })).toBeVisible({ timeout: 10_000 });
			if (file.mimeType.startsWith('image/')) {
				await expect(page.locator('img[alt="pixel.png"]')).toBeVisible();
				await expect
					.poll(() => page.locator('img[alt="pixel.png"]').evaluate((img: HTMLImageElement) => img.naturalWidth))
					.toBeGreaterThan(0);
			}
			if (file.mimeType.startsWith('audio/')) {
				await expect(page.locator('audio')).toBeVisible();
				await expect
					.poll(() => page.locator('audio').evaluate((audio: HTMLAudioElement) => audio.readyState))
					.toBeGreaterThan(0);
			}
			seen.push(...(await collected(page, logged)));
		}

		for (const path of ['/revoke', '/downloads', '/pending']) {
			await page.goto(path);
			await page.waitForLoadState('networkidle');
			seen.push(...(await collected(page, logged)));
		}

		expect(seen).toEqual([]);
	});
});
