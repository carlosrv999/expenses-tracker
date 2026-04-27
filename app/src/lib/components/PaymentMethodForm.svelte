<script lang="ts">
	import IconPicker from './IconPicker.svelte';
	import type { PaymentMethod, PaymentMethodInput, PaymentMethodType } from '$lib/types';
	import { Save, X } from 'lucide-svelte';

	const TYPES: PaymentMethodType[] = [
		'cash',
		'credit_card',
		'debit_card',
		'yape',
		'plin',
		'bank_transfer'
	];

	let {
		initial,
		onSubmit,
		submitting = false,
		cancelHref
	}: {
		initial?: PaymentMethod | null;
		onSubmit: (input: PaymentMethodInput) => void | Promise<void>;
		submitting?: boolean;
		cancelHref: string;
	} = $props();

	/* svelte-ignore state_referenced_locally */
	let method_name = $state(initial?.method_name ?? '');
	/* svelte-ignore state_referenced_locally */
	let method_type = $state<PaymentMethodType>(initial?.method_type ?? 'cash');
	/* svelte-ignore state_referenced_locally */
	let icon = $state(initial?.icon ?? '');

	async function handleSubmit(event: Event) {
		event.preventDefault();
		const input: PaymentMethodInput = {
			method_name: method_name.trim(),
			method_type,
			icon: icon || null
		};
		await onSubmit(input);
	}
</script>

<form onsubmit={handleSubmit} class="card">
	<div class="grid-form">
		<div>
			<label for="pm-name">Name</label>
			<input id="pm-name" bind:value={method_name} required maxlength="100" />
		</div>
		<div>
			<label for="pm-type">Type</label>
			<select id="pm-type" bind:value={method_type}>
				{#each TYPES as t (t)}
					<option value={t}>{t}</option>
				{/each}
			</select>
		</div>
		<div>
			<label for="pm-icon">Icon</label>
			<IconPicker id="pm-icon" bind:value={icon} />
		</div>
	</div>
	<div class="form-actions">
		<a href={cancelHref}><button type="button" class="ghost"><X size={15} /> Cancel</button></a>
		<button type="submit" class="primary" disabled={submitting}>
			<Save size={15} /> {submitting ? 'Saving…' : 'Save'}
		</button>
	</div>
</form>
