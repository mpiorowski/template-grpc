import { UserRole } from "$lib/proto/proto/UserRole";
import { grpcSafe } from "$lib/server/safe";
import { usersService } from "$lib/server/grpc";
import { logger, perf } from "$lib/server/logger";
import { createMetadata } from "$lib/server/metadata";
import { redirect } from "@sveltejs/kit";

/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {
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
    let token = event.url.searchParams.get("token");
    if (token) {
        event.cookies.set("token", token, {
            path: "/",
            maxAge: 10,
            domain: "localhost",
        });
        throw redirect(302, "/articles");
    }

    token = event.cookies.get("token") ?? "";
    if (!token) {
        logger.info("No token");
        throw redirect(302, "/auth");
    }

    const metadata = createMetadata(token);
    /** @type {import("$lib/server/safe.types").Safe<import("$lib/proto/proto/AuthResponse").AuthResponse__Output>} */
    const auth = await new Promise((res) => {
        usersService.Auth({}, metadata, grpcSafe(res));
    });
    if (!auth.success || !auth.data.token || !auth.data.user) {
        logger.error("Error during auth");
        throw redirect(302, "/auth");
    }

    event.locals.user = auth.data.user;
    event.locals.token = auth.data.token;
    logger.debug(event.locals.user, "user");

    if (event.url.pathname === "/") {
        throw redirect(302, "/articles");
    }

    end();
    const response = await resolve(event);
    // max age is 30 days
    response.headers.append(
        "set-cookie",
        `token=${auth.data.token}; HttpOnly; SameSite=Lax; Secure; Max-Age=2592000; Path=/; Domain=localhost;`
    );
    return response;
}
