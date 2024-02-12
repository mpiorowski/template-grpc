import { safe } from "./safe";

/**
 * API wrapper
 * @param {string} url
 * @param {{
 *  method?: "GET" | "POST" | "PUT" | "DELETE"
 *  body?: string
 *  }} options
 *  @returns {Promise<import("./safe.types").Safe<T>>}
 *  @template T
 */
export default async function api(url, { method = "GET", body = undefined } = {}) {
    const res = await safe(
        fetch(url, {
            method,
            headers: {
                "content-type": "application/json",
            },
            body: body ? JSON.stringify(body) : null,
        }),
    );

    if (!res.success) {
        return { success: false, error: res.error };
    }

    // check if empty response
    if (res.data.status === 204) {
        const empty = /** @type {T} */ ({});
        return { success: true, data: empty };
    }
    // check if invalid response
    if (!res.data.headers.get("content-type")?.includes("application/json")) {
        return { success: false, error: "Response was not JSON" };
    }
    const data = await res.data.json();
    return { success: true, data };
}
