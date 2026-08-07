const backendBaseUrl = '/api';

export type PasteCreated = {
	id: string;
	delete_token: string;
	expires_at: string | null;
	created_at: string;
};

export type PasteMeta = {
	id: string;
	exists: boolean;
	burned: boolean;
	has_password: boolean;
	syntax: string;
	expires_at: string | null;
	created_at: string;
	burn_after_read: boolean;
};

export type PasteContent = {
	content: string;
};

type ApiErrorPayload = {
	error?: { message?: string };
};

async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
	const headers = new Headers(options.headers);
	if (!headers.has('Content-Type') && options.body && !(options.body instanceof FormData)) {
		headers.set('Content-Type', 'application/json');
	}

	const response = await fetch(`${backendBaseUrl}${path}`, { ...options, headers });

	if (!response.ok) {
		let payload: ApiErrorPayload | undefined;
		try {
			payload = (await response.json()) as ApiErrorPayload;
		} catch {
			payload = undefined;
		}
		throw new Error(payload?.error?.message || `Request failed with status ${response.status}`);
	}

	const text = await response.text();
	if (!text) return {} as T;
	return JSON.parse(text) as T;
}

export const backend = {
	createPaste(data: {
		content: string;
		burn_after_read: boolean;
		expires_in?: string;
		max_views?: number;
		has_password?: boolean;
		syntax?: string;
	}) {
		return apiFetch<PasteCreated>('/pastes', {
			method: 'POST',
			body: JSON.stringify(data),
		});
	},

	getPasteMeta(id: string) {
		return apiFetch<PasteMeta>(`/pastes/${id}`);
	},

	getPasteContent(id: string) {
		return apiFetch<PasteContent>(`/pastes/${id}/content`, {
			method: 'POST',
		});
	},

	deletePaste(id: string, deleteToken: string) {
		return apiFetch<void>(`/pastes/${id}`, {
			method: 'DELETE',
			headers: { 'X-Delete-Token': deleteToken },
		});
	},
};
