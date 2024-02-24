import { DIRECTUS_URL } from "$env/static/private";
import api from "$lib/server/api";
import { logger, perf } from "$lib/server/logger";

/**
 * @typedef {{
 * id: number;
 * status: string;
 * title: string;
 * content: string;
 * image: string;
 * }} Article
 */

/**
 * Get all articles
 * @returns {Promise<import("../safe.types").Safe<Article[]>>}
 */
export async function getAllArticles() {
    const end = perf("get_all_articles");
    /** @type {import("../safe.types").Safe<{data: Article[]}>} */
    const r = await api(DIRECTUS_URL + "/items/articles");
    if (!r.success) {
        logger.error(r.error, "Error getting all articles");
        return { success: false, error: r.error };
    }
    logger.debug(r, "get_all_articles");
    end();
    return { success: true, data: r.data.data };
}
