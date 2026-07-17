<script>
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Navbar from '$lib/components/Navbar.svelte';
	import { authState } from '$lib/auth.svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { sharedStore } from '$lib/store.svelte';

	let { children } = $props();

	// Navigation route guard to check login state
	$effect(() => {
		if (!authState.loading) {
			const path = page.url.pathname;
			if (!authState.isLoggedIn && !path.startsWith('/auth')) {
				goto('/auth/login');
			} else if (authState.isLoggedIn && path.startsWith('/auth')) {
				goto('/sales');
			}
		}
	});

	// Reactively fetch Postgres orders and sales upon authentication
	$effect(() => {
		if (authState.isLoggedIn) {
			sharedStore.init();
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<!-- Add font imports here -->
	<link rel="preconnect" href="https://fonts.googleapis.com" />
	<link rel="preconnect" href="https://fonts.gstatic.com" />
	<link
		href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap"
		rel="stylesheet"
	/>
</svelte:head>

{#if authState.loading}
	<div class="flex h-screen w-screen items-center justify-center bg-background text-on-surface">
		<div class="animate-pulse font-body-md text-body-lg">Loading Ledger...</div>
	</div>
{:else}
	<Navbar isLoggedIn={authState.isLoggedIn} />
	{@render children()}
{/if}
