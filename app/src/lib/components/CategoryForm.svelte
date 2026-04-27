<script lang="ts">
	import IconPicker from './IconPicker.svelte';
	import type { Category, CategoryInput } from '$lib/types';
	import { Save, X } from 'lucide-svelte';

	let {
		initial,
		categories = [],
		onSubmit,
		submitting = false,
		cancelHref
	}: {
		initial?: Category | null;
		categories?: Category[];
		onSubmit: (input: CategoryInput) => void | Promise<void>;
		submitting?: boolean;
		cancelHref: string;
	} = $props();

	/* svelte-ignore state_referenced_locally */
	let category_name = $state(initial?.category_name ?? '');
	/* svelte-ignore state_referenced_locally */
	let parent_category_id = $state<number | ''>(initial?.parent_category_id ?? '');
	/* svelte-ignore state_referenced_locally */
	let icon = $state(initial?.icon ?? '');
	/* svelte-ignore state_referenced_locally */
	let color = $state(initial?.color ?? '#4f46e5');
	/* svelte-ignore state_referenced_locally */
	let useColor = $state(Boolean(initial?.color));

	let parentOptions = $derived(
		categories.filter((c) => c.category_id !== initial?.category_id)
	);

	async function handleSubmit(event: Event) {
		event.preventDefault();
		const input: CategoryInput = {
			category_name: category_name.trim(),
			parent_category_id: parent_category_id === '' ? null : Number(parent_category_id),
			icon: icon || null,
			color: useColor ? color : null
		};
		await onSubmit(input);
	}
</script>

<form onsubmit={handleSubmit} class="card">
	<div class="grid-form">
		<div>
			<label for="cat-name">Name</label>
			<input id="cat-name" bind:value={category_name} required maxlength="200" />
		</div>
		<div>
			<label for="cat-parent">Parent category</label>
			<select id="cat-parent" bind:value={parent_category_id}>
				<option value="">— None (top-level) —</option>
				{#each parentOptions as opt (opt.category_id)}
					<option value={opt.category_id}>{opt.category_name}</option>
				{/each}
			</select>
		</div>
		<div>
			<label for="cat-icon">Icon</label>
			<IconPicker id="cat-icon" bind:value={icon} />
		</div>
		<div>
			<label for="cat-color-toggle">Color</label>
			<div class="color-row">
				<input
					id="cat-color-toggle"
					type="checkbox"
					bind:checked={useColor}
					class="checkbox"
				/>
				<input type="color" bind:value={color} disabled={!useColor} class="color" />
				<input type="text" bind:value={color} disabled={!useColor} maxlength="7" class="hex" />
			</div>
		</div>
	</div>
	<div class="form-actions">
		<a href={cancelHref}><button type="button" class="ghost"><X size={15} /> Cancel</button></a>
		<button type="submit" class="primary" disabled={submitting}>
			<Save size={15} /> {submitting ? 'Saving…' : 'Save'}
		</button>
	</div>
</form>

<style>
	.color-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	.checkbox {
		width: auto;
	}
	.color {
		width: 48px;
		padding: 0;
		height: 38px;
	}
	.hex {
		flex: 1;
		min-width: 0;
	}
</style>
