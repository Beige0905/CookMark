<script lang="ts">
	import { goto } from '$app/navigation';
	import Link from 'lucide-svelte/icons/link';

	const todayRecipe = {
		title: '된장찌개',
		subtitle: '깊고 구수한 한 그릇',
		time: '25분',
		tag: '한식'
	};

	const recipes = [
		{
			id: 1,
			title: '아보카도 토스트',
			tags: ['브런치', '10분']
		},
		{
			id: 2,
			title: '스테이크 샐러드',
			tags: ['건강식', '20분']
		},
		{
			id: 3,
			title: '토마토 파스타',
			tags: ['양식', '15분']
		}
	];

	let linkInput = $state('');

	function handleLinkSubmit() {
		if (!linkInput.trim()) return;
		alert(`링크 저장 예정: ${linkInput}`);
		linkInput = '';
	}

	function navigateToRecipe(id: number) {
		goto(`/recipes/${id}`);
	}
</script>

<!-- URL Input Bar -->
<div class="mb-10">
	<div
		class="flex items-center gap-3 rounded-2xl border border-stone-200 bg-stone-50/80 px-5 py-4 transition-all focus-within:border-[#7C9A7E] focus-within:bg-white focus-within:ring-4 focus-within:ring-[#7C9A7E]/5"
	>
		<Link class="h-5 w-5 shrink-0 text-stone-400" />
		<input
			bind:value={linkInput}
			onkeydown={(e) => e.key === 'Enter' && handleLinkSubmit()}
			type="url"
			placeholder="유튜브·인스타 링크 붙여넣기..."
			class="flex-1 bg-transparent text-base text-stone-700 outline-none placeholder:text-stone-400"
		/>
		{#if linkInput}
			<button
				onclick={handleLinkSubmit}
				class="shrink-0 rounded-xl bg-[#7C9A7E] px-4 py-1.5 text-sm font-semibold text-white"
				>저장</button
			>
		{/if}
	</div>
</div>

<div class="grid grid-cols-1 gap-10 lg:grid-cols-2">
	<!-- Left Side of Grid: Hero & Stats(mobile) -->
	<div class="space-y-10">
		<!-- Hero -->
		<div>
			<p class="mb-4 text-[11px] font-semibold uppercase tracking-widest text-stone-400">
				오늘의 추천
			</p>
			<div class="relative h-64 overflow-hidden rounded-[2rem] shadow-lg lg:h-80 bg-gradient-to-br from-[#7C9A7E] to-[#4A6B4C]">
				<div class="absolute inset-0 bg-gradient-to-t from-black/50 via-transparent to-transparent"></div>
				<div class="absolute top-5 left-6">
					<span
						class="rounded-full bg-[#7C9A7E] px-3 py-1 text-xs font-semibold text-white shadow-sm"
						>{todayRecipe.tag}</span
					>
				</div>
				<div class="absolute top-5 right-6">
					<span
						class="rounded-full bg-white/20 px-3 py-1 text-xs font-medium text-white backdrop-blur-md"
						>⏱ {todayRecipe.time}</span
					>
				</div>
				<div class="absolute bottom-0 left-0 right-0 p-8">
					<p class="mb-1 text-sm font-medium text-white/70">{todayRecipe.subtitle}</p>
					<h2 class="text-3xl font-bold leading-tight text-white">{todayRecipe.title}</h2>
				</div>
			</div>
		</div>

		<!-- Stats (Mobile Only) -->
		<div class="grid grid-cols-2 gap-4 lg:hidden">
			<div class="rounded-2xl bg-[#F5F2EC] p-5">
				<p class="mb-1 text-[11px] font-medium text-stone-500 uppercase">총 레시피</p>
				<p class="text-3xl font-bold text-stone-800">12</p>
			</div>
			<div class="rounded-2xl bg-[#EEF3EE] p-5">
				<p class="mb-1 text-[11px] font-medium text-stone-500 uppercase">이번 달 요리</p>
				<p class="text-3xl font-bold text-stone-800">5<span class="ml-1 text-sm font-normal text-stone-400">회</span></p>
			</div>
		</div>
	</div>

	<!-- Right Side of Grid: Recipe List -->
	<div>
		<div class="mb-5 flex items-center justify-between">
			<h2 class="text-lg font-bold text-stone-900">내 레시피 보관함</h2>
			<button class="text-sm font-semibold text-[#7C9A7E] hover:underline">전체 보기 →</button>
		</div>
		<div class="flex flex-col divide-y divide-stone-100">
			{#each recipes as recipe}
				<div 
					role="button"
					tabindex="0"
					onclick={() => navigateToRecipe(recipe.id)}
					onkeydown={(e) => e.key === 'Enter' && navigateToRecipe(recipe.id)}
					class="group flex items-center gap-4 py-4 cursor-pointer"
				>
					<div class="min-w-0 flex-1">
						<p
							class="truncate text-base font-semibold text-stone-800 transition-colors group-hover:text-[#7C9A7E]"
						>
							{recipe.title}
						</p>
						<div class="mt-1.5 flex flex-wrap gap-1.5">
							{#each recipe.tags as tag}
								<span
									class="rounded-full bg-stone-100 px-2.5 py-0.5 text-[11px] font-medium text-stone-500 transition-colors group-hover:bg-[#EBF0EB] group-hover:text-[#6A8C6C]"
									>{tag}</span
								>
							{/each}
						</div>
					</div>
					<button
						onclick={(e) => { e.stopPropagation(); alert('만들었어!'); }}
						class="shrink-0 whitespace-nowrap rounded-xl bg-stone-50 px-3 py-2 text-xs font-bold text-[#7C9A7E] opacity-0 transition-all hover:bg-[#7C9A7E] hover:text-white group-hover:opacity-100"
						>만들었어! 🔥</button
					>
				</div>
			{/each}
		</div>
	</div>
</div>
