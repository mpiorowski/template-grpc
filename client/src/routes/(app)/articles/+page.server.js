import { PUBLIC_DIRECTUS_URL } from "$env/static/public";
import api from "$lib/server/api";
import { logger, perf } from "$lib/server/logger";
import { error } from "@sveltejs/kit";

/**
 * @typedef {{
 * id: number;
 * status: string;
 * title: string;
 * description: string;
 * cover: string;
 * }} Article
 */

/** @type {import('./$types').PageServerLoad} */
export async function load() {
    const end = perf("get_all_articles");
    /** @type {import("$lib/server/safe.types").Safe<{data: Article[]}>} */
    const articles = await api(PUBLIC_DIRECTUS_URL + "/items/articles");
    if (!articles.success) {
        logger.error(`Error getting all articles: ${articles.error}`);
        return error(500, articles.error);
    }
    end();
    return {
        articles: articles.data.data,
    };
}
