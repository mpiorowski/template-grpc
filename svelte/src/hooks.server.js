import { redirect } from "@sveltejs/kit";

/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {

    event.locals.user = {
        email: "mat@gmail.com",
        avatar: "https://avatars.githubusercontent.com/u/1?v=4"
    };

    if (event.url.pathname === "/") {
        throw redirect(302, "/users");
    }
    const response = await resolve(event);
    return response;
}
