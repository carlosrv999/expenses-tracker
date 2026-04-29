<script lang="ts">
	import { onMount } from 'svelte';
	import { categoriesApi, expensesApi, paymentMethodsApi, tagsApi } from '$lib/api';
	import { formatAmount, formatDate } from '$lib/format';
	import type { Category, Expense, PaymentMethod, Tag } from '$lib/types';
	import { Receipt, FolderTree, Tags as TagsIcon, CreditCard, Plus } from 'lucide-svelte';
	import { SvelteMap } from 'svelte/reactivity';
	import { resolve } from '$app/paths';
	import Chart from 'chart.js/auto';

	// ── Core data (all-time: stats, recent list, totals) ────────────────────────
	let expenses: Expense[] = $state([]);
	let totalExpenses = $state(0);
	let categories: Category[] = $state([]);
	let tags: Tag[] = $state([]);
	let paymentMethods: PaymentMethod[] = $state([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	// ── Pie-chart filter state ───────────────────────────────────────────────────
	const NOW = new Date();
	const CURRENT_YEAR = NOW.getFullYear();
	const CURRENT_MONTH = NOW.getMonth() + 1; // 1-based
	const FIRST_YEAR = 2025;

	const MONTH_NAMES = [
		'January',
		'February',
		'March',
		'April',
		'May',
		'June',
		'July',
		'August',
		'September',
		'October',
		'November',
		'December'
	];

	// Available years for the dropdown
	const yearOptions: Array<{ label: string; value: number | null }> = [
		{ label: 'All', value: null },
		...Array.from({ length: CURRENT_YEAR - FIRST_YEAR + 1 }, (_, i) => {
			const y = CURRENT_YEAR - i;
			return { label: String(y), value: y };
		})
	];

	let selectedYear = $state<number | null>(null); // null = All
	let selectedMonth = $state<number | null>(null); // null = All, 1–12 otherwise

	// Month options depend on selected year
	let monthOptions = $derived.by(() => {
		if (selectedYear === null) return [];
		const maxMonth = selectedYear === CURRENT_YEAR ? CURRENT_MONTH : 12;
		return [
			{ label: 'All', value: null },
			...Array.from({ length: maxMonth }, (_, i) => ({
				label: MONTH_NAMES[i],
				value: i + 1
			}))
		];
	});

	// ── Pie-chart expenses (re-fetched when filters change) ──────────────────────
	let chartExpenses: Expense[] = $state([]);
	let chartLoading = $state(false);
	let chartError = $state<string | null>(null);

	/** Build start_date / end_date strings from the current filter selection. */
	function buildDateRange(
		year: number | null,
		month: number | null
	): { start_date?: string; end_date?: string } {
		if (year === null) return {};

		if (month !== null) {
			const lastDay = new Date(year, month, 0).getDate();
			return {
				start_date: `${year}-${String(month).padStart(2, '0')}-01`,
				end_date: `${year}-${String(month).padStart(2, '0')}-${lastDay}`
			};
		} else {
			return {
				start_date: `${year}-01-01`,
				end_date: `${year}-12-31`
			};
		}
	}

	/** Fetch all pages for the given filters (handles >500 expenses). */
	async function fetchAllChartExpenses(
		year: number | null,
		month: number | null
	): Promise<Expense[]> {
		const dateRange = buildDateRange(year, month);
		const PAGE = 500;
		let offset = 0;
		const all: Expense[] = [];

		while (true) {
			const page = await expensesApi.list({ ...dateRange, limit: PAGE, offset });
			const rows = page?.expenses ?? [];
			all.push(...rows);
			const total = page?.total_count ?? 0;
			offset += rows.length;
			if (all.length >= total || rows.length === 0) break;
		}
		return all;
	}

	// Re-fetch chart data whenever year or month changes.
	// Capture the values as local consts so the async closure uses the
	// snapshot at the time the effect ran, not a later value.
	$effect(() => {
		const year = selectedYear;
		const month = selectedYear === null ? null : selectedMonth;

		chartLoading = true;
		chartError = null;

		fetchAllChartExpenses(year, month)
			.then((rows) => {
				chartExpenses = rows;
			})
			.catch((err) => {
				chartError = err instanceof Error ? err.message : String(err);
			})
			.finally(() => {
				chartLoading = false;
			});
	});

	// ── Derived lookups ──────────────────────────────────────────────────────────
	let categoryById = $derived(new Map(categories.map((c) => [c.category_id, c])));
	let pmById = $derived(new Map(paymentMethods.map((p) => [p.payment_method_id, p])));

	// Totals / recent always use the all-time dataset
	let totalByCurrency = $derived.by(() => {
		const totals = new SvelteMap<string, number>();
		for (const e of expenses) {
			totals.set(e.currency, (totals.get(e.currency) ?? 0) + e.amount);
		}
		return [...totals.entries()];
	});

	let recent = $derived(expenses.slice(0, 5));

	// Primary currency derived from the FILTERED chart expenses
	let chartPrimaryCurrency = $derived.by(() => {
		const totals = new SvelteMap<string, number>();
		for (const e of chartExpenses) {
			totals.set(e.currency, (totals.get(e.currency) ?? 0) + e.amount);
		}
		const entries = [...totals.entries()];
		if (entries.length === 0) return null;
		return entries.reduce((max, curr) => (curr[1] > max[1] ? curr : max))[0];
	});

	// Filter chart expenses to primary currency for accurate pie
	let filteredChartExpenses = $derived.by(() =>
		chartPrimaryCurrency
			? chartExpenses.filter((e) => e.currency === chartPrimaryCurrency)
			: chartExpenses
	);

	// Spending by category (for pie)
	let categorySpending = $derived.by(() => {
		const map = new SvelteMap<number, { amount: number; name: string; color: string }>();
		for (const e of filteredChartExpenses) {
			const cat = categoryById.get(e.category_id);
			if (!cat) continue;
			const key = e.category_id;
			const current = map.get(key) ?? {
				amount: 0,
				name: cat.category_name,
				color: cat.color || '#64748b'
			};
			current.amount += e.amount;
			map.set(key, current);
		}

		const items = Array.from(map.values());
		const totalAmount = items.reduce((sum, item) => sum + item.amount, 0);

		const dataItems = items
			.map((item) => ({
				...item,
				percentage: totalAmount > 0 ? Math.round((item.amount / totalAmount) * 100) : 0
			}))
			.sort((a, b) => b.amount - a.amount);

		return { items: dataItems, totalAmount };
	});

	// ── Chart.js ─────────────────────────────────────────────────────────────────
	let pieCanvas: HTMLCanvasElement | undefined = $state();
	let chartInstance: Chart | null = null;

	$effect(() => {
		if (loading || chartLoading || !pieCanvas || categorySpending.items.length === 0) {
			if (chartInstance) {
				chartInstance.destroy();
				chartInstance = null;
			}
			return;
		}

		if (chartInstance) {
			chartInstance.destroy();
		}

		const ctx = pieCanvas.getContext('2d');
		if (!ctx) return;

		chartInstance = new Chart(ctx, {
			type: 'doughnut',
			data: {
				labels: categorySpending.items.map((item) => item.name),
				datasets: [
					{
						data: categorySpending.items.map((item) => item.amount),
						backgroundColor: categorySpending.items.map((item) => item.color),
						borderColor: '#ffffff',
						borderWidth: 4,
						hoverOffset: 25
					}
				]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				cutout: '68%',
				plugins: {
					legend: {
						position: 'bottom',
						labels: {
							padding: 24,
							usePointStyle: true,
							boxWidth: 12,
							font: { size: 13 }
						}
					},
					tooltip: {
						callbacks: {
							label: (context) => {
								const value = context.raw as number;
								const idx = context.dataIndex;
								const item = categorySpending.items[idx];
								const perc = item?.percentage ?? 0;
								return `${context.label}: ${formatAmount(
									value,
									chartPrimaryCurrency || 'PEN'
								)} (${perc}%)`;
							}
						}
					}
				}
			}
		});
	});

	// ── Initial load (all-time data for stats / recent) ──────────────────────────
	onMount(async () => {
		try {
			const [paginated, c, t, p] = await Promise.all([
				expensesApi.list({ limit: 500 }),
				categoriesApi.list(),
				tagsApi.list(),
				paymentMethodsApi.list()
			]);

			expenses = paginated?.expenses ?? [];
			totalExpenses = paginated?.total_count ?? 0;
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

		<!-- Spending by Category – Pie/Donut Chart -->
		<div class="card">
			<div class="chart-header">
				<h2>
					Spending by Category
					{#if chartPrimaryCurrency}
						<span class="muted">({chartPrimaryCurrency})</span>
					{/if}
				</h2>

				<div class="chart-filters">
					<select
						bind:value={selectedYear}
						onchange={() => {
							selectedMonth = null;
						}}
					>
						{#each yearOptions as opt (opt.value)}
							<option value={opt.value}>{opt.label}</option>
						{/each}
					</select>

					<select bind:value={selectedMonth} disabled={selectedYear === null}>
						{#if selectedYear === null}
							<option value={null}>Month</option>
						{:else}
							{#each monthOptions as opt (opt.value)}
								<option value={opt.value}>{opt.label}</option>
							{/each}
						{/if}
					</select>
				</div>
			</div>

			{#if !chartLoading && categorySpending.items.length > 0}
				<div class="chart-total">
					<span class="chart-total-label">Total</span>
					<span class="chart-total-value"
						>{formatAmount(categorySpending.totalAmount, chartPrimaryCurrency || 'PEN')}</span
					>
				</div>
			{/if}

			{#if chartLoading}
				<p class="muted chart-loading">Loading chart…</p>
			{:else if chartError}
				<div class="error">{chartError}</div>
			{:else if categorySpending.items.length === 0}
				<p class="empty">No expenses for the selected period.</p>
			{:else}
				<div class="pie-container">
					<canvas bind:this={pieCanvas}></canvas>
				</div>
			{/if}
		</div>

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

	/* Chart card header: title left, filters right */
	.chart-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		flex-wrap: wrap;
		gap: 0.75rem;
		margin-bottom: 0.5rem;
	}

	.chart-header h2 {
		margin: 0;
	}

	.chart-filters {
		display: flex;
		gap: 0.5rem;
	}

	.chart-filters select {
		padding: 0.3rem 0.6rem;
		border: 1px solid var(--border, #e2e8f0);
		border-radius: 6px;
		font-size: 0.85rem;
		background: var(--surface, #fff);
		color: var(--text, inherit);
		cursor: pointer;
	}

	.chart-filters select:disabled {
		opacity: 0.45;
		cursor: not-allowed;
	}

	.chart-loading {
		padding: 2rem 0;
		text-align: center;
	}

	.chart-total {
		display: flex;
		align-items: baseline;
		gap: 0.5rem;
		margin: 0.25rem 0 0.75rem;
	}

	.chart-total-label {
		font-size: 0.8rem;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.04em;
	}

	.chart-total-value {
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

	/* Consistent vertical spacing between all major dashboard cards */
	.card {
		margin-bottom: 1.25rem;
	}

	/* Pie chart container */
	.pie-container {
		position: relative;
		height: 340px;
		width: 100%;
		margin: 1rem 0;
	}

	.pie-container canvas {
		max-width: 100%;
	}

	.total-value {
		font-size: 1.4rem;
		font-weight: 700;
	}

	.totals-grid {
		display: flex;
		gap: 2rem;
		flex-wrap: wrap;
	}
</style>
