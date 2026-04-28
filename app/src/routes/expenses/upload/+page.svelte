<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import Papa from 'papaparse';
	import { categoriesApi, paymentMethodsApi, tagsApi } from '$lib/api';
	import type { Category, PaymentMethod, Tag } from '$lib/types';
	import { Upload, Save, X } from 'lucide-svelte';
	import { resolve } from '$app/paths';

	let categories: Category[] = $state([]);
	let paymentMethods: PaymentMethod[] = $state([]);
	let tags: Tag[] = $state([]);

	let file: File | null = $state(null);

	interface ParsedRow {
		expense_date: string;
		amount: number;
		category_id: number;
		payment_method_id: number;
		currency: string;
		merchant_name: string;
		description: string;
		tag_ids: number[];
	}

	let parsedRows: ParsedRow[] = $state([]);
	let step: 'select' | 'edit' | 'uploading' = $state('select');

	onMount(async () => {
		// Fetch actual data for our dropdowns
		const [c, p, t] = await Promise.all([
			categoriesApi.list(),
			paymentMethodsApi.list(),
			tagsApi.list()
		]);
		categories = c ?? [];
		paymentMethods = p ?? [];
		tags = t ?? [];
	});

	function handleFileChange(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files && target.files.length > 0) {
			file = target.files[0];
		}
	}

	function parseCsv() {
		if (!file) return;
		Papa.parse(file, {
			header: true,
			skipEmptyLines: true,
			complete: (results) => {
				// 1. Cast PapaParse output from unknown[] to an array of string Records
				const rawData = results.data as Record<string, string>[];

				// 2. Explicitly map every field to satisfy the ParsedRow interface perfectly
				parsedRows = rawData.map((row) => ({
					expense_date: row.expense_date || '',
					amount: Number(row.amount) || 0,
					category_id: Number(row.category_id) || 0,
					payment_method_id: Number(row.payment_method_id) || 0,
					currency: row.currency || '',
					merchant_name: row.merchant_name || '',
					description: row.description || '',
					tag_ids: row.tag_ids ? row.tag_ids.split(',').map(Number) : []
				}));

				step = 'edit';
			}
		});
	}

	async function submitCsv() {
		step = 'uploading';

		// 1. Convert edited JSON back to original CSV string format
		const csvFormatData = parsedRows.map((row) => ({
			...row,
			// Convert array of tag IDs back to comma separated string
			tag_ids: row.tag_ids.join(',')
		}));

		const csvString = Papa.unparse(csvFormatData);

		// 2. Create a new file blob
		const updatedFile = new File([csvString], file?.name || 'updated_expenses.csv', {
			type: 'text/csv'
		});

		// 3. Prepare FormData
		const formData = new FormData();
		formData.append('file', updatedFile);

		// 4. Send request to backend
		try {
			const res = await fetch('http://localhost:8080/api/v1/expenses/upload', {
				method: 'POST',
				headers: {
					Accept: 'application/json'
				},
				body: formData
			});

			if (res.ok) {
				goto(resolve('/expenses'));
			} else {
				alert('Failed to upload expenses');
				step = 'edit';
			}
		} catch (error) {
			console.error(error);
			alert('An error occurred during upload');
			step = 'edit';
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1>Bulk Upload Expenses</h1>
		{#if step === 'edit'}
			<div class="actions">
				<button class="secondary" onclick={() => (step = 'select')}>Cancel</button>
				<button class="primary" onclick={submitCsv}>
					<Save size={16} style="margin-right: 8px;" />
					Upload to Backend
				</button>
			</div>
		{/if}
	</div>

	{#if step === 'select'}
		<div class="upload-box">
			<input type="file" accept=".csv" onchange={handleFileChange} />
			<button disabled={!file} onclick={parseCsv}>
				<Upload size={16} style="margin-right: 8px;" />
				Parse & Edit File
			</button>
		</div>
	{:else if step === 'edit'}
		<div class="table-container">
			<table>
				<thead>
					<tr>
						<th>Date (ISO)</th>
						<th>Amount (Cents)</th>
						<th>Currency</th>
						<th>Merchant</th>
						<th>Category</th>
						<th>Payment Method</th>
						<th>Tags</th>
						<th>Description</th>
						<th></th>
					</tr>
				</thead>
				<tbody>
					{#each parsedRows as row, i (row)}
						<tr>
							<td><input type="text" bind:value={row.expense_date} /></td>
							<td><input type="number" bind:value={row.amount} /></td>
							<td><input type="text" style="width: 60px;" bind:value={row.currency} /></td>
							<td><input type="text" bind:value={row.merchant_name} /></td>

							<td>
								<select bind:value={row.category_id}>
									{#each categories as category, index (index)}
										<option value={category.category_id}>{category.category_name}</option>
									{/each}
								</select>
							</td>

							<td>
								<select bind:value={row.payment_method_id}>
									{#each paymentMethods as pm, index (index)}
										<option value={pm.payment_method_id}>{pm.method_name}</option>
									{/each}
								</select>
							</td>

							<td>
								<select multiple bind:value={row.tag_ids} style="height: 60px;">
									{#each tags as tag, index (index)}
										<option value={tag.tag_id}>{tag.tag_name}</option>
									{/each}
								</select>
							</td>

							<td><input type="text" bind:value={row.description} /></td>
							<td>
								<button class="danger" onclick={() => parsedRows.splice(i, 1)}>
									<X size={15} />
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{:else if step === 'uploading'}
		<div class="loading">
			<p>Uploading and processing expenses...</p>
		</div>
	{/if}
</div>

<style>
	.page-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	.upload-box {
		padding: 3rem;
		border: 2px dashed var(--border);
		border-radius: 8px;
		text-align: center;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		align-items: center;
	}

	.table-container {
		overflow-x: auto;
	}

	input,
	select {
		width: 100%;
		padding: 4px;
		border: 1px solid var(--border);
		border-radius: 4px;
	}

	.actions {
		display: flex;
		gap: 0.75rem;
	}

	.loading {
		padding: 4rem;
		text-align: center;
		color: var(--text-muted);
	}
</style>
