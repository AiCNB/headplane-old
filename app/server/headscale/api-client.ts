import { readFile } from 'node:fs/promises';
import type { TFunction } from 'i18next';
import { data } from 'react-router';
import { Agent, Dispatcher, request } from 'undici';
import { errors } from 'undici';
import log from '~/utils/log';
import { getTranslator } from '../i18n.server';
import ResponseError from './api-error';

function isNodeNetworkError(error: unknown): error is NodeJS.ErrnoException {
	const keys = Object.keys(error as Record<string, unknown>);
	return (
		typeof error === 'object' &&
		error !== null &&
		keys.includes('code') &&
		keys.includes('errno')
	);
}

type Translator = TFunction<'server'>;

function friendlyError(givenError: unknown, t: Translator) {
	let error: unknown = givenError;
	if (error instanceof AggregateError) {
		error = error.errors[0];
	}

	switch (true) {
		case error instanceof errors.BodyTimeoutError:
		case error instanceof errors.ConnectTimeoutError:
		case error instanceof errors.HeadersTimeoutError:
			return data(t('errors.headscale.timeout'), {
				statusText: t('status.requestTimeout'),
				status: 408,
			});

		case error instanceof errors.SocketError:
		case error instanceof errors.SecureProxyConnectionError:
		case error instanceof errors.ClientClosedError:
		case error instanceof errors.ClientDestroyedError:
		case error instanceof errors.RequestAbortedError:
			return data(t('errors.headscale.unreachable'), {
				statusText: t('status.serviceUnavailable'),
				status: 503,
			});

		case error instanceof errors.InvalidArgumentError:
		case error instanceof errors.InvalidReturnValueError:
		case error instanceof errors.NotSupportedError:
			return data(t('errors.headscale.unableRequest'), {
				statusText: t('status.internalServerError'),
				status: 500,
			});

		case error instanceof errors.HeadersOverflowError:
		case error instanceof errors.RequestContentLengthMismatchError:
		case error instanceof errors.ResponseContentLengthMismatchError:
		case error instanceof errors.ResponseExceededMaxSizeError:
			return data(t('errors.headscale.malformed'), {
				statusText: t('status.badGateway'),
				status: 502,
			});

		case isNodeNetworkError(error):
			if (error.code === 'ECONNREFUSED') {
				return data(t('errors.headscale.unreachable'), {
					statusText: t('status.serviceUnavailable'),
					status: 503,
				});
			}

			if (error.code === 'ENOTFOUND') {
				return data(t('errors.headscale.unreachable'), {
					statusText: t('status.serviceUnavailable'),
					status: 503,
				});
			}

			if (error.code === 'EAI_AGAIN') {
				return data(t('errors.headscale.unreachable'), {
					statusText: t('status.serviceUnavailable'),
					status: 503,
				});
			}

			if (error.code === 'ETIMEDOUT') {
				return data(t('errors.headscale.timeout'), {
					statusText: t('status.requestTimeout'),
					status: 408,
				});
			}

			if (error.code === 'ECONNRESET') {
				return data(t('errors.headscale.unreachable'), {
					statusText: t('status.serviceUnavailable'),
					status: 503,
				});
			}

			if (error.code === 'EPIPE') {
				return data(t('errors.headscale.unreachable'), {
					statusText: t('status.serviceUnavailable'),
					status: 503,
				});
			}

			if (error.code === 'ENETUNREACH') {
				return data(t('errors.headscale.unreachable'), {
					statusText: t('status.serviceUnavailable'),
					status: 503,
				});
			}

			if (error.code === 'ENETRESET') {
				return data(t('errors.headscale.unreachable'), {
					statusText: t('status.serviceUnavailable'),
					status: 503,
				});
			}

			return data(t('errors.headscale.unreachable'), {
				statusText: t('status.serviceUnavailable'),
				status: 503,
			});

		default:
			return data((error as Error).message ?? t('errors.generic.unknown'), {
				statusText: t('status.internalServerError'),
				status: 500,
			});
	}
}

