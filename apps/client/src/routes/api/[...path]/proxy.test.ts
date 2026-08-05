import { describe, it, expect, vi, afterEach } from 'vitest';
import { fallback } from './+server';

const proxy = (upstream: Response) => {
	vi.stubGlobal('fetch', vi.fn(async () => upstream));
	const handler = fallback as unknown as (event: unknown) => Promise<Response>;
	return handler({
		request: new Request('http://localhost/api/pastes', { method: 'GET' }),
		params: { path: 'pastes' },
		url: new URL('http://localhost/api/pastes'),
	});
};

afterEach(() => {
	vi.unstubAllGlobals();
});

describe('api proxy response headers', () => {
	it('strips content-encoding and content-length that undici already decoded', async () => {
		const body = '{"ok":true}';
		const upstream = new Response(body, {
			status: 200,
			headers: {
				'content-type': 'application/json',
				'content-encoding': 'gzip',
				'content-length': String(body.length),
			},
		});

		const response = await proxy(upstream);

		expect(response.headers.get('content-encoding')).toBeNull();
		expect(response.headers.get('content-length')).toBeNull();
		expect(response.headers.get('content-type')).toBe('application/json');
		expect(await response.text()).toBe(body);
	});
});
