import { test, expect, type Page } from '@playwright/test';

/*
 * muse paints dark from `@media (prefers-color-scheme: dark)` scoped to `:root:not(.light)`,
 * plus a `:root.dark` twin. So the two halves are not symmetric: `.dark` is only needed to
 * force dark on a light OS, while `.light` is the *only* thing that can force light on a dark
 * one. A toggle that writes just `.dark` looks completely correct until someone on a dark OS
 * asks for light and nothing happens.
 *
 * The same asymmetry already caused one silent bug: `highlight.ts` picked its shiki theme by
 * reading the `.dark` class, which an OS-dark visitor never has, so it served light syntax
 * colours onto a dark page. Nothing threw, so only a rendered check catches it.
 */

/* Computed colours come back as `oklch(...)`; the canvas is the cheapest honest converter. */
async function isDark(page: Page): Promise<boolean> {
	return page.evaluate(() => {
		const cv = document.createElement('canvas');
		cv.width = cv.height = 1;
		const cx = cv.getContext('2d')!;
		cx.fillStyle = getComputedStyle(document.body).backgroundColor;
		cx.fillRect(0, 0, 1, 1);
		const [r, g, b] = Array.from(cx.getImageData(0, 0, 1, 1).data);
		return 0.2126 * r + 0.7152 * g + 0.0722 * b < 128;
	});
}

test.describe('theme', () => {
	test.use({ colorScheme: 'dark' });

	test('an OS-dark visitor with no stored preference gets a dark page', async ({ page }) => {
		await page.goto('/');
		await page.locator('h1').waitFor();
		expect(await isDark(page)).toBe(true);
		await expect(page.locator('html')).not.toHaveClass(/light/);
	});

	test('forcing light on a dark OS actually lands', async ({ page }) => {
		await page.addInitScript(() => localStorage.setItem('capsule-theme', 'light'));
		await page.goto('/');
		await page.locator('h1').waitFor();
		await expect(page.locator('html')).toHaveClass(/light/);
		expect(await isDark(page)).toBe(false);
	});

	test('the toggle flips the page and survives a reload', async ({ page }) => {
		await page.goto('/');
		await page.locator('h1').waitFor();
		expect(await isDark(page)).toBe(true);

		await page.getByRole('button', { name: /toggle theme/i }).click();
		await expect(page.locator('html')).toHaveClass(/light/);
		expect(await isDark(page)).toBe(false);

		await page.reload();
		await page.locator('h1').waitFor();
		expect(await isDark(page)).toBe(false);
	});

	test('syntax highlighting follows the theme with no class set', async ({ page }) => {
		await page.goto('/');
		await page.locator('textarea').fill('const apiKey = "theme check";');
		await page.getByRole('checkbox', { name: 'Burn after opening' }).uncheck();
		await page.locator('select').last().selectOption('typescript');
		await page.getByRole('button', { name: /seal capsule/i }).click();
		await expect(page.locator('h1')).toContainText('Capsule sealed', { timeout: 15_000 });
		const url = await page.locator('input[readonly]').inputValue();

		await page.goto(url);
		await page.getByRole('button', { name: /open capsule/i }).click();
		await expect(page.locator('h1')).toContainText('Capsule opened', { timeout: 15_000 });
		await expect(page.locator('pre span').first()).toBeVisible();

		// Light ink on the dark page. Reading `.dark` instead of the theme store put the
		// light shiki theme here, which is legible only by accident.
		const readable = await page.evaluate(() => {
			const cv = document.createElement('canvas');
			cv.width = cv.height = 1;
			const cx = cv.getContext('2d')!;
			const lum = (css: string) => {
				cx.fillStyle = css;
				cx.fillRect(0, 0, 1, 1);
				const [r, g, b] = Array.from(cx.getImageData(0, 0, 1, 1).data);
				const f = (c: number) => {
					const s = c / 255;
					return s <= 0.03928 ? s / 12.92 : ((s + 0.055) / 1.055) ** 2.4;
				};
				return 0.2126 * f(r) + 0.7152 * f(g) + 0.0722 * f(b);
			};
			const token = document.querySelector('pre span')!;
			const [hi, lo] = [
				lum(getComputedStyle(token).color),
				lum(getComputedStyle(document.body).backgroundColor)
			].sort((a, b) => b - a);
			return (hi + 0.05) / (lo + 0.05);
		});

		expect(readable).toBeGreaterThan(4.5);
	});
});
