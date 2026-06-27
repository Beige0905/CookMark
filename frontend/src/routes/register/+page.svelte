<script lang="ts">
	import { goto } from '$app/navigation';

	let email = $state('');
	let password = $state('');
	let displayName = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleRegister(e: Event) {
		e.preventDefault();
		error = '';
		if (password.length < 8) {
			error = '비밀번호는 8자 이상이어야 합니다';
			return;
		}
		loading = true;
		try {
			const res = await fetch('/api/auth/register', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password, display_name: displayName })
			});
			if (res.ok) {
				goto('/login');
			} else {
				const text = await res.text();
				error = text.trim() || '회원가입에 실패했습니다';
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
			<p class="mt-1 text-sm text-stone-400">계정 만들기</p>
		</div>

		<form onsubmit={handleRegister} class="space-y-4">
			<div>
				<label for="name" class="mb-1.5 block text-sm font-semibold text-stone-700">이름</label>
				<input
					id="name"
					type="text"
					bind:value={displayName}
					required
					autocomplete="name"
					placeholder="홍길동"
					class="w-full rounded-2xl border border-stone-200 px-4 py-3 text-sm outline-none transition focus:border-[#7C9A7E] focus:ring-2 focus:ring-[#7C9A7E]/20"
				/>
			</div>

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
				<label for="password" class="mb-1.5 block text-sm font-semibold text-stone-700"
					>비밀번호 <span class="font-normal text-stone-400">(8자 이상)</span></label
				>
				<input
					id="password"
					type="password"
					bind:value={password}
					required
					autocomplete="new-password"
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
				{loading ? '가입 중...' : '회원가입'}
			</button>
		</form>

		<p class="mt-6 text-center text-sm text-stone-400">
			이미 계정이 있으신가요?
			<a href="/login" class="font-semibold text-[#7C9A7E] hover:underline">로그인</a>
		</p>
	</div>
</div>
