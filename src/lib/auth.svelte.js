import { auth } from '$lib/firebase';
import { onAuthStateChanged, signOut } from 'firebase/auth';

class AuthState {
	/** @type {import('firebase/auth').User | null} */
	user = $state(null);
	/** @type {{ id?: number, name?: string, email?: string, role?: string, company_code?: string } | null} */
	profile = $state(null);
	isLoggedIn = $derived(this.user !== null);
	role = $derived(this.profile?.role || 'Employee');
	isOwner = $derived(this.role === 'Owner');
	loading = $state(true);

	constructor() {
		if (typeof window !== 'undefined' && auth) {
			onAuthStateChanged(auth, async (firebaseUser) => {
				this.user = firebaseUser;
				if (firebaseUser) {
					await this.fetchProfile(firebaseUser.uid);
				} else {
					this.profile = null;
				}
				this.loading = false;
			});
		} else if (typeof window !== 'undefined') {
			this.loading = false;
		}
	}

	async fetchProfile(uid) {
		try {
			const res = await fetch(`http://localhost:8080/api/users/${uid}`);
			if (res.ok) {
				this.profile = await res.json();
			}
		} catch (err) {
			console.error('Failed to fetch user profile:', err);
		}
	}

	async logout() {
		if (auth) {
			await signOut(auth);
			this.user = null;
			this.profile = null;
		}
	}
}

export const authState = new AuthState();
