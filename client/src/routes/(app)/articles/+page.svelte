<script>
    import { PUBLIC_DIRECTUS_URL } from "$env/static/public";
    /** @type {import("./$types").PageData} */
    export let data;
    const articles = data.articles;
</script>

<ul class="max-w-4xl list-outside list-disc">
    <h2 class="mb-10 text-xl">Articles</h2>
    {#each articles as article}
        <li class="flex flex-col gap-4 rounded-xl border p-4">
            <h2>{article.title}</h2>
            <img
                src="{PUBLIC_DIRECTUS_URL}/assets/{article.cover}"
                alt={article.title}
                class="w-1/2"
            />
            <div>{article.description}</div>
        </li>
    {/each}
</ul>
