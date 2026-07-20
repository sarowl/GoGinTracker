<script>
	import { onMount } from 'svelte';
	import { sharedStore } from '$lib/store.svelte';
	import { authState } from '$lib/auth.svelte';

	// Time Horizon Filter: 'mtd' (Month-to-Date) | 'ytd' (Year-to-Date) | 'all' (All-Time)
	let timeHorizon = $state('ytd');

	onMount(async () => {
		await Promise.all([
			sharedStore.fetchSales(),
			sharedStore.fetchExpenses(),
			sharedStore.fetchOrders()
		]);
	});

	const currentYear = new Date().getFullYear();
	const currentMonth = new Date().getMonth();

	const monthsShort = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
	/** @type {Record<string, string>} */
	const categoryColors = {
		Utilities: '#3b82f6',   // Blue
		Supplies: '#10b981',    // Emerald
		Payroll: '#f59e0b',     // Amber
		'Delivery/Logistics': '#8b5cf6', // Purple
		Rent: '#ec4899',        // Pink
		Misc: '#64748b'         // Slate
	};

	// --- Filtered Data Computations ---
	let filteredSales = $derived(
		sharedStore.sales.filter((s) => {
			if (!s.date) return true;
			const d = new Date(s.date);
			if (isNaN(d.getTime())) return true;
			if (timeHorizon === 'mtd') {
				return d.getFullYear() === currentYear && d.getMonth() === currentMonth;
			} else if (timeHorizon === 'ytd') {
				return d.getFullYear() === currentYear;
			}
			return true;
		})
	);

	let filteredExpenses = $derived(
		sharedStore.expenses.filter((e) => {
			if (!e.expenseDate) return true;
			const d = new Date(e.expenseDate);
			if (isNaN(d.getTime())) return true;
			if (timeHorizon === 'mtd') {
				return d.getFullYear() === currentYear && d.getMonth() === currentMonth;
			} else if (timeHorizon === 'ytd') {
				return d.getFullYear() === currentYear;
			}
			return true;
		})
	);

	// --- KPI Calculations ---
	let totalRevenue = $derived(filteredSales.reduce((sum, s) => sum + (s.amount || 0), 0));
	let totalExpenses = $derived(filteredExpenses.reduce((sum, e) => sum + (e.amount || 0), 0));
	let netProfit = $derived(totalRevenue - totalExpenses);
	let profitMargin = $derived(totalRevenue > 0 ? (netProfit / totalRevenue) * 100 : 0);

	// Order Statistics
	let totalOrdersCount = $derived(sharedStore.orders.length + sharedStore.sales.length);
	let deliveredOrdersCount = $derived(sharedStore.sales.length);
	let pendingOrdersCount = $derived(sharedStore.orders.filter((o) => o.status === 'Pending').length);
	let expiredOrdersCount = $derived(sharedStore.orders.filter((o) => o.status === 'Expired').length);
	let fulfillmentRate = $derived(
		totalOrdersCount > 0 ? (deliveredOrdersCount / totalOrdersCount) * 100 : 0
	);

	// --- Monthly Trend Data (Jan-Dec for Current Year) ---
	let monthlyTrend = $derived(
		monthsShort.map((monthName, idx) => {
			const monthSales = sharedStore.sales.filter((s) => {
				if (!s.date) return false;
				const d = new Date(s.date);
				return d.getFullYear() === currentYear && d.getMonth() === idx;
			});
			const monthExpenses = sharedStore.expenses.filter((e) => {
				if (!e.expenseDate) return false;
				const d = new Date(e.expenseDate);
				return d.getFullYear() === currentYear && d.getMonth() === idx;
			});

			const rev = monthSales.reduce((sum, s) => sum + (s.amount || 0), 0);
			const exp = monthExpenses.reduce((sum, e) => sum + (e.amount || 0), 0);
			return {
				month: monthName,
				revenue: rev,
				expense: exp,
				profit: rev - exp
			};
		})
	);

	// Max value for scaling SVG chart bars
	let maxMonthlyVal = $derived(
		Math.max(
			1,
			...monthlyTrend.map((m) => Math.max(m.revenue, m.expense))
		)
	);

	// --- Expense Category Breakdown ---
	let categoryBreakdown = $derived.by(() => {
		/** @type {Record<string, number>} */
		const map = {};
		for (const exp of filteredExpenses) {
			const cat = exp.category || 'Misc';
			map[cat] = (map[cat] || 0) + (exp.amount || 0);
		}
		const total = Object.values(map).reduce((a, b) => a + b, 0);
		return Object.entries(map).map(([name, val]) => ({
			name,
			value: val,
			percentage: total > 0 ? (val / total) * 100 : 0,
			color: categoryColors[name] || '#64748b'
		}));
	});

	// --- Month-over-Month (MTM) Comparison ---
	let currentMonthRevenue = $derived(monthlyTrend[currentMonth]?.revenue || 0);
	let prevMonthRevenue = $derived(monthlyTrend[currentMonth > 0 ? currentMonth - 1 : 11]?.revenue || 0);
	let momRevenueGrowth = $derived(
		prevMonthRevenue > 0
			? ((currentMonthRevenue - prevMonthRevenue) / prevMonthRevenue) * 100
			: currentMonthRevenue > 0 ? 100 : 0
	);
