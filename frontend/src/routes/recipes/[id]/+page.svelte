<script lang="ts">
	import { page } from '$app/stores';
	import { getRecipe, type Recipe } from '$lib/api/recipes';

	let recipe: Recipe | null = $state(null);
	let error = $state('');

	$effect(() => {
		getRecipe(Number($page.params.id))
			.then((data) => (recipe = data))
			.catch((e) => (error = e.message));
	});
</script>

{#if error}
	<p>{error}</p>
{:else if recipe}
	<h1>{recipe.title}</h1>
	<p>{recipe.description}</p>
	<small>{new Date(recipe.created_at).toLocaleDateString('ko-KR')}</small>
{:else}
	<p>불러오는 중...</p>
{/if}
