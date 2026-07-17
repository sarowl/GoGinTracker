<script>
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { auth, googleProvider } from '$lib/firebase';
	import { createUserWithEmailAndPassword, updateProfile, signInWithPopup } from 'firebase/auth';

	// --- Form State using Svelte 5 Runes ---
	let fullName = $state('');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let showPassword = $state(false);
	let showConfirmPassword = $state(false);
	let isSubmitting = $state(false);
	let errorMessage = $state('');

	// Active track states to replicate label micro-interactions
	let activeField = $state('');

	// Derived state for client-side validation
	let passwordsMatch = $derived(confirmPassword.length === 0 || password === confirmPassword);
	let isFormValid = $derived(
		fullName.trim() !== '' &&
			email.includes('@') &&
			password.length >= 6 &&
			password === confirmPassword
	);

	/**
	 * @param {import('firebase/auth').User} firebaseUser
	 * @param {string} [name]
	 */
	async function syncUserWithBackend(firebaseUser, name) {
		try {
			const response = await fetch('http://localhost:8080/api/users', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					name: name || firebaseUser.displayName || firebaseUser.email?.split('@')[0] || 'User',
					email: firebaseUser.email || '',
					firebase_uid: firebaseUser.uid
				})
			});
			if (!response.ok) {
				console.error('Failed to sync user with backend:', await response.text());
			}
		} catch (err) {
			console.error('Error syncing user with backend:', err);
		}
	}

	/**
	 * @param {SubmitEvent} event
	 */
	async function handleSubmit(event) {
		event.preventDefault();
		if (!isFormValid || isSubmitting) return;
		isSubmitting = true;
		errorMessage = '';

		if (!auth) {
			errorMessage =
				'Firebase is not configured. Please supply VITE_FIREBASE_* environment variables.';
			isSubmitting = false;
			return;
		}

		const activeAuth = /** @type {import('firebase/auth').Auth} */ (auth);

		try {
			const userCredential = await createUserWithEmailAndPassword(activeAuth, email, password);
			await updateProfile(userCredential.user, { displayName: fullName });
			await syncUserWithBackend(userCredential.user, fullName);
			await goto(resolve('/sales'));
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : 'Could not create account.';
		} finally {
			isSubmitting = false;
		}
	}

	async function handleGoogleOAuth() {
		errorMessage = '';
		if (!auth || !googleProvider) {
			errorMessage =
				'Firebase is not configured. Please supply VITE_FIREBASE_* environment variables.';
			return;
		}
		const activeAuth = /** @type {import('firebase/auth').Auth} */ (auth);
		const activeProvider = /** @type {import('firebase/auth').GoogleAuthProvider} */ (
			googleProvider
		);
		try {
			const userCredential = await signInWithPopup(activeAuth, activeProvider);
			await syncUserWithBackend(
				userCredential.user,
				userCredential.user.displayName || undefined
			);
			await goto(resolve('/sales'));
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : 'Google authentication failed.';
		}
	}
</script>

<svelte:head>
	<title>Create account · Ledger</title>
	<meta
		name="description"
		content="Create an account on Ledger to empower your business with financial control."
	/>
</svelte:head>

