<script lang="ts">
	import { toast } from '$lib/toast.svelte';
	import { getRecipes } from '$lib/api/recipes';
	import { getPantry, addPantryItem, deletePantryItem, getRecommendations, type PantryItem, type RecommendResult } from '$lib/api/pantry';

	let items = $state<PantryItem[]>([]);
	let recommendations = $state<RecommendResult[]>([]);
	let inputName = $state('');
	let suggestNames = $state<string[]>([]);
	let isRecommending = $state(false);
	let showRecommendations = $state(false);

	async function loadAll() {
		const [pantryData, recipes] = await Promise.all([getPantry(), getRecipes()]);
		items = pantryData;

		const nameSet = new Set<string>();
		for (const recipe of recipes) {
			for (const ing of recipe.ingredients) {
				if (ing.name) nameSet.add(ing.name);
			}
		}
		suggestNames = [...nameSet].sort((a, b) => a.localeCompare(b, 'ko'));
	}

	loadAll().catch(() => toast.add('냉장고 데이터를 불러오지 못했습니다.', 'error'));

	async function handleAdd(event: SubmitEvent) {
		event.preventDefault();
		const name = inputName.trim();
		if (!name) return;
		if (items.some((i) => i.name === name)) {
			toast.add('이미 추가된 재료입니다.', 'error');
			return;
		}
		try {
			const item = await addPantryItem(name);
			items = [item, ...items];
			inputName = '';
			showRecommendations = false;
		} catch (e: any) {
			toast.add(e.message || '재료 추가 실패', 'error');
		}
	}

	async function handleDelete(id: number) {
		try {
			await deletePantryItem(id);
			items = items.filter((i) => i.id !== id);
			showRecommendations = false;
		} catch (e: any) {
			toast.add(e.message || '재료 삭제 실패', 'error');
		}
	}

	async function handleRecommend() {
		if (items.length === 0) {
			toast.add('냉장고에 재료를 먼저 추가해주세요.', 'error');
			return;
		}
		isRecommending = true;
		try {
			recommendations = await getRecommendations();
			showRecommendations = true;
			if (recommendations.length === 0) {
				toast.add('매칭되는 레시피가 없습니다.', 'error');
			}
		} catch (e: any) {
			toast.add(e.message || '추천 조회 실패', 'error');
		} finally {
			isRecommending = false;
		}
	}
</script>

<div class="space-y-8">
	<div>
		<h1 class="text-3xl font-bold text-stone-900">내 냉장고 🧊</h1>
		<p class="text-stone-500">지금 냉장고에 있는 재료들로 무엇을 만들 수 있을까요?</p>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<div class="col-span-1 md:col-span-2 space-y-6">
			<!-- 재료 입력 -->
			<form onsubmit={handleAdd} class="flex gap-3">
				<div class="relative flex-1">
					<input
						type="text"
						bind:value={inputName}
						placeholder="재료 이름 입력..."
						list="ingredient-suggestions"
						class="w-full rounded-2xl border border-stone-200 bg-white px-4 py-3 outline-none focus:border-[#7C9A7E] focus:ring-4 focus:ring-[#7C9A7E]/5 transition-all"
					/>
					<datalist id="ingredient-suggestions">
						{#each suggestNames as name}
							<option value={name}></option>
						{/each}
					</datalist>
				</div>
				<button
					type="submit"
					disabled={!inputName.trim()}
					class="rounded-2xl bg-[#7C9A7E] px-5 py-3 font-bold text-white disabled:opacity-40 transition-all hover:bg-[#6A8C6C] active:scale-95"
				>
					추가
				</button>
			</form>

			<!-- 보관 중인 재료 -->
			<div class="rounded-[2rem] bg-stone-50 p-8 border border-stone-100">
				<h2 class="text-xl font-bold text-stone-800 mb-6">보관 중인 재료</h2>
				{#if items.length === 0}
					<p class="text-stone-400 text-sm text-center py-4">재료를 추가해보세요.</p>
				{:else}
					<div class="flex flex-wrap gap-3">
						{#each items as item (item.id)}
							<span class="inline-flex items-center gap-2 rounded-2xl bg-white px-4 py-2 text-sm font-semibold text-stone-700 shadow-sm border border-stone-100">
								{item.name}
								<button
									onclick={() => handleDelete(item.id)}
									class="text-stone-300 hover:text-red-400 transition-colors leading-none"
									aria-label="{item.name} 삭제"
								>×</button>
							</span>
						{/each}
					</div>
				{/if}
			</div>

			<button
				onclick={handleRecommend}
				disabled={isRecommending || items.length === 0}
				class="w-full rounded-[2rem] bg-[#7C9A7E] py-6 text-xl font-bold text-white shadow-xl shadow-[#7C9A7E]/20 transition-all hover:scale-[1.02] active:scale-95 disabled:opacity-40 disabled:scale-100"
			>
				{isRecommending ? '추천 중...' : '🍳 이 재료들로 뭐 만들지?'}
			</button>

			<!-- 추천 결과 -->
			{#if showRecommendations && recommendations.length > 0}
				<div class="space-y-4">
					<h2 class="text-xl font-bold text-stone-800">만들 수 있는 레시피</h2>
					{#each recommendations as result (result.recipe.id)}
						<a
							href="/recipes/{result.recipe.id}"
							class="flex items-center justify-between rounded-2xl bg-white border border-stone-100 px-6 py-4 shadow-sm hover:border-[#7C9A7E] transition-colors"
						>
							<div>
								<p class="font-bold text-stone-800">{result.recipe.title}</p>
								<p class="text-sm text-stone-400 mt-0.5">{result.recipe.base_servings}인분 기준</p>
							</div>
							<div class="text-right shrink-0 ml-4">
								<p class="text-2xl font-bold text-[#7C9A7E]">{result.matched_count}<span class="text-base text-stone-400">/{result.total_count}</span></p>
								<p class="text-xs text-stone-400">재료 보유</p>
							</div>
						</a>
					{/each}
				</div>
			{/if}
		</div>

		<!-- 통계 사이드바 -->
		<div class="space-y-6">
			<div class="rounded-[2rem] bg-[#EEF3EE] p-6">
				<h3 class="font-bold text-[#6A8C6C] mb-4 text-sm uppercase tracking-wider">냉장고 현황</h3>
				<p class="text-3xl font-bold text-[#5A7A5C]">{items.length}<span class="text-base font-normal text-[#7C9A7E] ml-1">가지</span></p>
				<p class="text-sm text-[#7C9A7E] mt-1">재료 보관 중</p>
			</div>
		</div>
	</div>
</div>
