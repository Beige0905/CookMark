<script lang="ts">
	import type { Recipe } from '$lib/api/recipes';
	import Clock from 'lucide-svelte/icons/clock';
	import Utensils from 'lucide-svelte/icons/utensils';

	let { recipe }: { recipe: Recipe } = $props();
</script>

<a href="/recipes/{recipe.id}" class="group block">
	<article class="h-full bg-white rounded-3xl border border-stone-100 shadow-sm overflow-hidden transition-all duration-300 hover:shadow-md hover:border-[#7C9A7E]/30 group-active:scale-[0.98]">
<div class="p-6">
			<div class="flex items-center gap-2 mb-2">
				<span class="px-2 py-0.5 bg-[#EEF3EE] text-[#5A7A5C] text-[10px] font-bold rounded-md uppercase tracking-wider">Recipe</span>
				<span class="text-[10px] text-stone-400 font-medium">{new Date(recipe.created_at).toLocaleDateString('ko-KR')}</span>
			</div>
			<h3 class="text-xl font-bold text-stone-900 mb-2 group-hover:text-[#7C9A7E] transition-colors line-clamp-1">{recipe.title}</h3>
			
			{#if recipe.description}
				<p class="text-stone-500 text-sm line-clamp-2 mb-4 leading-relaxed">{recipe.description}</p>
			{:else if recipe.instructions && recipe.instructions.length > 0}
				<p class="text-stone-500 text-sm line-clamp-2 mb-4 leading-relaxed">{recipe.instructions[0]}</p>
			{:else}
				<p class="text-stone-400 text-sm italic mb-4">작성된 설명이 없습니다.</p>
			{/if}

			<div class="flex items-center gap-4 pt-4 border-t border-stone-50">
				<div class="flex items-center gap-1.5 text-stone-400 text-xs">
					<Utensils size={14} />
					<span>재료 {recipe.ingredients?.length || 0}개</span>
				</div>
				<div class="flex items-center gap-1.5 text-stone-400 text-xs">
					<Clock size={14} />
					<span>{recipe.instructions?.length || 0} 단계</span>
				</div>
			</div>
		</div>
	</article>
</a>
