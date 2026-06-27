const BASE = '/api';

export type Ingredient = {
	name: string;
	amount_num?: number;
	unit?: string;
	note?: string;
	amount?: string; // 기존 호환성 및 수동 입력용
	scaling_type?: 'linear' | 'culinary';
	scaling_factor?: number;
};

export type Recipe = {
	id: number;
	title: string;
	description?: string;
	origin_url?: string;
	image_url?: string;
	base_servings: number;
	ingredients: Ingredient[];
	instructions: string[];
	created_at: string;
};

export type YouTubeExtractResult = {
	title: string;
	base_servings: number;
	ingredients: Ingredient[];
	image_url?: string;
};

export type ImageExtractResult = {
	title: string;
	base_servings: number;
	ingredients: Ingredient[];
	instructions: string[];
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

export async function extractFromYouTube(url: string): Promise<YouTubeExtractResult> {
	const res = await fetch(`${BASE}/youtube/extract`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ url })
	});
	if (!res.ok) {
		const msg = await res.text();
		throw new Error(msg || '재료 추출 실패');
	}
	return res.json();
}

export type RecipeNote = {
	recipe_id: number;
	memo: string;
	adjustments: Record<string, number>;
	updated_at: string;
};

export async function getRecipeNote(recipeID: number): Promise<RecipeNote> {
	const res = await fetch(`${BASE}/recipes/${recipeID}/note`);
	if (!res.ok) throw new Error('메모 조회 실패');
	return res.json();
}

export async function upsertRecipeNote(recipeID: number, note: Pick<RecipeNote, 'memo' | 'adjustments'>): Promise<void> {
	const res = await fetch(`${BASE}/recipes/${recipeID}/note`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(note)
	});
	if (!res.ok) throw new Error('메모 저장 실패');
}

export async function extractFromImage(file: File): Promise<ImageExtractResult> {
	const formData = new FormData();
	formData.append('image', file);

	const res = await fetch(`${BASE}/recipes/extract-image`, {
		method: 'POST',
		body: formData
	});
	if (!res.ok) {
		const msg = await res.text();
		throw new Error(msg || '이미지 분석 실패');
	}
	return res.json();
}
