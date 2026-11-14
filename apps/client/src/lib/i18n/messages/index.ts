import type { Catalogs, Message } from '@facile/muse';
import { create } from './create';
import { landing } from './landing';
import { pending } from './pending';
import { reveal } from './reveal';
import { revoke } from './revoke';
import { shell } from './shell';
import type { Row } from '../row';

const rows = { ...shell, ...create, ...landing, ...reveal, ...revoke, ...pending };

export type AppKey = keyof typeof rows;
type AppCatalog = Record<AppKey, Message>;

function column(index: number): AppCatalog {
	const entries = Object.entries(rows as Record<string, Row>).map(([key, row]) => [key, row[index]]);
	return Object.fromEntries(entries) as AppCatalog;
}

/** The app's catalogs: English is the source, and every row carries all four locales. */
export const catalogs: Catalogs<AppCatalog> = {
	en: column(0),
	fr: column(1),
	es: column(2),
	de: column(3)
};
