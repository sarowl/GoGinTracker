<script lang="ts">
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { fade } from 'svelte/transition';

	import { authState } from '$lib/auth.svelte';

	// --- Props using Svelte 5 Runes ---
	let { isLoggedIn = false, activeTab = '' } = $props();

	let isDropdownOpen = $state(false);

	function toggleDropdown() {
		isDropdownOpen = !isDropdownOpen;
	}

	function closeDropdown() {
		isDropdownOpen = false;
	}

	async function handleLogout() {
		isDropdownOpen = false;
		try {
			await authState.logout();
			await goto(resolve('/auth/login'));
		} catch (err) {
			console.error('Logout error:', err);
		}
	}
</script>

<svelte:window onclick={closeDropdown} />

<nav
	class="sticky top-0 z-50 w-full border-b border-outline-variant bg-surface/85 backdrop-blur-md transition-all duration-300"
>
	<div class="max-w-container-max mx-auto flex h-16 items-center justify-between px-gutter">
		<!-- Logo -->
		<div class="flex items-center gap-8">
			<a
				href={resolve('/')}
				class="font-headline-md text-headline-md tracking-tight text-primary transition-opacity hover:opacity-80"
			>
				Ledger
			</a>

			{#if isLoggedIn}
				<div class="hidden h-full items-center gap-6 md:flex" transition:fade={{ duration: 150 }}>
					<a
						href="/sales"
						class="relative flex h-16 items-center px-1 font-body-md text-sm font-medium transition-all duration-200 hover:text-primary {activeTab ===
						'sales'
							? 'border-b-2 border-primary text-primary'
							: 'text-on-surface-variant'}"
					>
						Sales
					</a>

					<a
						href="/expenses"
						class="relative flex h-16 items-center px-1 font-body-md text-sm font-medium transition-all duration-200 hover:text-primary {activeTab ===
						'expenses'
							? 'border-b-2 border-primary text-primary'
							: 'text-on-surface-variant'}"
					>
						Expenses
					</a>

					<a
						href="/orders"
						class="relative flex h-16 items-center px-1 font-body-md text-sm font-medium transition-all duration-200 hover:text-primary {activeTab ===
						'orders'
							? 'border-b-2 border-primary text-primary'
							: 'text-on-surface-variant'}"
					>
						Orders
					</a>
				</div>
			{/if}
		</div>

		<!-- Actions / Profile -->
		<div class="flex items-center gap-4">
			{#if isLoggedIn}
				<!-- Profile Dropdown Trigger -->
				<div class="relative">
					<button
						type="button"
						onclick={(e) => {
							e.stopPropagation();
							toggleDropdown();
						}}
						class="flex h-10 w-10 cursor-pointer items-center justify-center rounded-full border border-outline-variant bg-surface-container-low transition-all duration-200 hover:bg-surface-container-high focus:ring-2 focus:ring-primary/20 focus:outline-none"
						aria-label="User menu"
					>
						<!-- Avatar/Profile Icon SVG -->
						<svg
							class="h-5 w-5 text-on-surface"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
							/>
						</svg>
					</button>

					<!-- Dropdown Menu -->
					{#if isDropdownOpen}
						<div
							transition:fade={{ duration: 150 }}
							class="absolute right-0 z-50 mt-2 w-48 rounded-lg border border-outline-variant bg-surface p-1 shadow-lg backdrop-blur-md"
						>
							<a
								href="/settings"
								class="flex items-center gap-2 rounded px-3 py-2 text-sm font-medium text-on-surface transition-colors duration-150 hover:bg-surface-dim"
								onclick={closeDropdown}
							>
								<!-- Settings Icon -->
								<svg
									class="h-4 w-4 opacity-75"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
									/>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
									/>
								</svg>
								Settings
							</a>
							<button
								type="button"
								class="flex w-full cursor-pointer items-center gap-2 rounded px-3 py-2 text-left text-sm font-medium text-error transition-colors duration-150 hover:bg-error-container/20"
								onclick={handleLogout}
							>
								<!-- Logout Icon -->
								<svg
									class="h-4 w-4 opacity-75"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
									/>
								</svg>
								Log out
							</button>
						</div>
					{/if}
				</div>
			{/if}
		</div>
	</div>

	<!-- Mobile Navigation (Only shown when logged in) -->
	{#if isLoggedIn}
		<div
			class="flex items-center justify-around border-t border-outline-variant bg-surface/50 py-2 backdrop-blur-md md:hidden"
		>
			<a
				href="/sales"
				class="flex flex-col items-center text-xs font-medium transition-colors hover:text-primary {activeTab ===
				'sales'
					? 'text-primary'
					: 'text-on-surface-variant'}"
			>
				<span class="mt-0.5">Sales</span>
			</a>

			<a
				href="/expenses"
				class="flex flex-col items-center text-xs font-medium transition-colors hover:text-primary {activeTab ===
				'expenses'
					? 'text-primary'
					: 'text-on-surface-variant'}"
			>
				<span class="mt-0.5">Expenses</span>
			</a>

			<a
				href="/orders"
				class="flex flex-col items-center text-xs font-medium transition-colors hover:text-primary {activeTab ===
				'orders'
					? 'text-primary'
					: 'text-on-surface-variant'}"
			>
				<span class="mt-0.5">Orders</span>
			</a>
		</div>
	{/if}
</nav>
