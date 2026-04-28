<script lang="ts">
	import { goto } from '$app/navigation';
	import { paymentMethodsApi } from '$lib/api';
	import PaymentMethodForm from '$lib/components/PaymentMethodForm.svelte';
	import type { PaymentMethodInput } from '$lib/types';
	import { resolve } from '$app/paths';

	let submitting = $state(false);
	let error = $state<string | null>(null);

	async function handleSubmit(input: PaymentMethodInput) {
		submitting = true;
		error = null;
		try {
			await paymentMethodsApi.create(input);
			await goto(resolve('/payment-methods'));
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
			submitting = false;
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1>New payment method</h1>
	</div>
	{#if error}
		<div class="error">{error}</div>
	{/if}
	<PaymentMethodForm onSubmit={handleSubmit} {submitting} cancelHref={resolve('/payment-methods')} />
</div>
