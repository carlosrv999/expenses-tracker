<script lang="ts">
	import { resolveIcon, type IconComponent } from '$lib/icons';
	import { hexToTintBackground } from '$lib/color';
	import { Tag as TagFallback } from 'lucide-svelte';

	let {
		name,
		fallback,
		color,
		size = 18
	}: {
		name: string | null | undefined;
		fallback?: IconComponent;
		color?: string | null;
		size?: number;
	} = $props();

	const defaultFallback = TagFallback as unknown as IconComponent;
	let Icon = $derived(resolveIcon(name, fallback ?? defaultFallback));
	let bg = $derived(color ? hexToTintBackground(color) : undefined);
</script>

<span class="icon-bubble" style:background={bg} style:color={color ?? undefined}>
	<Icon {size} />
</span>
