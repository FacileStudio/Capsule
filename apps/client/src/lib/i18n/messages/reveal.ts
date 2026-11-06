import type { Row } from '../row';

/** Strings of the page that opens a capsule from its link. */
export const reveal = {
	'reveal.title': ['Capsule — Open', 'Capsule — Ouvrir', 'Capsule — Abrir', 'Capsule — Öffnen'],
	'reveal.description': [
		'Someone sent you an encrypted capsule.',
		'Quelqu’un vous a envoyé une capsule chiffrée.',
		'Alguien te ha enviado una cápsula cifrada.',
		'Jemand hat Ihnen eine verschlüsselte Kapsel geschickt.'
	],
	'reveal.noId': [
		'No capsule id in the URL.',
		'Aucun identifiant de capsule dans l’URL.',
		'No hay identificador de cápsula en la URL.',
		'Keine Kapsel-ID in der URL.'
	],
	'reveal.noKey': [
		'No decryption key found in URL. The link may be incomplete.',
		'Aucune clé de déchiffrement dans l’URL. Le lien est peut-être incomplet.',
		'No se encontró ninguna clave de descifrado en la URL. Puede que el enlace esté incompleto.',
		'Kein Entschlüsselungsschlüssel in der URL. Der Link ist möglicherweise unvollständig.'
	],
	'reveal.badFragment': [
		'Invalid URL fragment',
		'Fragment d’URL invalide',
		'Fragmento de URL no válido',
		'Ungültiges URL-Fragment'
	],
	'reveal.wrongPassword': [
		'Wrong password or corrupted link.',
		'Mot de passe erroné ou lien corrompu.',
		'Contraseña incorrecta o enlace dañado.',
		'Falsches Passwort oder beschädigter Link.'
	],
	'reveal.decryptFailed': [
		'Failed to decrypt. The key may be wrong or the data may be corrupted.',
		'Échec du déchiffrement. La clé est peut-être erronée ou les données corrompues.',
		'No se pudo descifrar. Puede que la clave sea incorrecta o que los datos estén dañados.',
		'Entschlüsselung fehlgeschlagen. Der Schlüssel ist möglicherweise falsch oder die Daten sind beschädigt.'
	],
	'reveal.sealedAgo': ['Sealed {when}', 'Scellée {when}', 'Sellada {when}', 'Versiegelt {when}'],
	'reveal.badgePassword': [
		'Password protected',
		'Protégée par mot de passe',
		'Protegida con contraseña',
		'Passwortgeschützt'
	],
	'reveal.badgeBurn': [
		'Burns after opening',
		'Se détruit à l’ouverture',
		'Se destruye al abrir',
		'Wird nach dem Öffnen gelöscht'
	],
	'reveal.emptyHeading': [
		'This capsule is empty',
		'Cette capsule est vide',
		'Esta cápsula está vacía',
		'Diese Kapsel ist leer'
	],
	'reveal.emptyText': [
		'It was already opened, expired, or never existed.',
		'Elle a déjà été ouverte, a expiré ou n’a jamais existé.',
		'Ya se abrió, caducó o nunca existió.',
		'Sie wurde bereits geöffnet, ist abgelaufen oder hat nie existiert.'
	],
	'reveal.sealNew': [
		'Seal a new capsule',
		'Sceller une nouvelle capsule',
		'Sellar una cápsula nueva',
		'Neue Kapsel versiegeln'
	],
	'reveal.errorHeading': [
		'Something went wrong',
		'Une erreur est survenue',
		'Algo salió mal',
		'Etwas ist schiefgelaufen'
	],
	'reveal.back': ['Back to Capsule', 'Retour à Capsule', 'Volver a Capsule', 'Zurück zu Capsule'],
	'reveal.passwordHeading': [
		'Password required',
		'Mot de passe requis',
		'Contraseña necesaria',
		'Passwort erforderlich'
	],
	'reveal.passwordText': [
		'This capsule is password-protected. Enter the password to decrypt it.',
		'Cette capsule est protégée par mot de passe. Saisissez-le pour la déchiffrer.',
		'Esta cápsula está protegida con contraseña. Introdúcela para descifrarla.',
		'Diese Kapsel ist passwortgeschützt. Geben Sie das Passwort zum Entschlüsseln ein.'
	],
	'reveal.passwordPlaceholder': [
		'Enter password...',
		'Saisissez le mot de passe...',
		'Introduce la contraseña...',
		'Passwort eingeben ...'
	],
	'reveal.unlock': ['Unlock capsule', 'Déverrouiller la capsule', 'Desbloquear cápsula', 'Kapsel entsperren'],
	'reveal.readyHeading': [
		'A capsule for you',
		'Une capsule pour vous',
		'Una cápsula para ti',
		'Eine Kapsel für Sie'
	],
	'reveal.readyBurn': [
		'This capsule will self-destruct after you open it.',
		'Cette capsule se détruira après son ouverture.',
		'Esta cápsula se autodestruirá después de abrirla.',
		'Diese Kapsel zerstört sich nach dem Öffnen selbst.'
	],
	'reveal.open': ['Open capsule', 'Ouvrir la capsule', 'Abrir cápsula', 'Kapsel öffnen'],
	'reveal.opening': ['Opening...', 'Ouverture...', 'Abriendo...', 'Wird geöffnet ...'],
	'reveal.openingLabel': ['Opening', 'Ouverture', 'Abriendo', 'Öffnen'],
	'reveal.openedHeading': ['Capsule opened', 'Capsule ouverte', 'Cápsula abierta', 'Kapsel geöffnet'],
	'reveal.destroyed': [
		'Content destroyed from server. This is the only time you’ll see it.',
		'Contenu détruit sur le serveur. C’est la seule fois où vous le verrez.',
		'Contenido destruido en el servidor. Esta es la única vez que lo verás.',
		'Inhalt auf dem Server gelöscht. Sie sehen ihn nur dieses eine Mal.'
	],
	'reveal.copy': ['Copy', 'Copier', 'Copiar', 'Kopieren'],
	'reveal.copied': ['Copied', 'Copié', 'Copiado', 'Kopiert']
} satisfies Record<string, Row>;
