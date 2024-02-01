import { TARGET } from "$env/static/private";
import pino from "pino";

export const logger = pino({
    transport: {
        target: "pino-pretty",
        options: {
            colorize: true,
        },
    },
    level: TARGET === "production" ? "info" : "debug",
});

/**
 * Measure the performance
 * @param {string} name - The name of the performance measurement
 * @returns {() => void} - The end function
 */
export function perf(name) {
    if (TARGET === "production") {
        return () => {
            // do nothing
        };
    }
    const start = performance.now();

    /**
     * End the performance measurement
     * @returns {void}
     */
    function end() {
        const duration = performance.now() - start;
        logger.info(`${name}: ${duration.toFixed(4)}ms`);
    }

    return end;
}
