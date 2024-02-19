import { sendContact } from "$lib/server/services/contact_service";
import { fail } from "@sveltejs/kit";

/** @type {import('./$types').Actions} */
export const actions = {
    contact: async ({ request }) => {
        const form = await request.formData();
        const rest = await sendContact(form);
        if (!rest.success) {
            return fail(500, { error: rest.error });
        }
        return {
            body: rest.data,
        };
    },
};
