import { UserRole } from "$lib/proto/proto/UserRole";
import { authService, grpcSafe } from "$lib/server/grpc";
import { logger, perf } from "$lib/server/logger";
import { createMetadata } from "$lib/server/metadata";
import { redirect } from "@sveltejs/kit";
import { building } from "$app/environment";

/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {
    if (building) {
        return await resolve(event);
    }

    const end = perf("auth");
    event.locals.user = {
        id: "",
        created: "",
        updated: "",
        deleted: "",
        email: "",
        avatar: "",
        role: UserRole.ROLE_UNSET,
        sub: "",
        subscription_id: "",
        subscription_end: "",
        subscription_check: "",
        subscription_active: false,
    };

    if (event.url.pathname === "/auth") {
        event.cookies.set("token", "", {
            path: "/",
            maxAge: 0,
        });
        return await resolve(event);
    }

    /**
     * Check if the user is coming from the oauth flow
     * If so, set a temporary cookie with the token
     * On the next request, the new token will be used
     */
    let token = event.url.pathname.includes("/token/")
        ? event.url.pathname.split("/token/")[1]
        : "";
    if (token) {
        event.cookies.set("token", token, {
            path: "/",
            maxAge: 10,
        });
        throw redirect(302, "/profile");
    }

    token = event.cookies.get("token") ?? "";
    if (!token) {
        logger.info("No token");
        throw redirect(302, "/auth");
    }

    const metadata = createMetadata(token);
    /** @type {import("$lib/server/safe").Safe<import("$lib/proto/proto/AuthResponse").AuthResponse__Output>} */
    const auth = await new Promise((res) => {
        authService.Auth({}, metadata, grpcSafe(res));
    });
    if (!auth.success) {
        logger.error(`Error during auth: ${auth.error}`);
        throw redirect(302, "/auth?error=unauthorized");
    }
    if (!auth.data.token || !auth.data.user) {
        logger.error("Error during auth");
        throw redirect(302, "/auth?error=unauthorized");
    }

    event.locals.user = auth.data.user;
    event.locals.token = auth.data.token;
    // logger.debug(event.locals.user, "user");

    /**
     * Last check to make sure we have a user
     */
    if (!event.locals.user.id) {
        logger.error("No user found");
        throw redirect(303, "/auth");
    }

    if (event.url.pathname === "/") {
        throw redirect(302, "/profile");
    }

    end();
    const response = await resolve(event);
    // max age is 7 days
    response.headers.append(
        "set-cookie",
        `token=${auth.data.token}; HttpOnly; SameSite=Lax; Secure; Max-Age=604800; Path=/`,
    );
    return response;
}
