import { DIRECTUS_URL } from "$env/static/private";
import api from "$lib/server/api";
import { logger, perf } from "$lib/server/logger";

/**
 * @typedef {{
 * username: string;
 * about: string;
 * first_name: string;
 * last_name: string;
 * email: string;
 * country: string;
 * street_address: string;
 * city: string;
 * state: string;
 * zip: string;
 * email_notifications: string[];
 * push_notification: string;
 * avatar: File;
 * cover: File;
 * }} User
 */

/**
 * Get users
 * @returns {Promise<import("./safe.types").Safe<User[]>>}
 */
export async function getAllUsers() {
    const end = perf("get_all_users");
    /** @type {import("./safe.types").Safe<{data: User[]}>} */
    const data = await api(DIRECTUS_URL + "/items/users");
    if (!data.success) {
        logger.error("Error getting users", data.error);
        return { success: false, error: "Error getting users" };
    }
    end();
    logger.debug(data, "get_all_users");
    return { success: true, data: data.data.data };
}

/**
 * Insert user
 * @param {User} user
 * @returns {Promise<import("./safe.types").Safe<{data: User}>>}
 */
export async function insertUser(user) {
    const end = perf("insert_user");
    const data = await api(DIRECTUS_URL + "/items/user", {
        method: "POST",
        body: JSON.stringify(user),
    });
    if (!data.success) {
        logger.error("Error creating user", data.error);
        return { success: false, error: "Error creating user" };
    }
    end();
    logger.debug(data, "create_user");
    return { success: true, data: data.data };
}
