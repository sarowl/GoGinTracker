<script>
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { authState } from '$lib/auth.svelte';

	// Owner mode tab state: 'create' | 'join'
	let ownerMode = $state('create');

	// Form binds
	let companyName = $state('');
	let companyCode = $state('');
	let isSubmitting = $state(false);
	let errorMessage = $state('');
	let successMessage = $state('');

	/**
	 * Create Company (Owner only)
	 * @param {SubmitEvent} event
	 */
	async function handleCreateCompany(event) {
		event.preventDefault();
		if (!companyName.trim() || isSubmitting) return;

		isSubmitting = true;
		errorMessage = '';
		successMessage = '';

		const userUid = authState.user?.uid || '';
		if (!userUid) {
			errorMessage = 'User session not found. Please log in again.';
			isSubmitting = false;
			return;
		}

		try {
			const res = await fetch('http://localhost:8080/api/companies/create', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					name: companyName.trim(),
					ownerUid: userUid
				})
			});

			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || 'Failed to create company');
			}

			if (data.user) {
				authState.profile = data.user;
			}

			successMessage = `Company "${data.company.name}" created successfully! Code: ${data.company.code}`;
			setTimeout(async () => {
				await goto(resolve('/sales'));
			}, 1200);
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : 'Error creating company.';
		} finally {
			isSubmitting = false;
		}
	}

	/**
	 * Join Company (Owner or Employee)
	 * @param {SubmitEvent} event
	 */
	async function handleJoinCompany(event) {
		event.preventDefault();
		if (!companyCode.trim() || isSubmitting) return;

		isSubmitting = true;
		errorMessage = '';
		successMessage = '';

		const userUid = authState.user?.uid || '';
		if (!userUid) {
			errorMessage = 'User session not found. Please log in again.';
			isSubmitting = false;
			return;
		}

		try {
			const res = await fetch('http://localhost:8080/api/companies/join', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					code: companyCode.trim().toUpperCase(),
					userUid
				})
			});

			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || 'Invalid Company Code');
			}

			if (data.user) {
				authState.profile = data.user;
			}

			successMessage = `Successfully joined "${data.company.name}"!`;
			setTimeout(async () => {
				await goto(resolve('/orders'));
			}, 1200);
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : 'Error joining company.';
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Company Setup · Ledger</title>
	<meta name="description" content="Create or join a business organization on Ledger." />
</svelte:head>

