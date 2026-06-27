<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';

	const displayName = $derived(auth.user?.display_name ?? '...');
	const email = $derived(auth.user?.email ?? '');
	const initial = $derived(displayName[0]?.toUpperCase() ?? '?');
</script>

<div class="max-w-2xl space-y-10">
	<div>
		<h1 class="text-3xl font-bold text-stone-900">설정 ⚙️</h1>
		<p class="text-stone-500">앱 환경과 계정 정보를 관리하세요.</p>
	</div>

	<div class="space-y-4">
		<h2 class="px-4 text-xs font-bold uppercase tracking-widest text-stone-400">계정</h2>
		<div class="overflow-hidden rounded-[2rem] border border-stone-100 bg-white">
			<div class="flex items-center justify-between border-b border-stone-50 p-6">
				<div class="flex items-center gap-4">
					<div
						class="flex h-12 w-12 items-center justify-center rounded-full bg-[#EBF0EB] font-bold text-[#6A8C6C]"
					>
						{initial}
					</div>
					<div>
						<p class="font-bold text-stone-800">{displayName} 님</p>
						<p class="text-sm text-stone-400">{email}</p>
					</div>
				</div>
			</div>
			<button
				onclick={() => auth.logout()}
				class="w-full p-6 text-left font-bold text-red-500 transition-colors hover:bg-red-50"
			>
				로그아웃
			</button>
		</div>
	</div>

	<div class="space-y-4">
		<h2 class="px-4 text-xs font-bold uppercase tracking-widest text-stone-400">앱 설정</h2>
		<div class="overflow-hidden rounded-[2rem] border border-stone-100 bg-white">
			<div class="flex items-center justify-between border-b border-stone-50 p-6">
				<span class="font-bold text-stone-700">다크 모드</span>
				<div class="relative h-6 w-11 cursor-pointer rounded-full bg-stone-200 p-1">
					<div class="h-4 w-4 rounded-full bg-white shadow-sm"></div>
				</div>
			</div>
			<div class="flex items-center justify-between border-b border-stone-50 p-6">
				<span class="font-bold text-stone-700">푸시 알림</span>
				<div class="relative h-6 w-11 cursor-pointer rounded-full bg-[#7C9A7E] p-1">
					<div class="ml-auto h-4 w-4 rounded-full bg-white shadow-sm"></div>
				</div>
			</div>
			<div class="p-6">
				<p class="text-xs text-stone-400">버전 1.0.0 (beta)</p>
			</div>
		</div>
	</div>
</div>
