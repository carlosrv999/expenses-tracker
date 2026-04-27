<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { categoriesApi } from '$lib/api';
	import CategoryForm from '$lib/components/CategoryForm.svelte';
	import type { Category, CategoryInput } from '$lib/types';
	import { Trash2 } from 'lucide-svelte';

	let id = $derived(Number(page.params.id));
	let category = $state<Category | null>(null);
	let categories: Category[] = $state([]);
	let submitting = $state(false);
	let loading = $state(true);
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			const [cat, list] = await Promise.all([categoriesApi.get(id), categoriesApi.list()]);
			category = cat;
			categories = list ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			loading = false;
		}
	});

	async function handleSubmit(input: CategoryInput) {
		submitting = true;
		error = null;
		try {
			await categoriesApi.update(id, input);
			await goto('/categories');
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
			submitting = false;
		}
	}

	async function remove() {
		if (!category) return;
		if (!confirm(`Delete category "${category.category_name}"?`)) return;
		try {
			await categoriesApi.remove(id);
			await goto('/categories');
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1>Edit category</h1>
		{#if category}
			<button class="danger" onclick={remove}><Trash2 size={15} /> Delete</button>
		{/if}
	</div>
	{#if error}
		<div class="error">{error}</div>
	{/if}
	{#if loading}
		<p class="muted">Loading…</p>
	{:else if category}
		<CategoryForm
			initial={category}
			{categories}
			onSubmit={handleSubmit}
			{submitting}
			cancelHref="/categories"
		/>
	{:else}
		<p class="muted">Category not found.</p>
	{/if}
</div>
