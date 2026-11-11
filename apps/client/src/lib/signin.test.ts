import { beforeEach, describe, expect, it } from 'vitest';
import { rememberReturn, takeReturn } from './signin';

describe('sign-in return', () => {
	beforeEach(() => {
		const store = new Map<string, string>();
		globalThis.sessionStorage = {
			getItem: (key: string) => store.get(key) ?? null,
			setItem: (key: string, value: string) => void store.set(key, value),
			removeItem: (key: string) => void store.delete(key)
		} as Storage;
	});

	it('returns a confirm page path once', () => {
		rememberReturn('/pending?action=abc_DEF-1');
		expect(takeReturn()).toBe('/pending?action=abc_DEF-1');
		expect(takeReturn()).toBeNull();
	});

	it('drops anything but a confirm page path', () => {
		rememberReturn('//evil.example/pending');
		expect(takeReturn()).toBeNull();
		sessionStorage.setItem('capsule-return', 'https://evil.example');
		expect(takeReturn()).toBeNull();
	});
});
