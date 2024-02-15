<script>
    import { page } from "$app/stores";
    import Drawer from "$lib/ui/Drawer.svelte";
    import LoaderIcon from "$lib/icons/LoaderIcon.svelte";
    import Toast from "$lib/ui/toast.svelte";
    import { toastStore } from "$lib/ui/toast.store.js";
    import Avatar from "./avatar.svelte";
    import Breadcrumbs from "./breadcrumbs.svelte";
    import "../app.css";
    import Nav from "./nav.svelte";
    import { navigating } from "$app/stores";

    /** @type {import("./$types").LayoutData} */
    export let data;

    let open = false;
    $: current = $page.url.pathname.split("/")[1];

    /** @type {boolean} */
    let isNavigating = false;

    /** @type {number} */
    let t;
    $: if ($navigating) {
        t = setTimeout(() => {
            isNavigating = true;
        }, 500);
    } else {
        clearTimeout(t);
        isNavigating = false;
    }
</script>

{#if isNavigating}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center bg-black opacity-50"
    >
        <span class="sr-only">Loading...</span>
        <LoaderIcon />
    </div>
{/if}

{#if open}
    <Drawer {open} close={() => (open = false)} position="left">
        <Nav close={() => (open = false)} />
    </Drawer>
{/if}

<div class="min-h-full bg-gray-900 font-poppins text-gray-50 antialiased">
    <!-- Static sidebar for mobile -->
    <div
        class="sticky top-0 z-30 flex items-center gap-x-6 border-b border-white/5 bg-gray-900 px-4 py-4 shadow-sm sm:px-6 lg:hidden"
    >
        <button
            type="button"
            class="-m-2.5 p-2.5 text-gray-400 lg:hidden"
            on:click={() => (open = true)}
        >
            <span class="sr-only">Open sidebar</span>
            <svg
                class="h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                aria-hidden="true"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
                />
            </svg>
        </button>
        <div class="flex-1 text-sm font-semibold leading-6 text-white">
            {current?.replace(/^\w/, (c) => c.toUpperCase())}
        </div>
        <Avatar email={data.email} avatarUrl={data.avatar} />
    </div>

    <!-- Static sidebar for desktop -->
    <div
        class="hidden lg:fixed lg:inset-y-0 lg:z-40 lg:flex lg:w-72 lg:flex-col"
    >
        <!-- Sidebar component, swap this element with another sidebar if you like -->
        <Nav />
    </div>

    <!-- Your content -->
    <main class="lg:pl-72">
        <header
            class="hidden items-center justify-between border-b border-white/5 px-4 py-2 sm:px-6 sm:py-4 lg:flex lg:px-8"
        >
            <Breadcrumbs />
            <Avatar email={data.email} avatarUrl={data.avatar} />
        </header>

        <div class="p-6 sm:p-8 lg:p-10">
            <slot />
        </div>
    </main>

    <!-- Global notification live region, render this permanently at the end of the document -->
    <div
        aria-live="assertive"
        class="pointer-events-none fixed inset-0 z-50 flex items-end px-4 py-6 sm:items-start sm:p-6"
    >
        <div class="flex w-full flex-col items-center space-y-4 sm:items-end">
            {#each $toastStore as toast}
                <Toast {toast} />
            {/each}
        </div>
    </div>
</div>
