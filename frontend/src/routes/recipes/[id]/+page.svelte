<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getRecipe, getRecipeNote, upsertRecipeNote, deleteRecipe, type Recipe, type Ingredient } from '$lib/api/recipes';
	import { toast } from '$lib/toast.svelte';
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import Users from 'lucide-svelte/icons/users';
	import Pencil from 'lucide-svelte/icons/pencil';
	import Trash2 from 'lucide-svelte/icons/trash-2';

	let recipe: Recipe | null = $state(null);
	let servingSize = $state(1);
	let error = $state('');
	let activeTab = $state<'ingredients' | 'steps'>('ingredients');

	let editMode = $state(false);
	let savedMemo = $state('');
	let savedAdjustments = $state<Record<string, number>>({});
	let pendingMemo = $state('');
	let pendingAdjustments = $state<Record<string, number>>({});
	let saving = $state(false);
	let deleting = $state(false);

	const recipeId = $derived(Number($page.params.id));

	async function handleDelete() {
		if (!confirm('이 레시피를 삭제할까요?')) return;
		deleting = true;
		try {
			await deleteRecipe(recipeId);
			toast.add('레시피가 삭제됐습니다.', 'success');
			goto('/recipes');
		} catch (e: any) {
			toast.add(e.message || '삭제 중 오류가 발생했습니다.', 'error');
			deleting = false;
		}
	}

	let hasAdjustments = $derived(Object.keys(savedAdjustments).some(k => savedAdjustments[k] !== 0));

	function enterEditMode() {
		pendingMemo = savedMemo;
		pendingAdjustments = { ...savedAdjustments };
		editMode = true;
	}

	function cancelEdit() {
		editMode = false;
	}

	async function saveNote() {
		if (!recipe) return;
		saving = true;
		try {
			await upsertRecipeNote(recipe.id, { memo: pendingMemo, adjustments: pendingAdjustments });
			savedMemo = pendingMemo;
			savedAdjustments = { ...pendingAdjustments };
			editMode = false;
		} finally {
			saving = false;
		}
	}

	function resetAdjustments() {
		savedAdjustments = {};
		savedMemo = '';
		upsertRecipeNote(recipe!.id, { memo: '', adjustments: {} });
	}

	function adjustIngredient(name: string, delta: number) {
		const current = pendingAdjustments[name] ?? 0;
		const next = Math.max(-80, Math.min(100, current + delta));
		pendingAdjustments = { ...pendingAdjustments, [name]: next };
	}

	const SPOON_UNITS = new Set(['스푼', '큰술', '작은술', 'T', 't', 'tbsp', 'tsp', '컵', 'cup']);
	const WEIGHT_VOLUME_UNITS = new Set(['g', 'ml', 'L', 'l', 'kg']);
	const COUNT_UNITS = new Set(['개', '장', '줄기', '쪽', '알', '마리', '봉지', '묶음', '조각', '덩어리', '토막', '가닥', '뿌리']);

	function roundByUnit(amount: number, unit: string): number {
		const u = unit.trim();
		if (SPOON_UNITS.has(u)) return Math.round(amount * 4) / 4;
		if (WEIGHT_VOLUME_UNITS.has(u)) {
			if (amount <= 10) return Math.round(amount);
			if (amount <= 100) return Math.round(amount / 5) * 5;
			return Math.round(amount / 10) * 10;
		}
		if (COUNT_UNITS.has(u)) return Math.ceil(amount);
		return Math.round(amount * 10) / 10;
	}

	function formatAmount(amount: number, unit: string): string {
		const u = unit.trim();
		if (!SPOON_UNITS.has(u)) return `${amount}${unit}`;
		const whole = Math.floor(amount);
		const frac = amount - whole;
		let fracStr = '';
		if (Math.abs(frac - 0.25) < 0.01) fracStr = '¼';
		else if (Math.abs(frac - 0.5) < 0.01) fracStr = '½';
		else if (Math.abs(frac - 0.75) < 0.01) fracStr = '¾';
		const wholeStr = whole > 0 ? `${whole}` : '';
		return `${wholeStr}${fracStr || (whole === 0 ? '0' : '')}${unit}`;
	}

	// 인분 수에 따른 재료 계산
	let activeIngredients = $derived.by(() => {
		if (!recipe) return [];
		return recipe.ingredients.map(ing => scaleIngredient(ing, servingSize));
	});

	function scaleIngredient(ing: Ingredient, size: number): Ingredient {
		if (ing.amount_num) {
			const factor = ing.scaling_factor ?? 1.0;
			const scaled = ing.amount_num * (1 + (size - 1) * factor);
			const adj = savedAdjustments[ing.name] ?? 0;
			const adjusted = scaled * (1 + adj / 100);
			return {
				...ing,
				amount_num: roundByUnit(adjusted, ing.unit || ''),
				amount: undefined
			};
		}
		return ing;
	}

	function getDisplayAmount(ing: Ingredient): string {
		if (ing.amount) return ing.amount;
		if (ing.amount_num !== undefined && ing.amount_num !== null) {
			return formatAmount(ing.amount_num, ing.unit || '');
		}
		return ing.note || '';
	}

	$effect(() => {
		const id = recipeId;
		getRecipe(id)
			.then((data) => (recipe = data))
			.catch((e) => (error = e.message));
		getRecipeNote(id)
			.then((note) => {
				savedMemo = note.memo;
				savedAdjustments = note.adjustments;
			})
			.catch(() => {});
	});
