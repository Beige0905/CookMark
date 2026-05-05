<script lang="ts">
	import type { Recipe } from '$lib/api/recipes';

	let { recipe }: { recipe: Recipe } = $props();
</script>

<a href="/recipes/{recipe.id}">
	<article>
		<h2>{recipe.title}</h2>
		<p>{recipe.description}</p>
	</article>
</a>
