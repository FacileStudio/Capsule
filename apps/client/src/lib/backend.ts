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

	deletePaste(id: string, deleteToken: string) {
		return apiFetch<void>(`/pastes/${id}`, {
			method: 'DELETE',
			headers: { 'X-Delete-Token': deleteToken },
		});
	},
};
