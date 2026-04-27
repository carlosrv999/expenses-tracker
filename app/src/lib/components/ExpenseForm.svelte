<script lang="ts">
	import type {
		Category,
		Expense,
		ExpenseInput,
		PaymentMethod,
		Tag
	} from '$lib/types';
	import {
		fromDateTimeLocal,
		toDateTimeLocal,
		toMajorAmount,
		toMinorAmount
	} from '$lib/format';
	import { Save, X } from 'lucide-svelte';

	let {
		initial,
		categories,
		paymentMethods,
		tags,
		onSubmit,
		submitting = false,
		cancelHref
	}: {
		initial?: Expense | null;
		categories: Category[];
		paymentMethods: PaymentMethod[];
		tags: Tag[];
		onSubmit: (input: ExpenseInput) => void | Promise<void>;
		submitting?: boolean;
		cancelHref: string;
	} = $props();

	/* svelte-ignore state_referenced_locally */
	let category_id = $state<number | ''>(initial?.category_id ?? '');
	/* svelte-ignore state_referenced_locally */
	let payment_method_id = $state<number | ''>(initial?.payment_method_id ?? '');
	/* svelte-ignore state_referenced_locally */
	let currency = $state(initial?.currency ?? 'PEN');
	/* svelte-ignore state_referenced_locally */
	let amountMajor = $state<string>(
		initial ? String(toMajorAmount(initial.amount).toFixed(2)) : ''
	);
	/* svelte-ignore state_referenced_locally */
	let expense_date = $state(toDateTimeLocal(initial?.expense_date));
	/* svelte-ignore state_referenced_locally */
	let merchant_name = $state(initial?.merchant_name ?? '');
	/* svelte-ignore state_referenced_locally */
	let description = $state(initial?.description ?? '');
	let initialTagIds = $derived(
		initial?.tag_ids ?? initial?.tags?.map((t) => t.tag_id) ?? []
	);
	/* svelte-ignore state_referenced_locally */
	let selectedTags = $state<Set<number>>(new Set(initialTagIds));

	function toggleTag(id: number) {
		const next = new Set(selectedTags);
		if (next.has(id)) next.delete(id);
		else next.add(id);
		selectedTags = next;
	}

	async function handleSubmit(event: Event) {
		event.preventDefault();
		if (category_id === '' || payment_method_id === '') return;
		const input: ExpenseInput = {
			category_id: Number(category_id),
			payment_method_id: Number(payment_method_id),
			currency: currency.toUpperCase(),
			amount: toMinorAmount(amountMajor),
			expense_date: fromDateTimeLocal(expense_date),
			merchant_name: merchant_name.trim() || null,
			description: description.trim() || null,
			tag_ids: [...selectedTags]
		};
		await onSubmit(input);
	}
</script>

<form onsubmit={handleSubmit} class="card">
	<div class="grid-form">
		<div>
			<label for="exp-category">Category *</label>
			<select id="exp-category" bind:value={category_id} required>
				<option value="" disabled>— Select —</option>
				{#each categories as c (c.category_id)}
					<option value={c.category_id}>{c.category_name}</option>
				{/each}
			</select>
		</div>
		<div>
			<label for="exp-pm">Payment method *</label>
			<select id="exp-pm" bind:value={payment_method_id} required>
				<option value="" disabled>— Select —</option>
				{#each paymentMethods as p (p.payment_method_id)}
					<option value={p.payment_method_id}>{p.method_name}</option>
				{/each}
			</select>
		</div>
		<div>
			<label for="exp-amount">Amount *</label>
			<input
				id="exp-amount"
				type="number"
				step="0.01"
				min="0.01"
				bind:value={amountMajor}
				required
			/>
		</div>
		<div>
			<label for="exp-currency">Currency</label>
			<input id="exp-currency" bind:value={currency} maxlength="3" required />
		</div>
		<div>
			<label for="exp-date">Date *</label>
			<input id="exp-date" type="datetime-local" bind:value={expense_date} required />
		</div>
		<div>
			<label for="exp-merchant">Merchant</label>
			<input id="exp-merchant" bind:value={merchant_name} maxlength="200" />
		</div>
		<div style="grid-column: 1 / -1;">
			<label for="exp-desc">Description</label>
			<textarea id="exp-desc" bind:value={description} rows="3"></textarea>
		</div>
		{#if tags.length > 0}
			<div style="grid-column: 1 / -1;">
				<span class="label-as-span">Tags</span>
				<div class="tag-list">
					{#each tags as t (t.tag_id)}
						<label class="tag-chip" class:selected={selectedTags.has(t.tag_id)}>
							<input
								type="checkbox"
								checked={selectedTags.has(t.tag_id)}
								onchange={() => toggleTag(t.tag_id)}
							/>
							{#if t.color}
								<span class="color-dot" style:background={t.color}></span>
							{/if}
							<span>{t.tag_name}</span>
						</label>
					{/each}
				</div>
			</div>
		{/if}
	</div>
	<div class="form-actions">
		<a href={cancelHref}><button type="button" class="ghost"><X size={15} /> Cancel</button></a>
		<button type="submit" class="primary" disabled={submitting}>
			<Save size={15} /> {submitting ? 'Saving…' : 'Save'}
		</button>
	</div>
</form>

<style>
	.label-as-span {
		display: block;
		margin-bottom: 0.4rem;
		font-weight: 500;
		color: var(--text-muted);
		font-size: 0.85rem;
	}
	.tag-list {
		display: flex;
		flex-wrap: wrap;
		gap: 0.4rem;
	}
	.tag-chip {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		padding: 0.35rem 0.7rem;
		border: 1px solid var(--border);
		border-radius: 999px;
		cursor: pointer;
		background: var(--surface);
		font-size: 0.85rem;
		margin: 0;
	}
	.tag-chip.selected {
		border-color: var(--primary);
		background: rgba(79, 70, 229, 0.08);
		color: var(--primary);
	}
	.tag-chip input {
		width: auto;
		margin: 0;
	}
</style>
