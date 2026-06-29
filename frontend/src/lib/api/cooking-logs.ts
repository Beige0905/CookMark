const BASE = '/api';

export type CookingLog = {
	id: number;
	recipe_id: number;
	comment?: string;
	cooked_at: string;
};

export async function getCookingLogs(recipeID: number): Promise<CookingLog[]> {
	const res = await fetch(`${BASE}/recipes/${recipeID}/logs`);
	if (!res.ok) throw new Error('요리 기록 조회 실패');
	return res.json();
}

export async function addCookingLog(recipeID: number, data: { comment?: string; cooked_at?: string; delete_pantry_ids?: number[] }): Promise<CookingLog> {
	const res = await fetch(`${BASE}/recipes/${recipeID}/logs`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(data)
	});
	if (!res.ok) throw new Error('요리 기록 추가 실패');
	return res.json();
}

export async function deleteCookingLog(recipeID: number, logID: number): Promise<void> {
	const res = await fetch(`${BASE}/recipes/${recipeID}/logs/${logID}`, { method: 'DELETE' });
	if (!res.ok) throw new Error('요리 기록 삭제 실패');
}
