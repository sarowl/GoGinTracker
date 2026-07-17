import { initializeApp, getApps, getApp } from 'firebase/app';
import { getAuth, GoogleAuthProvider } from 'firebase/auth';

const firebaseConfig = {
	apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
	authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
	projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
	storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
	messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
	appId: import.meta.env.VITE_FIREBASE_APP_ID
};

const isBrowser = typeof window !== 'undefined';

/** @type {import('firebase/app').FirebaseApp | undefined} */
let app;

/** @type {import('firebase/auth').Auth | undefined} */
let auth;

/** @type {import('firebase/auth').GoogleAuthProvider | undefined} */
let googleProvider;

if (isBrowser) {
	app = getApps().length === 0 ? initializeApp(firebaseConfig) : getApp();
	auth = getAuth(app);
	googleProvider = new GoogleAuthProvider();
}

export { auth, googleProvider };
