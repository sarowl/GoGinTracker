<script>
	import OrderModal from '$lib/components/OrderModal.svelte';
	import { sharedStore } from '$lib/store.svelte';
	import { authState } from '$lib/auth.svelte';

	// Modal display state
	let isModalOpen = $state(false);

	// Derived statistics
	let totalOrdersValue = $derived(sharedStore.orders.reduce((sum, order) => sum + order.totalCost, 0));
	let activeCount = $derived(sharedStore.orders.filter((o) => o.status !== 'Cancelled').length);
	let poSuggestion = $derived(20260000 + sharedStore.orders.length + 1);

	function openAddModal() {
		isModalOpen = true;
	}

	function closeAddModal() {
		isModalOpen = false;
	}

	/**
	 * @param {any} newOrder
	 */
	function handleAddOrder(newOrder) {
		const userUid = authState.user?.uid || '';
		sharedStore.addOrder({
			...newOrder,
			insertedBy: userUid
		});
		closeAddModal();
	}

	/**
	 * @param {string} id
	 */
	function removeOrder(id) {
		sharedStore.removeOrder(id);
	}

	/**
	 * @param {any} order
	 * @param {Event & { currentTarget: HTMLSelectElement }} event
	 */
	async function handleStatusChange(order, event) {
		const newStatus = event.currentTarget.value;
		const oldStatus = order.status;
		const userUid = authState.user?.uid || '';

		if (newStatus === 'Delivered') {
			const invoiceId = prompt(`Please enter the Sales Invoice ID for PO #${order.poNumber}:`);
			
			// User clicked cancel
			if (invoiceId === null) {
				event.currentTarget.value = oldStatus;
				return;
			}

			const trimmedId = invoiceId.trim();
			if (trimmedId === '') {
				alert("Sales Invoice ID is required to mark the order as Delivered.");
				event.currentTarget.value = oldStatus;
				return;
			}

			await sharedStore.updateOrderStatus(order.id, 'Delivered', trimmedId, userUid);
		} else {
			await sharedStore.updateOrderStatus(order.id, newStatus, undefined, userUid);
		}
	}
</script>

<svelte:head>
	<title>Orders — Ledger</title>
	<meta name="description" content="Manage your purchasing and sales orders." />
</svelte:head>

