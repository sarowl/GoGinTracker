<script>
	import ExpenseModal from '$lib/components/ExpenseModal.svelte';
	import { sharedStore } from '$lib/store.svelte';
	import { authState } from '$lib/auth.svelte';

	let isModalOpen = $state(false);
	let selectedCategory = $state('All');

	const categories = ['All', 'Utilities', 'Supplies', 'Payroll', 'Delivery/Logistics', 'Rent', 'Misc'];

	// Derived statistics
	let totalExpensesValue = $derived(
		sharedStore.expenses.reduce((sum, item) => sum + (item.amount || 0), 0)
	);

	let filteredExpenses = $derived(
		selectedCategory === 'All'
			? sharedStore.expenses
			: sharedStore.expenses.filter((item) => item.category === selectedCategory)
	);

	let expenseNoSuggestion = $derived(
		`EXP-${String(sharedStore.expenses.length + 1001).padStart(4, '0')}`
	);

	function openModal() {
		isModalOpen = true;
	}

	function closeModal() {
		isModalOpen = false;
	}

	/**
	 * @param {any} newExpense
	 */
	function handleAddExpense(newExpense) {
		const userUid = authState.user?.uid || 'System User';
		sharedStore.addExpense({
			...newExpense,
			recordedBy: userUid
		});
		closeModal();
	}

	/**
	 * @param {string|number} id
	 */
	function removeExpense(id) {
		if (!authState.isOwner) {
			alert('Employees are not permitted to delete expenses.');
			return;
		}
		if (confirm('Are you sure you want to delete this expense entry?')) {
			sharedStore.removeExpense(id, authState.role);
		}
	}
</script>

<svelte:head>
	<title>Expenses — Ledger</title>
	<meta name="description" content="Track and audit company operating expenses." />
</svelte:head>

