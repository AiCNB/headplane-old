import en from '../../locales/en/translation.json';

function getNested(obj: Record<string, unknown>, path: string[]): unknown {
	return path.reduce<unknown>((acc, key) => {
		if (acc && typeof acc === 'object') {
			return (acc as Record<string, unknown>)[key];
		}
		return undefined;
	}, obj);
}

export function t(key: string): string {
	const result = getNested(en as Record<string, unknown>, key.split('.'));
	return typeof result === 'string' ? result : key;
}
