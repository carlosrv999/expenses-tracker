<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { tagsApi } from '$lib/api';
	import TagForm from '$lib/components/TagForm.svelte';
	import type { Tag, TagInput } from '$lib/types';
	import { Trash2 } from 'lucide-svelte';

	let id = $derived(Number(page.params.id));
	let tag = $state<Tag | null>(null);
	let loading = $state(true);
	let submitting = $state(false);
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			tag = await tagsApi.get(id);
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			loading = false;
		}
	});

	async function handleSubmit(input: TagInput) {
		submitting = true;
		error = null;
		try {
			await tagsApi.update(id, input);
			await goto('/tags');
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
			submitting = false;
		}
	}

	async function remove() {
		if (!tag) return;
		if (!confirm(`Delete tag "${tag.tag_name}"?`)) return;
		try {
			await tagsApi.remove(id);
			await goto('/tags');
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1>Edit tag</h1>
		{#if tag}
			<button class="danger" onclick={remove}><Trash2 size={15} /> Delete</button>
		{/if}
	</div>
	{#if error}
		<div class="error">{error}</div>
	{/if}
	{#if loading}
		<p class="muted">Loading…</p>
	{:else if tag}
		<TagForm initial={tag} onSubmit={handleSubmit} {submitting} cancelHref="/tags" />
	{:else}
		<p class="muted">Tag not found.</p>
	{/if}
</div>
