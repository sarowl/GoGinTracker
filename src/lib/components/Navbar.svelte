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
						href="/dashboard"
						class="relative flex h-16 items-center px-1 font-body-md text-sm font-medium transition-all duration-200 hover:text-primary {activeTab ===
						'dashboard'
							? 'border-b-2 border-primary text-primary'
							: 'text-on-surface-variant'}"
					>
						Dashboard
					</a>

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
							class="absolute right-0 z-50 mt-2 w-52 rounded-lg border border-outline-variant bg-surface p-2 shadow-lg backdrop-blur-md"
						>
							<div class="border-b border-outline-variant px-3 py-2 text-xs">
								<p class="font-bold text-on-surface truncate">{authState.user?.displayName || authState.user?.email}</p>
								<div class="mt-1 flex items-center gap-1.5">
									<span class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider {authState.isOwner ? 'bg-amber-100 text-amber-900 border border-amber-300' : 'bg-blue-100 text-blue-900 border border-blue-300'}">
										{authState.isOwner ? '👑 Owner' : '💼 Employee'}
									</span>
									{#if authState.profile?.company_code}
										<span class="text-[10px] font-mono text-stone-500">({authState.profile.company_code})</span>
									{/if}
								</div>
							</div>
							<a
								href="/company"
								class="mt-1 flex items-center gap-2 rounded px-3 py-2 text-sm font-medium text-on-surface transition-colors duration-150 hover:bg-surface-dim"
								onclick={closeDropdown}
							>
								🏢 Company Setup
							</a>
							<a
								href="/settings"
								class="flex items-center gap-2 rounded px-3 py-2 text-sm font-medium text-on-surface transition-colors duration-150 hover:bg-surface-dim"
								onclick={closeDropdown}
							>
								⚙️ Settings
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
