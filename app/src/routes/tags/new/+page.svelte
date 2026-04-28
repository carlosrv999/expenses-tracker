<script lang="ts">
	import { goto } from '$app/navigation';
	import { tagsApi } from '$lib/api';
	import TagForm from '$lib/components/TagForm.svelte';
	import type { TagInput } from '$lib/types';
	import { resolve } from '$app/paths';

	let submitting = $state(false);
	let error = $state<string | null>(null);

	async function handleSubmit(input: TagInput) {
		submitting = true;
		error = null;
		try {
			await tagsApi.create(input);
			await goto(resolve('/tags'));
		} catch (err) {
			error = err instanceof Error ? err.message : String(err);
			submitting = false;
		}
	}
</script>

<div class="page">
	<div class="page-header">
		<h1>New tag</h1>
	</div>
	{#if error}
		<div class="error">{error}</div>
	{/if}
	<TagForm onSubmit={handleSubmit} {submitting} cancelHref={resolve('/tags')} />
</div>
