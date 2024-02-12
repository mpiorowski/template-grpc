import { DIRECTUS_URL } from "$env/static/private";
import api from "$lib/server/api";
import { logger, perf } from "$lib/server/logger";

/**
 * @typedef {{
 *   id: number
 *   status: string
 *   title: string
 * }} Article
 */

/**
 * Get all articles
 * @returns {Promise<import("../types").Safe<Article[]>>}
 */
export async function getAllArticles() {
    const end = perf("getAllArticles");
    /** @type {import("../types").Safe<{data: Article[]}>} */
    const data = await api(DIRECTUS_URL + "/items/articles");
    if (!data.success) {
        logger.error(`Failed to get articles: ${data.error}`);
        return { success: false, error: data.error };
    }
    end();
    return { success: true, data: data.data.data };
}
