<script lang="ts">
	import type { PageData } from './$types';
	import RecipeCard from '$lib/components/RecipeCard.svelte';
	import Search from 'lucide-svelte/icons/search';

	let { data }: { data: PageData } = $props();

	let searchQuery = $state('');
	let sortOrder = $state<'newest' | 'oldest' | 'alpha'>('newest');

	let filtered = $derived.by(() => {
		let result = data.recipes;

		if (searchQuery.trim()) {
			const q = searchQuery.toLowerCase();
			result = result.filter(
				(r) =>
					r.title.toLowerCase().includes(q) ||
					r.ingredients?.some((ing) => ing.name.toLowerCase().includes(q))
			);
		}

		return [...result].sort((a, b) => {
			if (sortOrder === 'newest') return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
			if (sortOrder === 'oldest') return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
			return a.title.localeCompare(b.title, 'ko');
		});
	});
</script>

<div class="space-y-8">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold text-stone-900">내 레시피 보관함</h1>
			<p class="text-stone-500">총 {data.recipes.length}개의 레시피가 저장되어 있어요.</p>
		</div>
		<a href="/recipes/new" class="rounded-2xl bg-[#7C9A7E] px-6 py-3 font-bold text-white shadow-lg shadow-[#7C9A7E]/20 transition-transform hover:scale-105 active:scale-95">
			+ 새 레시피 추가
		</a>
	</div>

	<!-- Search & Filter -->
	<div class="flex flex-wrap gap-3">
		<div class="flex-1 min-w-[240px] relative">
			<Search class="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-stone-400" />
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="레시피 이름, 재료 검색..."
				class="w-full rounded-2xl border border-stone-200 bg-white py-3 pl-12 pr-4 outline-none focus:border-[#7C9A7E] focus:ring-4 focus:ring-[#7C9A7E]/5 transition-all"
			/>
		</div>
		<select
			bind:value={sortOrder}
			class="rounded-2xl border border-stone-200 bg-white px-4 py-3 outline-none focus:border-[#7C9A7E]"
		>
			<option value="newest">최신순</option>
			<option value="oldest">오래된순</option>
			<option value="alpha">가나다순</option>
		</select>
	</div>

	<!-- Recipe Grid -->
	{#if filtered.length === 0}
		<div class="flex flex-col items-center justify-center py-24 text-stone-400">
			{#if searchQuery}
				<p class="text-lg font-medium">검색 결과가 없습니다.</p>
				<p class="text-sm mt-1">다른 키워드로 검색해보세요.</p>
			{:else}
				<p class="text-lg font-medium">아직 저장된 레시피가 없어요.</p>
				<a href="/recipes/new" class="mt-4 text-[#7C9A7E] font-bold hover:underline">첫 레시피 추가하기 →</a>
			{/if}
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			{#each filtered as recipe (recipe.id)}
				<RecipeCard {recipe} />
			{/each}
		</div>
	{/if}
</div>
