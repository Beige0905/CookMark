<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getRecipe, deleteRecipe, type Recipe, type Ingredient } from '$lib/api/recipes';
	import { getCookingLogs, addCookingLog, deleteCookingLog, type CookingLog } from '$lib/api/cooking-logs';
	import { getPantryMatchesForRecipe, type PantryItem } from '$lib/api/pantry';
	import { toast } from '$lib/toast.svelte';
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import Users from 'lucide-svelte/icons/users';
	import Pencil from 'lucide-svelte/icons/pencil';
	import Trash2 from 'lucide-svelte/icons/trash-2';

	let recipe = $state<Recipe | null>(null);
	let servingSize = $state(1);
	let error = $state('');
	let deleting = $state(false);

	let logs = $state<CookingLog[]>([]);
	let addingLog = $state(false);

	// 만들었어 모달
	let modalOpen = $state(false);
	let logComment = $state('');
	let logDate = $state('');
	let pantryMatches = $state<PantryItem[]>([]);
	let selectedPantryIDs = $state<Set<number>>(new Set());
	let loadingMatches = $state(false);

	const recipeId = $derived(Number($page.params.id));

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

	let activeIngredients = $derived.by(() => {
		if (!recipe) return [];
		return recipe.ingredients.map((ing) => scaleIngredient(ing, servingSize));
	});

	function scaleIngredient(ing: Ingredient, size: number): Ingredient {
		if (ing.amount_num) {
			const factor = ing.scaling_factor ?? 1.0;
			const scaled = ing.amount_num * (1 + (size - 1) * factor);
			return { ...ing, amount_num: roundByUnit(scaled, ing.unit || ''), amount: undefined };
		}
		return ing;
	}

