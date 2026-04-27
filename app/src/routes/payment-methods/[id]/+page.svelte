<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { paymentMethodsApi } from '$lib/api';
	import PaymentMethodForm from '$lib/components/PaymentMethodForm.svelte';
	import type { PaymentMethod, PaymentMethodInput } from '$lib/types';
	import { Trash2 } from 'lucide-svelte';

	let id = $derived(Number(page.params.id));
	let method = $state<PaymentMethod | null>(null);
	let loading = $state(true);
	let submitting = $state(false);
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			method = await paymentMethodsApi.get(id);
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			loading = false;
		}
	});

	async function handleSubmit(input: PaymentMethodInput) {
		submitting = true;
		error = null;
		try {
			await paymentMethodsApi.update(id, input);
			await goto('/payment-methods');
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
			submitting = false;
		}
	}

	async function remove() {
		if (!method) return;
		if (!confirm(`Delete payment method "${method.method_name}"?`)) return;
		try {
			await paymentMethodsApi.remove(id);
			await goto('/payment-methods');
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1>Edit payment method</h1>
		{#if method}
			<button class="danger" onclick={remove}><Trash2 size={15} /> Delete</button>
		{/if}
	</div>
	{#if error}
		<div class="error">{error}</div>
	{/if}
	{#if loading}
		<p class="muted">Loading…</p>
	{:else if method}
		<PaymentMethodForm
			initial={method}
			onSubmit={handleSubmit}
			{submitting}
			cancelHref="/payment-methods"
		/>
	{:else}
		<p class="muted">Payment method not found.</p>
	{/if}
</div>
