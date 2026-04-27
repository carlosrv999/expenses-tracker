<script lang="ts">
	import { AVAILABLE_ICONS, ICON_MAP } from '$lib/icons';

	let { value = $bindable(''), id }: { value?: string | null; id?: string } = $props();
</script>

<select {id} bind:value>
	<option value="">— No icon —</option>
	{#each AVAILABLE_ICONS as name (name)}
		<option value={name}>{name}</option>
	{/each}
</select>

{#if value && ICON_MAP[value]}
	{@const Icon = ICON_MAP[value]}
	<div class="preview">
		<span class="icon-bubble"><Icon size={18} /></span>
		<span class="muted">{value}</span>
	</div>
{/if}

<style>
	.preview {
		margin-top: 0.5rem;
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
	}
</style>
