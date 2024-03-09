import { PUBLIC_DIRECTUS_URL } from "$env/static/public";
import { getValue } from "$lib/helpers";
import api from "$lib/server/api";
import { logger, perf } from "$lib/server/logger";
import { fail } from "@sveltejs/kit";

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

/** @type {import('./$types').Actions} */
export const actions = {
    contact: async ({ request }) => {
        const end = perf("send_contact");
        const form = await request.formData();
        /** @type {Omit<Contact, "id">} */
        const data = {
            first_name: getValue(form, "first_name"),
            last_name: getValue(form, "last_name"),
            email: getValue(form, "email"),
            phone: getValue(form, "phone"),
            message: getValue(form, "message"),
        };
        /** @type {import("$lib/server/safe").Safe<{data: Contact}>} */
        const r = await api(PUBLIC_DIRECTUS_URL + "/items/contact", {
            method: "POST",
            body: data,
        });
        if (!r.success) {
            logger.error(r.error, "Error sending contact form");
            throw fail(500, { error: "Error sending contact form" });
        }
        end();
        return { success: true };
    },
};
