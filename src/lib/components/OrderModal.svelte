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
  let { 
    onClose = () => {}, 
    onAddOrder = () => {},
    poSuggestion = 0
  } = $props();

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
      alert("PO number must consist of exactly 8 digits.");
      return;
    }
    
    const quantity = Number(newQuantity) || 0;
    const cost = Number(newCost) || 0;
    
    const newOrder = {
      poNumber: Number(newPo),
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
  class="fixed inset-0 bg-black/60 backdrop-blur-sm z-[100] flex items-center justify-center p-gutter"
>
  <div 
    transition:scale={{ duration: 150, start: 0.95 }}
    class="bg-surface border border-outline-variant w-full max-w-lg rounded-2xl shadow-xl p-lg md:p-xl space-y-md relative text-on-surface"
  >
    <header class="flex justify-between items-center pb-md border-b border-outline-variant">
      <h2 class="font-headline-md text-headline-md text-primary tracking-tight">Create Purchase Order</h2>
      <button 
        onclick={onClose}
        class="text-on-surface-variant hover:text-on-surface transition-colors cursor-pointer p-1"
        aria-label="Close modal"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </header>

    <form onsubmit={handleSubmit} class="space-y-md font-body-md text-sm">
      <!-- PO Number & Ship To -->
      <div class="grid grid-cols-2 gap-md">
        <div class="flex flex-col gap-xs">
          <label for="poNumber" class="font-label-md text-label-md text-on-surface-variant">PO Number</label>
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
            class="w-full h-10 px-md bg-surface-container-low/40 border border-outline-variant rounded focus:outline-none focus:ring-2 focus:ring-primary/40 transition-all font-body-md"
          />
        </div>
        <div class="flex flex-col gap-xs">
          <label for="shipTo" class="font-label-md text-label-md text-on-surface-variant">Ship To</label>
          <input 
            id="shipTo"
            type="text" 
            bind:value={newShipTo} 
            placeholder="Client Name"
            required
            class="w-full h-10 px-md bg-surface-container-low/40 border border-outline-variant rounded focus:outline-none focus:ring-2 focus:ring-primary/40 transition-all font-body-md"
          />
        </div>
      </div>

      <!-- Description -->
      <div class="flex flex-col gap-xs">
        <label for="description" class="font-label-md text-label-md text-on-surface-variant">Description</label>
        <input 
          id="description"
          type="text" 
          bind:value={newDesc} 
          placeholder="Items description"
          required
          class="w-full h-10 px-md bg-surface-container-low/40 border rounded focus:outline-none focus:ring-2 focus:ring-primary/40 transition-all font-body-md"
        />
      </div>

      <!-- Quantity & Unit Cost -->
      <div class="grid grid-cols-2 gap-md">
        <div class="flex flex-col gap-xs">
          <label for="quantity" class="font-label-md text-label-md text-on-surface-variant">Quantity</label>
          <input 
            id="quantity"
            type="number" 
            min="1"
            bind:value={newQuantity} 
            required
            class="w-full h-10 px-md bg-surface-container-low/40 border rounded focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary transition-all font-body-md"
          />
        </div>
        <div class="flex flex-col gap-xs">
          <label for="cost" class="font-label-md text-label-md text-on-surface-variant">Unit Cost (₱)</label>
          <input 
            id="cost"
            type="number" 
            step="0.01"
            min="0"
            bind:value={newCost} 
            required
            class="w-full h-10 px-md bg-surface-container-low/40 border rounded focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary transition-all font-body-md"
          />
        </div>
      </div>

      <!-- Cancel Date & Status -->
      <div class="grid grid-cols-2 gap-md">
        <div class="flex flex-col gap-xs">
          <label for="cancelDate" class="font-label-md text-label-md text-on-surface-variant">Cancel Date</label>
          <input 
            id="cancelDate"
            type="date" 
            bind:value={newCancelDate} 
            required
            class="w-full h-10 px-md bg-surface-container-low/40 border rounded focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary transition-all font-body-md"
          />
        </div>
        <div class="flex flex-col gap-xs">
          <label for="status" class="font-label-md text-label-md text-on-surface-variant">Status</label>
          <select 
            id="status"
            bind:value={newStatus} 
            class="w-full h-10 px-md bg-surface-container-low/40 border rounded focus:outline-none focus:ring-2 focus:ring-primary/40 focus:border-primary transition-all font-body-md"
          >
            <option value="Pending">Pending</option>
            <option value="Cancelled">Cancelled</option>
          </select>
        </div>
      </div>

      <!-- Calculated Live Preview -->
      <div class="bg-surface-dim p-md rounded border border-outline-variant flex justify-between items-center">
        <span class="font-label-md text-on-surface-variant uppercase">Estimated Total:</span>
        <span class="font-headline-md text-primary font-bold">₱{(newQuantity * newCost).toLocaleString('en-US', { minimumFractionDigits: 2 })}</span>
      </div>

      <!-- Submit actions -->
      <footer class="flex justify-end gap-sm pt-md border-t border-outline-variant">
        <button 
          type="button" 
          onclick={onClose}
          class="h-10 px-md rounded border border-outline-variant hover:bg-surface-dim transition-colors cursor-pointer"
        >
          Cancel
        </button>
        <button 
          type="submit"
          class="h-10 px-lg rounded bg-primary text-on-primary font-semibold hover:brightness-110 active:scale-[0.98] transition-all cursor-pointer"
        >
          Create Order
        </button>
      </footer>
    </form>
  </div>
</div>