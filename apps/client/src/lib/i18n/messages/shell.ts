import type { Row } from '../row';

/** Strings of the page chrome: header, footer and theme toggle. */
export const shell = {
	'shell.seal': ['Seal', 'Sceller', 'Sellar', 'Versiegeln'],
	'shell.revoke': ['Revoke', 'Révoquer', 'Revocar', 'Widerrufen'],
	'shell.api': ['API reference', 'Référence de l’API', 'Referencia de la API', 'API-Referenz'],
	'shell.theme': ['Toggle theme', 'Changer de thème', 'Cambiar de tema', 'Design wechseln'],
	'shell.tagline': [
		'Share a secret that disappears once read.',
		'Partagez un secret qui disparaît après lecture.',
		'Comparte un secreto que desaparece tras leerlo.',
		'Teilen Sie ein Geheimnis, das nach dem Lesen verschwindet.'
	],
	'shell.facile': ['Facile Studio', 'Facile Studio', 'Facile Studio', 'Facile Studio']
} satisfies Record<string, Row>;
