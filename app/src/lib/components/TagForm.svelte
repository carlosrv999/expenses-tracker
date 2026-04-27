<script lang="ts">
	import IconPicker from './IconPicker.svelte';
	import type { Tag, TagInput } from '$lib/types';
	import { Save, X } from 'lucide-svelte';

	let {
		initial,
		onSubmit,
		submitting = false,
		cancelHref
	}: {
		initial?: Tag | null;
		onSubmit: (input: TagInput) => void | Promise<void>;
		submitting?: boolean;
		cancelHref: string;
	} = $props();

	/* svelte-ignore state_referenced_locally */
	let tag_name = $state(initial?.tag_name ?? '');
	/* svelte-ignore state_referenced_locally */
	let icon = $state(initial?.icon ?? '');
	/* svelte-ignore state_referenced_locally */
	let color = $state(initial?.color ?? '#16a34a');
	/* svelte-ignore state_referenced_locally */
	let useColor = $state(Boolean(initial?.color));

	async function handleSubmit(event: Event) {
		event.preventDefault();
		const input: TagInput = {
			tag_name: tag_name.trim(),
			icon: icon || null,
			color: useColor ? color : null
		};
		await onSubmit(input);
	}
</script>

<form onsubmit={handleSubmit} class="card">
	<div class="grid-form">
		<div>
			<label for="tag-name">Name</label>
			<input id="tag-name" bind:value={tag_name} required maxlength="100" />
		</div>
		<div>
			<label for="tag-icon">Icon</label>
			<IconPicker id="tag-icon" bind:value={icon} />
		</div>
		<div>
			<label for="tag-color-toggle">Color</label>
			<div class="color-row">
				<input id="tag-color-toggle" type="checkbox" bind:checked={useColor} class="checkbox" />
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
