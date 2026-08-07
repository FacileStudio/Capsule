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
	return toBase64Url(new Uint8Array(raw));
}

export async function importKey(encoded: string): Promise<CryptoKey> {
	const raw = fromBase64Url(encoded);
	return crypto.subtle.importKey(
		'raw',
		raw,
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

	return btoa(String.fromCharCode(...combined));
}

export async function decrypt(data: string, key: CryptoKey): Promise<string> {
	const combined = Uint8Array.from(atob(data), (c) => c.charCodeAt(0));

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
): Promise<{ encryptedKey: string; salt: string; iv: string }> {
	const salt = crypto.getRandomValues(new Uint8Array(32));
	const iv = crypto.getRandomValues(new Uint8Array(12));
	const wrappingKey = await deriveWrappingKey(password, salt);

	const rawKey = await crypto.subtle.exportKey('raw', contentKey);
	const wrapped = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, wrappingKey, rawKey);

	return {
		encryptedKey: toBase64Url(new Uint8Array(wrapped)),
		salt: toBase64Url(salt),
		iv: toBase64Url(iv),
	};
}

export async function unwrapContentKey(
	encryptedKey: string,
	salt: string,
	iv: string,
	password: string
): Promise<CryptoKey> {
	const wrappingKey = await deriveWrappingKey(password, fromBase64Url(salt));
	const rawKey = await crypto.subtle.decrypt(
		{ name: 'AES-GCM', iv: fromBase64Url(iv) },
		wrappingKey,
		fromBase64Url(encryptedKey)
	);

	return crypto.subtle.importKey(
		'raw',
		rawKey,
		{ name: 'AES-GCM', length: 256 },
		false,
		['decrypt']
	);
}

export function buildPasswordFragment(encryptedKey: string, salt: string, iv: string): string {
	return `${encryptedKey}.${salt}.${iv}`;
}

export function parsePasswordFragment(
	fragment: string
): { encryptedKey: string; salt: string; iv: string } | null {
	const parts = fragment.split('.');
	if (parts.length !== 3) return null;
	return { encryptedKey: parts[0], salt: parts[1], iv: parts[2] };
}

function toBase64Url(bytes: Uint8Array): string {
	const binary = String.fromCharCode(...bytes);
	return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

// Uint8Array is generic over its backing buffer since TypeScript 5.7, and the
// bare form widens to ArrayBufferLike — which includes SharedArrayBuffer and so
// is not a BufferSource. Every WebCrypto call here takes the result, so the
// annotation belongs at the source rather than at each call site.
function fromBase64Url(str: string): Uint8Array<ArrayBuffer> {
	const padded = str.replace(/-/g, '+').replace(/_/g, '/');
	const binary = atob(padded);
	return Uint8Array.from(binary, (c) => c.charCodeAt(0));
}
