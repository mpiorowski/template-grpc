/**
 * @param {Promise<T> | (() => T)} promiseOrFunc
 * @returns {Promise<import("./safe.types").Safe<T>> | import("./safe.types").Safe<T>}
 * @template T
 * @public
 */
export function safe(promiseOrFunc) {
    if (promiseOrFunc instanceof Promise) {
        return safeAsync(promiseOrFunc);
    }
    return safeSync(promiseOrFunc);
}

/**
 * @param {Promise<T>} promise
 * @returns {Promise<import("./safe.types").Safe<T>>}
 * @template T
 * @private
 */
async function safeAsync(promise) {
    try {
        const data = await promise;
        return { data, success: true };
    } catch (e) {
        if (e instanceof Error) {
            return { success: false, error: e.message };
        }
        return { success: false, error: "Something went wrong" };
    }
}

/**
 * @param {() => T} func
 * @returns {import("./safe.types").Safe<T>}
 * @template T
 * @private
 */
function safeSync(func) {
    try {
        const data = func();
        return { data, success: true };
    } catch (e) {
        if (e instanceof Error) {
            return { success: false, error: e.message };
        }
        return { success: false, error: "Something went wrong" };
    }
}

/**
 * Callback function for handling gRPC responses safely.
 *
 * @template T - The type of data expected in the response.
 *
 * @param {(value: import("./safe.types").Safe<T>) => void} res - The callback function to handle the response.
 * @returns {(err: import("@grpc/grpc-js").ServiceError | null, data: T | undefined) => void} - A callback function to be used with gRPC response handling.
 */
export function grpcSafe(res) {
    /**
     * Handles the gRPC response and calls the provided callback function safely.
     *
     * @param {import("@grpc/grpc-js").ServiceError | null} err - The error, if any, returned in the response.
     * @param {T | undefined} data - The data returned in the response.
     */
    return (err, data) => {
        if (err) {
            if (err.code === 3) {
                let fields = [];
                try {
                    fields = JSON.parse(err.details);
                } catch (e) {
                    return res({
                        success: false,
                        error: err?.message || "Something went wrong",
                    });
                }

                return res({
                    success: false,
                    error: "Invalid argument",
                    fields: fields,
                });
            }
            return res({
                success: false,
                error: err?.message || "Something went wrong",
            });
        }
        if (!data) {
            return res({
                success: false,
                error: "No data returned",
            });
        }
        res({ data, success: true });
    };
}
