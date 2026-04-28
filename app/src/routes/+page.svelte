<script lang="ts">
	import { onMount } from 'svelte';
	import { categoriesApi, expensesApi, paymentMethodsApi, tagsApi } from '$lib/api';
	import { formatAmount, formatDate } from '$lib/format';
	import type { Category, Expense, PaymentMethod, Tag } from '$lib/types';
	import { Receipt, FolderTree, Tags as TagsIcon, CreditCard, Plus } from 'lucide-svelte';
	import { SvelteMap } from 'svelte/reactivity';
	import { resolve } from '$app/paths';

	let expenses: Expense[] = $state([]);
	let totalExpenses = $state(0); // ← NEW: real total from backend
	let categories: Category[] = $state([]);
	let tags: Tag[] = $state([]);
	let paymentMethods: PaymentMethod[] = $state([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	let categoryById = $derived(new Map(categories.map((c) => [c.category_id, c])));
	let pmById = $derived(new Map(paymentMethods.map((p) => [p.payment_method_id, p])));

	let totalByCurrency = $derived.by(() => {
		const totals = new SvelteMap<string, number>();
		for (const e of expenses) {
			totals.set(e.currency, (totals.get(e.currency) ?? 0) + e.amount);
		}
		return [...totals.entries()];
	});

	let recent = $derived(expenses.slice(0, 5));

	onMount(async () => {
		try {
			const [paginated, c, t, p] = await Promise.all([
				expensesApi.list({ limit: 100 }), // still limited for "Recent expenses"
				categoriesApi.list(),
				tagsApi.list(),
				paymentMethodsApi.list()
			]);

			expenses = paginated?.expenses ?? [];
			totalExpenses = paginated?.total_count ?? 0; // ← THIS IS THE FIX
			categories = c ?? [];
			tags = t ?? [];
			paymentMethods = p ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			loading = false;
		}
	});
</script>

<div class="page">
	<div class="page-header">
		<h1>Dashboard</h1>
		<a class="primary-link" href={resolve('/expenses/new')}>
			<button class="primary"><Plus size={16} /> New expense</button>
		</a>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<p class="muted">Loading…</p>
	{:else}
		<div class="stats">
			<a href={resolve('/expenses')} class="stat card">
				<div class="stat-icon"><Receipt size={20} /></div>
				<div>
					<div class="stat-label">Expenses</div>
					<div class="stat-value">{totalExpenses}</div>
				</div>
			</a>
			<a href={resolve('/categories')} class="stat card">
				<div class="stat-icon"><FolderTree size={20} /></div>
				<div>
					<div class="stat-label">Categories</div>
					<div class="stat-value">{categories.length}</div>
				</div>
			</a>
			<a href={resolve('/tags')} class="stat card">
				<div class="stat-icon"><TagsIcon size={20} /></div>
				<div>
					<div class="stat-label">Tags</div>
					<div class="stat-value">{tags.length}</div>
				</div>
			</a>
			<a href={resolve('/payment-methods')} class="stat card">
				<div class="stat-icon"><CreditCard size={20} /></div>
				<div>
					<div class="stat-label">Payment methods</div>
					<div class="stat-value">{paymentMethods.length}</div>
				</div>
			</a>
		</div>

		{#if totalByCurrency.length > 0}
			<div class="card totals">
				<h2>Totals</h2>
				<div class="totals-grid">
					{#each totalByCurrency as [currency, amount] (currency)}
						<div class="total">
							<div class="muted">{currency}</div>
							<div class="total-value">{formatAmount(amount, currency)}</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<div class="card">
			<div class="recent-header">
				<h2>Recent expenses</h2>
				<a href={resolve('/expenses')}>View all →</a>
			</div>
			{#if recent.length === 0}
				<p class="empty">
					No expenses yet. <a href={resolve('/expenses/new')}>Create the first one</a>.
				</p>
			{:else}
				<table>
					<thead>
						<tr>
							<th>Date</th>
							<th>Merchant</th>
							<th>Category</th>
							<th>Payment</th>
							<th style="text-align:right">Amount</th>
						</tr>
					</thead>
					<tbody>
						{#each recent as e (e.expense_id)}
							<tr>
								<td>{formatDate(e.expense_date)}</td>
								<td>{e.merchant_name ?? '—'}</td>
								<td>{categoryById.get(e.category_id)?.category_name ?? '—'}</td>
								<td>{pmById.get(e.payment_method_id)?.method_name ?? '—'}</td>
								<td style="text-align:right">{formatAmount(e.amount, e.currency)}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/if}
		</div>
	{/if}
</div>

<style>
	.primary-link {
		text-decoration: none;
	}

	.stats {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1rem;
		margin-bottom: 1.25rem;
	}

	.stat {
		display: flex;
		align-items: center;
		gap: 0.85rem;
		text-decoration: none;
		color: inherit;
		transition: transform 0.1s ease;
	}

	.stat:hover {
		transform: translateY(-1px);
		text-decoration: none;
	}

	.stat-icon {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 40px;
		height: 40px;
		border-radius: 10px;
		background: rgba(79, 70, 229, 0.1);
		color: var(--primary);
	}

	.stat-label {
		font-size: 0.8rem;
		color: var(--text-muted);
	}

	.stat-value {
		font-size: 1.5rem;
		font-weight: 700;
	}

	.totals {
		margin-bottom: 1.25rem;
	}

	.totals-grid {
		display: flex;
		gap: 2rem;
		flex-wrap: wrap;
	}

	.total-value {
		font-size: 1.4rem;
		font-weight: 700;
	}

	.recent-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 0.5rem;
	}

	.recent-header h2 {
		margin: 0;
	}
</style>
