<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { categoriesApi, expensesApi, paymentMethodsApi, tagsApi } from '$lib/api';
	import ExpenseForm from '$lib/components/ExpenseForm.svelte';
	import type { Category, ExpenseInput, PaymentMethod, Tag } from '$lib/types';
	import { resolve } from '$app/paths';

	let categories: Category[] = $state([]);
	let paymentMethods: PaymentMethod[] = $state([]);
	let tags: Tag[] = $state([]);
	let loading = $state(true);
	let submitting = $state(false);
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			const [c, p, t] = await Promise.all([
				categoriesApi.list(),
				paymentMethodsApi.list(),
				tagsApi.list()
			]);
			categories = c ?? [];
			paymentMethods = p ?? [];
			tags = t ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			loading = false;
		}
	});

	async function handleSubmit(input: ExpenseInput) {
		submitting = true;
		error = null;
		try {
			await expensesApi.create(input);
			await goto(resolve('/expenses'));
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
			submitting = false;
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1>New expense</h1>
	</div>
	{#if error}
		<div class="error">{error}</div>
	{/if}
	{#if loading}
		<p class="muted">Loading…</p>
	{:else if categories.length === 0 || paymentMethods.length === 0}
		<div class="card">
			<p>
				You need at least one <a href={resolve('/categories/new')}>category</a> and one
				<a href={resolve('/payment-methods/new')}>payment method</a> before creating an expense.
			</p>
		</div>
	{:else}
		<ExpenseForm
			{categories}
			{paymentMethods}
			{tags}
			onSubmit={handleSubmit}
			{submitting}
			cancelHref={resolve('/expenses')}
		/>
	{/if}
</div>
