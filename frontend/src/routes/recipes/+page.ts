import type { PageLoad } from './$types';
import type { Recipe } from '$lib/api/recipes';

export const load: PageLoad = async ({ fetch }) => {
	const res = await fetch('/api/recipes');
	if (!res.ok) throw new Error('레시피 목록 조회 실패');
	const recipes: Recipe[] = await res.json();
	return { recipes };
};