<div class="bg-background text-on-surface p-xl relative min-h-screen overflow-hidden">
	<!-- Ambient Background Glow -->
	<div class="pointer-events-none absolute inset-0 z-0 overflow-hidden opacity-40">
		<div class="glow-blob -top-24 -left-24"></div>
		<div class="glow-blob-2 right-0 bottom-0"></div>
	</div>

	<div class="max-w-container-max relative z-10 mx-auto space-y-md">
		<!-- Header Title -->
		<header class="mb-xl flex flex-col gap-md md:flex-row md:items-center md:justify-between">
			<div>
				<h1 class="font-headline-lg text-headline-lg text-primary tracking-tight">
					Expenses Ledger
				</h1>
				<p class="font-body-md text-body-md text-on-surface-variant">
					Log, categorize, and audit business operating expenses and vendor payouts.
				</p>
			</div>
			<div>
				<button
					onclick={openModal}
					class="px-lg bg-primary text-on-primary font-body-md shadow-md hover:brightness-110 active:scale-[0.98] inline-flex h-10 cursor-pointer items-center justify-center gap-2 rounded text-sm font-semibold transition-all"
				>
					<svg
						class="h-4 w-4"
						fill="none"
						stroke="currentColor"
						stroke-width="2.5"
						viewBox="0 0 24 24"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
					</svg>
					Log Expense
				</button>
			</div>
		</header>

		<!-- Metrics Cards Grid -->
		<div class="grid grid-cols-1 gap-md md:grid-cols-3">
			<div class="glass-card space-y-sm rounded-xl p-lg">
				<p class="font-label-md text-label-md text-on-surface-variant uppercase tracking-wider">
					Total Expenses Outflow
				</p>
				<p class="font-headline-lg text-headline-lg text-primary tracking-tight">
					₱{totalExpensesValue.toLocaleString('en-US', {
						minimumFractionDigits: 2,
						maximumFractionDigits: 2
					})}
				</p>
				<span class="inline-flex items-center text-xs font-semibold text-stone-500"
					>Sum of all logged company payouts</span
				>
			</div>
			<div class="glass-card space-y-sm rounded-xl p-lg">
				<p class="font-label-md text-label-md text-on-surface-variant uppercase tracking-wider">
					Expense Transactions
				</p>
				<p class="font-headline-lg text-headline-lg text-primary tracking-tight">
					{sharedStore.expenses.length}
				</p>
				<span class="inline-flex items-center text-xs font-semibold text-stone-500"
					>Total recorded expense entries</span
				>
			</div>
			<div class="glass-card space-y-sm rounded-xl p-lg">
				<p class="font-label-md text-label-md text-on-surface-variant uppercase tracking-wider">
					Average Expense Value
				</p>
				<p class="font-headline-lg text-headline-lg text-primary tracking-tight">
					₱{sharedStore.expenses.length
						? (totalExpensesValue / sharedStore.expenses.length).toLocaleString('en-US', {
								minimumFractionDigits: 2,
								maximumFractionDigits: 2
							})
						: '0.00'}
				</p>
				<span class="inline-flex items-center text-xs font-semibold text-stone-500"
					>Average cost per logged expense</span
				>
			</div>
		</div>

		<!-- Main Expenses Container -->
		<div class="glass-card space-y-md rounded-xl p-lg">
			<!-- Category Filter Tabs -->
			<div class="border-outline-variant flex flex-wrap items-center gap-xs border-b pb-sm">
				{#each categories as cat}
					<button
						onclick={() => (selectedCategory = cat)}
						class="px-md py-1.5 font-body-md rounded-lg text-xs font-semibold transition-all cursor-pointer {selectedCategory ===
						cat
							? 'bg-primary text-on-primary shadow-sm'
							: 'text-on-surface-variant hover:bg-surface-dim/60 hover:text-on-surface'}"
					>
						{cat}
					</button>
				{/each}
			</div>

			<!-- Expenses Table -->
			<div class="border-outline-variant overflow-x-auto rounded-lg border">
				<table class="w-full border-collapse text-left">
					<thead>
						<tr class="border-outline-variant bg-surface-dim border-b">
							<th
								class="p-md font-label-md text-label-md text-on-surface-variant min-w-[140px] uppercase tracking-wider whitespace-nowrap"
								>Vendor Invoice #</th
							>
							<th
								class="p-md font-label-md text-label-md text-on-surface-variant min-w-[130px] uppercase tracking-wider whitespace-nowrap"
								>Date</th
							>
							<th
								class="p-md font-label-md text-label-md text-on-surface-variant min-w-[140px] uppercase tracking-wider whitespace-nowrap"
								>Category</th
							>
							<th
								class="p-md font-label-md text-label-md text-on-surface-variant min-w-[160px] uppercase tracking-wider whitespace-nowrap"
								>Vendor / Supplier</th
							>
							<th
								class="p-md font-label-md text-label-md text-on-surface-variant min-w-[220px] uppercase tracking-wider whitespace-nowrap"
								>Description</th
							>
							<th
								class="p-md font-label-md text-label-md text-on-surface-variant min-w-[130px] uppercase tracking-wider whitespace-nowrap"
								>Payment</th
							>
							<th
								class="p-md font-label-md text-label-md text-on-surface-variant min-w-[140px] text-right uppercase tracking-wider whitespace-nowrap"
								>Amount</th
							>
							<th
								class="p-md font-label-md text-label-md text-on-surface-variant min-w-[90px] text-center uppercase tracking-wider whitespace-nowrap"
								>Actions</th
							>
						</tr>
					</thead>
					<tbody class="divide-outline-variant bg-surface/50 divide-y">
						{#if filteredExpenses.length === 0}
							<tr>
								<td
									colspan="8"
									class="p-xl text-on-surface-variant text-center font-body-md whitespace-nowrap"
								>
									No expenses logged {selectedCategory !== 'All'
										? `in category "${selectedCategory}"`
										: 'yet'}. Click "Log Expense" above to record an expense.
								</td>
							</tr>
						{:else}
							{#each filteredExpenses as item (item.id)}
								<tr class="hover:bg-surface-dim/40 transition-colors">
									<td
										class="p-md font-body-md text-primary text-sm font-semibold whitespace-nowrap"
										>{item.vendorInvoiceNo || '—'}</td
									>
									<td class="p-md font-body-md text-on-surface-variant text-sm whitespace-nowrap">
										{item.expenseDate ? item.expenseDate.split('T')[0] : 'N/A'}
									</td>
									<td class="p-md font-body-md text-sm whitespace-nowrap">
										<span
											class="inline-flex items-center rounded-full border border-stone-200 bg-stone-100 px-2.5 py-0.5 text-xs font-semibold text-stone-700"
										>
											{item.category}
										</span>
									</td>
									<td class="p-md font-body-md text-on-surface text-sm whitespace-nowrap">
										{item.vendor || '—'}
									</td>
									<td
										class="p-md text-on-surface-variant font-body-md max-w-[220px] truncate text-sm"
										title={item.description}
									>
										{item.description}
									</td>
									<td class="p-md font-body-md text-on-surface-variant text-sm whitespace-nowrap">
										{item.paymentMethod || 'Cash'}
									</td>
									<td
										class="p-md font-body-md text-right text-sm font-semibold text-stone-900 whitespace-nowrap"
									>
										₱{(item.amount || 0).toLocaleString('en-US', {
											minimumFractionDigits: 2,
											maximumFractionDigits: 2
										})}
									</td>
									<td class="p-md text-center whitespace-nowrap">
										{#if authState.isOwner}
											<button
												onclick={() => removeExpense(item.id)}
												class="hover:bg-red-50 inline-flex cursor-pointer items-center justify-center rounded-full p-1.5 text-red-600 transition-all hover:text-red-800"
												aria-label="Delete expense"
												title="Delete expense"
											>
												<svg
													class="h-4.5 w-4.5"
													fill="none"
													stroke="currentColor"
													stroke-width="2"
													viewBox="0 0 24 24"
													xmlns="http://www.w3.org/2000/svg"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
													/>
												</svg>
											</button>
										{:else}
											<span class="inline-flex items-center gap-1 text-xs text-stone-400 font-medium" title="Only Owners can delete expenses">
												🔒 Read Only
											</span>
										{/if}
									</td>
								</tr>
							{/each}
						{/if}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</div>

{#if isModalOpen}
	<ExpenseModal
		expenseNoSuggestion={expenseNoSuggestion}
		onClose={closeModal}
		onAddExpense={handleAddExpense}
	/>
{/if}
