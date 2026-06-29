import type { PageServerLoad } from './$types';
import type { Recipe } from '$lib/api/recipes';
import type { RecommendResult } from '$lib/api/pantry';

export type HeroRecipe = {
	id: number;
	title: string;
	image_url?: string;
	matchedCount?: number;
	totalCount?: number;
};

export const load: PageServerLoad = async ({ fetch }) => {
	const [recipesRes, recommendRes] = await Promise.all([
		fetch('/api/recipes'),
		fetch('/api/pantry/recommend')
	]);

	const recipes: Recipe[] = recipesRes.ok ? await recipesRes.json() : [];
	const recommendations: RecommendResult[] = recommendRes.ok ? await recommendRes.json() : [];

	const recentRecipes = [...recipes]
		.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
		.slice(0, 3);

	let heroRecipe: HeroRecipe | null = null;

	if (recommendations.length > 0) {
		const top = recommendations[0];
		heroRecipe = {
			id: top.recipe.id,
			title: top.recipe.title,
			image_url: top.recipe.image_url,
			matchedCount: top.matched_count,
			totalCount: top.total_count
		};
	} else if (recipes.length > 0) {
		const r = recipes[Math.floor(Math.random() * recipes.length)];
		heroRecipe = { id: r.id, title: r.title, image_url: r.image_url };
	}

	return { recentRecipes, heroRecipe, recipeCount: recipes.length };
};
