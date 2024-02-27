import { DIRECTUS_URL } from "$env/static/private";
import { getValue } from "$lib/helpers";
import api from "$lib/server/api";
import { logger, perf } from "$lib/server/logger";

/**
 * @typedef {{
 * id: number;
 * first_name: string;
 * last_name: string;
 * email: string;
 * phone: string;
 * message: string;
 * }} Contact
 */

/**
 * Send a contact form
 * @param {FormData} form
 * @returns {Promise<import("../safe.types").Safe<Contact>>}
 */
export async function sendContact(form) {
    const end = perf("send_contact");
    /** @type {Omit<Contact, "id">} */
    const data = {
        first_name: getValue(form, "first_name"),
        last_name: getValue(form, "last_name"),
        email: getValue(form, "email"),
        phone: getValue(form, "phone"),
        message: getValue(form, "message"),
    };
    /** @type {import("../safe.types").Safe<{data: Contact}>} */
    const r = await api(DIRECTUS_URL + "/items/contact", {
        method: "POST",
        body: data,
    });
    if (!r.success) {
        logger.error(r.error, "Error sending contact form");
        return { success: false, error: "Error sending contact form" };
    }
    logger.debug(r.data, "send_contact");
    end();
    return { success: true, data: r.data.data };
}
