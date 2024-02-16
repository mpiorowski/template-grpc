import { DIRECTUS_URL } from "$env/static/private";
import api from "$lib/server/api";
import { logger, perf } from "$lib/server/logger";

/**
 * @typedef {{
 *   id: number
 *   status: string
 *   title: string
 *   content: string
 *   image: string
 * }} Article
 */

/**
 * Get all articles
 * @returns {Promise<import("../safe.types").Safe<Article[]>>}
 */
export async function getAllArticles() {
    const end = perf("getAllArticles");
    /** @type {import("../safe.types").Safe<{data: Article[]}>} */
    const data = await api(DIRECTUS_URL + "/items/articles");
    if (!data.success) {
        logger.error(`Failed to get articles: ${data.error}`);
        return { success: false, error: data.error };
    }
    logger.debug(data, "articles");
    end();
    return { success: true, data: data.data.data };
}
