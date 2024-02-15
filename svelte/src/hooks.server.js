import { redirect } from "@sveltejs/kit";

/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {

    event.locals.user = {
        email: "mateuszpiorowski@gmail.com",
        avatar: "https://picsum.photos/200/300",
    };

    if (event.url.pathname === "/") {
        throw redirect(302, "/users");
    }
    const response = await resolve(event);
    return response;
}
