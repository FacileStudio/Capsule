import type { Row } from '../row';

/** Strings of the landing hero, the how-it-works demo and the three demo tiles. */
export const landing = {
	'landing.heroTitle': [
		'Send a password. It’s gone once read.',
		'Envoyez un mot de passe, il disparaît après lecture.',
		'Envía una contraseña y desaparece tras leerla.',
		'Passwort senden, nach dem Lesen ist es weg.'
	],
	'landing.heroBody': [
		'A link that opens once, encrypted before it leaves your device. No account needed.',
		'Un lien qui ne s’ouvre qu’une fois, chiffré avant même de quitter votre appareil. Aucun compte nécessaire.',
		'Un enlace que solo se abre una vez, cifrado antes de salir de tu dispositivo. Sin cuenta.',
		'Ein Link, der sich nur einmal öffnet, verschlüsselt, bevor er Ihr Gerät verlässt. Kein Konto nötig.'
	],
	'landing.options': ['Options', 'Options', 'Opciones', 'Optionen'],
	'landing.nav.cli': ['CLI', 'CLI', 'CLI', 'CLI'],
	'downloads.title': ['Downloads · Capsule', 'Téléchargements · Capsule', 'Descargas · Capsule', 'Downloads · Capsule'],
	'downloads.heading': ['Get the capsule CLI', 'Obtenir la CLI capsule', 'Obtén la CLI capsule', 'Die capsule-CLI holen'],
	'landing.downloadBody': [
		'The capsule CLI encrypts on your machine and prints the link.',
		'La CLI capsule chiffre sur votre machine et affiche le lien.',
		'La CLI capsule cifra en tu equipo e imprime el enlace.',
		'Die capsule-CLI verschlüsselt auf Ihrem Rechner und gibt den Link aus.'
	],
	'landing.script': ['Install script', 'Script d’installation', 'Script de instalación', 'Installationsskript'],
	'landing.demoTitle':['How it works', 'Comment ça marche', 'Cómo funciona', 'So funktioniert es'],
	'landing.demoLabel': [
		'A Wi-Fi password is sealed, its link copied, then the guest opens it once and it is gone.',
		'Un code Wi-Fi est scellé, son lien copié, puis l’invité l’ouvre une fois et il disparaît.',
		'Se sella una clave Wi-Fi, se copia su enlace y el invitado la abre una vez antes de que desaparezca.',
		'Ein WLAN-Passwort wird versiegelt, der Link kopiert, der Gast öffnet es einmal und es ist weg.'
	],
	'landing.demoSecret': ['Wi-Fi: linden-sun-42', 'Wi-Fi : tilleul-soleil-42', 'Wi-Fi: tilo-sol-42', 'WLAN: linde-sonne-42'],
	'landing.tilesTitle': [
		'For anything that should not linger',
		'Pour tout ce qui ne doit pas traîner',
		'Para todo lo que no debe quedarse',
		'Für alles, was nicht liegen bleiben soll'
	],
	'landing.tile.share.title': [
		'A password, a Wi-Fi code',
		'Un mot de passe, un code Wi-Fi',
		'Una contraseña, una clave Wi-Fi',
		'Ein Passwort, ein WLAN-Code'
	],
	'landing.tile.share.body': [
		'A password for a colleague, a Wi-Fi code for a guest. They need no account either.',
		'Un mot de passe pour un collègue, un code Wi-Fi pour un invité. Eux non plus n’ont pas besoin de compte.',
		'Una contraseña para un compañero, una clave Wi-Fi para un invitado. Tampoco necesitan cuenta.',
		'Ein Passwort für eine Kollegin, ein WLAN-Code für einen Gast. Auch sie brauchen kein Konto.'
	],
	'landing.tile.share.label': [
		'Two links are sent: a password to a colleague and a Wi-Fi code to a guest.',
		'Deux liens partent : un mot de passe pour une collègue et un code Wi-Fi pour un invité.',
		'Salen dos enlaces: una contraseña para una compañera y una clave Wi-Fi para un invitado.',
		'Zwei Links gehen raus: ein Passwort an eine Kollegin und ein WLAN-Code an einen Gast.'
	],
	'landing.tile.share.lea': ['Léa, colleague', 'Léa, collègue', 'Léa, compañera', 'Léa, Kollegin'],
	'landing.tile.share.leaWhat': [
		'Shared account password',
		'Mot de passe du compte partagé',
		'Contraseña de la cuenta compartida',
		'Passwort des gemeinsamen Kontos'
	],
	'landing.tile.share.guest': ['Guest', 'Invité', 'Invitado', 'Gast'],
	'landing.tile.share.guestWhat': ['Wi-Fi code', 'Code Wi-Fi', 'Clave Wi-Fi', 'WLAN-Code'],
	'landing.tile.once.title': [
		'Opens once, then it is gone',
		'S’ouvre une fois, puis disparaît',
		'Se abre una vez y desaparece',
		'Öffnet sich einmal, dann ist es weg'
	],
	'landing.tile.once.body': [
		'After the first read the link leads nowhere. Otherwise it expires on the date you chose.',
		'Après la première lecture, le lien ne mène plus à rien. Sinon, il expire à la date choisie.',
		'Tras la primera lectura el enlace ya no lleva a nada. Si no, caduca en la fecha que elijas.',
		'Nach dem ersten Lesen führt der Link ins Leere. Sonst läuft er zum gewählten Termin ab.'
	],
	'landing.tile.once.label': [
		'The capsule is opened once, and a second visit finds it empty.',
		'La capsule est ouverte une fois, et une seconde visite la trouve vide.',
		'La cápsula se abre una vez y una segunda visita la encuentra vacía.',
		'Die Kapsel wird einmal geöffnet, ein zweiter Besuch findet sie leer.'
	],
	'landing.tile.extras.title': [
		'Add a password or a file',
		'Ajoutez un mot de passe ou un fichier',
		'Añade una contraseña o un archivo',
		'Passwort oder Datei hinzufügen'
	],
	'landing.tile.extras.body': [
		'A password you give in person, or a file up to 10 MB, encrypted too.',
		'Un mot de passe à donner de vive voix, ou un fichier jusqu’à 10 Mo, chiffré lui aussi.',
		'Una contraseña que das en persona, o un archivo de hasta 10 MB, también cifrado.',
		'Ein Passwort, das Sie persönlich weitergeben, oder eine Datei bis 10 MB, ebenfalls verschlüsselt.'
	],
	'landing.tile.extras.label': [
		'A password is set on the capsule and a PDF is attached.',
		'Un mot de passe protège la capsule et un PDF y est joint.',
		'Se pone una contraseña a la cápsula y se adjunta un PDF.',
		'Die Kapsel bekommt ein Passwort und ein PDF wird angehängt.'
	]
} satisfies Record<string, Row>;
