<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { toast } from '$lib/toast.svelte';
	import { getRecipe, updateRecipe, type Ingredient } from '$lib/api/recipes';
	import ArrowLeft from 'lucide-svelte/icons/arrow-left';
	import Loader2 from 'lucide-svelte/icons/loader-2';
	import Plus from 'lucide-svelte/icons/plus';
	import Minus from 'lucide-svelte/icons/minus';

	const id = $derived(Number($page.params.id));

	let title = $state('');
	let baseServings = $state(1);
	let originUrl = $state('');
	let ingredients = $state<{ name: string; amount: string }[]>([]);
	let instructions = $state<string[]>(['']);
	let isLoading = $state(true);
	let isSubmitting = $state(false);

	function getDisplayAmount(ing: Ingredient): string {
		if (ing.amount) return ing.amount;
		if (ing.amount_num !== undefined && ing.amount_num !== null) {
			return `${ing.amount_num}${ing.unit ?? ''}`;
		}
		return ing.note ?? '';
	}

	$effect(() => {
		getRecipe(id)
			.then((recipe) => {
				title = recipe.title;
				baseServings = recipe.base_servings;
				originUrl = recipe.origin_url ?? '';
				ingredients = recipe.ingredients.map((ing) => ({
					name: ing.name,
					amount: getDisplayAmount(ing)
				}));
				instructions = recipe.instructions?.length ? recipe.instructions : [''];
			})
			.catch(() => toast.add('레시피를 불러오지 못했습니다.', 'error'))
			.finally(() => (isLoading = false));
	});

	function addIngredient() {
		ingredients = [...ingredients, { name: '', amount: '' }];
	}

	function removeIngredient(i: number) {
		ingredients = ingredients.filter((_, idx) => idx !== i);
	}

	function addInstruction() {
		instructions = [...instructions, ''];
	}

	function removeInstruction(i: number) {
		if (instructions.length > 1) {
			instructions = instructions.filter((_, idx) => idx !== i);
		} else {
			instructions[0] = '';
		}
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		isSubmitting = true;
		try {
			await updateRecipe(id, {
				title,
				origin_url: originUrl || undefined,
				base_servings: baseServings,
				ingredients: ingredients
					.filter((ing) => ing.name.trim())
					.map((ing) => ({ name: ing.name, amount: ing.amount })),
				instructions: instructions.map((s) => s.trim()).filter(Boolean)
			});
			toast.add('레시피가 수정됐습니다!', 'success');
			goto(`/recipes/${id}`);
		} catch (e: any) {
			toast.add(e.message || '수정 중 오류가 발생했습니다.', 'error');
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>레시피 수정 | Cookmark</title>
</svelte:head>

<div class="max-w-2xl mx-auto pb-20">
	<header class="mb-10">
		<a href="/recipes/{id}" class="inline-flex items-center gap-2 text-sm font-semibold text-stone-400 hover:text-[#7C9A7E] transition-colors mb-4">
			<ArrowLeft class="h-4 w-4" />
			레시피로 돌아가기
		</a>
		<h1 class="text-3xl font-bold text-stone-900 tracking-tight">레시피 수정</h1>
	</header>

	{#if isLoading}
		<div class="flex justify-center py-20">
			<Loader2 class="animate-spin h-8 w-8 text-[#7C9A7E]" />
		</div>
	{:else}
		<form onsubmit={handleSubmit} class="space-y-8 bg-white p-8 rounded-3xl border border-stone-100 shadow-sm">
			<!-- 제목 & 인분 -->
			<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
				<div class="md:col-span-2 space-y-3">
					<label class="block text-sm font-bold text-stone-700" for="title">레시피 이름</label>
					<input
						id="title"
						type="text"
						bind:value={title}
						required
						class="w-full px-5 py-4 rounded-2xl border-2 border-stone-100 focus:border-[#7C9A7E] focus:outline-none transition-colors bg-stone-50/50"
					/>
				</div>
				<div class="space-y-3">
					<label class="block text-sm font-bold text-stone-700">기준 인분</label>
					<div class="flex items-center gap-1">
						<button
							type="button"
							onclick={() => (baseServings = Math.max(1, baseServings - 1))}
							class="p-4 rounded-2xl border-2 border-stone-100 hover:border-stone-200 text-stone-400"
						>
							<Minus size={18} />
						</button>
						<div class="flex-1 text-center font-bold text-lg text-stone-800 bg-stone-50/50 py-3.5 rounded-2xl border-2 border-stone-100">
							{baseServings}
						</div>
						<button
							type="button"
							onclick={() => (baseServings = Math.min(10, baseServings + 1))}
							class="p-4 rounded-2xl border-2 border-stone-100 hover:border-stone-200 text-stone-400"
						>
							<Plus size={18} />
						</button>
					</div>
				</div>
			</div>

			<!-- 출처 URL -->
			<div class="space-y-3">
				<label class="block text-sm font-bold text-stone-700" for="originUrl">출처 링크 (선택)</label>
				<input
					id="originUrl"
					type="url"
					bind:value={originUrl}
					placeholder="https://..."
					class="w-full px-5 py-4 rounded-2xl border-2 border-stone-100 focus:border-[#7C9A7E] focus:outline-none transition-colors bg-stone-50/50"
				/>
			</div>

			<!-- 재료 -->
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<span class="block text-sm font-bold text-stone-700">재료</span>
					<button type="button" onclick={addIngredient} class="text-xs font-bold text-[#7C9A7E] hover:text-[#5A7A5C]">
						+ 재료 추가
					</button>
				</div>
				<div class="space-y-2">
					{#each ingredients as ingredient, i}
						<div class="flex gap-2 items-center">
							<input
								type="text"
								bind:value={ingredient.name}
								placeholder="재료명"
								class="flex-1 px-4 py-3 rounded-xl border-2 border-stone-100 focus:border-[#7C9A7E] focus:outline-none transition-colors bg-stone-50/50 text-sm"
							/>
							<input
								type="text"
								bind:value={ingredient.amount}
								placeholder="용량"
								class="w-40 px-4 py-3 rounded-xl border-2 border-stone-100 focus:border-[#7C9A7E] focus:outline-none transition-colors bg-stone-50/50 text-sm"
							/>
							<button
								type="button"
								onclick={() => removeIngredient(i)}
								class="p-2 text-stone-300 hover:text-red-400 transition-colors rounded-lg"
							>✕</button>
						</div>
					{/each}
				</div>
			</div>

			<!-- 조리법 -->
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<span class="block text-sm font-bold text-stone-700">조리법</span>
					<button type="button" onclick={addInstruction} class="text-xs font-bold text-[#7C9A7E] hover:text-[#5A7A5C]">
						+ 단계 추가
					</button>
				</div>
				<div class="space-y-3">
					{#each instructions as step, i}
						<div class="flex gap-3 items-start group">
							<div class="mt-3 flex-shrink-0 w-6 h-6 rounded-full bg-stone-100 flex items-center justify-center text-xs font-bold text-stone-500 group-focus-within:bg-[#7C9A7E] group-focus-within:text-white transition-colors">
								{i + 1}
							</div>
							<textarea
								bind:value={instructions[i]}
								placeholder={`${i + 1}단계`}
								rows="2"
								class="flex-1 px-4 py-3 rounded-xl border-2 border-stone-100 focus:border-[#7C9A7E] focus:outline-none transition-colors bg-stone-50/50 text-sm resize-none"
							></textarea>
							<button
								type="button"
								onclick={() => removeInstruction(i)}
								class="mt-2 p-2 text-stone-300 hover:text-red-400 transition-colors rounded-lg"
							>✕</button>
						</div>
					{/each}
				</div>
			</div>

			<div class="pt-4">
				<button
					type="submit"
					disabled={isSubmitting || !title}
					class="w-full py-4 bg-[#7C9A7E] hover:bg-[#6A8C6C] disabled:bg-stone-300 text-white rounded-2xl font-bold text-lg transition-all shadow-lg shadow-[#7C9A7E]/20 active:scale-[0.98]"
				>
					{#if isSubmitting}
						<span class="flex items-center justify-center gap-2">
							<Loader2 class="animate-spin h-5 w-5" />
							저장 중...
						</span>
					{:else}
						수정 저장하기
					{/if}
				</button>
			</div>
		</form>
	{/if}
</div>
