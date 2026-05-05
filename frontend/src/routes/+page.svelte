<script lang="ts">
	const todayRecipe = {
		title: '된장찌개',
		subtitle: '깊고 구수한 한 그릇',
		time: '25분',
		tag: '한식',
		image: 'https://images.unsplash.com/photo-1498654896293-37aaa4f85252?w=800&auto=format&fit=crop&q=80'
	};

	const recipes = [
		{
			id: 1,
			title: '아보카도 토스트',
			tags: ['브런치', '10분'],
			image: 'https://images.unsplash.com/photo-1541519227354-08fa5d50c820?w=200&auto=format&fit=crop&q=80'
		},
		{
			id: 2,
			title: '까르보나라',
			tags: ['파스타', '20분'],
			image: 'https://images.unsplash.com/photo-1621996346565-e3dbc646d9a9?w=200&auto=format&fit=crop&q=80'
		},
		{
			id: 3,
			title: '그린 샐러드',
			tags: ['샐러드', '5분'],
			image: 'https://images.unsplash.com/photo-1512621776951-a57141f2eefd?w=200&auto=format&fit=crop&q=80'
		},
		{
			id: 4,
			title: '블루베리 팬케이크',
			tags: ['브런치', '15분'],
			image: 'https://images.unsplash.com/photo-1565299543923-37dd37887442?w=200&auto=format&fit=crop&q=80'
		},
		{
			id: 5,
			title: '된장찌개',
			tags: ['한식', '25분'],
			image: 'https://images.unsplash.com/photo-1498654896293-37aaa4f85252?w=200&auto=format&fit=crop&q=80'
		},
		{
			id: 6,
			title: '토마토 리조또',
			tags: ['이탈리안', '30분'],
			image: 'https://images.unsplash.com/photo-1476124369491-e7addf5db371?w=200&auto=format&fit=crop&q=80'
		}
	];

	let linkInput = $state('');
	let activeTab = $state('home');

	function handleLinkSubmit() {
		if (!linkInput.trim()) return;
		alert(`링크 저장 예정: ${linkInput}`);
		linkInput = '';
	}
</script>

