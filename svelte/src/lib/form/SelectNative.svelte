<script>
    /** @type {string} */
    export let name;
    /** @type {string} */
    export let label;
    /** @type {string | number} */
    export let value;
    /** @type {string} */
    export let error = "";
    /** @type {string} */
    export let helper = "";
</script>

<div class="mb-2">
    <label for={name} class="block text-sm font-medium leading-6 text-gray-900">
        {label}
    </label>
    <select
        bind:value
        id={name}
        {name}
        class="mt-2 block w-full rounded-md border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6"
        aria-describedby="{name}-description"
    >
        <slot />
    </select>
    <p
        id="{name}-description"
        class="inline-block text-xs leading-6
        {error ? 'text-red-600' : 'text-gray-500'}"
    >
        {error || helper}
    </p>
</div>
