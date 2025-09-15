import { type Server, createServer } from 'node:http';
import type { AddressInfo } from 'node:net';
import { afterAll, beforeAll, describe, expect, it } from 'vitest';

import {
	detectLocale,
	getTranslator,
	initServerI18n,
} from '../../app/server/i18n.server';

let server: Server;
let baseUrl: string;

describe('i18n HTTP integration', () => {
	beforeAll(async () => {
		await initServerI18n();

		server = createServer((req, res) => {
			const header = req.headers['accept-language'];
			const locale = detectLocale(
				typeof header === 'string'
					? header
					: Array.isArray(header)
						? header.join(',')
						: null,
			);

			const t = getTranslator(locale);
			const message = t('status.requestTimeout');

			res.statusCode = 200;
			res.setHeader('Content-Type', 'text/plain; charset=utf-8');
			res.end(message);
		});

		await new Promise<void>((resolve) => {
			server.listen(0, () => {
				const address = server.address() as AddressInfo;
				baseUrl = `http://127.0.0.1:${address.port}`;
				resolve();
			});
		});
	});

	afterAll(async () => {
		await new Promise<void>((resolve, reject) => {
			server.close((err) => {
				if (err) {
					reject(err);
					return;
				}

				resolve();
			});
		});
	});

	it('returns English content by default', async () => {
		const response = await fetch(baseUrl);
		expect(await response.text()).toBe('Request Timeout');
	});

	it('switches to Chinese when Accept-Language prefers zh', async () => {
		const response = await fetch(baseUrl, {
			headers: { 'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8' },
		});

		expect(await response.text()).toBe('请求超时');
	});
});
