<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import Home from 'lucide-svelte/icons/house';
	import BookOpen from 'lucide-svelte/icons/book-open';
	import Refrigerator from 'lucide-svelte/icons/refrigerator';
	import Settings from 'lucide-svelte/icons/settings';
	import ToastContainer from '$lib/components/ToastContainer.svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import { theme } from '$lib/stores/theme.svelte';

	const isActive = (path: string) => {
		if (path === '/') return $page.url.pathname === '/';
		return $page.url.pathname.startsWith(path);
	};

	const isAuthPage = $derived(
		$page.url.pathname.startsWith('/login') || $page.url.pathname.startsWith('/register')
	);

	onMount(() => {
		theme.init();
		if (!isAuthPage) {
			auth.fetchMe();
		}
	});

	const userInitial = $derived(auth.user?.display_name?.[0]?.toUpperCase() ?? '?');
</script>

{#if isAuthPage}
	<slot />
{:else}
<!-- Outer: warm ivory bg, creates depth on desktop -->
<div class="min-h-screen bg-[#EDE9E2] dark:bg-[#111111]">
	<!-- Centered app container: Mobile-first, becomes wider on desktop -->
	<div class="mx-auto min-h-screen max-w-xl bg-white shadow-[0_0_60px_rgba(0,0,0,0.08)] lg:max-w-5xl dark:bg-[#1c1c1c] dark:shadow-none">
		<div class="flex flex-col lg:flex-row">

			<!-- Left Sidebar: Hidden on mobile, sticky on desktop -->
			<aside class="hidden w-80 shrink-0 border-r border-stone-100 p-8 lg:sticky lg:top-0 lg:block lg:h-screen dark:border-stone-800">
				<div class="flex flex-col h-full">
					<!-- Logo/Header -->
					<a href="/" class="flex items-baseline gap-2 mb-12 hover:opacity-80 transition-opacity">
						<span class="text-2xl font-bold tracking-tight text-stone-900 dark:text-stone-100">Cookmark</span>
						<span class="text-xs font-medium text-stone-400 dark:text-stone-500">나의 레시피 보관함</span>
					</a>

					<!-- Sidebar Nav -->
					<nav class="flex flex-col gap-2 mb-12">
						<a
							href="/"
							class="flex items-center gap-3 px-4 py-3 rounded-2xl transition-colors {isActive('/') ? 'bg-[#7C9A7E] text-white' : 'text-stone-500 hover:bg-stone-50 dark:text-stone-400 dark:hover:bg-stone-800'}"
						>
							<Home size={20} strokeWidth={isActive('/') ? 2.5 : 2} />
							<span class="font-semibold">홈</span>
						</a>
						<a
							href="/recipes"
							class="flex items-center gap-3 px-4 py-3 rounded-2xl transition-colors {isActive('/recipes') ? 'bg-[#7C9A7E] text-white' : 'text-stone-500 hover:bg-stone-50 dark:text-stone-400 dark:hover:bg-stone-800'}"
						>
							<BookOpen size={20} strokeWidth={isActive('/recipes') ? 2.5 : 2} />
							<span class="font-semibold">레시피</span>
						</a>
						<a
							href="/fridge"
							class="flex items-center gap-3 px-4 py-3 rounded-2xl transition-colors {isActive('/fridge') ? 'bg-[#7C9A7E] text-white' : 'text-stone-500 hover:bg-stone-50 dark:text-stone-400 dark:hover:bg-stone-800'}"
						>
							<Refrigerator size={20} strokeWidth={isActive('/fridge') ? 2.5 : 2} />
							<span class="font-semibold">냉장고</span>
						</a>
						<a
							href="/settings"
							class="flex items-center gap-3 px-4 py-3 rounded-2xl transition-colors {isActive('/settings') ? 'bg-[#7C9A7E] text-white' : 'text-stone-500 hover:bg-stone-50 dark:text-stone-400 dark:hover:bg-stone-800'}"
						>
							<Settings size={20} strokeWidth={isActive('/settings') ? 2.5 : 2} />
							<span class="font-semibold">설정</span>
						</a>
					</nav>

					<!-- Stats in Sidebar for Desktop -->
					<div class="mt-auto space-y-4">
						<div class="rounded-2xl bg-[#F5F2EC] dark:bg-stone-800 p-5">
							<p class="mb-1 text-xs font-medium text-stone-500 dark:text-stone-400 uppercase tracking-wider">총 레시피</p>
							<p class="text-4xl font-bold text-stone-800 dark:text-stone-100">12</p>
						</div>
						<div class="rounded-2xl bg-[#EEF3EE] dark:bg-stone-800 p-5">
							<p class="mb-1 text-xs font-medium text-stone-500 dark:text-stone-400 uppercase tracking-wider">이번 달 요리</p>
							<p class="text-4xl font-bold text-stone-800 dark:text-stone-100">5<span class="ml-1 text-sm font-normal text-stone-400">회</span></p>
						</div>
					</div>
				</div>
			</aside>

			<!-- Main Content Area -->
			<main class="flex-1 px-5 pb-28 lg:px-10 lg:py-12">
				<!-- Header (Mobile Only) -->
				<div class="flex items-center justify-between pt-12 pb-6 lg:hidden">
					<a href="/" class="flex items-baseline gap-2">
						<span class="text-lg font-bold tracking-tight text-stone-900 dark:text-stone-100">Cookmark</span>
						<span class="text-[11px] font-medium text-stone-400 dark:text-stone-500">나의 레시피 보관함</span>
					</a>
					<div class="flex h-8 w-8 items-center justify-center rounded-full bg-[#EBF0EB] dark:bg-stone-700 text-xs font-bold text-[#6A8C6C] dark:text-stone-300">{userInitial}</div>
				</div>

				<slot />
			</main>
		</div>
	</div>
</div>

<!-- Bottom Nav (Mobile Only) -->
<nav class="fixed bottom-0 left-0 right-0 z-50 border-t border-stone-100 bg-white/90 backdrop-blur-md lg:hidden dark:border-stone-800 dark:bg-[#1c1c1c]/90">
	<div class="flex items-center justify-around px-8 py-4">
		<a href="/" class="flex flex-col items-center gap-1.5 transition-colors {isActive('/') ? 'text-[#7C9A7E]' : 'text-stone-400 dark:text-stone-500'}">
			<Home size={24} strokeWidth={isActive('/') ? 2.5 : 2} />
			<span class="text-[10px] font-bold">홈</span>
		</a>
		<a href="/recipes" class="flex flex-col items-center gap-1.5 transition-colors {isActive('/recipes') ? 'text-[#7C9A7E]' : 'text-stone-400 dark:text-stone-500'}">
			<BookOpen size={24} strokeWidth={isActive('/recipes') ? 2.5 : 2} />
			<span class="text-[10px] font-bold">레시피</span>
		</a>
		<a href="/fridge" class="flex flex-col items-center gap-1.5 transition-colors {isActive('/fridge') ? 'text-[#7C9A7E]' : 'text-stone-400 dark:text-stone-500'}">
			<Refrigerator size={24} strokeWidth={isActive('/fridge') ? 2.5 : 2} />
			<span class="text-[10px] font-bold">냉장고</span>
		</a>
		<a href="/settings" class="flex flex-col items-center gap-1.5 transition-colors {isActive('/settings') ? 'text-[#7C9A7E]' : 'text-stone-400 dark:text-stone-500'}">
			<Settings size={24} strokeWidth={isActive('/settings') ? 2.5 : 2} />
			<span class="text-[10px] font-bold">설정</span>
		</a>
	</div>
</nav>
{/if}

<ToastContainer />
