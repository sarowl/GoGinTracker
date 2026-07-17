import { auth } from '$lib/firebase';
import { onAuthStateChanged, signOut } from 'firebase/auth';

class AuthState {
	/** @type {import('firebase/auth').User | null} */
	user = $state(null);
	isLoggedIn = $derived(this.user !== null);
	loading = $state(true);

	constructor() {
		if (typeof window !== 'undefined' && auth) {
			onAuthStateChanged(auth, (firebaseUser) => {
				this.user = firebaseUser;
				this.loading = false;
			});
		} else if (typeof window !== 'undefined') {
			// In the browser, but auth is not initialized
			this.loading = false;
		}
	}

	async logout() {
		if (auth) {
			await signOut(auth);
		}
	}
}

export const authState = new AuthState();
