/**
 * API wrapper
 * @param {string} url
 * @param {{
 *  method: string,
 *  body?: string
 *  }} options
 *  @returns {Promise<T>}
 *  @template T
 */
export default async function api(url, { method = "GET", body = undefined }) {
    const res = await fetch(url, {
        method,
        headers: {
            "content-type": "application/json",
        },
        body: body ? JSON.stringify(body) : null,
    });
    if (!res.ok) {
        // server error
        return Promise.reject(res);
    }
    return await res.json();
}
