import { defineConfig } from 'vitest/config';

const ENV_ID = '$env/dynamic/private';
const ENV_RESOLVED = '\0$env/dynamic/private';

const sveltekitEnvStub = {
	name: 'sveltekit-env-stub',
	resolveId(id: string) {
		return id === ENV_ID ? ENV_RESOLVED : null;
	},
	load(id: string) {
		return id === ENV_RESOLVED ? 'export const env = {};' : null;
	},
};

export default defineConfig({
	plugins: [sveltekitEnvStub],
	test: {
		include: ['src/**/*.test.ts'],
		environment: 'node',
	},
});
