<script>
  import { sharedStore } from '$lib/store.svelte';

  // Derived statistics using sharedStore.sales
  let totalSalesValue = $derived(
    sharedStore.sales.reduce((sum, item) => sum + item.amount, 0)
  );

  let uniqueCustomers = $derived(
    new Set(sharedStore.sales.map(item => item.customer)).size
  );
</script>

<svelte:head>
  <title>Sales — Ledger</title>
  <meta name="description" content="View and audit business sales invoices." />
</svelte:head>

<div class="min-h-screen bg-background text-on-surface p-xl relative overflow-hidden">
  <!-- Subtle Ambient Background Element with Glow Blobs -->
  <div class="absolute inset-0 z-0 opacity-40 pointer-events-none overflow-hidden">
    <div class="glow-blob -top-24 -left-24"></div>
    <div class="glow-blob-2 bottom-0 right-0"></div>
  </div>

  <div class="max-w-container-max mx-auto relative z-10 space-y-md">
    <!-- Header Title -->
    <header class="mb-xl">
      <h1 class="font-headline-lg text-headline-lg text-primary tracking-tight">Sales Ledger</h1>
      <p class="font-body-md text-body-md text-on-surface-variant">View all completed contract revenues transferred automatically from delivered orders.</p>
    </header>

    <!-- Metrics Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-md">
      <div class="glass-card rounded-xl p-lg space-y-sm">
        <p class="font-label-md text-label-md text-on-surface-variant uppercase tracking-wider">Total Sales Revenue</p>
        <p class="font-headline-lg text-headline-lg text-primary tracking-tight">
          ₱{totalSalesValue.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
        </p>
        <span class="inline-flex items-center text-xs font-semibold text-stone-500">Delivered contracts payout</span>
      </div>
      <div class="glass-card rounded-xl p-lg space-y-sm">
        <p class="font-label-md text-label-md text-on-surface-variant uppercase tracking-wider">Invoiced Customers</p>
        <p class="font-headline-lg text-headline-lg text-primary tracking-tight">{uniqueCustomers}</p>
        <span class="inline-flex items-center text-xs font-semibold text-stone-500">Unique entities invoiced</span>
      </div>
      <div class="glass-card rounded-xl p-lg space-y-sm">
        <p class="font-label-md text-label-md text-on-surface-variant uppercase tracking-wider">Average Invoice Value</p>
        <p class="font-headline-lg text-headline-lg text-primary tracking-tight">
          ₱{sharedStore.sales.length 
            ? (totalSalesValue / sharedStore.sales.length).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) 
            : '0.00'}
        </p>
        <span class="inline-flex items-center text-xs font-semibold text-stone-500">Value per completed sale</span>
      </div>
    </div>

    <!-- Sales Table Container -->
    <div class="glass-card rounded-xl p-lg space-y-md">
      <div class="overflow-x-auto border border-outline-variant rounded-lg">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-surface-dim border-b border-outline-variant">
              <th class="p-md font-label-md text-label-md text-on-surface-variant uppercase tracking-wider whitespace-nowrap min-w-[120px]">Invoice ID</th>
              <th class="p-md font-label-md text-label-md text-on-surface-variant uppercase tracking-wider whitespace-nowrap min-w-[140px]">Ref PO</th>
              <th class="p-md font-label-md text-label-md text-on-surface-variant uppercase tracking-wider whitespace-nowrap min-w-[180px]">Customer</th>
              <th class="p-md font-label-md text-label-md text-on-surface-variant uppercase tracking-wider whitespace-nowrap min-w-[240px]">Description</th>
              <th class="p-md font-label-md text-label-md text-on-surface-variant uppercase tracking-wider whitespace-nowrap min-w-[140px]">Invoice Date</th>
              <th class="p-md font-label-md text-label-md text-on-surface-variant uppercase tracking-wider text-right whitespace-nowrap min-w-[140px]">Revenue</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant bg-surface/50">
            {#if sharedStore.sales.length === 0}
              <tr>
                <td colspan="6" class="p-xl text-center font-body-md text-on-surface-variant whitespace-nowrap">
                  No sales logged yet. Deliver an active order in the Registry to log your first sale.
                </td>
              </tr>
            {:else}
              {#each sharedStore.sales as item (item.id)}
                <tr class="hover:bg-surface-dim/40 transition-colors">
                  <td class="p-md font-body-md text-sm font-semibold text-primary whitespace-nowrap">{item.id}</td>
                  <td class="p-md font-body-md text-sm text-stone-600 whitespace-nowrap">#{item.refPo}</td>
                  <td class="p-md font-body-md text-sm text-on-surface whitespace-nowrap">{item.customer}</td>
                  <td class="p-md font-body-md text-sm text-on-surface-variant max-w-[240px] truncate" title={item.description}>{item.description}</td>
                  <td class="p-md font-body-md text-sm text-on-surface-variant whitespace-nowrap">{item.date}</td>
                  <td class="p-md font-body-md text-sm font-semibold text-right text-stone-900 whitespace-nowrap">
                    ₱{item.amount.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
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
