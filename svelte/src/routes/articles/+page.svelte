<script>
    /** @type {import("./$types").PageData} */
    export let data;
    const articles = data.articles;
</script>

<ul class="mt-6 max-w-4xl list-outside list-disc">
    <h2 class="text-xl">Articles</h2>
    {#each articles as article}
        <li>
            {article.title}
            <div class="rounded-xl border p-4">{@html article.content}</div>
        </li>
    {/each}
</ul>
