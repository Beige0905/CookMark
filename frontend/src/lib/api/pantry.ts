const BASE = '/api';

export type PantryItem = {
	id: number;
	name: string;
	created_at: string;
};

export type RecommendResult = {
	recipe: {
		id: number;
		title: string;
		image_url?: string;
		base_servings: number;
		created_at: string;
	};
	matched_count: number;
	total_count: number;
};

export async function getPantry(): Promise<PantryItem[]> {
	const res = await fetch(`${BASE}/pantry`);
	if (!res.ok) throw new Error('냉장고 재료 조회 실패');
	return res.json();
}

export async function addPantryItem(name: string): Promise<PantryItem> {
	const res = await fetch(`${BASE}/pantry`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ name })
	});
	if (!res.ok) throw new Error('재료 추가 실패');
	return res.json();
}

export async function deletePantryItem(id: number): Promise<void> {
	const res = await fetch(`${BASE}/pantry/${id}`, { method: 'DELETE' });
	if (!res.ok) throw new Error('재료 삭제 실패');
}

export async function getRecommendations(): Promise<RecommendResult[]> {
	const res = await fetch(`${BASE}/pantry/recommend`);
	if (!res.ok) throw new Error('추천 조회 실패');
	return res.json();
}
