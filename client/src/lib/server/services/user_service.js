import { DIRECTUS_URL } from "$env/static/private";
import { getAllValues, getFile, getValue } from "$lib/helpers";
import api from "$lib/server/api";
import { logger, perf } from "$lib/server/logger";

/**
 * @typedef {{
 * id: string
 * date_updated: string
 * active: string
 * username: string
 * about: string
 * first_name: string
 * last_name: string
 * email: string
 * country: string
 * street_address: string
 * city: string
 * state: string
 * zip: string
 * email_notifications: string[]
 * push_notification: string
 * resume: File
 * cover: File
 * position: string
 * skills: string
 * }} User
 */

/**
 * @type {User}
 */
export const emptyUser = {
    id: "",
    date_updated: "",
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
    resume: new File([], ""),
    cover: new File([], ""),
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
    const r = await api(DIRECTUS_URL + "/items/users");
    if (!r.success) {
        logger.error(r.error, "Error getting all users");
        return { success: false, error: "Error getting users" };
    }
    const users = r.data.data.sort((a, b) => {
        // date_updated is a string, so we need to convert it to a number
        return (
            new Date(b.date_updated).getTime() -
            new Date(a.date_updated).getTime()
        );
    });

    end();
    logger.debug(users, "get_all_users");
    return { success: true, data: users };
}

/**
 * Get user by id
 * @param {string} id
 * @returns {Promise<import("../safe.types").Safe<User>>}
 */
export async function getUserById(id) {
    const end = perf("get_user_by_id");
    /** @type {import("../safe.types").Safe<{data: User}>} */
    const r = await api(DIRECTUS_URL + "/items/users/" + id);
    if (!r.success) {
        logger.error(r.error, "Error getting user");
        return { success: false, error: "Error getting user" };
    }
    end();
    logger.debug(r, "get_user_by_id");
    return { success: true, data: r.data.data };
}

/**
 * Create user
 * @param {FormData} form_data
 * @param {string} id
 * @returns {Promise<import("../safe.types").Safe<User>>}
 */
export async function createUser(form_data, id) {
    const end = perf("create_user");

    /** @type {Omit<User, "id" | "date_updated">} */
    const user = {
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
        email_notifications: getAllValues(form_data, "email_notifications"),
        push_notification: getValue(form_data, "push_notification"),
        resume: getFile(form_data, "avatar"),
        cover: getFile(form_data, "cover"),
        position: getValue(form_data, "position"),
        skills: getValue(form_data, "skills"),
    };

    /** @type {import("../safe.types").Safe<{data: User}>} */
    let r;
    if (id !== "-1") {
        r = await api(DIRECTUS_URL + "/items/users/" + id, {
            method: "PATCH",
            body: user,
        });
    } else {
        r = await api(DIRECTUS_URL + "/items/users", {
            method: "POST",
            body: user,
        });
    }
    if (!r.success) {
        logger.error(r.error, "Error creating user");
        return { success: false, error: "Error creating user" };
    }
    end();
    logger.debug(r, "create_user");
    return { success: true, data: r.data.data };
}
