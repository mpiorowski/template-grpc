/**
 * Generate a unique ID for an element
 * @returns {string}
 */
export function generateId() {
    const timestamp = new Date().getTime().toString(36);
    const random = Math.random().toString(36).substring(2, 5); // Use a portion of the random number
    return timestamp + random;
}

/**
 * Check if an element is an HTMLElement
 * @param {Element | null | undefined} element
 * @returns {HTMLElement | undefined}
 */
export function checkElement(element) {
    if (!element || !(element instanceof HTMLElement)) {
        return undefined;
    }
    return element;
}

/**
 * Extract errors from a fields object
 * @param {Record<string, string[]> | undefined} fields
 * @param {string} field
 * @returns {string[]}
 */
export function extractErrors(fields, field) {
    if (!fields) {
        return [];
    }
    for (const [key, value] of Object.entries(fields)) {
        if (key === field) {
            return value;
        }
    }
    return [];
}
