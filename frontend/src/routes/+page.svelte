<script lang="ts">
	import { goto } from '$app/navigation';
	import { toast } from '$lib/toast.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let cookingInProgress = $state(new Set<number>());

	async function logCooking(recipeId: number, recipeTitle: string) {
		if (cookingInProgress.has(recipeId)) return;
		cookingInProgress = new Set([...cookingInProgress, recipeId]);
		try {
			const res = await fetch(`/api/recipes/${recipeId}/logs`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({})
			});
			if (!res.ok) throw new Error();
			toast.add(`${recipeTitle} 요리 완료! 🔥`);
		} catch {
			toast.add('기록에 실패했어요', 'error');
		} finally {
			cookingInProgress = new Set([...cookingInProgress].filter((id) => id !== recipeId));
		}
	}
</script>

<div class="grid grid-cols-1 gap-10 lg:grid-cols-2">
	<!-- Left Side: Hero & Stats(mobile) -->
	<div class="space-y-6">
		<!-- Hero -->
		<div>
			<p class="mb-4 text-[11px] font-semibold uppercase tracking-widest text-stone-400">
				오늘의 추천
			</p>
			{#if data.heroRecipe}
				<button
					onclick={() => goto(`/recipes/${data.heroRecipe!.id}`)}
					class="group relative h-64 w-full overflow-hidden rounded-[2rem] shadow-lg lg:h-80 text-left"
				>
					{#if data.heroRecipe.image_url}
						<img
							src={data.heroRecipe.image_url}
							alt={data.heroRecipe.title}
							class="absolute inset-0 h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
						/>
					{:else}
						<div
							class="absolute inset-0 bg-gradient-to-br from-[#7C9A7E] to-[#4A6B4C] transition-transform duration-500 group-hover:scale-105"
						></div>
					{/if}
					<div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent"></div>

					{#if data.heroRecipe.matchedCount !== undefined}
						<div class="absolute top-5 left-6">
							<span class="rounded-full bg-white/20 px-3 py-1 text-xs font-semibold text-white backdrop-blur-md">
								🧊 재료 {data.heroRecipe.matchedCount}/{data.heroRecipe.totalCount} 보유
							</span>
						</div>
					{/if}

					<div class="absolute bottom-0 left-0 right-0 p-8">
						<h2 class="text-2xl font-bold leading-tight text-white lg:text-3xl">
							{data.heroRecipe.title}
						</h2>
					</div>
				</button>
			{:else}
				<div
					class="flex h-64 items-center justify-center rounded-[2rem] bg-stone-100 dark:bg-stone-800 lg:h-80"
				>
					<div class="text-center">
						<p class="text-stone-400 dark:text-stone-500">아직 레시피가 없어요</p>
						<button
							onclick={() => goto('/recipes/new')}
							class="mt-3 rounded-xl bg-[#7C9A7E] px-4 py-2 text-sm font-semibold text-white"
						>
							첫 레시피 추가하기
						</button>
					</div>
				</div>
			{/if}
		</div>

		<!-- Stats (Mobile Only) -->
		<div class="lg:hidden">
			<div class="rounded-2xl bg-[#F5F2EC] p-5 dark:bg-stone-800">
				<p class="mb-1 text-[11px] font-medium uppercase text-stone-500 dark:text-stone-400">
					총 레시피
				</p>
				<p class="text-3xl font-bold text-stone-800 dark:text-stone-100">{data.recipeCount}</p>
			</div>
		</div>
	</div>

	<!-- Right Side: Recipe List -->
	<div>
		<div class="mb-5 flex items-center justify-between">
			<h2 class="text-lg font-bold text-stone-900 dark:text-stone-100">내 레시피 보관함</h2>
			<button
				onclick={() => goto('/recipes')}
				class="text-sm font-semibold text-[#7C9A7E] hover:underline"
			>
				전체 보기 →
			</button>
		</div>

		{#if data.recentRecipes.length === 0}
			<p class="py-8 text-center text-stone-400 dark:text-stone-500">레시피를 추가해보세요</p>
		{:else}
			<div class="flex flex-col divide-y divide-stone-100 dark:divide-stone-800">
				{#each data.recentRecipes as recipe}
					<div
						role="button"
						tabindex="0"
						onclick={() => goto(`/recipes/${recipe.id}`)}
						onkeydown={(e) => e.key === 'Enter' && goto(`/recipes/${recipe.id}`)}
						class="group flex items-center gap-4 py-4 cursor-pointer"
					>
						<div class="min-w-0 flex-1">
							<p
								class="truncate text-base font-semibold text-stone-800 transition-colors group-hover:text-[#7C9A7E] dark:text-stone-100 dark:group-hover:text-[#7C9A7E]"
							>
								{recipe.title}
							</p>
							<div class="mt-1.5 flex flex-wrap gap-1.5">
								<span
									class="rounded-full bg-stone-100 px-2.5 py-0.5 text-[11px] font-medium text-stone-500 transition-colors group-hover:bg-[#EBF0EB] group-hover:text-[#6A8C6C] dark:bg-stone-700 dark:text-stone-400"
								>
									재료 {recipe.ingredients?.length ?? 0}가지
								</span>
							</div>
						</div>
						<button
							onclick={(e) => {
								e.stopPropagation();
								logCooking(recipe.id, recipe.title);
							}}
							disabled={cookingInProgress.has(recipe.id)}
							class="shrink-0 whitespace-nowrap rounded-xl bg-stone-50 px-3 py-2 text-xs font-bold text-[#7C9A7E] opacity-0 transition-all hover:bg-[#7C9A7E] hover:text-white group-hover:opacity-100 disabled:opacity-50 dark:bg-stone-700"
						>
							{cookingInProgress.has(recipe.id) ? '저장 중...' : '만들었어! 🔥'}
						</button>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