</script>

<svelte:head>
	<title>Performance Dashboard · Ledger</title>
	<meta name="description" content="Interactive business performance dashboard with MTY and MTM analytics." />
</svelte:head>

<div class="bg-background font-body-md text-on-surface relative min-h-screen overflow-hidden px-gutter py-xl">
	<!-- Subtle Ambient Background Element with Glow Blobs -->
	<div class="pointer-events-none absolute inset-0 z-0 overflow-hidden opacity-40">
		<div class="glow-blob -top-24 -left-24"></div>
		<div class="glow-blob-2 right-0 bottom-0"></div>
	</div>

	<div class="max-w-container-max relative z-10 mx-auto space-y-xl">
		<!-- TOP PART: USER ACCOUNT & COMPANY HEADER BANNER -->
		<header class="glass-card flex flex-col justify-between gap-lg rounded-2xl p-lg md:p-xl md:flex-row md:items-center">
			<div class="flex items-center gap-md">
				<!-- Account Avatar Badge -->
				<div class="from-primary to-secondary shadow-primary/20 bg-linear-to-tr flex h-14 w-14 items-center justify-center rounded-2xl text-xl font-bold text-on-primary shadow-lg">
					{(authState.user?.displayName || authState.user?.email || 'U').charAt(0).toUpperCase()}
				</div>
				<div>
					<div class="flex items-center gap-2">
						<h1 class="font-headline-md text-headline-md text-on-surface tracking-tight">
							{authState.user?.displayName || authState.profile?.name || 'Account User'}
						</h1>
						<span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-bold uppercase tracking-wider {authState.isOwner ? 'bg-amber-100 text-amber-900 border border-amber-300' : 'bg-blue-100 text-blue-900 border border-blue-300'}">
							{authState.isOwner ? '👑 Owner' : '💼 Employee'}
						</span>
					</div>
					<p class="font-body-md text-body-md text-on-surface-variant flex items-center gap-2 mt-0.5">
						<span>{authState.user?.email || authState.profile?.email || 'user@company.com'}</span>
						{#if authState.profile?.company_code}
							<span class="font-mono text-xs font-bold text-stone-500 bg-surface-container-low px-2 py-0.5 rounded border border-outline-variant">
								Company Code: {authState.profile.company_code}
							</span>
						{/if}
					</p>
				</div>
			</div>

			<!-- Time Horizon Filter Control (MTM / MTY / All Time) -->
			<div class="border-outline-variant bg-surface-container-low/60 flex items-center rounded-xl border p-1 self-start md:self-auto">
				<button
					type="button"
					onclick={() => (timeHorizon = 'mtd')}
					class="font-body-md cursor-pointer rounded-lg px-3 py-1.5 text-xs font-semibold transition-all {timeHorizon === 'mtd'
						? 'bg-primary text-on-primary shadow-sm'
						: 'text-on-surface-variant hover:text-on-surface'}"
				>
					MTD (Month-to-Date)
				</button>
				<button
					type="button"
					onclick={() => (timeHorizon = 'ytd')}
					class="font-body-md cursor-pointer rounded-lg px-3 py-1.5 text-xs font-semibold transition-all {timeHorizon === 'ytd'
						? 'bg-primary text-on-primary shadow-sm'
						: 'text-on-surface-variant hover:text-on-surface'}"
				>
					MTY / YTD (Year-to-Date)
				</button>
				<button
					type="button"
					onclick={() => (timeHorizon = 'all')}
					class="font-body-md cursor-pointer rounded-lg px-3 py-1.5 text-xs font-semibold transition-all {timeHorizon === 'all'
						? 'bg-primary text-on-primary shadow-sm'
						: 'text-on-surface-variant hover:text-on-surface'}"
				>
					All Time
				</button>
			</div>
		</header>

		<!-- FINANCIAL KPI CARDS GRID -->
		<div class="grid grid-cols-1 gap-lg sm:grid-cols-2 lg:grid-cols-4 lg:gap-xl">
			<!-- Total Revenue Card -->
			<div class="glass-card flex flex-col justify-between space-y-sm rounded-xl p-lg">
				<div class="flex items-center justify-between">
					<p class="font-label-md text-label-md tracking-wider text-on-surface-variant uppercase">
						{timeHorizon === 'mtd' ? 'MTD Revenue' : timeHorizon === 'ytd' ? 'MTY / YTD Revenue' : 'Total Revenue'}
					</p>
					<div class="flex h-9 w-9 items-center justify-center rounded-lg bg-emerald-100 text-emerald-700">
						<svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
				</div>
				<p class="font-headline-lg text-headline-lg font-bold tracking-tight text-emerald-600">
					₱{totalRevenue.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
				</p>
				<div class="flex items-center gap-1.5 text-xs">
					<span class="font-semibold text-emerald-700">
						{momRevenueGrowth >= 0 ? `+${momRevenueGrowth.toFixed(1)}%` : `${momRevenueGrowth.toFixed(1)}%`}
					</span>
					<span class="text-on-surface-variant">vs previous month</span>
				</div>
			</div>

			<!-- Total Expenses Card -->
			<div class="glass-card flex flex-col justify-between space-y-sm rounded-xl p-lg">
				<div class="flex items-center justify-between">
					<p class="font-label-md text-label-md tracking-wider text-on-surface-variant uppercase">
						{timeHorizon === 'mtd' ? 'MTD Expenses' : timeHorizon === 'ytd' ? 'MTY / YTD Expenses' : 'Total Expenses'}
					</p>
					<div class="flex h-9 w-9 items-center justify-center rounded-lg bg-amber-100 text-amber-700">
						<svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" d="M15 12H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
				</div>
				<p class="font-headline-lg text-headline-lg font-bold tracking-tight text-amber-600">
					₱{totalExpenses.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
				</p>
				<div class="text-xs text-on-surface-variant">
					<span>Logged across {filteredExpenses.length} operational entries</span>
				</div>
			</div>

			<!-- Net Profit Card -->
			<div class="glass-card flex flex-col justify-between space-y-sm rounded-xl p-lg">
				<div class="flex items-center justify-between">
					<p class="font-label-md text-label-md tracking-wider text-on-surface-variant uppercase">
						Net Profit
					</p>
					<div class="flex h-9 w-9 items-center justify-center rounded-lg {netProfit >= 0 ? 'bg-blue-100 text-blue-700' : 'bg-red-100 text-red-700'}">
						<svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
						</svg>
					</div>
				</div>
				<p class="font-headline-lg text-headline-lg font-bold tracking-tight {netProfit >= 0 ? 'text-primary' : 'text-red-600'}">
					₱{netProfit.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
				</p>
				<div class="flex items-center gap-1.5 text-xs">
					<span class="font-semibold {profitMargin >= 0 ? 'text-emerald-600' : 'text-red-600'}">
						{profitMargin.toFixed(1)}% Margin
					</span>
					<span class="text-on-surface-variant">net yield</span>
				</div>
			</div>

			<!-- Order Fulfillment Rate Card -->
			<div class="glass-card flex flex-col justify-between space-y-sm rounded-xl p-lg">
				<div class="flex items-center justify-between">
					<p class="font-label-md text-label-md tracking-wider text-on-surface-variant uppercase">
						Fulfillment Rate
					</p>
					<div class="flex h-9 w-9 items-center justify-center rounded-lg bg-purple-100 text-purple-700">
						<svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
				</div>
				<p class="font-headline-lg text-headline-lg font-bold tracking-tight text-purple-600">
					{fulfillmentRate.toFixed(1)}%
				</p>
				<div class="text-xs text-on-surface-variant">
					<span>{deliveredOrdersCount} delivered of {totalOrdersCount} total orders</span>
				</div>
			</div>
		</div>

		<!-- CHARTS SECTION: VISUAL AIDS & GRAPHS -->
		<div class="grid grid-cols-1 gap-xl lg:grid-cols-3">
			<!-- MAIN BAR CHART: MONTHLY REVENUE VS EXPENSES (MTY / MTM TREND) -->
			<div class="glass-card space-y-md rounded-2xl p-lg md:p-xl lg:col-span-2">
				<div class="flex flex-col justify-between gap-sm sm:flex-row sm:items-center">
					<div>
						<h2 class="font-title-lg text-title-lg text-on-surface font-bold">
							Monthly Performance ({currentYear})
						</h2>
						<p class="font-body-md text-xs text-on-surface-variant">
							Month-to-Month (MTM) Revenue vs Expense visual breakdown
						</p>
					</div>
					<div class="flex items-center gap-md text-xs">
						<div class="flex items-center gap-1.5">
							<span class="h-3 w-3 rounded-full bg-emerald-500"></span>
							<span class="text-on-surface-variant font-medium">Revenue</span>
						</div>
						<div class="flex items-center gap-1.5">
							<span class="h-3 w-3 rounded-full bg-amber-500"></span>
							<span class="text-on-surface-variant font-medium">Expenses</span>
						</div>
					</div>
				</div>

				<!-- SVG CHART CONTAINER -->
				<div class="relative h-64 w-full pt-4">
					<div class="flex h-full items-end gap-2 border-b border-outline-variant pb-2">
						{#each monthlyTrend as m (m.month)}
							{@const revHeight = maxMonthlyVal > 0 ? (m.revenue / maxMonthlyVal) * 100 : 0}
							{@const expHeight = maxMonthlyVal > 0 ? (m.expense / maxMonthlyVal) * 100 : 0}
							<div class="group relative flex flex-1 items-end justify-center gap-1 h-full">
								<!-- Hover Tooltip -->
								<div class="pointer-events-none absolute -top-12 z-20 hidden rounded-lg bg-stone-900 px-2.5 py-1 text-[10px] text-white shadow-lg group-hover:block whitespace-nowrap">
									<p class="font-bold">{m.month} {currentYear}</p>
									<p class="text-emerald-400">Rev: ₱{m.revenue.toLocaleString()}</p>
									<p class="text-amber-400">Exp: ₱{m.expense.toLocaleString()}</p>
								</div>

								<!-- Revenue Bar -->
								<div
									style="height: {Math.max(revHeight, 4)}%"
									class="bg-linear-to-t w-full max-w-3.5 rounded-t from-emerald-600 to-emerald-400 transition-all duration-300 hover:brightness-110"
								></div>

								<!-- Expense Bar -->
								<div
									style="height: {Math.max(expHeight, 4)}%"
									class="bg-linear-to-t w-full max-w-3.5 rounded-t from-amber-600 to-amber-400 transition-all duration-300 hover:brightness-110"
								></div>
							</div>
						{/each}
					</div>

					<!-- Month X-Axis Labels -->
					<div class="flex justify-between pt-2">
						{#each monthsShort as month (month)}
							<span class="flex-1 text-center text-[11px] font-medium text-on-surface-variant">{month}</span>
						{/each}
					</div>
				</div>
			</div>

			<!-- DONUT CHART: EXPENSE CATEGORY BREAKDOWN -->
			<div class="glass-card space-y-md flex flex-col justify-between rounded-2xl p-lg">
				<div>
					<h2 class="font-title-lg text-title-lg text-on-surface font-bold">
						Expenses by Category
					</h2>
					<p class="font-body-md text-xs text-on-surface-variant">
						Distribution of operational spending
					</p>
				</div>

				{#if categoryBreakdown.length === 0}
					<div class="flex h-48 items-center justify-center text-xs text-on-surface-variant">
						No expense entries recorded for this timeframe.
					</div>
				{:else}
					<!-- Donut Chart Visual SVG -->
					<div class="flex items-center justify-center py-2">
						<svg class="h-44 w-44 -rotate-90 transform" viewBox="0 0 36 36">
							<!-- Donut Background -->
							<circle cx="18" cy="18" r="15.915" fill="none" class="stroke-surface-container-high" stroke-width="3.8" />
							
							{#each categoryBreakdown as cat, index (cat.name)}
								{@const cumulativePct = categoryBreakdown.slice(0, index).reduce((acc, c) => acc + c.percentage, 0)}
								<circle
									cx="18"
									cy="18"
									r="15.915"
									fill="none"
									stroke={cat.color}
									stroke-width="3.8"
									stroke-dasharray="{cat.percentage} {100 - cat.percentage}"
									stroke-dashoffset={-cumulativePct}
									class="transition-all duration-500 hover:stroke-width-5 cursor-pointer"
								/>
							{/each}
						</svg>
					</div>

					<!-- Legend Items -->
					<div class="space-y-xs pt-2">
						{#each categoryBreakdown as cat (cat.name)}
							<div class="flex items-center justify-between text-xs">
								<div class="flex items-center gap-2">
									<span class="h-2.5 w-2.5 rounded-full" style="background-color: {cat.color}"></span>
									<span class="text-on-surface font-medium">{cat.name}</span>
								</div>
								<div class="flex items-center gap-2">
									<span class="font-mono text-on-surface-variant font-semibold">₱{cat.value.toLocaleString()}</span>
									<span class="text-[10px] text-stone-500">({cat.percentage.toFixed(1)}%)</span>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>

		<!-- SECONDARY VISUAL AIDS: ORDER LIFECYCLE & MONTHLY SUMMARY -->
		<div class="grid grid-cols-1 gap-xl md:grid-cols-2">
			<!-- Order Lifecycle Status Distribution -->
			<div class="glass-card space-y-md rounded-2xl p-lg md:p-xl">
				<h2 class="font-title-lg text-title-lg text-on-surface font-bold">
					Order Pipeline Status
				</h2>

				<div class="grid grid-cols-3 gap-md">
					<div class="border-emerald-200 bg-emerald-50/30 flex flex-col items-center rounded-xl border p-md text-center">
						<span class="text-2xl font-bold text-emerald-700">{deliveredOrdersCount}</span>
						<span class="text-xs font-semibold text-emerald-800">Delivered</span>
					</div>
					<div class="border-stone-200 bg-stone-50/50 flex flex-col items-center rounded-xl border p-md text-center">
						<span class="text-2xl font-bold text-stone-700">{pendingOrdersCount}</span>
						<span class="text-xs font-semibold text-stone-800">Pending</span>
					</div>
					<div class="border-red-200 bg-red-50/30 flex flex-col items-center rounded-xl border p-md text-center">
						<span class="text-2xl font-bold text-red-700">{expiredOrdersCount}</span>
						<span class="text-xs font-semibold text-red-800">Expired</span>
					</div>
				</div>

				<!-- Progress Bar Visual Aid -->
				<div class="space-y-xs pt-2">
					<div class="flex justify-between text-xs text-on-surface-variant">
						<span>Lifecycle Completion Ratio</span>
						<span class="font-bold">{fulfillmentRate.toFixed(1)}% Completed</span>
					</div>
					<div class="h-3 w-full overflow-hidden rounded-full bg-surface-container-high flex">
						{#if totalOrdersCount > 0}
							<div style="width: {(deliveredOrdersCount / totalOrdersCount) * 100}%" class="bg-emerald-500 h-full"></div>
							<div style="width: {(pendingOrdersCount / totalOrdersCount) * 100}%" class="bg-stone-400 h-full"></div>
							<div style="width: {(expiredOrdersCount / totalOrdersCount) * 100}%" class="bg-red-500 h-full"></div>
						{/if}
					</div>
				</div>
			</div>

			<!-- Executive Performance Summary -->
			<div class="glass-card space-y-md flex flex-col justify-between rounded-2xl p-lg">
				<div>
					<h2 class="font-title-lg text-title-lg text-on-surface font-bold">
						Executive Performance Summary
					</h2>
					<p class="font-body-md text-xs text-on-surface-variant mt-1">
						Quick business metrics for management audit
					</p>
				</div>

				<div class="space-y-sm text-sm">
					<div class="flex justify-between border-b border-outline-variant/50 pb-2">
						<span class="text-on-surface-variant">Active User Profile:</span>
						<span class="font-semibold text-on-surface">{authState.user?.displayName || authState.profile?.name || 'Account User'}</span>
					</div>
					<div class="flex justify-between border-b border-outline-variant/50 pb-2">
						<span class="text-on-surface-variant">User Assigned Role:</span>
						<span class="font-semibold text-primary">{authState.role}</span>
					</div>
					<div class="flex justify-between border-b border-outline-variant/50 pb-2">
						<span class="text-on-surface-variant">Company Code:</span>
						<span class="font-mono font-bold text-stone-700">{authState.profile?.company_code || 'N/A'}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-on-surface-variant">Delivered Revenue Total:</span>
						<span class="font-bold text-emerald-600">₱{totalRevenue.toLocaleString('en-US', { minimumFractionDigits: 2 })}</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