<div class="flex min-h-screen flex-col bg-background font-body-md text-on-surface">
	<!-- Main Content Canvas -->
	<main class="relative flex flex-grow items-center justify-center overflow-hidden py-xl px-gutter">
		<!-- Subtle Ambient Background Element with Glow Blobs -->
		<div class="pointer-events-none absolute inset-0 z-0 overflow-hidden opacity-50">
			<div class="glow-blob -top-24 -left-24"></div>
			<div class="glow-blob-2 right-0 bottom-0"></div>
		</div>

		<!-- Registration Card -->
		<div
			class="glass-card relative z-10 w-full max-w-md rounded-2xl p-lg transition-all duration-300 md:p-xl"
		>
			<div class="mb-xl text-center">
				<div class="mb-md flex justify-center">
					<div
						class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-tr from-primary to-secondary shadow-lg shadow-primary/20"
					>
						<svg
							class="h-6 w-6 text-white"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							xmlns="http://www.w3.org/2000/svg"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
					</div>
				</div>
				<h1 class="mb-xs font-headline-lg text-headline-lg tracking-tight text-on-surface">
					Create Account
				</h1>
				<p class="font-body-md text-body-md text-on-surface-variant">
					Empower your business with precise financial control.
				</p>
			</div>

			{#if errorMessage}
				<div class="px-md py-sm mb-md rounded-lg border border-error/30 bg-error-container">
					<p class="font-body-md text-body-md text-on-error-container">{errorMessage}</p>
				</div>
			{/if}

			<form onsubmit={handleSubmit} class="space-y-md">
				<!-- Full Name -->
				<div class="flex flex-col gap-xs">
					<label
						class="font-label-md text-label-md transition-colors duration-200 {activeField ===
						'fullName'
							? 'text-secondary'
							: 'text-on-surface-variant'}"
						for="full_name"
					>
						Full Name
					</label>
					<div
						class="group relative w-full transition-transform duration-150"
						class:scale-[1.01]={activeField === 'fullName'}
					>
						<input
							bind:value={fullName}
							onfocus={() => (activeField = 'fullName')}
							onblur={() => (activeField = '')}
							class="px-md h-10 w-full rounded border border-outline-variant bg-surface-container-low/40 font-body-md text-body-md text-on-surface placeholder-on-surface-variant/40 transition-all focus:border-primary focus:ring-2 focus:ring-primary/40 focus:outline-none"
							id="full_name"
							placeholder="John Doe"
							type="text"
							required
						/>
					</div>
				</div>

				<!-- Email -->
				<div class="flex flex-col gap-xs">
					<label
						class="font-label-md text-label-md transition-colors duration-200 {activeField ===
						'email'
							? 'text-secondary'
							: 'text-on-surface-variant'}"
						for="email"
					>
						Email Address
					</label>
					<div
						class="group relative w-full transition-transform duration-150"
						class:scale-[1.01]={activeField === 'email'}
					>
						<input
							bind:value={email}
							onfocus={() => (activeField = 'email')}
							onblur={() => (activeField = '')}
							class="px-md h-10 w-full rounded border border-outline-variant bg-surface-container-low/40 font-body-md text-body-md text-on-surface placeholder-on-surface-variant/40 transition-all focus:border-primary focus:ring-2 focus:ring-primary/40 focus:outline-none"
							id="email"
							placeholder="name@company.com"
							type="email"
							required
						/>
					</div>
				</div>

				<!-- Password -->
				<div class="flex flex-col gap-xs">
					<label
						class="font-label-md text-label-md transition-colors duration-200 {activeField ===
						'password'
							? 'text-secondary'
							: 'text-on-surface-variant'}"
						for="password"
					>
						Password
					</label>
					<div
						class="group relative w-full transition-transform duration-150"
						class:scale-[1.01]={activeField === 'password'}
					>
						<input
							bind:value={password}
							onfocus={() => (activeField = 'password')}
							onblur={() => (activeField = '')}
							class="pl-md h-10 w-full rounded border border-outline-variant bg-surface-container-low/40 pr-12 font-body-md text-body-md text-on-surface placeholder-on-surface-variant/40 transition-all focus:border-primary focus:ring-2 focus:ring-primary/40 focus:outline-none"
							id="password"
							placeholder="••••••••"
							type={showPassword ? 'text' : 'password'}
							minlength="6"
							required
						/>
						<button
							type="button"
							onclick={() => (showPassword = !showPassword)}
							class="absolute top-1/2 right-3 -translate-y-1/2 cursor-pointer p-1 text-on-surface-variant transition-colors hover:text-on-surface"
							aria-label={showPassword ? 'Hide password' : 'Show password'}
						>
							{#if showPassword}
								<svg
									class="h-5 w-5 opacity-70"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88"
									/>
								</svg>
							{:else}
								<svg
									class="h-5 w-5 opacity-70"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
									/>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
									/>
								</svg>
							{/if}
						</button>
					</div>
				</div>

				<!-- Confirm Password -->
				<div class="flex flex-col gap-xs">
					<div class="flex items-center justify-between">
						<label
							class="font-label-md text-label-md transition-colors duration-200 {activeField ===
							'confirm_password'
								? 'text-secondary'
								: 'text-on-surface-variant'}"
							for="confirm_password"
						>
							Confirm Password
						</label>
						{#if confirmPassword && !passwordsMatch}
							<span class="text-xs font-medium text-error">Passwords do not match</span>
						{/if}
					</div>
					<div
						class="group relative w-full transition-transform duration-150"
						class:scale-[1.01]={activeField === 'confirm_password'}
					>
						<input
							bind:value={confirmPassword}
							onfocus={() => (activeField = 'confirm_password')}
							onblur={() => (activeField = '')}
							class="pl-md h-10 w-full rounded border bg-surface-container-low/40 pr-12 font-body-md text-body-md text-on-surface placeholder-on-surface-variant/40 transition-all focus:ring-2 focus:outline-none {passwordsMatch
								? 'border-outline-variant focus:border-primary focus:ring-primary/40'
								: 'border-error focus:ring-error/40'}"
							id="confirm_password"
							placeholder="••••••••"
							type={showConfirmPassword ? 'text' : 'password'}
							required
						/>
						<button
							type="button"
							onclick={() => (showConfirmPassword = !showConfirmPassword)}
							class="absolute top-1/2 right-3 -translate-y-1/2 cursor-pointer p-1 text-on-surface-variant transition-colors hover:text-on-surface"
							aria-label={showConfirmPassword ? 'Hide password' : 'Show password'}
						>
							{#if showConfirmPassword}
								<svg
									class="h-5 w-5 opacity-70"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88"
									/>
								</svg>
							{:else}
								<svg
									class="h-5 w-5 opacity-70"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
									/>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
									/>
								</svg>
							{/if}
						</button>
					</div>
				</div>

				<!-- Primary Action -->
				<button
					class="mt-lg relative flex h-12 w-full cursor-pointer items-center justify-center gap-2 overflow-hidden rounded-lg bg-gradient-to-r from-[#1c1917] to-[#44403c] font-headline-md text-body-lg text-on-primary shadow-lg shadow-primary/25 transition-all duration-200 hover:shadow-primary/35 hover:brightness-110 active:scale-[0.98] disabled:pointer-events-none disabled:opacity-40 disabled:shadow-none"
					type="submit"
					disabled={!isFormValid || isSubmitting}
				>
					{#if isSubmitting}
						<svg
							class="h-5 w-5 animate-spin text-white"
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
						>
							<circle
								class="opacity-25"
								cx="12"
								cy="12"
								r="10"
								stroke="currentColor"
								stroke-width="4"
							></circle>
							<path
								class="opacity-75"
								fill="currentColor"
								d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
							></path>
						</svg>
						<span>Creating account…</span>
					{:else}
						<span>Sign Up</span>
					{/if}
				</button>

				<!-- Divider -->
				<div class="py-md flex items-center gap-md">
					<div class="flex-grow border-t border-outline-variant"></div>
					<span class="font-label-md text-label-md text-outline">OR</span>
					<div class="flex-grow border-t border-outline-variant"></div>
				</div>

				<!-- OAuth Action -->
				<button
					onclick={handleGoogleOAuth}
					class="flex h-12 w-full cursor-pointer items-center justify-center gap-md rounded-lg border border-outline-variant bg-surface-container-low/40 font-body-md text-body-md text-on-surface transition-all hover:bg-surface-container-low active:scale-[0.99]"
					type="button"
				>
					<svg class="h-5 w-5" viewBox="0 0 24 24">
						<path
							d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
							fill="#4285F4"
						></path>
						<path
							d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
							fill="#34A853"
						></path>
						<path
							d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z"
							fill="#FBBC05"
						></path>
						<path
							d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 1.2-4.53z"
							fill="#EA4335"
						></path>
					</svg>
					Continue with Google
				</button>
			</form>

			<div class="mt-xl text-center">
				<p class="font-body-md text-body-md text-on-surface-variant">
					Already have an account?
					<a
						class="font-bold text-secondary transition-all hover:text-secondary/80 hover:underline"
						href={resolve('/auth/login')}>Sign in</a
					>
				</p>
			</div>
		</div>
	</main>

	<!-- Footer -->
	<footer
		class="bottom-0 w-full border-t border-outline-variant bg-surface-container-low py-4 dark:border-outline dark:bg-surface-container-lowest"
	></footer>
</div>
