import { fromBase64, toBase64 } from './envelope';

const PBKDF2_ITERATIONS = 600_000;

export async function generateKey(): Promise<CryptoKey> {
	return crypto.subtle.generateKey(
		{ name: 'AES-GCM', length: 256 },
		true,
		['encrypt', 'decrypt']
	);
}

export async function exportKey(key: CryptoKey): Promise<string> {
	const raw = await crypto.subtle.exportKey('raw', key);
	return toBase64(new Uint8Array(raw), true);
}

export async function importKey(encoded: string): Promise<CryptoKey> {
	return crypto.subtle.importKey(
		'raw',
		fromBase64(encoded),
		{ name: 'AES-GCM', length: 256 },
		false,
		['decrypt']
	);
}

export async function encrypt(plaintext: string, key: CryptoKey): Promise<string> {
	const iv = crypto.getRandomValues(new Uint8Array(12));
	const encoded = new TextEncoder().encode(plaintext);

	const ciphertext = await crypto.subtle.encrypt(
		{ name: 'AES-GCM', iv },
		key,
		encoded
	);

	const combined = new Uint8Array(iv.length + ciphertext.byteLength);
	combined.set(iv, 0);
	combined.set(new Uint8Array(ciphertext), iv.length);

	return toBase64(combined);
}

export async function decrypt(data: string, key: CryptoKey): Promise<string> {
	const combined = fromBase64(data);

	const iv = combined.slice(0, 12);
	const ciphertext = combined.slice(12);

	const decrypted = await crypto.subtle.decrypt(
		{ name: 'AES-GCM', iv },
		key,
		ciphertext
	);

	return new TextDecoder().decode(decrypted);
}

async function deriveWrappingKey(password: string, salt: Uint8Array<ArrayBuffer>): Promise<CryptoKey> {
	const encoded = new TextEncoder().encode(password);
	const baseKey = await crypto.subtle.importKey('raw', encoded, 'PBKDF2', false, [
		'deriveKey',
	]);

	return crypto.subtle.deriveKey(
		{
			name: 'PBKDF2',
			salt,
			iterations: PBKDF2_ITERATIONS,
			hash: 'SHA-256',
		},
		baseKey,
		{ name: 'AES-GCM', length: 256 },
		false,
		['encrypt', 'decrypt']
	);
}

export async function wrapContentKey(
	contentKey: CryptoKey,
	password: string
): Promise<{ encryptedKey: string; salt: string; iv: string; fragment: string }> {
	const salt = crypto.getRandomValues(new Uint8Array(32));
	const iv = crypto.getRandomValues(new Uint8Array(12));
	const wrappingKey = await deriveWrappingKey(password, salt);

	const rawKey = await crypto.subtle.exportKey('raw', contentKey);
	const wrapped = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, wrappingKey, rawKey);

	const encryptedKey = toBase64(new Uint8Array(wrapped), true);
	const saltStr = toBase64(salt, true);
	const ivStr = toBase64(iv, true);

	return {
		encryptedKey,
		salt: saltStr,
		iv: ivStr,
		fragment: `${encryptedKey}.${saltStr}.${ivStr}`,
	};
}

export async function unwrapContentKey(
	fragmentOrKey: string,
	password: string,
	salt?: string,
	iv?: string
): Promise<CryptoKey> {
	let encryptedKey = fragmentOrKey;
	let saltVal = salt;
	let ivVal = iv;
	if (!saltVal || !ivVal) {
		const parts = fragmentOrKey.split('.');
		if (parts.length !== 3) throw new Error('Invalid URL fragment');
		[encryptedKey, saltVal, ivVal] = parts;
	}
	const wrappingKey = await deriveWrappingKey(password, fromBase64(saltVal));
	const rawKey = await crypto.subtle.decrypt(
		{ name: 'AES-GCM', iv: fromBase64(ivVal) },
		wrappingKey,
		fromBase64(encryptedKey)
	);

	return crypto.subtle.importKey(
		'raw',
		rawKey,
		{ name: 'AES-GCM', length: 256 },
		false,
		['decrypt']
	);
}
