<script>
	import OrderModal from '$lib/components/OrderModal.svelte';
	import { sharedStore } from '$lib/store.svelte';
	import { authState } from '$lib/auth.svelte';

	// Modal display state
	let isModalOpen = $state(false);

	// Derived statistics
	let ordersList = $derived(
		sharedStore.orders.filter((o) => o.status !== 'Delivered')
	);
	let totalOrdersValue = $derived(
		ordersList.reduce((sum, order) => sum + order.totalCost, 0)
	);
	let activeCount = $derived(ordersList.filter((o) => o.status === 'Pending').length);
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
		if (!authState.isOwner) {
			alert('Employees are not permitted to delete orders.');
			return;
		}
		sharedStore.removeOrder(id, authState.role);
	}

	/**
	 * @param {any} order
	 * @param {HTMLSelectElement} [selectEl]
	 */
	async function triggerDeliver(order, selectEl) {
		const userUid = authState.user?.uid || '';
		const invoiceId = prompt(`Please enter the Sales Invoice ID for PO #${order.poNumber}:`);

		if (invoiceId === null) {
			if (selectEl) selectEl.value = order.status;
			return;
		}

		const trimmedId = invoiceId.trim();
		if (trimmedId === '') {
			alert('Sales Invoice ID is required to mark the order as Delivered.');
			if (selectEl) selectEl.value = order.status;
			return;
		}

		await sharedStore.updateOrderStatus(order.id, 'Delivered', trimmedId, userUid, authState.role);
	}

	/**
	 * @param {any} order
	 * @param {Event & { currentTarget: HTMLSelectElement }} event
	 */
	async function handleStatusChange(order, event) {
		const newStatus = event.currentTarget.value;
		const userUid = authState.user?.uid || '';

		if (newStatus === 'Delivered') {
			await triggerDeliver(order, event.currentTarget);
		} else {
			await sharedStore.updateOrderStatus(order.id, newStatus, undefined, userUid, authState.role);
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
					₱{ordersList.length
						? (totalOrdersValue / ordersList.length).toLocaleString('en-US', {
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
								class="p-md min-w-25 font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>PO Number</th
							>
							<th
								class="p-md min-w-[180px] font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Ship To</th
							>
							<th
								class="p-md min-w-[240px] font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Description</th
							>
							<th
								class="p-md min-w-[100px] text-right font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Quantity</th
							>
							<th
								class="p-md min-w-[120px] text-right font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Unit Cost</th
							>
							<th
								class="p-md min-w-[260px] pr-16 text-right font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Total Cost</th
							>
							<th
								class="p-md min-w-[200px] pl-8 font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Cancel Date</th
							>
							<th
								class="p-md min-w-[120px] font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Status</th
							>
							<th
								class="p-md min-w-[100px] text-center font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Actions</th
							>
						</tr>
					</thead>
					<tbody class="divide-y divide-outline-variant bg-surface/50">
						{#if ordersList.length === 0}
							<tr>
								<td
									colspan="9"
									class="p-xl text-center font-body-md whitespace-nowrap text-on-surface-variant"
								>
									No active or pending orders. Click "Add Orders" above to log your first transaction.
								</td>
							</tr>
						{:else}
							{#each ordersList as order (order.id)}
								<tr class="transition-colors hover:bg-surface-dim/40">
									<td class="p-md font-body-md text-sm font-semibold whitespace-nowrap text-primary"
										>{order.poNumber}</td
									>
									<td class="p-md font-body-md text-sm whitespace-nowrap text-on-surface"
										>{order.shipTo}</td
									>
									<td
										class="p-md max-w-[240px] truncate font-body-md text-sm text-on-surface-variant"
										title={order.description}>{order.description}</td
									>
									<td class="p-md text-right font-body-md text-sm whitespace-nowrap"
										>{order.quantity}</td
									>
									<td class="p-md text-right font-body-md text-sm whitespace-nowrap"
										>₱{order.cost.toLocaleString('en-US', {
											minimumFractionDigits: 2,
											maximumFractionDigits: 2
										})}</td
									>
									<td
										class="p-md pr-16 text-right font-body-md text-sm font-semibold whitespace-nowrap"
										>₱{order.totalCost.toLocaleString('en-US', {
											minimumFractionDigits: 2,
											maximumFractionDigits: 2
										})}</td
									>
									<td
										class="p-md pl-8 font-body-md text-sm whitespace-nowrap text-on-surface-variant"
										>{order.cancelDate}</td
									>
									<td class="p-md font-body-md text-sm whitespace-nowrap">
										<select
											value={order.cancelDate &&
											order.cancelDate !== 'N/A' &&
											order.cancelDate < new Date().toISOString().split('T')[0] &&
											order.status === 'Pending'
												? 'Expired'
												: order.status}
											onchange={(e) => handleStatusChange(order, e)}
											class="cursor-pointer rounded border border-outline-variant bg-surface-container-low/40 px-2.5 py-1 font-body-md text-xs font-semibold transition-all focus:border-primary focus:ring-1 focus:ring-primary/30 focus:outline-none {order.status ===
												'Expired' ||
											(order.cancelDate &&
												order.cancelDate !== 'N/A' &&
												order.cancelDate < new Date().toISOString().split('T')[0])
												? 'border-red-200 bg-red-50/20 text-red-700'
												: 'border-stone-200 text-stone-600'}"
										>
											<option value="Pending">Pending</option>
											<option value="Expired">Expired</option>
											<option value="Delivered">Delivered</option>
										</select>
									</td>
									<td class="p-md text-center whitespace-nowrap">
										{#if authState.isOwner}
											<button
												onclick={() => removeOrder(order.id)}
												class="inline-flex cursor-pointer items-center justify-center rounded-full p-1.5 text-red-600 transition-all hover:bg-red-50 hover:text-red-800"
												aria-label="Delete order"
												title="Delete order"
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
											<span
												class="inline-flex items-center gap-1 text-xs font-medium text-stone-400"
												title="Only Owners can delete orders"
											>
												
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
	<OrderModal {poSuggestion} onClose={closeAddModal} onAddOrder={handleAddOrder} />
{/if}
