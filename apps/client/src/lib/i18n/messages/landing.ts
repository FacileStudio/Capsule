import type { Row } from '../row';

/** Strings of the public landing around the seal form. */
export const landing = {
	'landing.heroTitle': [
		'Share a secret that burns after reading',
		'Partagez un secret qui se détruit à la lecture',
		'Comparte un secreto que se destruye al leerlo',
		'Teilen Sie ein Geheimnis, das sich nach dem Lesen löscht'
	],
	'landing.heroBody': [
		'Encrypted in your browser. The server never sees the content.',
		'Chiffré dans votre navigateur. Le serveur ne voit jamais le contenu.',
		'Cifrado en tu navegador. El servidor nunca ve el contenido.',
		'Im Browser verschlüsselt. Der Server sieht den Inhalt nie.'
	],
	'landing.cli': ['Install the CLI', 'Installer la CLI', 'Instalar la CLI', 'CLI installieren'],
	'landing.nav.cli': ['CLI', 'CLI', 'CLI', 'CLI'],
	'landing.demoTitle': ['How it works', 'Comment ça marche', 'Cómo funciona', 'So funktioniert es'],
	'landing.demoLabel': [
		'A password is sealed, its link copied, then the recipient opens it once and it is gone.',
		'Un mot de passe est scellé, son lien copié, puis le destinataire l’ouvre une fois et il disparaît.',
		'Se sella una contraseña, se copia su enlace y el destinatario la abre una vez antes de que desaparezca.',
		'Ein Passwort wird versiegelt, der Link kopiert, der Empfänger öffnet es einmal und es ist weg.'
	],
	'landing.encrypted.title': [
		'Encrypted in your browser',
		'Chiffré dans votre navigateur',
		'Cifrado en tu navegador',
		'Im Browser verschlüsselt'
	],
	'landing.encrypted.body': [
		'AES-256-GCM runs on your device. The key stays in the link and never reaches the server.',
		'AES-256-GCM tourne sur votre appareil. La clé reste dans le lien et n’atteint jamais le serveur.',
		'AES-256-GCM se ejecuta en tu dispositivo. La clave va en el enlace y nunca llega al servidor.',
		'AES-256-GCM läuft auf Ihrem Gerät. Der Schlüssel bleibt im Link und erreicht nie den Server.'
	],
	'landing.once.title': ['Opens once', 'S’ouvre une fois', 'Se abre una vez', 'Öffnet sich einmal'],
	'landing.once.body': [
		'A capsule burns on first read, or expires when you choose.',
		'Une capsule se détruit à la première lecture, ou expire quand vous le décidez.',
		'Una cápsula se destruye al primer uso, o caduca cuando tú decidas.',
		'Eine Kapsel löscht sich beim ersten Lesen oder läuft ab, wann Sie wollen.'
	],
	'landing.extras.title': [
		'Passwords and files',
		'Mots de passe et fichiers',
		'Contraseñas y archivos',
		'Passwörter und Dateien'
	],
	'landing.extras.body': [
		'Add a password, or attach a file up to 10 MB.',
		'Ajoutez un mot de passe, ou joignez un fichier jusqu’à 10 Mo.',
		'Añade una contraseña o adjunta un archivo de hasta 10 MB.',
		'Fügen Sie ein Passwort hinzu oder hängen Sie eine Datei bis 10 MB an.'
	],
	'landing.downloadTitle': [
		'Seal from the terminal',
		'Scellez depuis le terminal',
		'Sella desde la terminal',
		'Aus dem Terminal versiegeln'
	],
	'landing.downloadBody': [
		'The capsule CLI encrypts on your machine and prints the link.',
		'La CLI capsule chiffre sur votre machine et affiche le lien.',
		'La CLI capsule cifra en tu equipo e imprime el enlace.',
		'Die capsule-CLI verschlüsselt auf Ihrem Rechner und gibt den Link aus.'
	],
	'landing.script': ['Install script', 'Script d’installation', 'Script de instalación', 'Installationsskript'],
	'landing.ctaTitle': [
		'Seal your first capsule',
		'Scellez votre première capsule',
		'Sella tu primera cápsula',
		'Versiegeln Sie Ihre erste Kapsel'
	],
	'landing.ctaBody': [
		'No account needed.',
		'Aucun compte requis.',
		'No necesitas cuenta.',
		'Kein Konto nötig.'
	]
} satisfies Record<string, Row>;
