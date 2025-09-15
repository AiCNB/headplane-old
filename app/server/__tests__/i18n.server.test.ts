import { beforeAll, describe, expect, it } from 'vitest';

import { getTranslator, initServerI18n } from '../i18n.server';

describe('server translator', () => {
	beforeAll(async () => {
		await initServerI18n();
	});

	it('returns Chinese text for status.requestTimeout', () => {
		const t = getTranslator('zh');
		expect(t('status.requestTimeout')).toBe('请求超时');
	});

	it('returns Chinese text for errors.headscale.timeout', () => {
		const t = getTranslator('zh');
		expect(t('errors.headscale.timeout')).toBe('等待 Headscale API 响应超时');
	});
});
