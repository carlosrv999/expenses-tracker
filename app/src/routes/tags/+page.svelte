<script lang="ts">
	import { onMount } from 'svelte';
	import { tagsApi } from '$lib/api';
	import type { Tag } from '$lib/types';
	import IconBadge from '$lib/components/IconBadge.svelte';
	import { Tag as TagIcon, Pencil, Trash2, Plus, Tags } from 'lucide-svelte';
	import type { IconComponent } from '$lib/icons';

	const tagFallback = TagIcon as unknown as IconComponent;

	let tags: Tag[] = $state([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	onMount(load);

	async function load() {
		loading = true;
		error = null;
		try {
			tags = (await tagsApi.list()) ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			loading = false;
		}
	}

	async function remove(id: number, name: string) {
		if (!confirm(`Delete tag "${name}"?`)) return;
		try {
			await tagsApi.remove(id);
			await load();
		} catch (err) {
			alert(err instanceof Error ? err.message : String(err));
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1><Tags size={22} /> Tags</h1>
		<a href="/tags/new"><button class="primary"><Plus size={16} /> New tag</button></a>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<p class="muted">Loading…</p>
	{:else if tags.length === 0}
		<div class="empty card">
			No tags yet. <a href="/tags/new">Create the first one</a>.
		</div>
	{:else}
		<table>
			<thead>
				<tr>
					<th>Name</th>
					<th>Color</th>
					<th style="width: 1%"></th>
				</tr>
			</thead>
			<tbody>
				{#each tags as t (t.tag_id)}
					<tr>
						<td>
							<span class="icon-cell">
								<IconBadge name={t.icon} fallback={tagFallback} color={t.color} />
								<a href={`/tags/${t.tag_id}`}>{t.tag_name}</a>
							</span>
						</td>
						<td>
							{#if t.color}
								<span class="color-dot" style:background={t.color}></span>
								<span class="muted">{t.color}</span>
							{:else}
								—
							{/if}
						</td>
						<td>
							<div class="actions">
								<a href={`/tags/${t.tag_id}`}><button title="Edit"><Pencil size={15} /></button></a>
								<button class="danger" title="Delete" onclick={() => remove(t.tag_id, t.tag_name)}><Trash2 size={15} /></button>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>
