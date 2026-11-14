import type { Row } from '../row';

/** Strings of the landing blocks below the demos: features, suite, pricing, developers and the closing call. */
export const blocks = {
	'blocks.features.title': [
		'Encrypted before it leaves',
		'Chiffré avant de partir',
		'Cifrado antes de salir',
		'Verschlüsselt, bevor es losgeht'
	],
	'blocks.features.line': [
		'Your device encrypts the secret. The server only keeps text nobody can read.',
		'Votre appareil chiffre le secret. Le serveur ne garde qu’un texte illisible.',
		'Tu dispositivo cifra el secreto. El servidor solo guarda un texto ilegible.',
		'Ihr Gerät verschlüsselt das Geheimnis. Der Server behält nur unlesbaren Text.'
	],
	'blocks.features.device.title': ['On your device', 'Sur votre appareil', 'En tu dispositivo', 'Auf Ihrem Gerät'],
	'blocks.features.device.text': [
		'The key stays in the link and never reaches the server.',
		'La clé reste dans le lien et n’arrive jamais sur le serveur.',
		'La clave va en el enlace y nunca llega al servidor.',
		'Der Schlüssel bleibt im Link und erreicht nie den Server.'
	],
	'blocks.features.expiry.title': ['You pick the expiry', 'Expiration au choix', 'Tú eliges la caducidad', 'Ablauf nach Wahl'],
	'blocks.features.expiry.text': [
		'One hour, one day, one week or one month.',
		'Une heure, un jour, une semaine ou un mois.',
		'Una hora, un día, una semana o un mes.',
		'Eine Stunde, ein Tag, eine Woche oder ein Monat.'
	],
	'blocks.features.revoke.title': ['Revoke at any time', 'Révocable à tout moment', 'Revocable en cualquier momento', 'Jederzeit widerrufbar'],
	'blocks.features.revoke.text': [
		'Keep the delete token and erase the capsule before anyone reads it.',
		'Gardez le jeton de suppression et effacez la capsule avant sa lecture.',
		'Guarda el token de eliminación y borra la cápsula antes de que alguien la lea.',
		'Behalten Sie den Löschtoken und löschen Sie die Kapsel, bevor jemand sie liest.'
	],
	'blocks.features.account.title': ['No account needed', 'Aucun compte nécessaire', 'Sin cuenta', 'Kein Konto nötig'],
	'blocks.features.account.text': [
		'Not for you, not for the person who opens the link.',
		'Ni pour vous, ni pour la personne qui ouvre le lien.',
		'Ni para ti ni para quien abre el enlace.',
		'Weder für Sie noch für die Person, die den Link öffnet.'
	],
	'blocks.with.bulle': [
		'Checks whether a capsule was opened and revokes it.',
		'Vérifie si une capsule a été lue et la révoque.',
		'Comprueba si una cápsula se abrió y la revoca.',
		'Prüft, ob eine Kapsel geöffnet wurde, und widerruft sie.'
	],
	'blocks.with.line': [
		'Bulle never sees the content and asks before it revokes. Bulle needs a Pro or Studio plan.',
		'Bulle ne voit jamais le contenu et demande votre accord avant de révoquer. Bulle demande une offre Pro ou Studio.',
		'Bulle nunca ve el contenido y te pide permiso antes de revocar. Bulle requiere un plan Pro o Studio.',
		'Bulle sieht den Inhalt nie und fragt vor dem Widerruf nach. Bulle braucht den Tarif Pro oder Studio.'
	],
	'blocks.pricing.title': ['Free, and in every plan', 'Gratuit, et dans chaque offre', 'Gratis, y en cada plan', 'Kostenlos und in jedem Tarif'],
	'blocks.pricing.free': [
		'Free, no account needed. Capsule is part of the Facile suite.',
		'Gratuit, sans compte. Capsule fait partie de la suite Facile.',
		'Gratis y sin cuenta. Capsule forma parte de la suite Facile.',
		'Kostenlos, ohne Konto. Capsule gehört zur Facile-Suite.'
	],
	'blocks.pricing.paid': [
		'Capsule, plus Bulle to check and revoke your capsules',
		'Capsule, et Bulle pour vérifier et révoquer vos capsules',
		'Capsule, y Bulle para comprobar y revocar tus cápsulas',
		'Capsule, dazu Bulle zum Prüfen und Widerrufen Ihrer Kapseln'
	],
	'blocks.dev.line': [
		'A CLI and an API to seal and revoke from your scripts. Open source, and you can run it yourself.',
		'Un CLI et une API pour sceller et révoquer depuis vos scripts. Open source, et installable chez vous.',
		'Un CLI y una API para sellar y revocar desde tus scripts. Código abierto, y puedes instalarlo tú mismo.',
		'Eine CLI und eine API zum Versiegeln und Widerrufen aus Skripten. Open Source, auch selbst betreibbar.'
	],
	'blocks.dev.cli': [
		'Encrypts on your machine and prints the link. Every platform on the downloads page.',
		'Chiffre sur votre machine et affiche le lien. Toutes les plateformes sur la page de téléchargement.',
		'Cifra en tu equipo e imprime el enlace. Todas las plataformas en la página de descargas.',
		'Verschlüsselt auf Ihrem Rechner und gibt den Link aus. Alle Plattformen auf der Download-Seite.'
	],
	'blocks.dev.api': [
		'Seal, read the status and revoke over a documented REST API.',
		'Scellez, lisez l’état et révoquez avec une API REST documentée.',
		'Sella, consulta el estado y revoca con una API REST documentada.',
		'Versiegeln, Status lesen und widerrufen über eine dokumentierte REST-API.'
	],
	'blocks.dev.mcp': [
		'Agents check and revoke a capsule at /mcp, and never read its content.',
		'Les agents vérifient et révoquent une capsule sur /mcp, sans jamais lire son contenu.',
		'Los agentes comprueban y revocan una cápsula en /mcp, sin leer nunca su contenido.',
		'Agenten prüfen und widerrufen eine Kapsel über /mcp und lesen nie den Inhalt.'
	],
	'blocks.dev.source': [
		'The code is open. The README explains how to run Capsule yourself.',
		'Le code est ouvert. Le README explique comment installer Capsule chez vous.',
		'El código es abierto. El README explica cómo instalar Capsule tú mismo.',
		'Der Code ist offen. Das README erklärt, wie Sie Capsule selbst betreiben.'
	],
	'blocks.cta.title': [
		'One Facile account, the whole suite',
		'Un compte Facile, toute la suite',
		'Una cuenta Facile, toda la suite',
		'Ein Facile-Konto, die ganze Suite'
	],
	'blocks.cta.body': [
		'Capsule works without an account. A free Facile account opens the other apps: files, calendar, quotes, signatures.',
		'Capsule marche sans compte. Un compte Facile gratuit ouvre les autres apps : fichiers, agenda, devis, signatures.',
		'Capsule funciona sin cuenta. Una cuenta Facile gratuita abre las demás apps: archivos, agenda, presupuestos, firmas.',
		'Capsule funktioniert ohne Konto. Ein kostenloses Facile-Konto öffnet die anderen Apps: Dateien, Kalender, Angebote, Signaturen.'
	],
	'blocks.cta.seal': ['Seal a secret', 'Sceller un secret', 'Sellar un secreto', 'Geheimnis versiegeln'],
	'blocks.footer.developers': ['Developers', 'Développeurs', 'Desarrolladores', 'Entwickler'],
	'blocks.footer.suite': ['Facile suite', 'Suite Facile', 'Suite Facile', 'Facile-Suite']
} satisfies Record<string, Row>;
