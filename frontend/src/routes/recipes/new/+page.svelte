<script lang="ts">
	import { goto } from '$app/navigation';
	import { toast } from '$lib/toast.svelte';
	import { extractFromYouTube, extractFromImage, createRecipe, type Ingredient } from '$lib/api/recipes';
	import Loader2 from 'lucide-svelte/icons/loader-2';
	import ArrowLeft from 'lucide-svelte/icons/arrow-left';
	import Camera from 'lucide-svelte/icons/camera';
	import Pencil from 'lucide-svelte/icons/pencil';
	import YouTube from '$lib/components/icons/YouTube.svelte';
	import Plus from 'lucide-svelte/icons/plus';
	import Minus from 'lucide-svelte/icons/minus';

	let platform = $state('youtube');
	let url = $state('');
	let title = $state('');
	let instructions = $state<string[]>(['']);
	let ingredients = $state<Ingredient[]>([]);
	let originalIngredients = $state<Ingredient[]>([]);
	let baseServings = $state(1);
	let currentServings = $state(1);
	let isSubmitting = $state(false);
	let isExtracting = $state(false);

	// 플랫폼 변경 시 폼 초기화
	$effect(() => {
		if (platform) {
			url = '';
			title = '';
			instructions = [''];
			ingredients = [];
			originalIngredients = [];
			baseServings = 1;
			currentServings = 1;
		}
	});

	function formatAmount(ing: Ingredient, servings: number, base: number): string {
		if (!ing.amount_num || !ing.unit) return ing.note || '';
		
		let amountNum = ing.amount_num;
		const scalingType = ing.scaling_type || 'linear';
		const factor = ing.scaling_factor || 1.0;

		if (servings !== base) {
			// 1. 1인분 기준 용량 산출
			let unitAmount = 0;
			if (scalingType === 'linear') {
				unitAmount = ing.amount_num / base;
			} else {
				unitAmount = ing.amount_num / (1 + (base - 1) * factor);
			}

			// 2. 목표 인분으로 확장
			if (scalingType === 'linear') {
				amountNum = unitAmount * servings;
			} else {
				amountNum = unitAmount * (1 + (servings - 1) * factor);
			}
		}

		// 소수점 처리 (0.5 단위로 반올림하거나 소수점 1자리까지)
		const rounded = Math.round(amountNum * 10) / 10;
		const stdPart = `${rounded}${ing.unit}`;
		
		return ing.note ? `${stdPart} (${ing.note})` : stdPart;
	}

	function updateScaledIngredients() {
		ingredients = originalIngredients.map(ing => ({
			...ing,
			amount: formatAmount(ing, currentServings, baseServings)
		}));
	}

	async function compressImage(file: File): Promise<Blob> {
		return new Promise((resolve, reject) => {
			const reader = new FileReader();
			reader.onload = (e) => {
				const img = new Image();
				img.onload = () => {
					const canvas = document.createElement('canvas');
					let width = img.width;
					let height = img.height;
					const max_size = 1200;

					if (width > height) {
						if (width > max_size) {
							height *= max_size / width;
							width = max_size;
						}
					} else {
						if (height > max_size) {
							width *= max_size / height;
							height = max_size;
						}
					}

					canvas.width = width;
					canvas.height = height;
					const ctx = canvas.getContext('2d');
					ctx?.drawImage(img, 0, 0, width, height);

					canvas.toBlob(
						(blob) => {
							if (blob) resolve(blob);
							else reject(new Error('압축 실패'));
						},
						'image/jpeg',
						0.7
					);
				};
				img.onerror = reject;
				img.src = e.target?.result as string;
			};
			reader.onerror = reject;
			reader.readAsDataURL(file);
		});
	}

	async function handleExtract() {
		if (!url) return;
		isExtracting = true;
		try {
			const result = await extractFromYouTube(url);
			if (!title) title = result.title;
			baseServings = result.base_servings || 1;
			currentServings = baseServings;
			originalIngredients = result.ingredients;
			updateScaledIngredients();

			if (ingredients.length === 0) {
				toast.add('영상 설명글에서 재료를 찾지 못했습니다.', 'error');
			} else {
				toast.add(`${baseServings}인분 기준 레시피 정보를 가져왔습니다!`, 'success');
			}
		} catch (error: any) {
			toast.add(error.message || '재료 추출 중 오류가 발생했습니다.', 'error');
		} finally {
			isExtracting = false;
		}
	}

	async function handleImageUpload(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		if (!file) return;

		isExtracting = true;
		try {
			const compressedBlob = await compressImage(file);
			const compressedFile = new File([compressedBlob], file.name, { type: 'image/jpeg' });
			
			const result = await extractFromImage(compressedFile);
			if (!title) title = result.title;
			baseServings = result.base_servings || 1;
			currentServings = baseServings;
			originalIngredients = result.ingredients;
			
			updateScaledIngredients();
			
			instructions = result.instructions.length > 0 ? result.instructions : [''];
			toast.add('사진에서 레시피 정보를 추출했습니다!', 'success');
		} catch (error: any) {
			toast.add(error.message || '사진 분석 중 오류가 발생했습니다.', 'error');
		} finally {
			isExtracting = false;
			target.value = '';
		}
	}

	function changeServings(delta: number) {
		const next = currentServings + delta;
		if (next < 1) return;
		currentServings = next;
		updateScaledIngredients();
	}

	function addIngredient() {
		ingredients = [...ingredients, { name: '', amount: '' }];
	}

	function removeIngredient(index: number) {
		ingredients = ingredients.filter((_, i) => i !== index);
	}

	function addInstruction() {
		instructions = [...instructions, ''];
	}

	function removeInstruction(index: number) {
		if (instructions.length > 1) {
			instructions = instructions.filter((_, i) => i !== index);
		} else {
			instructions[0] = '';
		}
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		isSubmitting = true;
		try {
			const parsedInstructions = platform !== 'youtube'
				? instructions
					.map((line) => line.trim())
					.filter((line) => line.length > 0)
				: [];

			await createRecipe({
				title,
				origin_url: platform === 'youtube' ? url : undefined,
				base_servings: currentServings, // 조절된 인분 수로 저장
				ingredients: ingredients.map(ing => ({
					...ing,
				})),
				instructions: parsedInstructions
			});

			toast.add('레시피가 성공적으로 등록되었습니다!', 'success');
			goto('/recipes');
		} catch (error: any) {
			toast.add(error.message || '등록 중 오류가 발생했습니다. 다시 시도해주세요.', 'error');
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>새 레시피 등록 | Cookmark</title>
</svelte:head>

<div class="max-w-2xl mx-auto pb-20">
	<header class="mb-10">
		<h1 class="text-3xl font-bold text-stone-900 tracking-tight mb-2">새 레시피 등록</h1>
		<p class="text-stone-500">레시피 링크를 저장하거나 사진을 찍어 자동으로 등록하세요.</p>
	</header>

	<form onsubmit={handleSubmit} class="space-y-8 bg-white p-8 rounded-3xl border border-stone-100 shadow-sm">
		<!-- Platform Selection -->
		<div class="space-y-3">
			<label class="block text-sm font-bold text-stone-700" for="platform">등록 방식</label>
			<div class="grid grid-cols-3 gap-3">
				<button
					type="button"
					onclick={() => (platform = 'youtube')}
					class="flex flex-col items-center justify-center gap-2 px-3 py-4 rounded-2xl border-2 transition-all {platform === 'youtube' ? 'border-[#7C9A7E] bg-[#EEF3EE] text-[#5A7A5C]' : 'border-stone-100 text-stone-400 hover:border-stone-200'}"
				>
					<YouTube class="w-6 h-6" />
					<span class="font-bold text-sm">YouTube</span>
				</button>
				<button
					type="button"
					onclick={() => (platform = 'photo')}
					class="flex flex-col items-center justify-center gap-2 px-3 py-4 rounded-2xl border-2 transition-all {platform === 'photo' ? 'border-[#7C9A7E] bg-[#EEF3EE] text-[#5A7A5C]' : 'border-stone-100 text-stone-400 hover:border-stone-200'}"
				>
					<Camera class="w-6 h-6" />
					<span class="font-bold text-sm">사진으로 추가</span>
				</button>
				<button
					type="button"
					onclick={() => (platform = 'manual')}
					class="flex flex-col items-center justify-center gap-2 px-3 py-4 rounded-2xl border-2 transition-all {platform === 'manual' ? 'border-[#7C9A7E] bg-[#EEF3EE] text-[#5A7A5C]' : 'border-stone-100 text-stone-400 hover:border-stone-200'}"
				>
					<Pencil class="w-6 h-6" />
					<span class="font-bold text-sm">직접 입력</span>
				</button>
			</div>
		</div>

		<!-- Input Area -->
		<div class="space-y-3">
			{#if platform === 'photo'}
				<label class="block text-sm font-bold text-stone-700" for="image">레시피 사진 분석</label>
				<div class="relative">
					<input
						id="image"
						type="file"
						accept="image/*"
						onchange={handleImageUpload}
						class="hidden"
						disabled={isExtracting}
					/>
					<label
						for="image"
						class="flex flex-col items-center justify-center w-full p-8 border-2 border-dashed border-stone-200 rounded-2xl hover:border-[#7C9A7E] hover:bg-stone-50 cursor-pointer transition-all {isExtracting ? 'opacity-50 cursor-not-allowed' : ''}"
					>
						{#if isExtracting}
							<Loader2 class="animate-spin h-8 w-8 text-[#7C9A7E] mb-2" />
							<span class="text-stone-500 font-medium">Gemini가 사진을 분석하고 있어요...</span>
						{:else}
							<Camera class="h-8 w-8 text-stone-300 mb-2" />
							<span class="text-stone-700 font-bold">사진 업로드 또는 촬영</span>
							<span class="text-stone-400 text-xs mt-1">종이 레시피, 인스타 스크린샷 등을 올려주세요</span>
						{/if}
					</label>
				</div>
			{:else if platform === 'youtube'}
				<label class="block text-sm font-bold text-stone-700" for="url">링크 주소</label>
				<input
					id="url"
					type="url"
					bind:value={url}
					required
					placeholder="https://youtube.com/..."
					class="w-full px-5 py-4 rounded-2xl border-2 border-stone-100 focus:border-[#7C9A7E] focus:outline-none transition-colors bg-stone-50/50"
				/>
				<button
					type="button"
					onclick={handleExtract}
					disabled={!url || isExtracting}
					class="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-[#EEF3EE] text-[#5A7A5C] font-bold text-sm hover:bg-[#dde8de] disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
				>
					{#if isExtracting}
						<Loader2 class="animate-spin h-4 w-4" />
						재료 가져오는 중...
					{:else}
						<YouTube class="w-4 h-4" />
						재료 가져오기
					{/if}
				</button>
			{/if}
		</div>

		<!-- Title & Servings -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			<div class="md:col-span-2 space-y-3">
				<label class="block text-sm font-bold text-stone-700" for="title">레시피 이름</label>
				<input
					id="title"
					type="text"
					bind:value={title}
					required
					placeholder="맛있는 김치찌개 레시피"
					class="w-full px-5 py-4 rounded-2xl border-2 border-stone-100 focus:border-[#7C9A7E] focus:outline-none transition-colors bg-stone-50/50"
				/>
			</div>
			<div class="space-y-3">
				<label class="block text-sm font-bold text-stone-700" for="servings">기준 인분</label>
				<div class="flex items-center gap-1">
					<button
						type="button"
						onclick={() => changeServings(-1)}
						class="p-4 rounded-2xl border-2 border-stone-100 hover:border-stone-200 text-stone-400"
					>
						<Minus size={18} />
					</button>
					<div class="flex-1 text-center font-bold text-lg text-stone-800 bg-stone-50/50 py-3.5 rounded-2xl border-2 border-stone-100">
						{currentServings}
					</div>
					<button
						type="button"
						onclick={() => changeServings(1)}
						class="p-4 rounded-2xl border-2 border-stone-100 hover:border-stone-200 text-stone-400"
					>
						<Plus size={18} />
					</button>
				</div>
			</div>
		</div>

		<!-- Ingredients -->
		<div class="space-y-3">
			<div class="flex items-center justify-between">
				<span class="block text-sm font-bold text-stone-700">재료</span>
				<button
					type="button"
					onclick={addIngredient}
					class="text-xs font-bold text-[#7C9A7E] hover:text-[#5A7A5C] transition-colors"
				>
					+ 재료 추가
				</button>
			</div>
			{#if ingredients.length > 0}
				<div class="space-y-2">
					{#each ingredients as ingredient, i}
						<div class="flex gap-2 items-center">
							<input
								type="text"
								bind:value={ingredient.name}
								placeholder="재료명"
								class="flex-1 px-4 py-3 rounded-xl border-2 border-stone-100 focus:border-[#7C9A7E] focus:outline-none transition-colors bg-stone-50/50 text-sm"
								aria-label={`${i + 1}번째 재료 이름`}
							/>
							<input
								type="text"
								bind:value={ingredient.amount}
								placeholder="용량"
								class="w-40 px-4 py-3 rounded-xl border-2 border-stone-100 focus:border-[#7C9A7E] focus:outline-none transition-colors bg-stone-50/50 text-sm"
								aria-label={`${i + 1}번째 재료 용량`}
							/>
							<button
								type="button"
								onclick={() => removeIngredient(i)}
								class="p-2 text-stone-300 hover:text-red-400 transition-colors rounded-lg"
								aria-label={`${i + 1}번째 재료 삭제`}
							>
								✕
							</button>
						</div>
					{/each}
				</div>
			{:else}
				<p class="text-sm text-stone-400 text-center py-4 border-2 border-dashed border-stone-100 rounded-2xl">
					재료가 없습니다. 사진을 분석하거나 직접 추가하세요.
				</p>
			{/if}
		</div>

		<!-- Instructions Input -->
		{#if platform !== 'youtube'}
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<span class="block text-sm font-bold text-stone-700">조리법</span>
					<button
						type="button"
						onclick={addInstruction}
						class="text-xs font-bold text-[#7C9A7E] hover:text-[#5A7A5C] transition-colors"
					>
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
								placeholder={`${i + 1}단계 설명을 입력하세요.`}
								rows="2"
								class="flex-1 px-4 py-3 rounded-xl border-2 border-stone-100 focus:border-[#7C9A7E] focus:outline-none transition-colors bg-stone-50/50 text-sm resize-none"
								aria-label={`${i + 1}단계 조리 방법`}
							></textarea>
							<button
								type="button"
								onclick={() => removeInstruction(i)}
								class="mt-2 p-2 text-stone-300 hover:text-red-400 transition-colors rounded-lg"
								aria-label={`${i + 1}단계 삭제`}
							>
								✕
							</button>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<!-- Submit Button -->
		<div class="pt-4">
			<button
				type="submit"
				disabled={isSubmitting || !title}
				class="w-full py-4 bg-[#7C9A7E] hover:bg-[#6A8C6C] disabled:bg-stone-300 text-white rounded-2xl font-bold text-lg transition-all shadow-lg shadow-[#7C9A7E]/20 active:scale-[0.98]"
			>
				{#if isSubmitting}
					<span class="flex items-center justify-center gap-2">
						<Loader2 class="animate-spin h-5 w-5" />
						등록 중...
					</span>
				{:else}
					레시피 저장하기
				{/if}
			</button>
		</div>
	</form>

	<div class="mt-8 text-center">
		<a href="/recipes" class="inline-flex items-center gap-2 text-sm font-bold text-stone-400 hover:text-stone-600 transition-colors underline underline-offset-4">
			<ArrowLeft size={16} />
			레시피 목록으로 돌아가기
		</a>
	</div>
</div>
