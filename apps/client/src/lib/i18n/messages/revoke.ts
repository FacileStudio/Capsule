import type { Row } from '../row';

/** Strings of the page that revokes a capsule with its delete token. */
export const revoke = {
	'revoke.title': [
		'Revoke — Capsule',
		'Révoquer — Capsule',
		'Revocar — Capsule',
		'Widerrufen — Capsule'
	],
	'revoke.failed': [
		'Failed to revoke capsule',
		'Impossible de révoquer la capsule',
		'No se pudo revocar la cápsula',
		'Kapsel konnte nicht widerrufen werden'
	],
	'revoke.doneHeading': ['Capsule revoked', 'Capsule révoquée', 'Cápsula revocada', 'Kapsel widerrufen'],
	'revoke.doneText': [
		'The capsule has been burned. Its content is permanently gone.',
		'La capsule a été détruite. Son contenu a définitivement disparu.',
		'La cápsula se ha destruido. Su contenido ha desaparecido para siempre.',
		'Die Kapsel wurde vernichtet. Ihr Inhalt ist endgültig gelöscht.'
	],
	'revoke.another': ['Revoke another', 'En révoquer une autre', 'Revocar otra', 'Weitere widerrufen'],
	'revoke.seal': ['Seal a capsule', 'Sceller une capsule', 'Sellar una cápsula', 'Kapsel versiegeln'],
	'revoke.heading': [
		'Revoke a capsule',
		'Révoquer une capsule',
		'Revocar una cápsula',
		'Kapsel widerrufen'
	],
	'revoke.intro': [
		'Burn a capsule before it’s opened. You’ll need the capsule URL or ID and the delete token you received when sealing.',
		'Détruisez une capsule avant son ouverture. Il vous faut l’URL ou l’identifiant de la capsule et le jeton de suppression reçu au scellement.',
		'Destruye una cápsula antes de que se abra. Necesitas la URL o el identificador de la cápsula y el token de eliminación que recibiste al sellarla.',
		'Vernichten Sie eine Kapsel, bevor sie geöffnet wird. Sie benötigen die URL oder ID der Kapsel und den Löschtoken, den Sie beim Versiegeln erhalten haben.'
	],
	'revoke.urlLabel': [
		'Capsule URL or ID',
		'URL ou identifiant de la capsule',
		'URL o identificador de la cápsula',
		'URL oder ID der Kapsel'
	],
	'revoke.tokenLabel': ['Delete token', 'Jeton de suppression', 'Token de eliminación', 'Löschtoken'],
	'revoke.tokenPlaceholder': [
		'Paste your delete token here...',
		'Collez votre jeton de suppression ici...',
		'Pega tu token de eliminación aquí...',
		'Fügen Sie Ihren Löschtoken hier ein ...'
	],
	'revoke.submit': ['Revoke capsule', 'Révoquer la capsule', 'Revocar cápsula', 'Kapsel widerrufen'],
	'revoke.revoking': ['Revoking...', 'Révocation...', 'Revocando...', 'Wird widerrufen ...'],
	'revoke.revokingLabel': ['Revoking', 'Révocation', 'Revocando', 'Widerrufen']
} satisfies Record<string, Row>;
