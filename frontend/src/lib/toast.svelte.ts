export type ToastType = 'success' | 'error' | 'info';

export interface Toast {
	id: number;
	message: string;
	type: ToastType;
}

class ToastState {
	toasts = $state<Toast[]>([]);

	add(message: string, type: ToastType = 'success', duration = 3000) {
		const id = Math.random();
		this.toasts.push({ id, message, type });
		
		if (duration > 0) {
			setTimeout(() => {
				this.remove(id);
			}, duration);
		}
	}

	remove(id: number) {
		this.toasts = this.toasts.filter((t) => t.id !== id);
	}
}

export const toast = new ToastState();