export async function createApiClient(base: string, certPath?: string) {
	if (!certPath) {
		return new ApiClient(new Agent(), base, getTranslator('en'));
	}

	try {
		log.debug('config', 'Loading certificate from %s', certPath);
		const data = await readFile(certPath, 'utf8');

		log.info('config', 'Using certificate from %s', certPath);
		return new ApiClient(
			new Agent({ connect: { ca: data.trim() } }),
			base,
			getTranslator('en'),
		);
	} catch (error) {
		log.error('config', 'Failed to load Headscale TLS cert: %s', error);
		log.debug('config', 'Error Details: %o', error);
		return new ApiClient(new Agent(), base, getTranslator('en'));
	}
}

export class ApiClient {
	private agent: Agent;
	private base: string;
	private translator: Translator;

	constructor(agent: Agent, base: string, translator: Translator) {
		this.agent = agent;
		this.base = base;
		this.translator = translator;
	}

	withTranslator(translator: Translator) {
		if (translator === this.translator) {
			return this;
		}

		return new ApiClient(this.agent, this.base, translator);
	}

	private async defaultFetch(
		url: string,
		options?: Partial<Dispatcher.RequestOptions>,
	) {
		const method = options?.method ?? 'GET';
		log.debug('api', '%s %s', method, url);

		try {
			const res = await request(new URL(url, this.base), {
				dispatcher: this.agent,
				headers: {
					...options?.headers,
					Accept: 'application/json',
					'User-Agent': `Headplane/${__VERSION__}`,
				},
				body: options?.body,
				method,
			});

			return res;
		} catch (error: unknown) {
			throw friendlyError(error, this.translator);
		}
	}

	async healthcheck() {
		try {
			const res = await request(new URL('/health', this.base), {
				dispatcher: this.agent,
				headers: {
					Accept: 'application/json',
					'User-Agent': `Headplane/${__VERSION__}`,
				},
			});

			return res.statusCode === 200;
		} catch (error) {
			log.debug('api', 'Healthcheck failed %o', error);
			return false;
		}
	}

	async get<T = unknown>(url: string, key: string) {
		const res = await this.defaultFetch(`/api/${url}`, {
			headers: {
				Authorization: `Bearer ${key}`,
			},
		});

		if (res.statusCode >= 400) {
			log.debug('api', 'GET %s failed with status %d', url, res.statusCode);
			throw new ResponseError(res.statusCode, await res.body.text());
		}

		return res.body.json() as Promise<T>;
	}

	async post<T = unknown>(url: string, key: string, body?: unknown) {
		const res = await this.defaultFetch(`/api/${url}`, {
			method: 'POST',
			body: body ? JSON.stringify(body) : undefined,
			headers: {
				Authorization: `Bearer ${key}`,
			},
		});

		if (res.statusCode >= 400) {
			log.debug('api', 'POST %s failed with status %d', url, res.statusCode);
			throw new ResponseError(res.statusCode, await res.body.text());
		}

		return res.body.json() as Promise<T>;
	}

	async put<T = unknown>(url: string, key: string, body?: unknown) {
		const res = await this.defaultFetch(`/api/${url}`, {
			method: 'PUT',
			body: body ? JSON.stringify(body) : undefined,
			headers: {
				Authorization: `Bearer ${key}`,
			},
		});

		if (res.statusCode >= 400) {
			log.debug('api', 'PUT %s failed with status %d', url, res.statusCode);
			throw new ResponseError(res.statusCode, await res.body.text());
		}

		return res.body.json() as Promise<T>;
	}

	async delete<T = unknown>(url: string, key: string) {
		const res = await this.defaultFetch(`/api/${url}`, {
			method: 'DELETE',
			headers: {
				Authorization: `Bearer ${key}`,
			},
		});

		if (res.statusCode >= 400) {
			log.debug('api', 'DELETE %s failed with status %d', url, res.statusCode);
			throw new ResponseError(res.statusCode, await res.body.text());
		}

		return res.body.json() as Promise<T>;
	}
}
