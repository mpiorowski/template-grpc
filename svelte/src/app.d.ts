// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
    namespace App {
        // interface Error {}
        interface Locals {
            user: {
                id: string;
                email: string;
                avatar: string;
                subscription_id: string;
                subscription_end: string;
                subscription_check: string;
                subscription_active: boolean;
            };
            token: string;
        }
        // interface PageData {}
        // interface PageState {}
        // interface Platform {}
    }
}

export { };
