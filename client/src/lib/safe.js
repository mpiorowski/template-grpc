/**
 * @typedef {{
 * data: T,
 * success: true,
 * }|{
 * error: true,
 * msg: string,
 * }} Safe<T>
 * @template T
 */

/**
 * @param {Promise<T> | (() => T)} promiseOrFunc
 * @returns {Promise<import("./safe").Safe<T>> | import("./safe").Safe<T>}
 * @template T
 */
export function safe(promiseOrFunc) {
    if (promiseOrFunc instanceof Promise) {
        return safeAsync(promiseOrFunc);
    }
    return safeSync(promiseOrFunc);
}

/**
 * @param {Promise<T>} promise
 * @returns {Promise<Safe<T>>}
 * @template T
 * @private
 */
async function safeAsync(promise) {
    try {
        const data = await promise;
        return { data, error: false };
    } catch (e) {
        if (e instanceof Error) {
            return { error: true, msg: e.message };
        }
        return { error: true, msg: "Something went wrong" };
    }
}

/**
 * @param {() => T} func
 * @returns {import("./safe").Safe<T>}
 * @template T
 * @private
 */
function safeSync(func) {
    try {
        const data = func();
        return { data, error: false };
    } catch (e) {
        if (e instanceof Error) {
            return { error: true, msg: e.message };
        }
        return { error: true, msg: "Something went wrong" };
    }
}
