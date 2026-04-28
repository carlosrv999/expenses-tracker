<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { categoriesApi } from '$lib/api';
	import CategoryForm from '$lib/components/CategoryForm.svelte';
	import type { Category, CategoryInput } from '$lib/types';
	import { resolve } from '$app/paths';

	let categories: Category[] = $state([]);
	let submitting = $state(false);
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			categories = (await categoriesApi.list()) ?? [];
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
		}
	});

	async function handleSubmit(input: CategoryInput) {
		submitting = true;
		error = null;
		try {
			await categoriesApi.create(input);
			await goto(resolve('/categories'));
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
			submitting = false;
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1>New category</h1>
	</div>
	{#if error}
		<div class="error">{error}</div>
	{/if}
	<CategoryForm
		{categories}
		onSubmit={handleSubmit}
		{submitting}
		cancelHref={resolve('/categories')}
	/>
</div>
