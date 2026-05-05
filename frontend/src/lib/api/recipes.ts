const BASE = '/api';

export type Recipe = {
	id: number;
	title: string;
	description: string;
	created_at: string;
};

export async function getRecipes(): Promise<Recipe[]> {
	const res = await fetch(`${BASE}/recipes`);
	if (!res.ok) throw new Error('레시피 목록 조회 실패');
	return res.json();
}

export async function getRecipe(id: number): Promise<Recipe> {
	const res = await fetch(`${BASE}/recipes/${id}`);
	if (!res.ok) throw new Error('레시피 조회 실패');
	return res.json();
}

export async function createRecipe(data: Omit<Recipe, 'id' | 'created_at'>): Promise<Recipe> {
	const res = await fetch(`${BASE}/recipes`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(data)
	});
	if (!res.ok) throw new Error('레시피 생성 실패');
	return res.json();
}
