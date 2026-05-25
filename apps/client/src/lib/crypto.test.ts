import { describe, it, expect } from 'vitest';
import { generateKey, exportKey, importKey, encrypt, decrypt } from './crypto';

describe('generateKey', () => {
	it('produces a CryptoKey', async () => {
		const key = await generateKey();
		expect(key).toBeDefined();
		expect(key.algorithm).toMatchObject({ name: 'AES-GCM', length: 256 });
		expect(key.usages).toContain('encrypt');
		expect(key.usages).toContain('decrypt');
	});

	it('produces unique keys each call', async () => {
		const k1 = await exportKey(await generateKey());
		const k2 = await exportKey(await generateKey());
		expect(k1).not.toBe(k2);
	});
});

describe('exportKey / importKey round-trip', () => {
	it('exports to a base64url string', async () => {
		const key = await generateKey();
		const exported = await exportKey(key);
		expect(typeof exported).toBe('string');
		expect(exported.length).toBeGreaterThan(0);
		expect(exported).not.toMatch(/[+/=]/);
	});

	it('import restores a usable key', async () => {
		const original = await generateKey();
		const exported = await exportKey(original);
		const restored = await importKey(exported);
		expect(restored.algorithm).toMatchObject({ name: 'AES-GCM', length: 256 });
		expect(restored.usages).toContain('decrypt');
	});

	it('exported key is deterministic for same CryptoKey', async () => {
		const key = await generateKey();
		const a = await exportKey(key);
		const b = await exportKey(key);
		expect(a).toBe(b);
	});
});

describe('encrypt / decrypt', () => {
	it('round-trips basic text', async () => {
		const key = await generateKey();
		const plaintext = 'Hello, Capsule!';
		const ciphertext = await encrypt(plaintext, key);
		const decrypted = await decrypt(ciphertext, key);
		expect(decrypted).toBe(plaintext);
	});

	it('round-trips empty string', async () => {
		const key = await generateKey();
		const ciphertext = await encrypt('', key);
		const decrypted = await decrypt(ciphertext, key);
		expect(decrypted).toBe('');
	});

	it('round-trips unicode and emoji', async () => {
		const key = await generateKey();
		const plaintext = '日本語テスト 🔐💊🎉 Ñoño';
		const ciphertext = await encrypt(plaintext, key);
		const decrypted = await decrypt(ciphertext, key);
		expect(decrypted).toBe(plaintext);
	});

	it('round-trips large content', async () => {
		const key = await generateKey();
		const plaintext = 'A'.repeat(100_000);
		const ciphertext = await encrypt(plaintext, key);
		const decrypted = await decrypt(ciphertext, key);
		expect(decrypted).toBe(plaintext);
	});

	it('round-trips multiline content', async () => {
		const key = await generateKey();
		const plaintext = 'line 1\nline 2\nline 3\n\ttabbed';
		const ciphertext = await encrypt(plaintext, key);
		const decrypted = await decrypt(ciphertext, key);
		expect(decrypted).toBe(plaintext);
	});

	it('produces different ciphertexts for same plaintext (IV uniqueness)', async () => {
		const key = await generateKey();
		const plaintext = 'same input';
		const c1 = await encrypt(plaintext, key);
		const c2 = await encrypt(plaintext, key);
		expect(c1).not.toBe(c2);
	});

	it('ciphertext is base64 encoded', async () => {
		const key = await generateKey();
		const ciphertext = await encrypt('test', key);
		expect(() => atob(ciphertext)).not.toThrow();
	});

	it('fails to decrypt with wrong key', async () => {
		const key1 = await generateKey();
		const key2 = await generateKey();
		const ciphertext = await encrypt('secret', key1);
		await expect(decrypt(ciphertext, key2)).rejects.toThrow();
	});

	it('fails to decrypt corrupted ciphertext', async () => {
		const key = await generateKey();
		const ciphertext = await encrypt('data', key);
		const corrupted = ciphertext.slice(0, -4) + 'AAAA';
		await expect(decrypt(corrupted, key)).rejects.toThrow();
	});
});

describe('full flow: generate → export → import → encrypt → decrypt', () => {
	it('simulates the real capsule workflow', async () => {
		const encryptKey = await generateKey();
		const keyString = await exportKey(encryptKey);

		const plaintext = 'Super secret paste content 🔑';
		const ciphertext = await encrypt(plaintext, encryptKey);

		const decryptKey = await importKey(keyString);
		const recovered = await decrypt(ciphertext, decryptKey);

		expect(recovered).toBe(plaintext);
	});
});
