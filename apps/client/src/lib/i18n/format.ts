const UNITS: [Intl.RelativeTimeFormatUnit, number][] = [
	['day', 86_400],
	['hour', 3_600],
	['minute', 60]
];

/** Formats a past instant as "3 hours ago" in the given locale; under a minute it is "now". */
export function timeAgo(value: string | number | Date, locale: string, now = Date.now()): string {
	const seconds = Math.max(Math.floor((now - new Date(value).getTime()) / 1000), 0);
	const [unit, size] = UNITS.find(([, length]) => seconds >= length) ?? ['second', Infinity];
	return new Intl.RelativeTimeFormat(locale, { numeric: 'auto' }).format(-Math.floor(seconds / size), unit);
}

/** Formats an instant as a medium date and a short time in the given locale. */
export function formatDateTime(value: string | number | Date, locale: string): string {
	return new Intl.DateTimeFormat(locale, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value));
}
