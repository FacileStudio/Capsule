import { describe, it, expect } from 'vitest';
import { serializePayload, parsePayload, fileToBlob, formatBytes, type CapsuleFile } from './envelope';

const sampleFile: CapsuleFile = {
	name: 'secret.png',
	type: 'image/png',
	size: 4,
	data: 'AQIDBA==',
};

describe('serializePayload / parsePayload', () => {
	it('keeps plain text unmodified when no file is attached', () => {
		const raw = serializePayload('hello world', null);
		expect(raw).toBe('hello world');
		expect(parsePayload(raw)).toEqual({ text: 'hello world', file: null });
	});

	it('round-trips a file-only capsule', () => {
		const raw = serializePayload('', sampleFile);
		expect(raw.startsWith('{"__capsule":1')).toBe(true);
		expect(parsePayload(raw)).toEqual({ text: '', file: sampleFile });
	});

	it('round-trips text and file together', () => {
		const raw = serializePayload('see attached', sampleFile);
		expect(parsePayload(raw)).toEqual({ text: 'see attached', file: sampleFile });
	});

	it('falls back to plain text on malformed envelope JSON', () => {
		const malformed = '{"__capsule":1,broken';
		expect(parsePayload(malformed)).toEqual({ text: malformed, file: null });
	});
});

describe('fileToBlob and formatBytes', () => {
	it('reconstructs a Blob with the expected size and type', async () => {
		const blob = fileToBlob(sampleFile);
		expect(blob.size).toBe(4);
		expect(blob.type).toBe('image/png');
		const bytes = new Uint8Array(await blob.arrayBuffer());
		expect(Array.from(bytes)).toEqual([1, 2, 3, 4]);
	});

	it('formats byte sizes', () => {
		expect(formatBytes(512)).toBe('512 B');
		expect(formatBytes(2048)).toBe('2.0 KB');
		expect(formatBytes(5 * 1024 * 1024)).toBe('5.0 MB');
	});
});
