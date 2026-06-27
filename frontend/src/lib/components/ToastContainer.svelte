<script lang="ts">
	import { toast } from '$lib/toast.svelte';
	import { fly } from 'svelte/transition';
	import CheckCircle from 'lucide-svelte/icons/check-circle';
	import AlertCircle from 'lucide-svelte/icons/alert-circle';
	import Info from 'lucide-svelte/icons/info';
	import X from 'lucide-svelte/icons/x';

	const icons = {
		success: CheckCircle,
		error: AlertCircle,
		info: Info
	};

	const colors = {
		success: 'bg-green-50 text-green-800 border-green-200',
		error: 'bg-red-50 text-red-800 border-red-200',
		info: 'bg-blue-50 text-blue-800 border-blue-200'
	};
</script>

<div class="fixed top-6 right-6 z-[100] flex flex-col gap-3 pointer-events-none">
	{#each toast.toasts as t (t.id)}
		<div
			in:fly={{ y: -20, duration: 300 }}
			out:fly={{ x: 20, duration: 300 }}
			class="flex items-center gap-3 px-4 py-3 rounded-2xl border shadow-lg max-w-sm pointer-events-auto {colors[t.type]}"
		>
			<svelte:component this={icons[t.type]} size={20} />
			<p class="text-sm font-bold flex-1">{t.message}</p>
			<button
				onclick={() => toast.remove(t.id)}
				class="p-1 hover:bg-black/5 rounded-full transition-colors"
			>
				<X size={16} />
			</button>
		</div>
	{/each}
</div>
