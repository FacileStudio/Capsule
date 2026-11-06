import type { Row } from '../row';

/** Strings of the page chrome: header, footer and theme toggle. */
export const shell = {
	'shell.seal': ['Seal', 'Sceller', 'Sellar', 'Versiegeln'],
	'shell.revoke': ['Revoke', 'Révoquer', 'Revocar', 'Widerrufen'],
	'shell.api': ['API reference', 'Référence de l’API', 'Referencia de la API', 'API-Referenz'],
	'shell.theme': ['Toggle theme', 'Changer de thème', 'Cambiar de tema', 'Design wechseln'],
	'shell.tagline': [
		'Secrets that burn after reading. The server only ever stores ciphertext.',
		'Des secrets qui se détruisent à la lecture. Le serveur ne stocke que du texte chiffré.',
		'Secretos que se destruyen al leerlos. El servidor solo almacena texto cifrado.',
		'Geheimnisse, die sich nach dem Lesen auflösen. Der Server speichert nur Chiffretext.'
	],
	'shell.facile': ['Facile Studio', 'Facile Studio', 'Facile Studio', 'Facile Studio']
} satisfies Record<string, Row>;
