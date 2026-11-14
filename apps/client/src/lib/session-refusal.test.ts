import { describe, expect, it } from 'vitest';
import { refusesOwnSession } from './session-refusal';

/** Shapes an error the way `@facile/client`'s ApiError carries it. */
function apiError(status: number, code: string, message: string) {
	return Object.assign(new Error(message), { name: 'ApiError', status, code });
}

describe('own-session refusal', () => {
	it('asks to sign in again only for the own-session refusal', () => {
		expect(refusesOwnSession(apiError(403, 'permission_denied', "only the person's own session can confirm actions"))).toBe(true);
		expect(refusesOwnSession(apiError(412, 'failed_precondition', 'This capsule changed'))).toBe(false);
		expect(refusesOwnSession(apiError(403, 'insufficient_scope', 'no'))).toBe(false);
		expect(refusesOwnSession(new Error('network'))).toBe(false);
		expect(refusesOwnSession(null)).toBe(false);
	});
});
