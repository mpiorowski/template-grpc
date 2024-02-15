import { emptyUser, getUserById, insertUser } from "$lib/server/services/user_service";
import { error } from "@sveltejs/kit";

/** @type {import('./$types').PageServerLoad} */
export async function load({ params }) {
    const user_id = params.user_id;
    if (!user_id || user_id === "-1") {
        return { user: emptyUser };
    }
    const user = await getUserById(user_id);
    if (!user.success) {
        return error(500, user.error);
    }
    return {
        user: user.data,
    };
}

/** @type {import('./$types').Actions} */
export const actions = {
    insert_user: async ({ request }) => {
        const form = await request.formData();
        return await insertUser(form);
    },
};
