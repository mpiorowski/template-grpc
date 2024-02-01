<script>
    /** @type {import("./$types").PageData} */
    export let data;
    const articles = data.articles;
</script>

<h1 class="text-2xl">Welcome to SvelteKit</h1>
<p class="mt-2 text-black">
    Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a>
    to read the documentation
</p>

<ul class="mt-6 list-decimal text-black">
    <h2 class="text-xl">Articles</h2>
    {#each articles as article}
        <li>
            {article.title}
        </li>
    {/each}
</ul>
