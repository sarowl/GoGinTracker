<script>
	import { fade, scale } from 'svelte/transition';
	import { onMount } from 'svelte';

	/**
	 * @type {{
	 *   onClose?: () => void,
	 *   onAddOrder?: (order: any) => void,
	 *   poSuggestion?: number
	 * }}
	 */
	let { onClose = () => {}, onAddOrder = () => {}, poSuggestion = 0 } = $props();

	// Local Form Binds
	let newPo = $state('');
	let newShipTo = $state('');
	let newDesc = $state('');
	let newQuantity = $state(1);
	let newCost = $state(0);
	let newCancelDate = $state('');
	let newStatus = $state('Pending');

	// Initialize once on mount to avoid warning and prevent inputs from being overridden on keystrokes
	onMount(() => {
		if (poSuggestion) {
			newPo = String(poSuggestion);
		}
	});

	/**
	 * @param {SubmitEvent} event
	 */
	function handleSubmit(event) {
		event.preventDefault();

		// Validate PO number is exactly 8 digits
		if (!/^[0-9]{8}$/.test(newPo)) {
			alert('PO number must consist of exactly 8 digits.');
			return;
		}

		const quantity = Number(newQuantity) || 0;
		const cost = Number(newCost) || 0;

		const newOrder = {
			poNumber: Number(newPo),
			customer: newShipTo,
			shipTo: newShipTo,
			description: newDesc,
			quantity,
			cost,
			totalCost: quantity * cost,
			cancelDate: newCancelDate || 'N/A',
			status: newStatus
		};

		onAddOrder(newOrder);
	}

	/**
	 * @param {any} e
	 */
	function filterNumericInput(e) {
		const target = /** @type {HTMLInputElement} */ (e.currentTarget);
		target.value = target.value.replace(/[^0-9]/g, '');
		newPo = target.value;
	}
</script>

<div
	transition:fade={{ duration: 150 }}
	class="p-gutter fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm"
>
	<div
		transition:scale={{ duration: 150, start: 0.95 }}
		class="relative space-y-md w-full max-w-lg rounded-2xl border border-outline-variant bg-surface p-lg text-on-surface shadow-xl md:p-xl"
	>
		<header class="pb-md flex items-center justify-between border-b border-outline-variant">
			<h2 class="font-headline-md text-headline-md tracking-tight text-primary">
				Create Purchase Order
			</h2>
			<button
				onclick={onClose}
				class="cursor-pointer p-1 text-on-surface-variant transition-colors hover:text-on-surface"
				aria-label="Close modal"
			>
				<svg class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>
		</header>

		<form onsubmit={handleSubmit} class="space-y-md font-body-md text-sm">
			<!-- PO Number & Ship To -->
			<div class="grid grid-cols-2 gap-md">
				<div class="flex flex-col gap-xs">
					<label for="poNumber" class="font-label-md text-label-md text-on-surface-variant"
						>PO Number</label
					>
					<input
						id="poNumber"
						type="text"
						bind:value={newPo}
						oninput={filterNumericInput}
						minlength="8"
						maxlength="8"
						pattern={'[0-9]{8}'}
						title="PO number must be exactly 8 digits"
						required
						class="px-md h-10 w-full rounded border border-outline-variant bg-surface-container-low/40 font-body-md transition-all focus:ring-2 focus:ring-primary/40 focus:outline-none"
					/>
				</div>
				<div class="flex flex-col gap-xs">
					<label for="shipTo" class="font-label-md text-label-md text-on-surface-variant"
						>Ship To</label
					>
					<input
						id="shipTo"
						type="text"
						bind:value={newShipTo}
						placeholder="Client Name"
						required
						class="px-md h-10 w-full rounded border border-outline-variant bg-surface-container-low/40 font-body-md transition-all focus:ring-2 focus:ring-primary/40 focus:outline-none"
					/>
				</div>
			</div>

			<!-- Description -->
			<div class="flex flex-col gap-xs">
				<label for="description" class="font-label-md text-label-md text-on-surface-variant"
					>Description</label
				>
				<input
					id="description"
					type="text"
					bind:value={newDesc}
					placeholder="Items description"
					required
					class="px-md h-10 w-full rounded border bg-surface-container-low/40 font-body-md transition-all focus:ring-2 focus:ring-primary/40 focus:outline-none"
				/>
			</div>

			<!-- Quantity & Unit Cost -->
			<div class="grid grid-cols-2 gap-md">
				<div class="flex flex-col gap-xs">
					<label for="quantity" class="font-label-md text-label-md text-on-surface-variant"
						>Quantity</label
					>
					<input
						id="quantity"
						type="number"
						min="1"
						bind:value={newQuantity}
						required
						class="px-md h-10 w-full rounded border bg-surface-container-low/40 font-body-md transition-all focus:border-primary focus:ring-2 focus:ring-primary/40 focus:outline-none"
					/>
				</div>
				<div class="flex flex-col gap-xs">
					<label for="cost" class="font-label-md text-label-md text-on-surface-variant"
						>Unit Cost (₱)</label
					>
					<input
						id="cost"
						type="number"
						step="0.01"
						min="0"
						bind:value={newCost}
						required
						class="px-md h-10 w-full rounded border bg-surface-container-low/40 font-body-md transition-all focus:border-primary focus:ring-2 focus:ring-primary/40 focus:outline-none"
					/>
				</div>
			</div>

			<!-- Cancel Date & Status -->
			<div class="grid grid-cols-2 gap-md">
				<div class="flex flex-col gap-xs">
					<label for="cancelDate" class="font-label-md text-label-md text-on-surface-variant"
						>Cancel Date</label
					>
					<input
						id="cancelDate"
						type="date"
						bind:value={newCancelDate}
						required
						class="px-md h-10 w-full rounded border bg-surface-container-low/40 font-body-md transition-all focus:border-primary focus:ring-2 focus:ring-primary/40 focus:outline-none"
					/>
				</div>
				<div class="flex flex-col gap-xs">
					<label for="status" class="font-label-md text-label-md text-on-surface-variant"
						>Status</label
					>
					<select
						id="status"
						bind:value={newStatus}
						class="px-md h-10 w-full rounded border bg-surface-container-low/40 font-body-md transition-all focus:border-primary focus:ring-2 focus:ring-primary/40 focus:outline-none"
					>
						<option value="Pending">Pending</option>
						<option value="Expired">Expired</option>
					</select>
				</div>
			</div>

			<!-- Calculated Live Preview -->
			<div
				class="p-md flex items-center justify-between rounded border border-outline-variant bg-surface-dim"
			>
				<span class="font-label-md text-on-surface-variant uppercase">Estimated Total:</span>
				<span class="font-headline-md font-bold text-primary"
					>₱{(newQuantity * newCost).toLocaleString('en-US', { minimumFractionDigits: 2 })}</span
				>
			</div>

			<!-- Submit actions -->
			<footer class="pt-md flex justify-end gap-sm border-t border-outline-variant">
				<button
					type="button"
					onclick={onClose}
					class="px-md h-10 cursor-pointer rounded border border-outline-variant transition-colors hover:bg-surface-dim"
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-lg h-10 cursor-pointer rounded bg-primary font-semibold text-on-primary transition-all hover:brightness-110 active:scale-[0.98]"
				>
					Create Order
				</button>
			</footer>
		</form>
	</div>
</div>
