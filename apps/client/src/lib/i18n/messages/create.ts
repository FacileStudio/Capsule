import type { Row } from '../row';

/** Strings of the page that seals a capsule, before and after sealing. */
export const create = {
	'create.title': [
		'Capsule — Ephemeral Encrypted Sharing',
		'Capsule — Partage chiffré et éphémère',
		'Capsule — Compartir cifrado y efímero',
		'Capsule — Verschlüsselt und flüchtig teilen'
	],
	'create.description': [
		'Share secrets that self-destruct. End-to-end encrypted, zero knowledge.',
		'Partagez des secrets qui se détruisent. Chiffrement de bout en bout, zéro connaissance.',
		'Comparte secretos que se autodestruyen. Cifrado de extremo a extremo, conocimiento cero.',
		'Teilen Sie Geheimnisse, die sich selbst zerstören. Ende-zu-Ende-verschlüsselt, Zero Knowledge.'
	],
	'create.heading': ['Seal a capsule', 'Sceller une capsule', 'Sellar una cápsula', 'Kapsel versiegeln'],
	'create.intro': [
		'Your content is encrypted in the browser. The server never sees the plaintext.',
		'Votre contenu est chiffré dans le navigateur. Le serveur ne voit jamais le texte en clair.',
		'Tu contenido se cifra en el navegador. El servidor nunca ve el texto sin cifrar.',
		'Ihr Inhalt wird im Browser verschlüsselt. Der Server sieht den Klartext nie.'
	],
	'create.placeholder': [
		'Paste your secret here...',
		'Collez votre secret ici...',
		'Pega tu secreto aquí...',
		'Fügen Sie Ihr Geheimnis hier ein ...'
	],
	'create.expires': ['Expires in', 'Expire dans', 'Caduca en', 'Läuft ab in'],
	'create.syntax': ['Syntax', 'Syntaxe', 'Sintaxis', 'Syntax'],
	'create.burn': [
		'Burn after opening',
		'Détruire après ouverture',
		'Destruir tras abrir',
		'Nach dem Öffnen löschen'
	],
	'create.password': [
		'Password protect',
		'Protéger par mot de passe',
		'Proteger con contraseña',
		'Mit Passwort schützen'
	],
	'create.passwordPlaceholder': [
		'Enter a password...',
		'Saisissez un mot de passe...',
		'Introduce una contraseña...',
		'Passwort eingeben ...'
	],
	'create.seal': ['Seal capsule', 'Sceller la capsule', 'Sellar cápsula', 'Kapsel versiegeln'],
	'create.sealing': ['Sealing...', 'Scellement...', 'Sellando...', 'Wird versiegelt ...'],
	'create.sealingLabel': ['Sealing', 'Scellement', 'Sellando', 'Versiegeln'],
	'create.sealFailed': [
		'Failed to seal capsule',
		'Impossible de sceller la capsule',
		'No se pudo sellar la cápsula',
		'Kapsel konnte nicht versiegelt werden'
	],
	'expiry.1h': ['1 hour', '1 heure', '1 hora', '1 Stunde'],
	'expiry.24h': ['24 hours', '24 heures', '24 horas', '24 Stunden'],
	'expiry.7d': ['7 days', '7 jours', '7 días', '7 Tage'],
	'expiry.30d': ['30 days', '30 jours', '30 días', '30 Tage'],
	'sealed.heading': ['Capsule sealed', 'Capsule scellée', 'Cápsula sellada', 'Kapsel versiegelt'],
	'sealed.burnOnce': [
		'Share this link — the capsule can only be opened once.',
		'Partagez ce lien : la capsule ne peut être ouverte qu’une fois.',
		'Comparte este enlace: la cápsula solo se puede abrir una vez.',
		'Teilen Sie diesen Link: Die Kapsel lässt sich nur einmal öffnen.'
	],
	'sealed.burnOncePassword': [
		'Share this link and the password — the capsule can only be opened once.',
		'Partagez ce lien et le mot de passe : la capsule ne peut être ouverte qu’une fois.',
		'Comparte este enlace y la contraseña: la cápsula solo se puede abrir una vez.',
		'Teilen Sie diesen Link und das Passwort: Die Kapsel lässt sich nur einmal öffnen.'
	],
	'sealed.expires': [
		'Share this link — it expires in {duration}.',
		'Partagez ce lien : il expire dans {duration}.',
		'Comparte este enlace: caduca en {duration}.',
		'Teilen Sie diesen Link: Er läuft in {duration} ab.'
	],
	'sealed.expiresPassword': [
		'Share this link and the password — it expires in {duration}.',
		'Partagez ce lien et le mot de passe : il expire dans {duration}.',
		'Comparte este enlace y la contraseña: caduca en {duration}.',
		'Teilen Sie diesen Link und das Passwort: Er läuft in {duration} ab.'
	],
	'sealed.copy': ['Copy link', 'Copier le lien', 'Copiar enlace', 'Link kopieren'],
	'sealed.copied': ['Copied', 'Copié', 'Copiado', 'Kopiert'],
	'sealed.revoked': ['Capsule revoked', 'Capsule révoquée', 'Cápsula revocada', 'Kapsel widerrufen'],
	'sealed.token': ['Delete token', 'Jeton de suppression', 'Token de eliminación', 'Löschtoken'],
	'sealed.tokenHelp': [
		'Copy it now — it is not stored anywhere you can read it back.',
		'Copiez-le maintenant : il n’est stocké nulle part où vous pourriez le relire.',
		'Cópialo ahora: no se guarda en ningún sitio donde puedas volver a leerlo.',
		'Kopieren Sie ihn jetzt: Er wird nirgends gespeichert, wo Sie ihn erneut lesen könnten.'
	],
	'sealed.revoke': ['Revoke', 'Révoquer', 'Revocar', 'Widerrufen'],
	'sealed.revoking': ['Revoking...', 'Révocation...', 'Revocando...', 'Wird widerrufen ...'],
	'sealed.another': [
		'Seal another capsule',
		'Sceller une autre capsule',
		'Sellar otra cápsula',
		'Weitere Kapsel versiegeln'
	]
} satisfies Record<string, Row>;
