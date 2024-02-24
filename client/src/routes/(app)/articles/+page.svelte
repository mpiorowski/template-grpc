<script>
    /** @type {import("./$types").PageData} */
    export let data;
    const articles = data.articles;
</script>

<ul class="mt-6 max-w-4xl list-outside list-disc">
    <h2 class="text-xl">Articles</h2>
    {#each articles as article}
        <li class="flex flex-col gap-4 rounded-xl border p-4">
            <h2>{article.title}</h2>
            <img
                src="http://0.0.0.0:8055/assets/{article.image}"
                alt={article.title}
                class="w-1/2"
            />
            <div>{@html article.content}</div>
        </li>
    {/each}
</ul>
