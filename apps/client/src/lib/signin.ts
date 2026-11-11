const RETURN_KEY = 'capsule-return';
const RETURN_PATTERN = /^\/pending(\?action=[A-Za-z0-9_-]{1,64})?$/;

/** The sign-in route. porte lands the browser on the home page once it is done. */
export const SIGN_IN_HREF = '/api/auth/oidc';

/** Remembers a confirm page path, so the home page can send the person back after sign-in. */
export function rememberReturn(path: string) {
	if (!RETURN_PATTERN.test(path)) return;
	try {
		sessionStorage.setItem(RETURN_KEY, path);
	} catch {
		return;
	}
}

/** Returns the remembered confirm page path once, or null. Anything but a confirm page path is dropped. */
export function takeReturn(): string | null {
	try {
		const path = sessionStorage.getItem(RETURN_KEY);
		sessionStorage.removeItem(RETURN_KEY);
		return path && RETURN_PATTERN.test(path) ? path : null;
	} catch {
		return null;
	}
}
