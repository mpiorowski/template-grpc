import { getAllArticles } from "$lib/server/services/articles_service";
import { error } from "@sveltejs/kit";

export const prerender = true;

/** @type {import('./$types').PageServerLoad} */
export async function load() {
    const articles = await getAllArticles();
    if (!articles.success) {
        return error(500, articles.error);
    }
    return {
        articles: articles.data,
    };
}
