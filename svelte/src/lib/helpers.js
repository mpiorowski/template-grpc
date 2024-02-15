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
export function getError(fields, field) {
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

/**
 * Get the string from a form field
 * @param {FormData} form
 * @param {string} key
 * @returns {string}
 */
export function getValue(form, key) {
    const value = form.get(key);
    if (!value || typeof value !== "string") {
        return "";
    }
    return value;
}

/**
 * Get the file from a form field
 * @param {FormData} form
 * @param {string} key
 * @returns {File}
 */
export function getFile(form, key) {
    const value = form.get(key);
    if (!value || !(value instanceof File)) {
        return new File([], "");
    }
    return value;
}
