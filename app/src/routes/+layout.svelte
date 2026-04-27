<script lang="ts">
	import '$lib/styles.css';
	import { page } from '$app/state';
	import { Wallet, FolderTree, Tags, CreditCard, Receipt, LayoutDashboard } from 'lucide-svelte';

	let { children } = $props();

	const navItems = [
		{ href: '/', label: 'Dashboard', icon: LayoutDashboard, exact: true },
		{ href: '/expenses', label: 'Expenses', icon: Receipt },
		{ href: '/categories', label: 'Categories', icon: FolderTree },
		{ href: '/tags', label: 'Tags', icon: Tags },
		{ href: '/payment-methods', label: 'Payment methods', icon: CreditCard }
	];

	function isActive(href: string, exact = false): boolean {
		const path = page.url.pathname;
		if (exact) return path === href;
		return path === href || path.startsWith(href + '/');
	}
</script>

<svelte:head>
	<title>Expenses Tracker</title>
</svelte:head>

<div class="app">
	<aside class="sidebar">
		<div class="brand">
			<Wallet size={22} />
			<span>Expenses</span>
		</div>
		<nav>
			{#each navItems as item (item.href)}
				{@const ItemIcon = item.icon}
				<a href={item.href} class="nav-item" class:active={isActive(item.href, item.exact)}>
					<ItemIcon size={18} />
					<span>{item.label}</span>
				</a>
			{/each}
		</nav>
	</aside>
	<main class="main">
		{@render children()}
	</main>
</div>

<style>
	.app {
		display: grid;
		grid-template-columns: 240px 1fr;
		min-height: 100vh;
	}

	.sidebar {
		background: var(--surface);
		border-right: 1px solid var(--border);
		padding: 1.25rem 0.75rem;
		position: sticky;
		top: 0;
		align-self: start;
		height: 100vh;
	}

	.brand {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-weight: 700;
		font-size: 1.05rem;
		padding: 0.5rem 0.75rem 1.25rem;
		color: var(--primary);
	}

	nav {
		display: flex;
		flex-direction: column;
		gap: 0.15rem;
	}

	.nav-item {
		display: flex;
		align-items: center;
		gap: 0.65rem;
		padding: 0.6rem 0.75rem;
		border-radius: var(--radius-sm);
		color: var(--text-muted);
		text-decoration: none;
		font-weight: 500;
		font-size: 0.9rem;
	}

	.nav-item:hover {
		background: var(--surface-2);
		color: var(--text);
		text-decoration: none;
	}

	.nav-item.active {
		background: rgba(79, 70, 229, 0.1);
		color: var(--primary);
	}

	.main {
		min-width: 0;
	}

	@media (max-width: 720px) {
		.app {
			grid-template-columns: 1fr;
		}
		.sidebar {
			position: static;
			height: auto;
			border-right: none;
			border-bottom: 1px solid var(--border);
		}
		nav {
			flex-direction: row;
			overflow-x: auto;
		}
		.nav-item span {
			display: none;
		}
	}
</style>