<div class="bg-background font-body-md text-on-surface flex min-h-screen flex-col">
	<main class="px-gutter py-xl relative flex flex-grow items-center justify-center overflow-hidden">
		<!-- Subtle Ambient Background Element with Glow Blobs -->
		<div class="pointer-events-none absolute inset-0 z-0 overflow-hidden opacity-50">
			<div class="glow-blob -top-24 -left-24"></div>
			<div class="glow-blob-2 right-0 bottom-0"></div>
		</div>

		<div class="glass-card p-lg md:p-xl relative z-10 w-full max-w-md rounded-2xl transition-all duration-300">
			<!-- Header Title -->
			<div class="mb-xl text-center">
				<div class="mb-md flex justify-center">
					<div class="from-primary to-secondary shadow-primary/20 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-tr shadow-lg">
						<span class="text-2xl">{authState.isOwner ? '👑' : '💼'}</span>
					</div>
				</div>
				<h1 class="text-headline-lg font-headline-lg mb-xs text-on-surface tracking-tight">
					{authState.isOwner ? 'Company Onboarding' : 'Join Business Organization'}
				</h1>
				<p class="text-body-md font-body-md text-on-surface-variant">
					{#if authState.isOwner}
						Set up your business organization or join an existing ledger workspace.
					{:else}
						Enter the unique company code provided by your business owner to continue.
					{/if}
				</p>
			</div>

			<!-- Success Alert -->
			{#if successMessage}
				<div class="border-success/30 bg-success-container mb-md px-md py-sm rounded-lg border">
					<p class="text-body-md font-body-md text-on-success-container font-semibold">{successMessage}</p>
				</div>
			{/if}

			<!-- Error Alert -->
			{#if errorMessage}
				<div class="border-error/30 bg-error-container mb-md px-md py-sm rounded-lg border">
					<p class="text-body-md font-body-md text-on-error-container">{errorMessage}</p>
				</div>
			{/if}

			{#if authState.isOwner}
				<!-- OWNER VIEW: Create vs Join Segmented Control -->
				<div class="border-outline-variant bg-surface-container-low/60 mb-lg flex rounded-lg border p-1">
					<button
						type="button"
						onclick={() => {
							ownerMode = 'create';
							errorMessage = '';
						}}
						class="font-body-md flex-1 cursor-pointer rounded-md py-2 text-center text-xs font-semibold transition-all {ownerMode === 'create'
							? 'bg-primary text-on-primary shadow-sm'
							: 'text-on-surface-variant hover:text-on-surface'}"
					>
						Create Company
					</button>
					<button
						type="button"
						onclick={() => {
							ownerMode = 'join';
							errorMessage = '';
						}}
						class="font-body-md flex-1 cursor-pointer rounded-md py-2 text-center text-xs font-semibold transition-all {ownerMode === 'join'
							? 'bg-primary text-on-primary shadow-sm'
							: 'text-on-surface-variant hover:text-on-surface'}"
					>
						Join Company
					</button>
				</div>

				{#if ownerMode === 'create'}
					<!-- OWNER CREATE FORM -->
					<form onsubmit={handleCreateCompany} class="space-y-md">
						<div class="flex flex-col gap-xs">
							<label for="company_name" class="font-label-md text-label-md text-on-surface-variant">
								Business / Company Name
							</label>
							<input
								id="company_name"
								type="text"
								bind:value={companyName}
								placeholder="e.g. Acme Logistics Corp"
								required
								class="bg-surface-container-low/40 border-outline-variant focus:border-primary focus:ring-primary/40 font-body-md text-body-md text-on-surface placeholder-on-surface-variant/40 h-10 w-full rounded border px-md transition-all focus:ring-2 focus:outline-none"
							/>
						</div>

						<button
							type="submit"
							disabled={!companyName.trim() || isSubmitting}
							class="from-primary to-secondary text-on-primary font-headline-md text-body-lg shadow-primary/25 hover:shadow-primary/35 active:scale-[0.98] mt-lg relative flex h-12 w-full cursor-pointer items-center justify-center gap-2 overflow-hidden rounded-lg bg-gradient-to-r font-semibold shadow-lg transition-all hover:brightness-110 disabled:pointer-events-none disabled:opacity-40"
						>
							{#if isSubmitting}
								<span>Creating Organization…</span>
							{:else}
								<span>Create Business Organization</span>
							{/if}
						</button>
					</form>
				{:else}
					<!-- OWNER JOIN FORM -->
					<form onsubmit={handleJoinCompany} class="space-y-md">
						<div class="flex flex-col gap-xs">
							<label for="company_code" class="font-label-md text-label-md text-on-surface-variant">
								Company Code
							</label>
							<input
								id="company_code"
								type="text"
								bind:value={companyCode}
								placeholder="e.g. CMP-948201"
								required
								class="bg-surface-container-low/40 border-outline-variant focus:border-primary focus:ring-primary/40 font-body-md text-body-md text-on-surface placeholder-on-surface-variant/40 h-10 w-full rounded border px-md uppercase transition-all focus:ring-2 focus:outline-none"
							/>
						</div>

						<button
							type="submit"
							disabled={!companyCode.trim() || isSubmitting}
							class="from-primary to-secondary text-on-primary font-headline-md text-body-lg shadow-primary/25 hover:shadow-primary/35 active:scale-[0.98] mt-lg relative flex h-12 w-full cursor-pointer items-center justify-center gap-2 overflow-hidden rounded-lg bg-gradient-to-r font-semibold shadow-lg transition-all hover:brightness-110 disabled:pointer-events-none disabled:opacity-40"
						>
							{#if isSubmitting}
								<span>Joining Organization…</span>
							{:else}
								<span>Join Existing Organization</span>
							{/if}
						</button>
					</form>
				{/if}
			{:else}
				<!-- EMPLOYEE VIEW: Only Company Code input form -->
				<form onsubmit={handleJoinCompany} class="space-y-md">
					<div class="flex flex-col gap-xs">
						<label for="employee_company_code" class="font-label-md text-label-md text-on-surface-variant">
							Company Code
						</label>
						<input
							id="employee_company_code"
							type="text"
							bind:value={companyCode}
							placeholder="e.g. CMP-948201"
							required
							class="bg-surface-container-low/40 border-outline-variant focus:border-primary focus:ring-primary/40 font-body-md text-body-md text-on-surface placeholder-on-surface-variant/40 h-11 w-full rounded border px-md text-center tracking-widest uppercase transition-all focus:ring-2 focus:outline-none"
						/>
						<span class="text-on-surface-variant mt-1 text-xs">Ask your business owner or administrator for the CMP code.</span>
					</div>

					<button
						type="submit"
						disabled={!companyCode.trim() || isSubmitting}
						class="from-primary to-secondary text-on-primary font-headline-md text-body-lg shadow-primary/25 hover:shadow-primary/35 active:scale-[0.98] mt-lg relative flex h-12 w-full cursor-pointer items-center justify-center gap-2 overflow-hidden rounded-lg bg-gradient-to-r font-semibold shadow-lg transition-all hover:brightness-110 disabled:pointer-events-none disabled:opacity-40"
					>
						{#if isSubmitting}
							<span>Validating Code…</span>
						{:else}
							<span>Join Company Ledger</span>
						{/if}
					</button>
				</form>
			{/if}
		</div>
	</main>
</div>