</script>

<div class="max-w-2xl">
	<a href="/" class="mb-8 inline-flex items-center gap-2 text-sm font-semibold text-stone-400 hover:text-[#7C9A7E] transition-colors">
		<ChevronLeft class="h-4 w-4" strokeWidth={2} />
		홈으로 돌아가기
	</a>

	{#if error}
		<div class="rounded-2xl bg-red-50 p-6 text-red-600">
			<p class="font-bold">오류가 발생했습니다</p>
			<p class="text-sm">{error}</p>
		</div>
	{:else if recipe}
		<div class="space-y-4">
			<div class="flex items-start justify-between gap-4">
				<div class="space-y-2">
					<h1 class="text-4xl font-bold text-stone-900">{recipe.title}</h1>
					<p class="text-stone-500">{recipe.description}</p>
				</div>
				<div class="flex items-center gap-2 shrink-0">
					<a
						href="/recipes/{recipeId}/edit"
						class="flex items-center gap-1.5 px-3 py-2 rounded-xl border border-stone-200 text-sm font-semibold text-stone-500 hover:border-[#7C9A7E] hover:text-[#7C9A7E] transition-colors"
					>
						<Pencil class="h-4 w-4" />
						수정
					</a>
					<button
						onclick={handleDelete}
						disabled={deleting}
						class="flex items-center gap-1.5 px-3 py-2 rounded-xl border border-stone-200 text-sm font-semibold text-stone-500 hover:border-red-300 hover:text-red-500 transition-colors disabled:opacity-40"
					>
						<Trash2 class="h-4 w-4" />
						삭제
					</button>
				</div>
			</div>

			<div class="flex items-center gap-4 text-sm text-stone-400">
				<span>📅 {new Date(recipe.created_at).toLocaleDateString('ko-KR')}</span>
				<span>👤 작성자: 나</span>
			</div>

			<!-- Sticky 탭 (모바일 전용) -->
			<div class="sticky top-0 z-10 -mx-5 bg-white px-5 border-b border-stone-100 md:hidden">
				<div class="flex">
					<button
						onclick={() => activeTab = 'ingredients'}
						class="flex-1 py-3 text-sm font-semibold transition-colors border-b-2 {activeTab === 'ingredients' ? 'border-[#7C9A7E] text-[#7C9A7E]' : 'border-transparent text-stone-400'}"
					>
						재료
					</button>
					<button
						onclick={() => activeTab = 'steps'}
						class="flex-1 py-3 text-sm font-semibold transition-colors border-b-2 {activeTab === 'steps' ? 'border-[#7C9A7E] text-[#7C9A7E]' : 'border-transparent text-stone-400'}"
					>
						요리방법
					</button>
				</div>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<!-- 재료 패널 -->
				<div class="rounded-3xl bg-white p-8 border border-stone-100 shadow-sm {activeTab !== 'ingredients' ? 'hidden md:block' : ''}">
					<div class="mb-6 flex items-center justify-between">
						<h2 class="text-xl font-bold text-stone-800">재료</h2>

						<!-- 인분 수 선택기 -->
						<div class="flex items-center gap-1 bg-stone-50 px-2 py-1 rounded-xl border border-stone-100">
							<Users class="w-3.5 h-3.5 text-stone-400" />
							<button
								onclick={() => servingSize = Math.max(1, servingSize - 1)}
								disabled={servingSize <= 1}
								class="w-6 h-6 flex items-center justify-center rounded-lg text-stone-400 hover:text-stone-600 hover:bg-stone-200 disabled:opacity-30 transition-all font-bold text-base leading-none"
							>−</button>
							<span class="text-xs font-bold text-stone-700 w-10 text-center">{servingSize}인분</span>
							<button
								onclick={() => servingSize = Math.min(10, servingSize + 1)}
								disabled={servingSize >= 10}
								class="w-6 h-6 flex items-center justify-center rounded-lg text-stone-400 hover:text-stone-600 hover:bg-stone-200 disabled:opacity-30 transition-all font-bold text-base leading-none"
							>+</button>
						</div>
					</div>

					<!-- 내 취향 조정 배너 -->
					{#if hasAdjustments && !editMode}
						<div class="mb-4 flex items-center justify-between rounded-xl bg-[#7C9A7E]/10 px-3 py-2 text-xs text-[#7C9A7E]">
							<span class="font-semibold">내 취향으로 조정됨</span>
							<button onclick={resetAdjustments} class="font-medium underline underline-offset-2 hover:opacity-70">
								원래 값으로 되돌리기
							</button>
						</div>
					{/if}

					<ul class="space-y-3">
						{#each activeIngredients as ingredient}
							<li class="pb-2 border-b border-stone-50 last:border-0">
								{#if editMode}
									<div class="flex items-center justify-between">
										<span class="font-medium text-stone-600">{ingredient.name}</span>
										<div class="flex items-center gap-2">
											<span class="text-xs text-stone-400 w-14 text-right">{getDisplayAmount(ingredient)}</span>
											<div class="flex items-center gap-0.5 bg-stone-50 rounded-lg border border-stone-100 px-1 py-0.5">
												<button
													onclick={() => adjustIngredient(ingredient.name, -10)}
													disabled={(pendingAdjustments[ingredient.name] ?? 0) <= -80}
													class="w-5 h-5 flex items-center justify-center text-stone-400 hover:text-stone-600 disabled:opacity-30 font-bold text-sm"
												>−</button>
												<span class="text-xs font-bold w-10 text-center {(pendingAdjustments[ingredient.name] ?? 0) !== 0 ? 'text-[#7C9A7E]' : 'text-stone-400'}">
													{(pendingAdjustments[ingredient.name] ?? 0) > 0 ? '+' : ''}{pendingAdjustments[ingredient.name] ?? 0}%
												</span>
												<button
													onclick={() => adjustIngredient(ingredient.name, 10)}
													disabled={(pendingAdjustments[ingredient.name] ?? 0) >= 100}
													class="w-5 h-5 flex items-center justify-center text-stone-400 hover:text-stone-600 disabled:opacity-30 font-bold text-sm"
												>+</button>
											</div>
										</div>
									</div>
								{:else}
									<div class="flex justify-between items-center text-stone-600">
										<span class="font-medium">{ingredient.name}</span>
										<span class="text-sm text-stone-400">{getDisplayAmount(ingredient)}</span>
									</div>
								{/if}
							</li>
						{/each}
					</ul>

					<!-- 편집 모드 -->
					{#if editMode}
						<div class="mt-6 space-y-3">
							<textarea
								bind:value={pendingMemo}
								placeholder="메모 (예: 좀 짰음, 다음엔 소금 덜)"
								rows="3"
								class="w-full rounded-xl border border-stone-100 bg-stone-50 px-3 py-2 text-sm text-stone-600 placeholder-stone-300 outline-none focus:border-[#7C9A7E] resize-none"
							></textarea>
							<div class="flex gap-2">
								<button
									onclick={cancelEdit}
									class="flex-1 rounded-xl border border-stone-100 py-2 text-sm font-semibold text-stone-400 hover:bg-stone-50 transition-colors"
								>취소</button>
								<button
									onclick={saveNote}
									disabled={saving}
									class="flex-1 rounded-xl bg-[#7C9A7E] py-2 text-sm font-semibold text-white hover:bg-[#6a8a6c] disabled:opacity-50 transition-colors"
								>{saving ? '저장 중...' : '저장'}</button>
							</div>
						</div>
					{:else}
						<button
							onclick={enterEditMode}
							class="mt-5 w-full rounded-xl border border-stone-100 py-2 text-xs font-semibold text-stone-400 hover:border-[#7C9A7E] hover:text-[#7C9A7E] transition-colors"
						>내 취향으로 조정</button>
						{#if savedMemo}
							<p class="mt-3 text-xs text-stone-400 italic">{savedMemo}</p>
						{/if}
					{/if}
				</div>

				<!-- 요리방법 패널 -->
				<div class="rounded-3xl bg-[#F5F2EC] p-8 {activeTab !== 'steps' ? 'hidden md:block' : ''}">
					<h2 class="mb-6 text-xl font-bold text-stone-800">요리방법</h2>
					<div class="space-y-6">
						{#if recipe.instructions && recipe.instructions.length > 0}
							{#each recipe.instructions as step, i}
								<div class="flex gap-4">
									<span class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-[#7C9A7E] text-xs font-bold text-white mt-0.5">
										{i + 1}
									</span>
									<p class="text-[15px] leading-relaxed text-stone-600">{step}</p>
								</div>
							{/each}
						{:else}
							<p class="leading-relaxed text-stone-400 italic">등록된 요리방법이 없습니다.</p>
						{/if}
					</div>
				</div>
			</div>
		</div>
	{:else}
		<div class="flex flex-col items-center justify-center py-20 text-stone-400">
			<div class="mb-4 h-12 w-12 animate-spin rounded-full border-4 border-stone-200 border-t-[#7C9A7E]"></div>
			<p>레시피를 불러오는 중입니다...</p>
		</div>
	{/if}
</div>
