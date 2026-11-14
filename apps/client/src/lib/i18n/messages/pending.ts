import type { Row } from '../row';

/** The page where the person confirms what an AI client asked Capsule to do. */
export const pending = {
	'pending.title': ['Pending actions · Capsule', 'Actions en attente · Capsule', 'Acciones pendientes · Capsule', 'Ausstehende Aktionen · Capsule'],
	'pending.heading': ['Pending actions', 'Actions en attente', 'Acciones pendientes', 'Ausstehende Aktionen'],
	'pending.intro': [
		'An AI client asked Capsule to do these. Nothing happens until you confirm.',
		'Un client IA a demandé ces actions à Capsule. Rien ne se passe tant que vous ne confirmez pas.',
		'Un cliente de IA pidió estas acciones a Capsule. No ocurre nada hasta que confirmes.',
		'Ein KI-Client hat Capsule um diese Aktionen gebeten. Nichts passiert, bevor Sie bestätigen.'
	],
	'pending.signInHeading': ['Sign in to continue', 'Connectez-vous pour continuer', 'Inicia sesión para continuar', 'Melden Sie sich an, um fortzufahren'],
	'pending.signInBody': [
		'Only you can confirm what an AI client asked for. Sign in with your Facile account and you will come back here.',
		'Vous seul pouvez confirmer ce qu’un client IA a demandé. Connectez-vous avec votre compte Facile pour revenir ici.',
		'Solo tú puedes confirmar lo que pidió un cliente de IA. Inicia sesión con tu cuenta Facile y volverás aquí.',
		'Nur Sie können bestätigen, worum ein KI-Client gebeten hat. Melden Sie sich mit Ihrem Facile-Konto an, dann kommen Sie hierher zurück.'
	],
	'pending.signIn': ['Sign in', 'Se connecter', 'Iniciar sesión', 'Anmelden'],
	'pending.signInAgain': ['Sign in again to confirm', 'Reconnectez-vous pour confirmer', 'Vuelve a iniciar sesión para confirmar', 'Melden Sie sich erneut an, um zu bestätigen'],
	'pending.signInAgainBody': [
		'Only a session you opened in this browser can confirm. Sign in again here, then confirm. The command line and the mobile or desktop apps cannot confirm.',
		'Seule une session ouverte dans ce navigateur peut confirmer. Reconnectez-vous ici, puis confirmez. La ligne de commande et les applications mobile ou de bureau ne peuvent pas confirmer.',
		'Solo una sesión abierta en este navegador puede confirmar. Vuelve a iniciar sesión aquí y luego confirma. La línea de comandos y las aplicaciones móviles o de escritorio no pueden confirmar.',
		'Nur eine in diesem Browser geöffnete Sitzung kann bestätigen. Melden Sie sich hier erneut an und bestätigen Sie dann. Die Kommandozeile und die Mobil- oder Desktop-Apps können nicht bestätigen.'
	],
	'pending.signInAgainButton': ['Sign in again', 'Se reconnecter', 'Iniciar sesión de nuevo', 'Erneut anmelden'],
	'pending.none': ['Nothing is waiting', 'Rien en attente', 'Nada pendiente', 'Nichts ausstehend'],
	'pending.noneBody': [
		'Actions an AI client asks for appear here until you confirm or dismiss them.',
		'Les actions demandées par un client IA apparaissent ici jusqu’à ce que vous les confirmiez ou les ignoriez.',
		'Las acciones que pide un cliente de IA aparecen aquí hasta que las confirmes o descartes.',
		'Aktionen, um die ein KI-Client bittet, erscheinen hier, bis Sie sie bestätigen oder verwerfen.'
	],
	'pending.kind.revoke': ['Revoke capsule', 'Révoquer la capsule', 'Revocar cápsula', 'Kapsel widerrufen'],
	'pending.askedBy': ['Asked by {client}', 'Demandé par {client}', 'Pedido por {client}', 'Angefragt von {client}'],
	'pending.sealed': ['Sealed {date}', 'Scellée le {date}', 'Sellada el {date}', 'Versiegelt am {date}'],
	'pending.expires': ['Expires {date}', 'Expire le {date}', 'Caduca el {date}', 'Läuft ab am {date}'],
	'pending.burns': ['Burns after reading', 'Détruite à la lecture', 'Se destruye al leerla', 'Löscht sich nach dem Lesen'],
	'pending.password': ['Password', 'Mot de passe', 'Contraseña', 'Passwort'],
	'pending.revokeNote': [
		'The capsule’s content is destroyed for good and its link stops working.',
		'Le contenu de la capsule est détruit définitivement et son lien ne fonctionne plus.',
		'El contenido de la cápsula se destruye para siempre y su enlace deja de funcionar.',
		'Der Inhalt der Kapsel wird endgültig vernichtet und ihr Link funktioniert nicht mehr.'
	],
	'pending.confirmRevoke': ['Revoke the capsule', 'Révoquer la capsule', 'Revocar la cápsula', 'Kapsel widerrufen'],
	'pending.dismiss': ['Dismiss', 'Ignorer', 'Descartar', 'Verwerfen'],
	'pending.changed': [
		'This capsule changed after the request. Ask the AI client again.',
		'Cette capsule a changé après la demande. Redemandez au client IA.',
		'Esta cápsula cambió después de la solicitud. Vuelve a pedírselo al cliente de IA.',
		'Diese Kapsel wurde nach der Anfrage geändert. Fragen Sie den KI-Client erneut.'
	],
	'pending.missing': [
		'This capsule was already opened, revoked or expired.',
		'Cette capsule a déjà été ouverte, révoquée ou a expiré.',
		'Esta cápsula ya se abrió, se revocó o caducó.',
		'Diese Kapsel wurde bereits geöffnet, widerrufen oder ist abgelaufen.'
	],
	'pending.revoked': ['Capsule revoked', 'Capsule révoquée', 'Cápsula revocada', 'Kapsel widerrufen'],
	'pending.dismissed': ['Request dismissed', 'Demande ignorée', 'Solicitud descartada', 'Anfrage verworfen'],
	'pending.loadFailed': [
		'Could not load the pending actions.',
		'Impossible de charger les actions en attente.',
		'No se pudieron cargar las acciones pendientes.',
		'Die ausstehenden Aktionen konnten nicht geladen werden.'
	],
	'pending.failed': [
		'That did not work. Try again.',
		'Cela n’a pas fonctionné. Réessayez.',
		'No funcionó. Inténtalo de nuevo.',
		'Das hat nicht geklappt. Versuchen Sie es erneut.'
	]
} satisfies Record<string, Row>;
