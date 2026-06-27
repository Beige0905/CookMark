import { redirect } from '@sveltejs/kit';
import type { Handle } from '@sveltejs/kit';

const PUBLIC_PATHS = ['/login', '/register'];

export const handle: Handle = async ({ event, resolve }) => {
	const accessToken = event.cookies.get('access_token');
	const isPublic = PUBLIC_PATHS.some((p) => event.url.pathname.startsWith(p));

	if (!isPublic && !accessToken) {
		throw redirect(302, '/login');
	}
	if (isPublic && accessToken) {
		throw redirect(302, '/');
	}

	return resolve(event);
};
