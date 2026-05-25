import { test, expect } from '@playwright/test';

test.describe('create page', () => {
	test('renders seal form', async ({ page }) => {
		await page.goto('/');
		await expect(page.locator('h1')).toContainText('Seal a capsule');
		await expect(page.locator('textarea')).toBeVisible();
		await expect(page.getByRole('button', { name: /seal capsule/i })).toBeVisible();
	});

	test('seal button disabled when textarea empty', async ({ page }) => {
		await page.goto('/');
		const btn = page.getByRole('button', { name: /seal capsule/i });
		await expect(btn).toBeDisabled();
	});

	test('seal button enabled when textarea has content', async ({ page }) => {
		await page.goto('/');
		await page.locator('textarea').fill('test content');
		const btn = page.getByRole('button', { name: /seal capsule/i });
		await expect(btn).toBeEnabled();
	});

	test('syntax selector has expected options', async ({ page }) => {
		await page.goto('/');
		const syntaxSelect = page.locator('select').last();
		const options = await syntaxSelect.locator('option').allTextContents();
		expect(options).toContain('plaintext');
		expect(options).toContain('javascript');
		expect(options).toContain('python');
		expect(options).toContain('go');
	});

	test('expiry selector has expected options', async ({ page }) => {
		await page.goto('/');
		const expirySelect = page.locator('select').first();
		const options = await expirySelect.locator('option').allTextContents();
		expect(options).toContain('1 hour');
		expect(options).toContain('24 hours');
		expect(options).toContain('7 days');
		expect(options).toContain('30 days');
	});

	test('burn after read checkbox exists and is checked by default', async ({ page }) => {
		await page.goto('/');
		const checkbox = page.locator('input[type="checkbox"]');
		await expect(checkbox).toBeChecked();
	});
});

test.describe('sealed state', () => {
	test('sealing a capsule shows success UI', async ({ page }) => {
		await page.goto('/');
		await page.locator('textarea').fill('my secret paste');
		await page.getByRole('button', { name: /seal capsule/i }).click();

		await expect(page.locator('h1')).toContainText('Capsule sealed', { timeout: 10_000 });
		await expect(page.locator('input[readonly]')).toBeVisible();
		await expect(page.getByRole('button', { name: /copy link/i })).toBeVisible();
	});

	test('capsule URL contains paste id and key fragment', async ({ page }) => {
		await page.goto('/');
		await page.locator('textarea').fill('url check');
		await page.getByRole('button', { name: /seal capsule/i }).click();

		await expect(page.locator('h1')).toContainText('Capsule sealed', { timeout: 10_000 });

		const urlInput = page.locator('input[readonly]');
		const capsuleUrl = await urlInput.inputValue();
		expect(capsuleUrl).toMatch(/\/cap_[0-9a-f]+#.+/);
	});

	test('delete token is shown', async ({ page }) => {
		await page.goto('/');
		await page.locator('textarea').fill('token check');
		await page.getByRole('button', { name: /seal capsule/i }).click();

		await expect(page.locator('h1')).toContainText('Capsule sealed', { timeout: 10_000 });
		await expect(page.locator('code')).toBeVisible();
		const token = await page.locator('code').textContent();
		expect(token).toBeTruthy();
		expect(token!.length).toBeGreaterThan(10);
	});

	test('seal another capsule button resets', async ({ page }) => {
		await page.goto('/');
		await page.locator('textarea').fill('reset check');
		await page.getByRole('button', { name: /seal capsule/i }).click();

		await expect(page.locator('h1')).toContainText('Capsule sealed', { timeout: 10_000 });

		await page.getByRole('button', { name: /seal another/i }).click();
		await expect(page.locator('h1')).toContainText('Seal a capsule');
		await expect(page.locator('textarea')).toHaveValue('');
	});
});

test.describe('reveal page', () => {
	test('non-existent paste shows empty state', async ({ page }) => {
		await page.goto('/cap_0000000000000000#fakekey');
		await expect(page.locator('h1')).toContainText('empty', { timeout: 10_000 });
	});

	test('missing key fragment shows error', async ({ page }) => {
		await page.goto('/cap_0000000000000000');
		await expect(page.locator('h1')).toContainText('Something went wrong', { timeout: 10_000 });
	});
});

test.describe('full create → reveal flow', () => {
	test('create paste, navigate to URL, reveal content', async ({ page }) => {
		await page.goto('/');
		const secret = 'E2E test secret ' + Date.now();
		await page.locator('textarea').fill(secret);

		const burnCheckbox = page.locator('input[type="checkbox"]');
		if (await burnCheckbox.isChecked()) {
			await burnCheckbox.uncheck();
		}

		await page.getByRole('button', { name: /seal capsule/i }).click();
		await expect(page.locator('h1')).toContainText('Capsule sealed', { timeout: 10_000 });

		const capsuleUrl = await page.locator('input[readonly]').inputValue();

		await page.goto(capsuleUrl);
		await expect(page.locator('h1')).toContainText('A capsule for you', { timeout: 10_000 });

		await page.getByRole('button', { name: /open capsule/i }).click();
		await expect(page.locator('h1')).toContainText('Capsule opened', { timeout: 10_000 });
		await expect(page.locator('pre')).toContainText(secret);
	});

	test('burn-after-read: second reveal fails', async ({ page, context }) => {
		await page.goto('/');
		await page.locator('textarea').fill('burn test ' + Date.now());
		await page.getByRole('button', { name: /seal capsule/i }).click();
		await expect(page.locator('h1')).toContainText('Capsule sealed', { timeout: 10_000 });

		const capsuleUrl = await page.locator('input[readonly]').inputValue();

		await page.goto(capsuleUrl);
		await expect(page.locator('h1')).toContainText('A capsule for you', { timeout: 10_000 });
		await expect(page.locator('text=Burns after opening')).toBeVisible();
		await page.getByRole('button', { name: /open capsule/i }).click();
		await expect(page.locator('h1')).toContainText('Capsule opened', { timeout: 10_000 });

		const page2 = await context.newPage();
		await page2.goto(capsuleUrl);
		await expect(page2.locator('h1')).toContainText('empty', { timeout: 10_000 });
	});
});
