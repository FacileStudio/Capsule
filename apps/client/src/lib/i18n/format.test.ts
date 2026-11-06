import { describe, expect, it } from 'vitest';
import { timeAgo } from './format';

const NOW = Date.UTC(2026, 10, 6, 12, 0, 0);

describe('timeAgo', () => {
	it('picks the largest whole unit', () => {
		expect(timeAgo(NOW - 30_000, 'en', NOW)).toBe('now');
		expect(timeAgo(NOW - 5 * 60_000, 'en', NOW)).toBe('5 minutes ago');
		expect(timeAgo(NOW - 3 * 3_600_000, 'en', NOW)).toBe('3 hours ago');
		expect(timeAgo(NOW - 2 * 86_400_000, 'en', NOW)).toBe('2 days ago');
	});

	it('follows the locale', () => {
		expect(timeAgo(NOW - 3 * 3_600_000, 'fr', NOW)).toBe('il y a 3 heures');
		expect(timeAgo(NOW - 3 * 3_600_000, 'es', NOW)).toBe('hace 3 horas');
		expect(timeAgo(NOW - 3 * 3_600_000, 'de', NOW)).toBe('vor 3 Stunden');
	});

	it('never reports the future', () => {
		expect(timeAgo(NOW + 60_000, 'en', NOW)).toBe('now');
	});
});
