import { redirect } from "@sveltejs/kit";

/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {
    if (event.url.pathname === "/") {
        throw redirect(302, "/form");
    }
    const response = await resolve(event);
    return response;
}
