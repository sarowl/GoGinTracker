<script>
	import { sharedStore } from '$lib/store.svelte';

	// Derived statistics using sharedStore.sales
	let totalSalesValue = $derived(sharedStore.sales.reduce((sum, item) => sum + item.amount, 0));

	let uniqueCustomers = $derived(new Set(sharedStore.sales.map((item) => item.customer)).size);
</script>

<svelte:head>
	<title>Sales — Ledger</title>
	<meta name="description" content="View and audit business sales invoices." />
</svelte:head>

<div class="relative min-h-screen overflow-hidden bg-background p-xl text-on-surface">
	<!-- Subtle Ambient Background Element with Glow Blobs -->
	<div class="pointer-events-none absolute inset-0 z-0 overflow-hidden opacity-40">
		<div class="glow-blob -top-24 -left-24"></div>
		<div class="glow-blob-2 right-0 bottom-0"></div>
	</div>

	<div class="max-w-container-max relative z-10 mx-auto space-y-md">
		<!-- Header Title -->
		<header class="mb-xl">
			<h1 class="font-headline-lg text-headline-lg tracking-tight text-primary">Sales Ledger</h1>
			<p class="font-body-md text-body-md text-on-surface-variant">
				View all completed contract revenues transferred automatically from delivered orders.
			</p>
		</header>

		<!-- Metrics Cards Grid -->
		<div class="grid grid-cols-1 gap-md md:grid-cols-3">
			<div class="glass-card space-y-sm rounded-xl p-lg">
				<p class="font-label-md text-label-md tracking-wider text-on-surface-variant uppercase">
					Total Sales Revenue
				</p>
				<p class="font-headline-lg text-headline-lg tracking-tight text-primary">
					₱{totalSalesValue.toLocaleString('en-US', {
						minimumFractionDigits: 2,
						maximumFractionDigits: 2
					})}
				</p>
				<span class="inline-flex items-center text-xs font-semibold text-stone-500"
					>Delivered contracts payout</span
				>
			</div>
			<div class="glass-card space-y-sm rounded-xl p-lg">
				<p class="font-label-md text-label-md tracking-wider text-on-surface-variant uppercase">
					Invoiced Customers
				</p>
				<p class="font-headline-lg text-headline-lg tracking-tight text-primary">
					{uniqueCustomers}
				</p>
				<span class="inline-flex items-center text-xs font-semibold text-stone-500"
					>Unique entities invoiced</span
				>
			</div>
			<div class="glass-card space-y-sm rounded-xl p-lg">
				<p class="font-label-md text-label-md tracking-wider text-on-surface-variant uppercase">
					Average Invoice Value
				</p>
				<p class="font-headline-lg text-headline-lg tracking-tight text-primary">
					₱{sharedStore.sales.length
						? (totalSalesValue / sharedStore.sales.length).toLocaleString('en-US', {
								minimumFractionDigits: 2,
								maximumFractionDigits: 2
							})
						: '0.00'}
				</p>
				<span class="inline-flex items-center text-xs font-semibold text-stone-500"
					>Value per completed sale</span
				>
			</div>
		</div>

		<!-- Sales Table Container -->
		<div class="glass-card space-y-md rounded-xl p-lg">
			<div class="overflow-x-auto rounded-lg border border-outline-variant">
				<table class="w-full border-collapse text-left">
					<thead>
						<tr class="border-b border-outline-variant bg-surface-dim">
							<th
								class="p-md min-w-[120px] font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Invoice ID</th
							>
							<th
								class="p-md min-w-[140px] font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Ref PO</th
							>
							<th
								class="p-md min-w-[180px] font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Customer</th
							>
							<th
								class="p-md min-w-[240px] font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Description</th
							>
							<th
								class="p-md min-w-[140px] font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Invoice Date</th
							>
							<th
								class="p-md min-w-[140px] text-right font-label-md text-label-md tracking-wider whitespace-nowrap text-on-surface-variant uppercase"
								>Revenue</th
							>
						</tr>
					</thead>
					<tbody class="divide-y divide-outline-variant bg-surface/50">
						{#if sharedStore.sales.length === 0}
							<tr>
								<td
									colspan="6"
									class="p-xl text-center font-body-md whitespace-nowrap text-on-surface-variant"
								>
									No sales logged yet. Deliver an active order in the Registry to log your first
									sale.
								</td>
							</tr>
						{:else}
							{#each sharedStore.sales as item (item.id)}
								<tr class="transition-colors hover:bg-surface-dim/40">
									<td class="p-md font-body-md text-sm font-semibold whitespace-nowrap text-primary"
										>{item.id}</td
									>
									<td class="p-md font-body-md text-sm whitespace-nowrap text-stone-600"
										>#{item.refPo}</td
									>
									<td class="p-md font-body-md text-sm whitespace-nowrap text-on-surface"
										>{item.customer}</td
									>
									<td
										class="p-md max-w-[240px] truncate font-body-md text-sm text-on-surface-variant"
										title={item.description}>{item.description}</td
									>
									<td class="p-md font-body-md text-sm whitespace-nowrap text-on-surface-variant"
										>{item.date}</td
									>
									<td
										class="p-md text-right font-body-md text-sm font-semibold whitespace-nowrap text-stone-900"
									>
										₱{item.amount.toLocaleString('en-US', {
											minimumFractionDigits: 2,
											maximumFractionDigits: 2
										})}
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
