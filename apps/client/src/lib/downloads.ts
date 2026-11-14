import type { CliCommand, DownloadEntry } from '@facile/muse';

const CLI_VERSION = '0.4.2';
const RELEASE = `https://github.com/FacileStudio/capsule-cli/releases/download/v${CLI_VERSION}`;

export const RELEASES_HREF = 'https://github.com/FacileStudio/capsule-cli/releases';
export const CHECKSUMS_HREF = `${RELEASE}/checksums.txt`;

function archive(platform: 'macos' | 'linux', os: string, arch: string, label: string, size: string): DownloadEntry {
	return {
		platform,
		href: `${RELEASE}/capsule-cli_${CLI_VERSION}_${os}_${arch}.tar.gz`,
		label,
		arch,
		version: CLI_VERSION,
		size,
		checksumHref: CHECKSUMS_HREF
	};
}

/** The capsule CLI archives attached to the pinned capsule-cli release. */
export const CLI_DOWNLOADS: DownloadEntry[] = [
	archive('macos', 'darwin', 'arm64', 'Apple Silicon', '4.0 MB'),
	archive('macos', 'darwin', 'amd64', 'Intel', '4.3 MB'),
	archive('linux', 'linux', 'amd64', 'x86_64', '4.3 MB'),
	archive('linux', 'linux', 'arm64', 'ARM64', '3.8 MB')
];

/** The install commands for the capsule CLI, with the script label translated by the caller. */
export function cliCommands(scriptLabel: string): CliCommand[] {
	return [
		{ label: scriptLabel, command: 'curl -fsSL https://raw.githubusercontent.com/FacileStudio/capsule-cli/main/install.sh | bash' },
		{ label: 'facile', command: 'facile install capsule' },
		{ label: 'Homebrew', command: 'brew install facilestudio/tap/capsule' }
	];
}
