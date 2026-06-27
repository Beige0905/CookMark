interface User {
	id: string;
	email: string;
	display_name: string;
	avatar_url?: string;
}

let user = $state<User | null>(null);
let loading = $state(true);

async function fetchMe() {
	try {
		const res = await fetch('/api/auth/me');
		if (res.ok) {
			user = await res.json();
		} else {
			user = null;
		}
	} catch {
		user = null;
	} finally {
		loading = false;
	}
}

async function logout() {
	await fetch('/api/auth/logout', { method: 'POST' });
	user = null;
	window.location.href = '/login';
}

export const auth = {
	get user() {
		return user;
	},
	get loading() {
		return loading;
	},
	fetchMe,
	logout
};
