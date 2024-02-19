import {
    emptyUser,
    getUserById,
    createUser,
} from "$lib/server/services/user_service";
import { error, fail, redirect } from "@sveltejs/kit";

/** @type {import('./$types').PageServerLoad} */
export async function load({ params }) {
    const user_id = params.user_id;
    if (!user_id || user_id === "-1") {
        return { user: { ...emptyUser, resume: "", cover: "" } };
    }
    const user = await getUserById(user_id);
    if (!user.success) {
        return error(500, user.error);
    }
    return {
        user: {
            ...user.data,
            resume: user.data.resume
                ? await user.data.resume.arrayBuffer()
                : "",
            cover: user.data.cover ? await user.data.cover.arrayBuffer() : "",
        },
    };
}

/** @type {import('./$types').Actions} */
export const actions = {
    create_user: async ({ request, params }) => {
        const form = await request.formData();
        const rest = await createUser(form, params.user_id);
        if (!rest.success) {
            return fail(500, { error: rest.error });
        }
        if (params.user_id === "-1") {
            throw redirect(302, `/users?success=created`);
        }
        return { user: rest.data };
    },
};