<!-- Outer: warm ivory bg, creates depth on desktop -->
<div class="min-h-screen bg-[#EDE9E2]">
	<!-- Centered app container -->
	<div class="mx-auto max-w-xl min-h-screen bg-white shadow-[0_0_60px_rgba(0,0,0,0.08)]">
		<div class="px-5 pb-28">

			<!-- Header -->
			<div class="flex items-center justify-between pt-12 pb-6">
				<div class="flex items-baseline gap-2">
					<span class="text-lg font-bold tracking-tight text-stone-900">recipejar</span>
					<span class="text-[11px] font-medium text-stone-400">나의 레시피 보관함</span>
				</div>
				<div
					class="flex h-8 w-8 items-center justify-center rounded-full bg-[#EBF0EB] text-xs font-bold text-[#6A8C6C]"
				>
					나
				</div>
			</div>

			<!-- URL Input Bar -->
			<div class="mb-8">
				<div
					class="flex items-center gap-3 rounded-2xl border border-stone-200 bg-stone-50/80 px-4 py-3 transition-colors focus-within:border-[#7C9A7E] focus-within:bg-white"
				>
					<svg
						class="h-4 w-4 shrink-0 text-stone-400"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						viewBox="0 0 24 24"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
						/>
					</svg>
					<input
						bind:value={linkInput}
						onkeydown={(e) => e.key === 'Enter' && handleLinkSubmit()}
						type="url"
						placeholder="유튜브·인스타 링크 붙여넣기..."
						class="flex-1 bg-transparent text-sm text-stone-700 outline-none placeholder:text-stone-400"
					/>
					{#if linkInput}
						<button
							onclick={handleLinkSubmit}
							class="shrink-0 rounded-xl bg-[#7C9A7E] px-3 py-1 text-xs font-semibold text-white"
						>
							저장
						</button>
					{/if}
				</div>
			</div>

			<!-- Hero: Magazine cover style -->
			<div class="mb-8">
				<p class="mb-3 text-[11px] font-semibold uppercase tracking-widest text-stone-400">
					오늘의 추천
				</p>
				<div class="relative h-48 overflow-hidden rounded-3xl">
					<img
						src={todayRecipe.image}
						alt={todayRecipe.title}
						class="h-full w-full object-cover"
					/>
					<!-- Bottom gradient -->
					<div
						class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent"
					></div>
					<!-- Top gradient for badge readability -->
					<div
						class="absolute inset-x-0 top-0 h-20 bg-gradient-to-b from-black/30 to-transparent"
					></div>
					<!-- Tag badge (top-left) -->
					<div class="absolute top-3.5 left-4">
						<span
							class="rounded-full bg-[#7C9A7E]/90 px-2.5 py-0.5 text-[11px] font-semibold text-white backdrop-blur-sm"
						>
							{todayRecipe.tag}
						</span>
					</div>
					<!-- Time badge (top-right) -->
					<div class="absolute top-3.5 right-4">
						<span
							class="rounded-full bg-white/20 px-2.5 py-0.5 text-[11px] font-medium text-white backdrop-blur-sm"
						>
							⏱ {todayRecipe.time}
						</span>
					</div>
					<!-- Bottom text overlay -->
					<div class="absolute bottom-0 left-0 right-0 p-4">
						<p class="mb-0.5 text-[11px] font-medium text-white/60">{todayRecipe.subtitle}</p>
						<h2 class="text-xl font-bold leading-snug text-white">{todayRecipe.title}</h2>
					</div>
				</div>
			</div>

			<!-- Stats Dashboard -->
			<div class="mb-8 grid grid-cols-2 gap-3">
				<div class="rounded-2xl bg-[#F5F2EC] p-4">
					<p class="mb-1.5 text-[11px] font-medium text-stone-500">총 레시피</p>
					<p class="text-3xl font-bold text-stone-800">12</p>
				</div>
				<div class="rounded-2xl bg-[#EEF3EE] p-4">
					<p class="mb-1.5 text-[11px] font-medium text-stone-500">이번 달 요리</p>
					<p class="text-3xl font-bold text-stone-800">
						5<span class="ml-1 text-sm font-normal text-stone-400">회</span>
					</p>
				</div>
			</div>

			<!-- Recipe List -->
			<div>
				<div class="mb-3 flex items-center justify-between">
					<h2 class="text-sm font-semibold text-stone-900">내 레시피 보관함</h2>
					<button class="text-xs font-medium text-[#7C9A7E]">전체 보기 →</button>
				</div>
				<div class="flex flex-col divide-y divide-stone-100">
					{#each recipes as recipe}
						<div class="group flex items-center gap-3 py-3">
							<!-- YouTube-style 16:9 thumbnail -->
							<div class="h-12 w-[86px] shrink-0 overflow-hidden rounded-lg bg-stone-100">
								<img
									src={recipe.image}
									alt={recipe.title}
									class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
								/>
							</div>
							<!-- Info -->
							<div class="min-w-0 flex-1">
								<p class="truncate text-sm font-medium text-stone-800">{recipe.title}</p>
								<div class="mt-1 flex flex-wrap gap-1">
									{#each recipe.tags as tag}
										<span
											class="rounded-full bg-stone-100 px-2 py-0.5 text-[10px] leading-none text-stone-500"
										>
											{tag}
										</span>
									{/each}
								</div>
							</div>
							<!-- Cook button -->
							<button
								class="shrink-0 whitespace-nowrap text-xs font-semibold text-[#7C9A7E] transition-colors hover:text-[#5A7D5C]"
							>
								만들었어! 🔥
							</button>
						</div>
					{/each}
				</div>
			</div>

		</div>
	</div>
</div>

<!-- Bottom Nav: fixed, centered, matching container width -->
<nav
	class="fixed bottom-0 left-1/2 z-50 w-full max-w-xl -translate-x-1/2 border-t border-stone-100 bg-white/90 backdrop-blur-md"
>
	<div class="flex items-center justify-around px-8 py-3">
		<button
			onclick={() => (activeTab = 'home')}
			class="flex flex-col items-center gap-1 transition-colors {activeTab === 'home'
				? 'text-[#7C9A7E]'
				: 'text-stone-400'}"
		>
			<svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
				<path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z" />
			</svg>
			<span class="text-[10px] font-medium">홈</span>
		</button>

		<button
			onclick={() => (activeTab = 'random')}
			class="flex flex-col items-center gap-1 transition-colors {activeTab === 'random'
				? 'text-[#7C9A7E]'
				: 'text-stone-400'}"
		>
			<svg
				class="h-5 w-5"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				viewBox="0 0 24 24"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
				/>
			</svg>
			<span class="text-[10px] font-medium">랜덤</span>
		</button>

		<button
			onclick={() => (activeTab = 'settings')}
			class="flex flex-col items-center gap-1 transition-colors {activeTab === 'settings'
				? 'text-[#7C9A7E]'
				: 'text-stone-400'}"
		>
			<svg
				class="h-5 w-5"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				viewBox="0 0 24 24"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
				/>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
				/>
			</svg>
			<span class="text-[10px] font-medium">설정</span>
		</button>
	</div>
</nav>