<div class="relative min-h-screen overflow-hidden bg-background p-xl text-on-surface">
	<!-- Subtle Ambient Background Element with Glow Blobs -->
	<div class="pointer-events-none absolute inset-0 z-0 overflow-hidden opacity-40">
		<div class="glow-blob -top-24 -left-24"></div>
		<div class="glow-blob-2 right-0 bottom-0"></div>
	</div>

	<div class="max-w-container-max relative z-10 mx-auto space-y-md">
		<!-- Header Title -->
		<header class="mb-xl flex flex-col gap-md md:flex-row md:items-center md:justify-between">
			<div>
				<h1 class="font-headline-lg text-headline-lg tracking-tight text-primary">
					Orders Registry
				</h1>
				<p class="font-body-md text-body-md text-on-surface-variant">
					Track, add, and audit purchase orders and shipment statuses.
				</p>
			</div>
			<div>
				<button
					onclick={openAddModal}
					class="px-lg inline-flex h-10 cursor-pointer items-center justify-center gap-2 rounded bg-primary font-body-md text-sm font-semibold text-on-primary shadow-md transition-all hover:brightness-110 active:scale-[0.98]"
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
					Add Orders
				</button>
			</div>
		</header>

		<!-- Metrics Cards Grid -->
		<div class="grid grid-cols-1 gap-md md:grid-cols-3">
			<div class="glass-card space-y-sm rounded-xl p-lg">
				<p class="font-label-md text-label-md tracking-wider text-on-surface-variant uppercase">
					Active POs
				</p>
				<p class="font-headline-lg text-headline-lg tracking-tight text-primary">{activeCount}</p>
				<span class="inline-flex items-center text-xs font-semibold text-stone-500"
					>Uncancelled orders</span
				>
			</div>
			<div class="glass-card space-y-sm rounded-xl p-lg">
				<p class="font-label-md text-label-md tracking-wider text-on-surface-variant uppercase">
					Total Contract Value
				</p>
				<p class="font-headline-lg text-headline-lg tracking-tight text-primary">
					₱{totalOrdersValue.toLocaleString('en-US', {
						minimumFractionDigits: 2,
						maximumFractionDigits: 2
					})}
				</p>
				<span class="inline-flex items-center text-xs font-semibold text-stone-500"
					>Sum of all total costs</span
				>
			</div>
			<div class="glass-card space-y-sm rounded-xl p-lg">
				<p class="font-label-md text-label-md tracking-wider text-on-surface-variant uppercase">
					Average Order Value
				</p>
				<p class="font-headline-lg text-headline-lg tracking-tight text-primary">
					₱{sharedStore.orders.length
						? (totalOrdersValue / sharedStore.orders.length).toLocaleString('en-US', {
								minimumFractionDigits: 2,
								maximumFractionDigits: 2
							})
						: '0.00'}
				</p>
				<span class="inline-flex items-center text-xs font-semibold text-stone-500"
					>Value per registered PO</span
				>
			</div>
		</div>

		<!-- Orders Data Table Container -->
		<div class="glass-card space-y-md rounded-xl p-lg">
			<div class="overflow-x-auto rounded-lg border border-outline-variant">
				<table class="w-full border-collapse text-left">
					<thead>
						<tr class="border-b border-outline-variant bg-surface-dim">
							<th
								class="p-md font-label-md text-label-md tracking-wider text-on-surface-variant uppercase whitespace-nowrap min-w-25"
								>PO Number</th
							>
							<th
								class="p-md font-label-md text-label-md tracking-wider text-on-surface-variant uppercase whitespace-nowrap min-w-[180px]"
								>Ship To</th
							>
							<th
								class="p-md font-label-md text-label-md tracking-wider text-on-surface-variant uppercase whitespace-nowrap min-w-[240px]"
								>Description</th
							>
							<th
								class="p-md text-right font-label-md text-label-md tracking-wider text-on-surface-variant uppercase whitespace-nowrap min-w-[100px]"
								>Quantity</th
							>
							<th
								class="p-md text-right font-label-md text-label-md tracking-wider text-on-surface-variant uppercase whitespace-nowrap min-w-[120px]"
								>Unit Cost</th
							>
							<th
								class="p-md pr-16 text-right font-label-md text-label-md tracking-wider text-on-surface-variant uppercase whitespace-nowrap min-w-[260px]"
								>Total Cost</th
							>
							<th
								class="p-md pl-8 font-label-md text-label-md tracking-wider text-on-surface-variant uppercase whitespace-nowrap min-w-[200px]"
								>Cancel Date</th
							>
							<th
								class="p-md font-label-md text-label-md tracking-wider text-on-surface-variant uppercase whitespace-nowrap min-w-[120px]"
								>Status</th
							>
							<th
								class="p-md text-center font-label-md text-label-md tracking-wider text-on-surface-variant uppercase whitespace-nowrap min-w-[100px]"
								>Actions</th
							>
						</tr>
					</thead>
					<tbody class="divide-y divide-outline-variant bg-surface/50">
						{#if sharedStore.orders.length === 0}
							<tr>
								<td colspan="9" class="p-xl text-center font-body-md text-on-surface-variant whitespace-nowrap">
									No orders registered yet. Click "Add Orders" above to log your first transaction.
								</td>
							</tr>
						{:else}
							{#each sharedStore.orders as order (order.id)}
								<tr class="transition-colors hover:bg-surface-dim/40">
									<td class="p-md font-body-md text-sm font-semibold text-primary whitespace-nowrap"
										>{order.poNumber}</td
									>
									<td class="p-md font-body-md text-sm text-on-surface whitespace-nowrap">{order.shipTo}</td>
									<td
										class="p-md max-w-[240px] truncate font-body-md text-sm text-on-surface-variant"
										title={order.description}>{order.description}</td
									>
									<td class="p-md text-right font-body-md text-sm whitespace-nowrap">{order.quantity}</td>
									<td class="p-md text-right font-body-md text-sm whitespace-nowrap">₱{order.cost.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
									<td class="p-md pr-16 text-right font-body-md text-sm font-semibold whitespace-nowrap"
										>₱{order.totalCost.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td
									>
									<td class="p-md pl-8 font-body-md text-sm text-on-surface-variant whitespace-nowrap"
										>{order.cancelDate}</td
									>
									<td class="p-md font-body-md text-sm whitespace-nowrap">
										{#if order.cancelDate && order.cancelDate !== 'N/A' && order.cancelDate < new Date().toISOString().split('T')[0]}
											<span class="inline-flex items-center rounded-full border border-red-100 bg-red-50 px-2.5 py-1 text-xs font-semibold text-red-800">
												Cancelled (Expired)
											</span>
										{:else}
											<select
												value={order.status}
												onchange={(e) => handleStatusChange(order, e)}
												class="cursor-pointer rounded border border-outline-variant bg-surface-container-low/40 px-2.5 py-1 font-body-md text-xs font-semibold focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary/30 transition-all {
													order.status === 'Pending' 
														? 'text-stone-600 border-stone-200' 
														: 'text-red-700 bg-red-50/20 border-red-200'
												}"
											>
												<option value="Pending">Pending</option>
												<option value="Cancelled">Cancelled</option>
												<option value="Delivered">Delivered</option>
											</select>
										{/if}
									</td>
									<td class="p-md text-center whitespace-nowrap">
										<button
											onclick={() => removeOrder(order.id)}
											class="cursor-pointer inline-flex items-center justify-center p-1.5 rounded-full hover:bg-red-50 text-red-600 hover:text-red-800 transition-all"
											aria-label="Delete order"
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
	<OrderModal
		poSuggestion={poSuggestion}
		onClose={closeAddModal}
		onAddOrder={handleAddOrder}
	/>
{/if}