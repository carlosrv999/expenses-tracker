<script lang="ts">
	import { onMount } from 'svelte';
	import {
		categoriesApi,
		expensesApi,
		paymentMethodsApi,
		tagsApi
	} from '$lib/api';
	import { formatAmount, formatDate } from '$lib/format';
	import type { Category, Expense, ExpenseListFilters, PaymentMethod, Tag } from '$lib/types';
	import { Receipt, Pencil, Trash2, Plus, Filter, X, Upload } from 'lucide-svelte';
	import { resolve } from '$app/paths';

	let expenses: Expense[] = $state([]);
	let categories: Category[] = $state([]);
	let paymentMethods: PaymentMethod[] = $state([]);
	let tags: Tag[] = $state([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	let filterCategory = $state<number | ''>('');
	let filterPm = $state<number | ''>('');
	let filterLimit = $state<number>(50);

	let categoryById = $derived(new Map(categories.map((c) => [c.category_id, c])));
	let pmById = $derived(new Map(paymentMethods.map((p) => [p.payment_method_id, p])));
	let tagById = $derived(new Map(tags.map((t) => [t.tag_id, t])));

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
			await load();
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
			loading = false;
		}
	});

	async function load() {
		loading = true;
		error = null;
		try {
			const filters: ExpenseListFilters = { limit: filterLimit };
			if (filterCategory !== '') filters.category_id = Number(filterCategory);
			if (filterPm !== '') filters.payment_method_id = Number(filterPm);
			expenses = (await expensesApi.list(filters)) ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			loading = false;
		}
	}

	function clearFilters() {
		filterCategory = '';
		filterPm = '';
		filterLimit = 50;
		load();
	}

	async function remove(e: Expense) {
		if (!confirm(`Delete expense #${e.expense_id}?`)) return;
		try {
			await expensesApi.remove(e.expense_id);
			await load();
		} catch (err) {
			alert(err instanceof Error ? err.message : String(err));
		}
	}

	function expenseTagIds(e: Expense): number[] {
		return e.tag_ids ?? e.tags?.map((t) => t.tag_id) ?? [];
	}
</script>

<div class="page">
	<div class="page-header">
		<h1><Receipt size={22} /> Expenses</h1>
		<a href={resolve('/expenses/upload')}>
			<button class="secondary"><Upload size={16} /> Bulk Upload CSV</button>
		</a>
		<a href={resolve('/expenses/new')}><button class="primary"><Plus size={16} /> New expense</button></a>
	</div>

	<div class="card filters">
		<div class="filter-grid">
			<div>
				<label for="f-category"><Filter size={13} /> Category</label>
				<select id="f-category" bind:value={filterCategory}>
					<option value="">All</option>
					{#each categories as c (c.category_id)}
						<option value={c.category_id}>{c.category_name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label for="f-pm">Payment method</label>
				<select id="f-pm" bind:value={filterPm}>
					<option value="">All</option>
					{#each paymentMethods as p (p.payment_method_id)}
						<option value={p.payment_method_id}>{p.method_name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label for="f-limit">Limit</label>
				<input id="f-limit" type="number" min="1" max="500" bind:value={filterLimit} />
			</div>
			<div class="filter-actions">
				<button onclick={load} class="primary">Apply</button>
				<button onclick={clearFilters} class="ghost"><X size={14} /> Clear</button>
			</div>
		</div>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<p class="muted">Loading…</p>
	{:else if expenses.length === 0}
		<div class="empty card">
			No expenses match. <a href={resolve('/expenses/new')}>Create one</a>.
		</div>
	{:else}
		<table>
			<thead>
				<tr>
					<th>Date</th>
					<th>Merchant</th>
					<th>Category</th>
					<th>Payment</th>
					<th>Tags</th>
					<th style="text-align:right">Amount</th>
					<th style="width: 1%"></th>
				</tr>
			</thead>
			<tbody>
				{#each expenses as e (e.expense_id)}
					<tr>
						<td>{formatDate(e.expense_date)}</td>
						<td>
							<a href={resolve(`/expenses/${e.expense_id}`)}>{e.merchant_name ?? `#${e.expense_id}`}</a>
							{#if e.description}
								<div class="muted small">{e.description}</div>
							{/if}
						</td>
						<td>{categoryById.get(e.category_id)?.category_name ?? '—'}</td>
						<td>{pmById.get(e.payment_method_id)?.method_name ?? '—'}</td>
						<td>
							{#each expenseTagIds(e) as tid (tid)}
								{@const tag = tagById.get(tid)}
								{#if tag}
									<span class="badge">
										{#if tag.color}<span class="color-dot" style:background={tag.color}></span>{/if}
										{tag.tag_name}
									</span>
								{/if}
							{/each}
						</td>
						<td style="text-align:right">{formatAmount(e.amount, e.currency)}</td>
						<td>
							<div class="actions">
								<a href={resolve(`/expenses/${e.expense_id}`)}><button title="Edit"><Pencil size={15} /></button></a>
								<button class="danger" title="Delete" onclick={() => remove(e)}><Trash2 size={15} /></button>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>

<style>
	.filters {
		margin-bottom: 1rem;
	}
	.filter-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
		gap: 0.85rem;
		align-items: end;
	}
	.filter-actions {
		display: flex;
		gap: 0.5rem;
	}
	.small {
		font-size: 0.8rem;
		margin-top: 0.15rem;
	}
</style>
