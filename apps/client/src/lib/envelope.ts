export const MAX_FILE_SIZE = 10 * 1024 * 1024;

export type CapsuleFile = {
	name: string;
	type: string;
	size: number;
	data: string;
};

export type CapsulePayload = {
	text: string;
	file: CapsuleFile | null;
};

type Envelope = {
	__capsule: 1;
	text: string;
	file?: CapsuleFile;
};

/** Encodes a byte buffer into base64 or base64url in chunks to avoid call stack limits. */
export function toBase64(bytes: Uint8Array, urlSafe = false): string {
	let binary = '';
	const chunk = 8192;
	for (let i = 0; i < bytes.length; i += chunk) {
		binary += String.fromCharCode(...bytes.subarray(i, i + chunk));
	}
	const b64 = btoa(binary);
	return urlSafe ? b64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '') : b64;
}

/** Decodes a base64 or base64url string into an ArrayBuffer-backed Uint8Array. */
export function fromBase64(str: string): Uint8Array<ArrayBuffer> {
	const padded = str.replace(/-/g, '+').replace(/_/g, '/');
	const binary = atob(padded);
	const bytes = new Uint8Array(binary.length);
	for (let i = 0; i < binary.length; i++) {
		bytes[i] = binary.charCodeAt(i);
	}
	return bytes;
}

/** Encodes a File into a CapsuleFile with base64-encoded bytes. */
export async function encodeFile(file: File): Promise<CapsuleFile> {
	const bytes = new Uint8Array(await file.arrayBuffer());
	return {
		name: file.name,
		type: file.type || 'application/octet-stream',
		size: file.size,
		data: toBase64(bytes),
	};
}

/** Serializes text and an optional attached file into the plaintext string to encrypt. */
export function serializePayload(text: string, file: CapsuleFile | null): string {
	if (!file) return text;
	const env: Envelope = { __capsule: 1, text, file };
	return JSON.stringify(env);
}

function isCapsuleFile(val: unknown): val is CapsuleFile {
	if (!val || typeof val !== 'object') return false;
	const f = val as Record<string, unknown>;
	return (
		typeof f.name === 'string' &&
		typeof f.type === 'string' &&
		typeof f.size === 'number' &&
		typeof f.data === 'string'
	);
}

/** Parses decrypted plaintext into text and an optional file, preserving plain-text compatibility. */
export function parsePayload(decrypted: string): CapsulePayload {
	if (!decrypted.startsWith('{"__capsule":1')) {
		return { text: decrypted, file: null };
	}
	try {
		const parsed = JSON.parse(decrypted) as Record<string, unknown>;
		if (parsed.__capsule !== 1 || typeof parsed.text !== 'string') {
			return { text: decrypted, file: null };
		}
		const file = isCapsuleFile(parsed.file) ? parsed.file : null;
		return { text: parsed.text, file };
	} catch {
		return { text: decrypted, file: null };
	}
}

/** Converts a CapsuleFile's base64 data back into a Blob. */
export function fileToBlob(file: CapsuleFile): Blob {
	return new Blob([fromBase64(file.data)], { type: file.type || 'application/octet-stream' });
}

/** Formats a byte count into a human-readable string (B, KB, MB, GB). */
export function formatBytes(bytes: number): string {
	if (!Number.isFinite(bytes) || bytes < 0) return '0 B';
	if (bytes < 1024) return `${Math.round(bytes)} B`;
	const units = ['KB', 'MB', 'GB'];
	let value = bytes / 1024;
	let unit = 0;
	while (value >= 1024 && unit < units.length - 1) {
		value /= 1024;
		unit += 1;
	}
	return `${value >= 10 ? Math.round(value) : value.toFixed(1)} ${units[unit]}`;
}
