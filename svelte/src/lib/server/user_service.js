import { DIRECTUS_URL } from "$env/static/private";
import api from "$lib/server/api";
import { logger, perf } from "$lib/server/logger";

/**
 * @typedef {{
 * username: string;
 * about: string;
 * country: string;
 * city: string;
 * radio: string;
 * checkbox: string[];
 * switch: boolean;
 * multi_select: string;
 * resume: File;
 * coverPhoto: File;
 * }} Form
 */

/**
 * Get form values
 * @returns {Promise<import("./safe.types").Safe<Form[]>>}
 */
export async function getAllForms() {
    const end = perf("get_all_forms");
    /** @type {import("./safe.types").Safe<{data: Form[]}>} */
    const data = await api(DIRECTUS_URL + "/items/form");
    if (!data.success) {
        logger.error("Error getting forms", data.error);
        return { success: false, error: "Error getting forms" };
    }
    end();
    logger.debug(data, "get_all_forms");
    return { success: true, data: data.data.data };
}

/**
 * Insert form values
 * @param {Form} form
 * @returns {Promise<import("./safe.types").Safe<{data: Form}>>}
 */
export async function insertForm(form) {
    const end = perf("insert_form");
    const data = await api(DIRECTUS_URL + "/items/form", {
        method: "POST",
        body: JSON.stringify(form),
    });
    if (!data.success) {
        logger.error("Error creating form", data.error);
        return { success: false, error: "Error creating form" };
    }
    end();
    logger.debug(data, "create_form");
    return { success: true, data: data.data };
}
