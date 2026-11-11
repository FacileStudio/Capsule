import { createApi } from '@facile/client';

const backendBaseUrl = '/api';

const apiFetch = createApi({ base: backendBaseUrl });

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

export type SessionUser = { email: string; name?: string; avatar_url?: string };

export type Invitations = { count: number; href: string };

export type PendingCapsule = {
	id: string;
	created_at: string;
	expires_at: string | null;
	burn_after_read: boolean;
	has_password: boolean;
};

export type PendingAction = {
	id: string;
	kind: string;
	client_id: string;
	summary: string;
	created_at: string;
	expires_at: string;
	state: 'ready' | 'changed' | 'missing';
	capsule?: PendingCapsule;
};

export const backend = {
	async sessionUser() {
		const { user } = await apiFetch<{ user: SessionUser | null }>('/auth/session');
		return user;
	},

	async invitations(): Promise<Invitations | null> {
		const answer = await apiFetch<{ count: number; href?: string }>('/me/invitations');
		return answer.href ? { count: answer.count, href: answer.href } : null;
	},

	logout() {
		return apiFetch<void>('/auth/logout', { method: 'POST' });
	},

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

	listPending() {
		return apiFetch<{ actions: PendingAction[] }>('/mcp/pending');
	},

	confirmPending(id: string) {
		return apiFetch<{ status: string }>(`/mcp/pending/${encodeURIComponent(id)}/confirm`, { method: 'POST' });
	},

	dismissPending(id: string) {
		return apiFetch<{ status: string }>(`/mcp/pending/${encodeURIComponent(id)}/dismiss`, { method: 'POST' });
	},

	deletePaste(id: string, deleteToken: string) {
		return apiFetch<void>(`/pastes/${id}`, {
			method: 'DELETE',
			headers: { 'X-Delete-Token': deleteToken },
		});
	},
};
