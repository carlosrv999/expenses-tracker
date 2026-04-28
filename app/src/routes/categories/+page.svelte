<script lang="ts">
	import { onMount } from 'svelte';
	import { categoriesApi } from '$lib/api';
	import type { Category } from '$lib/types';
	import IconBadge from '$lib/components/IconBadge.svelte';
	import { Folder, FolderTree, Pencil, Trash2, Plus } from 'lucide-svelte';
	import type { IconComponent } from '$lib/icons';
	import { resolve } from '$app/paths';

	const folderIcon = Folder as unknown as IconComponent;

	let categories: Category[] = $state([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	let parentName = $derived((id: number | null) => {
		if (id == null) return null;
		return categories.find((c) => c.category_id === id)?.category_name ?? `#${id}`;
	});

	onMount(load);

	async function load() {
		loading = true;
		error = null;
		try {
			categories = (await categoriesApi.list()) ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			loading = false;
		}
	}

	async function remove(id: number, name: string) {
		if (!confirm(`Delete category "${name}"?`)) return;
		try {
			await categoriesApi.remove(id);
			await load();
		} catch (err) {
			alert(err instanceof Error ? err.message : String(err));
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1><FolderTree size={22} /> Categories</h1>
		<a href={resolve('/categories/new')}><button class="primary"><Plus size={16} /> New category</button></a>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<p class="muted">Loading…</p>
	{:else if categories.length === 0}
		<div class="empty card">
			No categories yet. <a href={resolve('/categories/new')}>Create the first one</a>.
		</div>
	{:else}
		<table>
			<thead>
				<tr>
					<th>Name</th>
					<th>Parent</th>
					<th>Color</th>
					<th style="width: 1%"></th>
				</tr>
			</thead>
			<tbody>
				{#each categories as c (c.category_id)}
					<tr>
						<td>
							<span class="icon-cell">
								<IconBadge name={c.icon} fallback={folderIcon} color={c.color} />
								<a href={resolve(`/categories/${c.category_id}`)}>{c.category_name}</a>
							</span>
						</td>
						<td>{parentName(c.parent_category_id) ?? '—'}</td>
						<td>
							{#if c.color}
								<span class="color-dot" style:background={c.color}></span>
								<span class="muted">{c.color}</span>
							{:else}
								—
							{/if}
						</td>
						<td>
							<div class="actions">
								<a href={resolve(`/categories/${c.category_id}`)}><button title="Edit"><Pencil size={15} /></button></a>
								<button class="danger" title="Delete" onclick={() => remove(c.category_id, c.category_name)}><Trash2 size={15} /></button>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>
