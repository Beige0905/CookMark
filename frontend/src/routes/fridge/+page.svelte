<div class="space-y-8">
	<div>
		<h1 class="text-3xl font-bold text-stone-900">내 냉장고 🧊</h1>
		<p class="text-stone-500">지금 냉장고에 있는 재료들로 무엇을 만들 수 있을까요?</p>
	</div>

	<!-- Fridge Status -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<div class="col-span-1 md:col-span-2 space-y-6">
			<div class="rounded-[2rem] bg-stone-50 p-8 border border-stone-100">
				<div class="flex items-center justify-between mb-6">
					<h2 class="text-xl font-bold text-stone-800">보관 중인 재료</h2>
					<button class="text-sm font-bold text-[#7C9A7E] hover:underline">+ 재료 추가</button>
				</div>
				<div class="flex flex-wrap gap-3">
					{#each ['양파', '두부', '계란', '대파', '돼지고기'] as item}
						<span class="inline-flex items-center gap-2 rounded-2xl bg-white px-4 py-2 text-sm font-semibold text-stone-700 shadow-sm border border-stone-100 group">
							{item}
							<button class="text-stone-300 hover:text-red-400">×</button>
						</span>
					{/each}
				</div>
			</div>

			<button class="w-full rounded-[2rem] bg-[#7C9A7E] py-6 text-xl font-bold text-white shadow-xl shadow-[#7C9A7E]/20 transition-all hover:scale-[1.02] active:scale-95">
				🍳 이 재료들로 뭐 만들지?
			</button>
		</div>

		<div class="space-y-6">
			<div class="rounded-[2rem] bg-[#EEF3EE] p-6">
				<h3 class="font-bold text-[#6A8C6C] mb-2 text-sm uppercase tracking-wider">추천 필터</h3>
				<div class="space-y-3">
					<label class="flex items-center gap-3 cursor-pointer group">
						<input type="checkbox" checked class="h-5 w-5 rounded-lg border-stone-300 text-[#7C9A7E] focus:ring-[#7C9A7E]" />
						<span class="text-stone-700 font-medium">유통기한 임박 우선</span>
					</label>
					<label class="flex items-center gap-3 cursor-pointer group">
						<input type="checkbox" class="h-5 w-5 rounded-lg border-stone-300 text-[#7C9A7E] focus:ring-[#7C9A7E]" />
						<span class="text-stone-700 font-medium">메인 재료 포함</span>
					</label>
				</div>
			</div>
		</div>
	</div>
</div>
