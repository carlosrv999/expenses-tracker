<script lang="ts">
	import { onMount } from 'svelte';
	import { paymentMethodsApi } from '$lib/api';
	import type { PaymentMethod } from '$lib/types';
	import IconBadge from '$lib/components/IconBadge.svelte';
	import { CreditCard, Pencil, Trash2, Plus } from 'lucide-svelte';
	import type { IconComponent } from '$lib/icons';
	import { resolve } from '$app/paths';

	const ccFallback = CreditCard as unknown as IconComponent;

	let methods: PaymentMethod[] = $state([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	onMount(load);

	async function load() {
		loading = true;
		error = null;
		try {
			methods = (await paymentMethodsApi.list()) ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			loading = false;
		}
	}

	async function remove(id: number, name: string) {
		if (!confirm(`Delete payment method "${name}"?`)) return;
		try {
			await paymentMethodsApi.remove(id);
			await load();
		} catch (err) {
			alert(err instanceof Error ? err.message : String(err));
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1><CreditCard size={22} /> Payment methods</h1>
		<a href={resolve('/payment-methods/new')}><button class="primary"><Plus size={16} /> New method</button></a>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<p class="muted">Loading…</p>
	{:else if methods.length === 0}
		<div class="empty card">
			No payment methods yet. <a href={resolve('/payment-methods/new')}>Create the first one</a>.
		</div>
	{:else}
		<table>
			<thead>
				<tr>
					<th>Name</th>
					<th>Type</th>
					<th style="width: 1%"></th>
				</tr>
			</thead>
			<tbody>
				{#each methods as m (m.payment_method_id)}
					<tr>
						<td>
							<span class="icon-cell">
								<IconBadge name={m.icon} fallback={ccFallback} />
								<a href={resolve(`/payment-methods/${m.payment_method_id}`)}>{m.method_name}</a>
							</span>
						</td>
						<td><span class="badge">{m.method_type}</span></td>
						<td>
							<div class="actions">
								<a href={resolve(`/payment-methods/${m.payment_method_id}`)}><button title="Edit"><Pencil size={15} /></button></a>
								<button class="danger" title="Delete" onclick={() => remove(m.payment_method_id, m.method_name)}><Trash2 size={15} /></button>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>
