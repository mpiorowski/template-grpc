import { getAllUsers } from "$lib/server/services/user_service";
import { error } from "@sveltejs/kit";

/** @type {import('./$types').PageServerLoad} */
export async function load() {
    const users = await getAllUsers();
    if (!users.success) {
        return error(500, users.error);
    }
    return {
        users: users.data,
    };
}
