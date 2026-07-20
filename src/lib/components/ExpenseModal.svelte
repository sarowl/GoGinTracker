<script>
	import { fade, scale } from 'svelte/transition';
	import { onMount } from 'svelte';

	/**
	 * @type {{
	 *   onClose?: () => void,
	 *   onAddExpense?: (expense: any) => void,
	 *   expenseNoSuggestion?: string
	 * }}
	 */
	let {
		onClose = () => {},
		onAddExpense = () => {},
		expenseNoSuggestion = 'EXP-1001'
	} = $props();

	// Form State
	let expenseNo = $state('');
	let category = $state('Utilities');
	let vendor = $state('');
	let vendorInvoiceNo = $state('');
	let description = $state('');
	let amount = $state(0);
	let paymentMethod = $state('Bank Transfer');
	let receiptUrl = $state('');
	let expenseDate = $state(new Date().toISOString().split('T')[0]);

	const categories = ['Utilities', 'Supplies', 'Payroll', 'Delivery/Logistics', 'Rent', 'Misc'];
	const paymentMethods = ['Cash', 'Bank Transfer', 'Card', 'Check'];

	onMount(() => {
		if (expenseNoSuggestion) {
			expenseNo = expenseNoSuggestion;
		}
	});

	/**
	 * @param {SubmitEvent} event
	 */
	function handleSubmit(event) {
		event.preventDefault();

		if (Number(amount) <= 0) {
			alert('Please enter a valid expense amount.');
			return;
		}

		const newExpense = {
			expenseNo: expenseNo || expenseNoSuggestion,
			vendorInvoiceNo,
			category,
			vendor,
			description,
			amount: Number(amount) || 0,
			paymentMethod,
			receiptUrl,
			expenseDate
		};

		onAddExpense(newExpense);
	}
</script>

<div
	transition:fade={{ duration: 150 }}
	class="p-gutter fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm"
>
	<div
		transition:scale={{ duration: 150, start: 0.95 }}
		class="p-lg md:p-xl space-y-md border-outline-variant bg-surface text-on-surface relative w-full max-w-lg rounded-2xl border shadow-xl"
	>
		<header class="border-outline-variant flex items-center justify-between border-b pb-md">
			<h2 class="font-headline-md text-headline-md text-primary tracking-tight">
				Log Business Expense
			</h2>
			<button
				onclick={onClose}
				class="text-on-surface-variant hover:text-on-surface cursor-pointer p-1 transition-colors"
				aria-label="Close modal"
			>
				<svg class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>
		</header>

		<form onsubmit={handleSubmit} class="font-body-md space-y-md text-sm">
			<!-- Vendor Invoice # & Category -->
			<div class="grid grid-cols-2 gap-md">
				<div class="flex flex-col gap-xs">
					<label for="vendorInvoiceNo" class="font-label-md text-label-md text-on-surface-variant"
						>Vendor Invoice #</label
					>
					<input
						id="vendorInvoiceNo"
						type="text"
						bind:value={vendorInvoiceNo}
						placeholder="INV-99823"
						class="bg-surface-container-low/40 border-outline-variant focus:ring-primary/40 h-10 w-full rounded border px-md transition-all focus:ring-2 focus:outline-none"
					/>
				</div>
				<div class="flex flex-col gap-xs">
					<label for="category" class="font-label-md text-label-md text-on-surface-variant"
						>Category</label
					>
					<select
						id="category"
						bind:value={category}
						class="bg-surface-container-low/40 border-outline-variant focus:ring-primary/40 h-10 w-full rounded border px-md transition-all focus:ring-2 focus:outline-none"
					>
						{#each categories as cat}
							<option value={cat}>{cat}</option>
						{/each}
					</select>
				</div>
			</div>

			<!-- Vendor / Supplier -->
			<div class="flex flex-col gap-xs">
				<label for="vendor" class="font-label-md text-label-md text-on-surface-variant"
					>Vendor / Supplier</label
				>
				<input
					id="vendor"
					type="text"
					bind:value={vendor}
					placeholder="e.g. Meralco, Office Warehouse"
					required
					class="bg-surface-container-low/40 border-outline-variant focus:ring-primary/40 h-10 w-full rounded border px-md transition-all focus:ring-2 focus:outline-none"
				/>
			</div>

			<!-- Description -->
			<div class="flex flex-col gap-xs">
				<label for="description" class="font-label-md text-label-md text-on-surface-variant"
					>Description</label
				>
				<input
					id="description"
					type="text"
					bind:value={description}
					placeholder="Brief breakdown of line items or notes"
					required
					class="bg-surface-container-low/40 border-outline-variant focus:ring-primary/40 h-10 w-full rounded border px-md transition-all focus:ring-2 focus:outline-none"
				/>
			</div>

			<!-- Amount & Payment Method -->
			<div class="grid grid-cols-2 gap-md">
				<div class="flex flex-col gap-xs">
					<label for="amount" class="font-label-md text-label-md text-on-surface-variant"
						>Amount (₱)</label
					>
					<input
						id="amount"
						type="number"
						step="0.01"
						min="0.01"
						bind:value={amount}
						required
						class="bg-surface-container-low/40 border-outline-variant focus:ring-primary/40 h-10 w-full rounded border px-md transition-all focus:ring-2 focus:outline-none"
					/>
				</div>
				<div class="flex flex-col gap-xs">
					<label for="paymentMethod" class="font-label-md text-label-md text-on-surface-variant"
						>Payment Method</label
					>
					<select
						id="paymentMethod"
						bind:value={paymentMethod}
						class="bg-surface-container-low/40 border-outline-variant focus:ring-primary/40 h-10 w-full rounded border px-md transition-all focus:ring-2 focus:outline-none"
					>
						{#each paymentMethods as pm}
							<option value={pm}>{pm}</option>
						{/each}
					</select>
				</div>
			</div>

			<!-- Expense Date & Receipt URL -->
			<div class="grid grid-cols-2 gap-md">
				<div class="flex flex-col gap-xs">
					<label for="expenseDate" class="font-label-md text-label-md text-on-surface-variant"
						>Expense Date</label
					>
					<input
						id="expenseDate"
						type="date"
						bind:value={expenseDate}
						required
						class="bg-surface-container-low/40 border-outline-variant focus:ring-primary/40 h-10 w-full rounded border px-md transition-all focus:ring-2 focus:outline-none"
					/>
				</div>
				<div class="flex flex-col gap-xs">
					<label for="receiptUrl" class="font-label-md text-label-md text-on-surface-variant"
						>Receipt Attachment URL</label
					>
					<input
						id="receiptUrl"
						type="url"
						bind:value={receiptUrl}
						placeholder="https://..."
						class="bg-surface-container-low/40 border-outline-variant focus:ring-primary/40 h-10 w-full rounded border px-md transition-all focus:ring-2 focus:outline-none"
					/>
				</div>
			</div>

			<!-- Form Actions -->
			<footer class="border-outline-variant flex justify-end gap-sm border-t pt-md">
				<button
					type="button"
					onclick={onClose}
					class="border-outline-variant hover:bg-surface-dim h-10 cursor-pointer rounded border px-md transition-colors"
				>
					Cancel
				</button>
				<button
					type="submit"
					class="bg-primary text-on-primary h-10 cursor-pointer rounded px-lg font-semibold transition-all hover:brightness-110 active:scale-[0.98]"
				>
					Record Expense
				</button>
			</footer>
		</form>
	</div>
</div>