function getYouTubeThumbnail(url: string): string {
		try {
			let videoId = '';
			if (url.includes('youtu.be/')) {
				videoId = url.split('youtu.be/')[1]?.split(/[?#]/)[0] ?? '';
			} else {
				const u = new URL(url);
				if (u.pathname.startsWith('/shorts/')) {
					videoId = u.pathname.split('/shorts/')[1]?.split('/')[0] ?? '';
				} else {
					videoId = u.searchParams.get('v') ?? '';
				}
			}
			return videoId ? `https://img.youtube.com/vi/${videoId}/hqdefault.jpg` : '';
		} catch {
			return '';
		}
	}

	function getDisplayAmount(ing: Ingredient): string {
		if (ing.amount) return ing.amount;
		if (ing.amount_num !== undefined && ing.amount_num !== null) {
			return formatAmount(ing.amount_num, ing.unit || '');
		}
		return ing.note || '';
	}

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

	async function openLogModal() {
		modalOpen = true;
		logComment = '';
		logDate = '';
		loadingMatches = true;
		try {
			const matches = await getPantryMatchesForRecipe(recipeId);
			pantryMatches = matches;
			selectedPantryIDs = new Set(matches.map((m) => m.id));
		} catch {
			pantryMatches = [];
			selectedPantryIDs = new Set();
		} finally {
			loadingMatches = false;
		}
	}

	function togglePantryItem(id: number) {
		const next = new Set(selectedPantryIDs);
		if (next.has(id)) next.delete(id);
		else next.add(id);
		selectedPantryIDs = next;
	}

	function selectAll() {
		selectedPantryIDs = new Set(pantryMatches.map((m) => m.id));
	}

	function deselectAll() {
		selectedPantryIDs = new Set();
	}

	async function handleAddLog(event: SubmitEvent) {
		event.preventDefault();
		addingLog = true;
		try {
			const log = await addCookingLog(recipeId, {
				comment: logComment || undefined,
				cooked_at: logDate || undefined,
				delete_pantry_ids: selectedPantryIDs.size > 0 ? [...selectedPantryIDs] : undefined
			});
			logs = [log, ...logs];
			modalOpen = false;
			toast.add(
				selectedPantryIDs.size > 0
					? `요리 기록이 추가됐습니다! 냉장고 재료 ${selectedPantryIDs.size}개를 소비했습니다.`
					: '요리 기록이 추가됐습니다!',
				'success'
			);
		} catch (e: any) {
			toast.add(e.message || '기록 추가 실패', 'error');
		} finally {
			addingLog = false;
		}
	}

	async function handleDeleteLog(logId: number) {
		try {
			await deleteCookingLog(recipeId, logId);
			logs = logs.filter((l) => l.id !== logId);
		} catch (e: any) {
			toast.add(e.message || '기록 삭제 실패', 'error');
		}
	}

	$effect(() => {
		const id = recipeId;
		getRecipe(id)
			.then((data) => {
				recipe = data;
				servingSize = data.base_servings;
			})
			.catch((e) => (error = e.message));
		getCookingLogs(id)
			.then((data) => (logs = data))
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

			<div class="space-y-6">
				<!-- 재료 패널 -->
				<div class="rounded-3xl bg-white p-8 border border-stone-100 shadow-sm">
					<div class="mb-6 flex items-center justify-between">
						<h2 class="text-xl font-bold text-stone-800">재료</h2>
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
					<ul class="space-y-3">
						{#each activeIngredients as ingredient}
							<li class="pb-2 border-b border-stone-50 last:border-0">
								<div class="flex justify-between items-center text-stone-600">
									<span class="font-medium">{ingredient.name}</span>
									<span class="text-sm text-stone-400">{getDisplayAmount(ingredient)}</span>
								</div>
							</li>
						{/each}
					</ul>
				</div>

				<!-- 요리방법 패널 -->
				<div class="rounded-3xl bg-[#F5F2EC] p-8">
					<h2 class="mb-6 text-xl font-bold text-stone-800">요리방법</h2>
					<div class="space-y-6">
						{#if recipe.origin_url}
							{@const thumbnail = recipe.image_url || getYouTubeThumbnail(recipe.origin_url)}
							<a
								href={recipe.origin_url}
								target="_blank"
								rel="noopener noreferrer"
								class="group block overflow-hidden rounded-2xl"
							>
								{#if thumbnail}
									<div class="relative aspect-video w-full">
										<img
											src={thumbnail}
											alt={recipe.title}
											class="h-full w-full object-cover transition-opacity group-hover:opacity-90"
										/>
										<div class="absolute inset-0 flex items-center justify-center">
											<div class="flex items-center gap-2 rounded-full bg-black/70 px-5 py-2.5 text-sm font-bold text-white transition-transform group-hover:scale-105">
												<svg class="h-4 w-4 fill-white" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
												YouTube에서 보기
											</div>
										</div>
									</div>
								{:else}
									<div class="flex items-center gap-3 rounded-2xl border-2 border-stone-200 bg-white px-5 py-4 transition-colors hover:border-[#7C9A7E]">
										<svg class="h-5 w-5 fill-red-500 shrink-0" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
										<span class="font-bold text-stone-700">YouTube에서 보기</span>
									</div>
								{/if}
							</a>
						{/if}
						{#if recipe.instructions && recipe.instructions.length > 0}
							<div class="space-y-4 {recipe.origin_url ? 'mt-6' : ''}">
								{#each recipe.instructions as step, i}
									<div class="flex gap-4">
										<span class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-[#7C9A7E] text-xs font-bold text-white mt-0.5">
											{i + 1}
										</span>
										<p class="text-[15px] leading-relaxed text-stone-600">{step}</p>
									</div>
								{/each}
							</div>
						{:else if !recipe.origin_url}
							<p class="leading-relaxed text-stone-400 italic">등록된 요리방법이 없습니다.</p>
						{/if}
					</div>
				</div>
			</div>

			<!-- 요리 기록 -->
			<div class="rounded-3xl bg-white border border-stone-100 shadow-sm p-8 space-y-6">
				<div class="flex items-center justify-between">
					<h2 class="text-xl font-bold text-stone-800">요리 기록</h2>
					<button
						onclick={openLogModal}
						class="rounded-xl bg-[#7C9A7E] px-4 py-2 text-sm font-bold text-white hover:bg-[#6A8C6C] transition-colors"
					>
						오늘 만들었어! 🍳
					</button>
				</div>

				{#if logs.length === 0}
					<p class="text-center text-sm text-stone-400 py-4">아직 요리 기록이 없습니다.</p>
				{:else}
					<ul class="space-y-3">
						{#each logs as log (log.id)}
							<li class="flex items-start justify-between gap-4 border-b border-stone-50 pb-3 last:border-0 last:pb-0">
								<div class="space-y-1">
									<p class="text-xs font-semibold text-stone-400">
										{new Date(log.cooked_at).toLocaleDateString('ko-KR')}
									</p>
									{#if log.comment}
										<p class="text-sm text-stone-600">{log.comment}</p>
									{/if}
								</div>
								<button
									onclick={() => handleDeleteLog(log.id)}
									class="text-stone-300 hover:text-red-400 transition-colors text-lg leading-none shrink-0"
									aria-label="기록 삭제"
								>×</button>
							</li>
						{/each}
					</ul>
				{/if}
			</div>
		</div>
	{:else}
		<div class="flex flex-col items-center justify-center py-20 text-stone-400">
			<div class="mb-4 h-12 w-12 animate-spin rounded-full border-4 border-stone-200 border-t-[#7C9A7E]"></div>
			<p>레시피를 불러오는 중입니다...</p>
		</div>
	{/if}
</div>

<!-- 만들었어 모달 -->
{#if modalOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
		role="presentation"
		onclick={(e) => { if (e.target === e.currentTarget) modalOpen = false; }}
		onkeydown={(e) => { if (e.key === 'Escape') modalOpen = false; }}
	>
		<div class="w-full max-w-md rounded-3xl bg-white p-8 shadow-xl space-y-6">
			<h3 class="text-lg font-bold text-stone-800">오늘 만들었어! 🍳</h3>

			<form onsubmit={handleAddLog} class="space-y-5">
				<!-- 냉장고 재료 -->
				{#if loadingMatches}
					<div class="text-sm text-stone-400">냉장고 재료 확인 중...</div>
				{:else if pantryMatches.length > 0}
					<div class="space-y-3">
						<div class="flex items-center justify-between">
							<p class="text-sm font-semibold text-stone-600">소비한 재료 선택</p>
							<div class="flex gap-2 text-xs font-semibold">
								<button type="button" onclick={selectAll} class="text-[#7C9A7E] hover:underline">전체선택</button>
								<span class="text-stone-300">|</span>
								<button type="button" onclick={deselectAll} class="text-stone-400 hover:underline">전체해제</button>
							</div>
						</div>
						<ul class="space-y-2 max-h-40 overflow-y-auto">
							{#each pantryMatches as item (item.id)}
								<li>
									<label class="flex items-center gap-3 cursor-pointer group">
										<input
											type="checkbox"
											checked={selectedPantryIDs.has(item.id)}
											onchange={() => togglePantryItem(item.id)}
											class="h-4 w-4 rounded border-stone-300 accent-[#7C9A7E]"
										/>
										<span class="text-sm text-stone-700 group-hover:text-stone-900">{item.name}</span>
									</label>
								</li>
							{/each}
						</ul>
					</div>
				{/if}

				<!-- 날짜 -->
				<div class="space-y-1">
					<label for="log-date" class="text-sm font-semibold text-stone-600">날짜</label>
					<input
						id="log-date"
						type="date"
						bind:value={logDate}
						class="w-full rounded-xl border border-stone-200 px-3 py-2 text-sm outline-none focus:border-[#7C9A7E] transition-colors"
					/>
				</div>

				<!-- 코멘트 -->
				<div class="space-y-1">
					<label for="log-comment" class="text-sm font-semibold text-stone-600">메모 <span class="font-normal text-stone-400">(선택)</span></label>
					<textarea
						id="log-comment"
						bind:value={logComment}
						placeholder="오늘 어땠나요?"
						rows="2"
						class="w-full rounded-xl border border-stone-200 px-4 py-3 text-sm outline-none focus:border-[#7C9A7E] transition-colors resize-none placeholder-stone-300"
					></textarea>
				</div>

				<!-- 버튼 -->
				<div class="flex gap-3 pt-1">
					<button
						type="button"
						onclick={() => (modalOpen = false)}
						class="flex-1 rounded-xl border border-stone-200 py-2.5 text-sm font-semibold text-stone-500 hover:border-stone-300 transition-colors"
					>
						취소
					</button>
					<button
						type="submit"
						disabled={addingLog}
						class="flex-1 rounded-xl bg-[#7C9A7E] py-2.5 text-sm font-bold text-white disabled:opacity-40 hover:bg-[#6A8C6C] transition-colors"
					>
						{addingLog ? '저장 중...' : '저장'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
