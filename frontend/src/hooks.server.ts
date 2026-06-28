import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
import type { Cookies, Handle } from '@sveltejs/kit';

const PUBLIC_PATHS = ['/login', '/register'];
const BACKEND_URL = env.API_URL ?? 'http://localhost:8080';

export const handle: Handle = async ({ event, resolve }) => {
	const accessToken = event.cookies.get('access_token');
	const isPublic = PUBLIC_PATHS.some((p) => event.url.pathname.startsWith(p));

	if (!accessToken) {
		const refreshToken = event.cookies.get('refresh_token');

		if (refreshToken && !isPublic) {
			const refreshed = await tryRefresh(event.cookies, refreshToken);
			if (refreshed) {
				return resolve(event);
			}
		}

		if (!isPublic) {
			throw redirect(302, '/login');
		}
	} else if (isPublic) {
		throw redirect(302, '/');
	}

	return resolve(event);
};

async function tryRefresh(cookies: Cookies, refreshToken: string): Promise<boolean> {
	try {
		const res = await fetch(`${BACKEND_URL}/api/auth/refresh`, {
			method: 'POST',
			headers: { Cookie: `refresh_token=${refreshToken}` }
		});

		if (!res.ok) return false;

		const setCookie = res.headers.get('set-cookie');
		if (!setCookie) return false;

		const match = setCookie.match(/access_token=([^;]+)/);
		if (!match) return false;

		cookies.set('access_token', match[1], {
			path: '/',
			httpOnly: true,
			sameSite: 'lax',
			maxAge: 15 * 60
		});
		return true;
	} catch {
		return false;
	}
}
