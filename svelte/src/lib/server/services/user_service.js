import { DIRECTUS_URL } from "$env/static/private";
import { getFile, getValue } from "$lib/helpers";
import api from "$lib/server/api";
import { logger, perf } from "$lib/server/logger";

/**
 * @typedef {{
 * id: string;
 * active: string;
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
 * resume: File;
 * cover: File;
 * position: string;
 * skills: string;
 * }} User
 */

/**
 * @type {Omit<User, "resume" | "cover">}
 */
export const emptyUser = {
    id: "",
    active: "off",
    username: "",
    about: "",
    first_name: "",
    last_name: "",
    email: "",
    country: "",
    street_address: "",
    city: "",
    state: "",
    zip: "",
    email_notifications: [],
    push_notification: "",
    position: "",
    skills: "",
};

/**
 * Get users
 * @returns {Promise<import("../safe.types").Safe<User[]>>}
 */
export async function getAllUsers() {
    const end = perf("get_all_users");
    /** @type {import("../safe.types").Safe<{data: User[]}>} */
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
 * Get user by id
 * @param {string} id
 * @returns {Promise<import("../safe.types").Safe<User>>}
 */
export async function getUserById(id) {
    const end = perf("get_user_by_id");
    /** @type {import("../safe.types").Safe<{data: User}>} */
    const data = await api(DIRECTUS_URL + "/items/user/" + id);
    if (!data.success) {
        logger.error("Error getting user", data.error);
        return { success: false, error: "Error getting user" };
    }
    end();
    logger.debug(data, "get_user_by_id");
    return { success: true, data: data.data.data };
}

/**
 * Insert user
 * @param {FormData} form_data
 * @returns {Promise<import("../safe.types").Safe<{data: User}>>}
 */
export async function insertUser(form_data) {
    const end = perf("insert_user");

    /** @type {User} */
    const user = {
        id: "",
        active: getValue(form_data, "active"),
        username: getValue(form_data, "username"),
        about: getValue(form_data, "about"),
        first_name: getValue(form_data, "first_name"),
        last_name: getValue(form_data, "last_name"),
        email: getValue(form_data, "email"),
        country: getValue(form_data, "country"),
        street_address: getValue(form_data, "street_address"),
        city: getValue(form_data, "city"),
        state: getValue(form_data, "state"),
        zip: getValue(form_data, "zip"),
        email_notifications: [],
        push_notification: getValue(form_data, "push_notification"),
        resume: getFile(form_data, "avatar"),
        cover: getFile(form_data, "cover"),
        position: getValue(form_data, "position"),
        skills: getValue(form_data, "skills"),
    };

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
