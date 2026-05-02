<script lang="ts">
	import { onMount } from 'svelte';
	import { categoriesApi, expensesApi, paymentMethodsApi } from '$lib/api';
	import { formatAmount, formatDate } from '$lib/format';
	import type { Category, Expense, ExpenseListFilters, PaymentMethod, Tag } from '$lib/types';
	import {
		Receipt,
		Pencil,
		Trash2,
		Plus,
		Funnel,
		Upload,
		ChevronLeft,
		ChevronRight
	} from 'lucide-svelte';
	import { SvelteMap } from 'svelte/reactivity';
	import { resolve } from '$app/paths';

	let expenses: Expense[] = $state([]);
	let totalCount: number = $state(0);

	let categories: Category[] = $state([]);
	let paymentMethods: PaymentMethod[] = $state([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	// Filter & Pagination States
	let filterCategory = $state<number | ''>('');
	let filterPm = $state<number | ''>('');
	let filterStartDate = $state<string>(''); // ← NEW
	let filterEndDate = $state<string>(''); // ← NEW
	let filterLimit = $state<number>(20);
	let offset = $state<number>(0);

	let categoryById = $derived(new SvelteMap(categories.map((c) => [c.category_id, c])));
	let pmById = $derived(new SvelteMap(paymentMethods.map((p) => [p.payment_method_id, p])));

	// Pagination Math
	let currentPage = $derived(Math.floor(offset / filterLimit) + 1);
	let totalPages = $derived(Math.ceil(totalCount / filterLimit) || 1);

	// Return full Tag objects from the expense (what the backend returns)
	function getExpenseTags(e: Expense): Tag[] {
		return Array.isArray(e.tags) ? e.tags : [];
	}

	async function fetchExpenses() {
		loading = true;
		try {
			const filters: ExpenseListFilters = {
				limit: filterLimit,
				offset: offset,
				relations: ['category', 'payment_method'] // ← NEW: ask backend to embed full objects
			};
			if (filterCategory !== '') filters.category_id = filterCategory;
			if (filterPm !== '') filters.payment_method_id = filterPm;
			if (filterStartDate) filters.start_date = filterStartDate;
			if (filterEndDate) filters.end_date = filterEndDate;

			const paginatedResult = await expensesApi.list(filters);

			if (paginatedResult) {
				expenses = paginatedResult.expenses || [];
				totalCount = paginatedResult.total_count || 0;
			} else {
				expenses = [];
				totalCount = 0;
			}
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			loading = false;
		}
	}

	function applyFilters() {
		offset = 0; // Reset to page 1 on filter/limit change
		fetchExpenses();
	}

	function prevPage() {
		if (offset > 0) {
			offset = Math.max(0, offset - filterLimit);
			fetchExpenses();
		}
	}

	function nextPage() {
		if (offset + filterLimit < totalCount) {
			offset += filterLimit;
			fetchExpenses();
		}
	}

	async function remove(e: Expense) {
		if (!confirm('Are you sure you want to delete this expense?')) return;
		try {
			await expensesApi.remove(e.expense_id);
			await fetchExpenses();
		} catch (err: unknown) {
			alert('Failed to delete: ' + (err instanceof Error ? err.message : String(err)));
		}
	}

	onMount(async () => {
		try {
			const [c, p] = await Promise.all([categoriesApi.list(), paymentMethodsApi.list()]);
			categories = c ?? [];
			paymentMethods = p ?? [];

			await fetchExpenses();
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			loading = false;
		}
	});
</script>

<div class="page">
	<div class="page-header">
		<h1>Expenses</h1>
		<div class="actions">
			<a href={resolve('/expenses/upload')}>
				<button class="secondary"><Upload size={16} /> Bulk Upload CSV</button>
			</a>
			<a href={resolve('/expenses/new')}>
				<button class="primary"><Plus size={16} /> New Expense</button>
			</a>
		</div>
	</div>

	<div class="filters card">
		<div class="filter-bar">
			<!-- Left filters -->
			<div class="filter-group main-filters">
				<Funnel size={16} style="opacity: 0.6; flex-shrink: 0;" />

				<select bind:value={filterCategory} onchange={applyFilters} class="filter-select">
					<option value="">All Categories</option>
					{#each categories as c (c.category_id)}
						<option value={c.category_id}>{c.category_name}</option>
					{/each}
				</select>

				<select bind:value={filterPm} onchange={applyFilters} class="filter-select">
					<option value="">All Payment Methods</option>
					{#each paymentMethods as p (p.payment_method_id)}
						<option value={p.payment_method_id}>{p.method_name}</option>
					{/each}
				</select>
			</div>

			<!-- Date range -->
			<div class="date-range">
				<div class="date-field">
					<label for="start-date" class="date-label">Desde</label>
					<input
						type="date"
						bind:value={filterStartDate}
						onchange={applyFilters}
						class="date-input"
					/>
				</div>
				<div class="date-field">
					<label for="end-date" class="date-label">Hasta</label>
					<input
						type="date"
						bind:value={filterEndDate}
						onchange={applyFilters}
						class="date-input"
					/>
				</div>
			</div>

			<!-- Show selector (pushed to the right) -->
			<div class="filter-group show-group">
				<span class="muted small">Show:</span>
				<select bind:value={filterLimit} onchange={applyFilters} class="filter-select">
					<option value={20}>20</option>
					<option value={50}>50</option>
					<option value={100}>100</option>
				</select>
			</div>
		</div>
	</div>

	<div class="card">
		{#if loading && expenses.length === 0}
			<div class="blank-state">
				<p>Loading expenses...</p>
			</div>
		{:else if error}
			<div class="blank-state error">
				<p>{error}</p>
			</div>
		{:else if expenses.length === 0}
			<div class="blank-state">
				<Receipt size={40} style="margin-bottom: 1rem; opacity: 0.5;" />
				<p>No expenses found.</p>
			</div>
		{:else}
			<div
				class="table-container"
				style="border-bottom-left-radius: 0; border-bottom-right-radius: 0;"
			>
				<table>
					<thead>
						<tr>
							<th>Date</th>
							<th>Merchant</th>
							<th>Category</th>
							<th>Payment Method</th>
							<th>Tags</th>
							<th style="text-align:right">Amount</th>
							<th></th>
						</tr>
					</thead>
					<tbody>
						{#each expenses as e (e.expense_id)}
							<tr>
								<td>{formatDate(e.expense_date)}</td>
								<td>
									<a href={resolve(`/expenses/${e.expense_id}`)}
										>{e.merchant_name ?? `#${e.expense_id}`}</a
									>
									{#if e.description}
										<div class="muted small">{e.description}</div>
									{/if}
								</td>
								<td
									>{e.category?.category_name ??
										categoryById.get(e.category_id)?.category_name ??
										'—'}</td
								>
								<td
									>{e.payment_method?.method_name ??
										pmById.get(e.payment_method_id)?.method_name ??
										'—'}</td
								>
								<td>
									{#each getExpenseTags(e) as tag (tag.tag_id)}
										<span class="badge">
											{#if tag.color}<span class="color-dot" style:background={tag.color}
												></span>{/if}
											{tag.tag_name}
										</span>
									{/each}
								</td>
								<td style="text-align:right">{formatAmount(e.amount, e.currency)}</td>
								<td>
									<div class="actions">
										<a href={resolve(`/expenses/${e.expense_id}`)}>
											<button title="Edit" class="secondary"><Pencil size={15} /></button>
										</a>
										<button class="danger" title="Delete" onclick={() => remove(e)}>
											<Trash2 size={15} />
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<div class="pagination">
				<div class="page-info">
					Showing {offset + 1} to {Math.min(offset + filterLimit, totalCount)} of {totalCount} records
				</div>
				<div class="page-controls">
					<button class="secondary" disabled={offset === 0} onclick={prevPage}>
						<ChevronLeft size={16} /> Prev
					</button>
					<span class="page-count">Page {currentPage} of {totalPages}</span>
					<button
						class="secondary"
						disabled={offset + filterLimit >= totalCount}
						onclick={nextPage}
					>
						Next <ChevronRight size={16} />
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	/* Keep all your existing styles – they are already excellent */
	.filters {
		margin-bottom: 1.5rem;
	}

	.filter-bar {
		display: flex;
		align-items: flex-end;
		gap: 1.25rem;
		flex-wrap: wrap;
	}

	.filter-group {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.filters select {
		padding: 0.4rem 0.75rem;
		border-radius: var(--radius-sm);
		border: 1px solid var(--border);
		background: var(--surface-2);
		color: var(--text);
		font-size: 0.9rem;
	}

	.main-filters {
		gap: 0.75rem;
	}

	.filter-select {
		padding: 0.45rem 0.85rem;
		border-radius: var(--radius-sm);
		border: 1px solid var(--border);
		background: var(--surface-2);
		color: var(--text);
		font-size: 0.95rem;
		min-width: 80px;
		height: 42px;
	}

	.table-container {
		overflow-x: auto;
	}

	.actions {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	/* Pagination Styles */
	.pagination {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.25rem;
		background: var(--surface);
		border: 1px solid var(--border);
		border-top: none;
		border-radius: 0 0 var(--radius) var(--radius);
		flex-wrap: wrap;
		gap: 1rem;
	}

	.page-info {
		font-size: 0.9rem;
		color: var(--text-muted);
	}

	.page-controls {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.page-count {
		font-size: 0.9rem;
		font-weight: 500;
		color: var(--text);
	}

	.page-controls button {
		display: flex;
		align-items: center;
		gap: 0.25rem;
	}

	.date-range {
		display: flex;
		align-items: flex-end;
		gap: 1rem;
	}

	.date-field {
		display: flex;
		flex-direction: column;
		gap: 0.15rem;
	}

	.date-label {
		font-size: 0.75rem;
		color: var(--text-muted);
		white-space: nowrap;
		margin-bottom: 2px;
	}

	.date-input {
		padding: 0.45rem 0.85rem;
		border-radius: var(--radius-sm);
		border: 1px solid var(--border);
		background: var(--surface-2);
		color: var(--text);
		font-size: 0.95rem;
		width: 160px;
		height: 42px;
	}

	.show-group {
		margin-left: auto;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	/* Responsive: stack on very small screens */
	@media (max-width: 640px) {
		.filter-bar {
			flex-direction: column;
			align-items: stretch;
		}
		.show-group {
			margin-left: 0;
			justify-content: flex-end;
		}
	}
</style>
