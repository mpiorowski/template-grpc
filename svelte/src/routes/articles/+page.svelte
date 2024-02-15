<script>
    /** @type {import("./$types").PageData} */
    export let data;
    const articles = data.articles;
</script>

<ul class="mt-6 list-inside list-disc">
    <h2 class="text-xl text-gray-900">Articles</h2>
    {#each articles as article}
        <li>
            {article.title}
        </li>
    {/each}
</ul>
