import { getAllForms } from "$lib/server/form_service";
import { error } from "@sveltejs/kit";

export const prerender = true;

/** @type {import('./$types').PageServerLoad} */
export async function load() {
    const forms = await getAllForms();
    if (!forms.success) {
        return error(500, forms.error);
    }
    return {
        forms: forms.data,
    };
}

/** @type {import('./$types').Actions} */
export const actions = {
    insert_form: async ({ locals, request }) => {
        const form = await request.formData();



        return { id };
    },
};
