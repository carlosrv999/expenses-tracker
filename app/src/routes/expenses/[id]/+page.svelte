<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { categoriesApi, expensesApi, paymentMethodsApi, tagsApi } from '$lib/api';
	import ExpenseForm from '$lib/components/ExpenseForm.svelte';
	import type { Category, Expense, ExpenseInput, PaymentMethod, Tag } from '$lib/types';
	import { Trash2 } from 'lucide-svelte';
	import { resolve } from '$app/paths';

	let id = $derived(Number(page.params.id));
	let expense = $state<Expense | null>(null);
	let categories: Category[] = $state([]);
	let paymentMethods: PaymentMethod[] = $state([]);
	let tags: Tag[] = $state([]);
	let loading = $state(true);
	let submitting = $state(false);
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			const [exp, c, p, t] = await Promise.all([
				expensesApi.get(id),
				categoriesApi.list(),
				paymentMethodsApi.list(),
				tagsApi.list()
			]);
			expense = exp;
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
			await expensesApi.update(id, input);
			await goto(resolve('/expenses'));
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
			submitting = false;
		}
	}

	async function remove() {
		if (!expense) return;
		if (!confirm(`Delete expense #${expense.expense_id}?`)) return;
		try {
			await expensesApi.remove(id);
			await goto(resolve('/expenses'));
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1>Edit expense</h1>
		{#if expense}
			<button class="danger" onclick={remove}><Trash2 size={15} /> Delete</button>
		{/if}
	</div>
	{#if error}
		<div class="error">{error}</div>
	{/if}
	{#if loading}
		<p class="muted">Loading…</p>
	{:else if expense}
		<ExpenseForm
			initial={expense}
			{categories}
			{paymentMethods}
			{tags}
			onSubmit={handleSubmit}
			{submitting}
			cancelHref={resolve('/expenses')}
		/>
	{:else}
		<p class="muted">Expense not found.</p>
	{/if}
</div>
