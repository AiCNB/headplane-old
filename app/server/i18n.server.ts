import i18next, { type i18n, type TFunction } from 'i18next';
import enServer from '../../locales/en/server.json';
import zhServer from '../../locales/zh/server.json';

const supportedLocales = ['en', 'zh'] as const;
type SupportedLocale = (typeof supportedLocales)[number];

let instance: i18n | undefined;

export async function initServerI18n() {
	if (instance) {
		return instance;
	}

	const created = i18next.createInstance();
	await created.init({
		fallbackLng: 'en',
		supportedLngs: [...supportedLocales],
		defaultNS: 'server',
		resources: {
			en: { server: enServer },
			zh: { server: zhServer },
		},
		interpolation: { escapeValue: false },
	});

	instance = created;
	return instance;
}

function normalizeLocale(locale: string): SupportedLocale {
	const base = locale.toLowerCase().split('-')[0] as SupportedLocale;
	return supportedLocales.includes(base) ? base : 'en';
}

export function getTranslator(locale?: string): TFunction<'server'> {
	if (!instance) {
		throw new Error('Server i18n has not been initialized');
	}

	if (!locale) {
		return instance.getFixedT('en', 'server');
	}

	return instance.getFixedT(normalizeLocale(locale), 'server');
}

export function detectLocale(acceptLanguage: string | null): SupportedLocale {
	if (!acceptLanguage) {
		return 'en';
	}

	const languages = acceptLanguage
		.split(',')
		.map((entry) => {
			const [lang, ...params] = entry.trim().split(';');
			let quality = 1;

			for (const param of params) {
				const [key, value] = param.split('=');
				if (key?.trim() === 'q') {
					const parsed = Number.parseFloat(value ?? '1');
					if (!Number.isNaN(parsed)) {
						quality = parsed;
					}
				}
			}

			return { lang, quality };
		})
		.sort((a, b) => b.quality - a.quality);

	for (const { lang } of languages) {
		if (!lang) {
			continue;
		}

		const base = lang.split('-')[0] as SupportedLocale;
		if (supportedLocales.includes(base)) {
			return base;
		}
	}

	return 'en';
}
