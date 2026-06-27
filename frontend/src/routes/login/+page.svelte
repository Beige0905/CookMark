<script lang="ts">
	import { goto } from '$app/navigation';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleLogin(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;
		try {
			const res = await fetch('/api/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password })
			});
			if (res.ok) {
				goto('/');
			} else {
				const text = await res.text();
				error = text.trim() || '로그인에 실패했습니다';
			}
		} catch {
			error = '서버와 연결할 수 없습니다';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-[#EDE9E2]">
	<div class="w-full max-w-sm rounded-3xl bg-white p-8 shadow-[0_0_40px_rgba(0,0,0,0.08)]">
		<div class="mb-8 text-center">
			<h1 class="text-2xl font-bold text-stone-900">Cookmark</h1>
			<p class="mt-1 text-sm text-stone-400">나의 레시피 보관함</p>
		</div>

		<form onsubmit={handleLogin} class="space-y-4">
			<div>
				<label for="email" class="mb-1.5 block text-sm font-semibold text-stone-700">이메일</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					required
					autocomplete="email"
					placeholder="you@example.com"
					class="w-full rounded-2xl border border-stone-200 px-4 py-3 text-sm outline-none transition focus:border-[#7C9A7E] focus:ring-2 focus:ring-[#7C9A7E]/20"
				/>
			</div>

			<div>
				<label for="password" class="mb-1.5 block text-sm font-semibold text-stone-700">비밀번호</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					required
					autocomplete="current-password"
					placeholder="••••••••"
					class="w-full rounded-2xl border border-stone-200 px-4 py-3 text-sm outline-none transition focus:border-[#7C9A7E] focus:ring-2 focus:ring-[#7C9A7E]/20"
				/>
			</div>

			{#if error}
				<p class="rounded-xl bg-red-50 px-4 py-2.5 text-sm text-red-600">{error}</p>
			{/if}

			<button
				type="submit"
				disabled={loading}
				class="w-full rounded-2xl bg-[#7C9A7E] py-3 text-sm font-bold text-white transition hover:bg-[#6A8C6C] disabled:opacity-60"
			>
				{loading ? '로그인 중...' : '로그인'}
			</button>
		</form>

		<p class="mt-6 text-center text-sm text-stone-400">
			계정이 없으신가요?
			<a href="/register" class="font-semibold text-[#7C9A7E] hover:underline">회원가입</a>
		</p>
	</div>
</div>
