import type { HighlighterCore } from 'shiki';
import { isDark } from './theme.svelte';

let highlighter: HighlighterCore | null = null;

const SUPPORTED_LANGS = [
	'javascript',
	'typescript',
	'python',
	'go',
	'rust',
	'json',
	'yaml',
	'sql',
	'bash',
	'html',
	'css',
	'markdown'
] as const;

async function getHighlighter(): Promise<HighlighterCore> {
	if (highlighter) return highlighter;

	const { createHighlighterCore } = await import('shiki/core');
	const { createJavaScriptRegexEngine } = await import('shiki/engine/javascript');

	highlighter = await createHighlighterCore({
		themes: [
			import('shiki/themes/github-dark.mjs'),
			import('shiki/themes/github-light.mjs')
		],
		langs: [
			import('shiki/langs/javascript.mjs'),
			import('shiki/langs/typescript.mjs'),
			import('shiki/langs/python.mjs'),
			import('shiki/langs/go.mjs'),
			import('shiki/langs/rust.mjs'),
			import('shiki/langs/json.mjs'),
			import('shiki/langs/yaml.mjs'),
			import('shiki/langs/sql.mjs'),
			import('shiki/langs/bash.mjs'),
			import('shiki/langs/html.mjs'),
			import('shiki/langs/css.mjs'),
			import('shiki/langs/markdown.mjs')
		],
		engine: createJavaScriptRegexEngine()
	});

	return highlighter;
}

export async function highlight(code: string, lang: string): Promise<string> {
	if (!lang || lang === 'plaintext' || !SUPPORTED_LANGS.includes(lang as any)) {
		return '';
	}

	const hl = await getHighlighter();

	return hl.codeToHtml(code, {
		lang,
		themes: {
			dark: 'github-dark',
			light: 'github-light'
		},
		defaultColor: isDark() ? 'dark' : 'light'
	});
}
