<script>
    import Drawer from "$lib/ui/Drawer.svelte";
    import LoaderIcon from "$lib/icons/LoaderIcon.svelte";
    import Toast from "$lib/ui/toast.svelte";
    import { toastStore } from "$lib/ui/toast.store.js";
    import Avatar from "./avatar.svelte";
    import "../app.css";
    import Nav from "./nav.svelte";
    import { navigating } from "$app/stores";

    /** @type {import("./$types").LayoutData} */
    export let data;

    let open = false;

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
    <!-- Static sidebar for desktop -->
    <div
        class="hidden lg:fixed lg:inset-y-0 lg:z-40 lg:flex lg:w-72 lg:flex-col"
    >
        <!-- Sidebar component, swap this element with another sidebar if you like -->
        <Nav />
    </div>

    <!-- Your content -->
    <main class="lg:pl-72">
        <div
            class="sticky top-0 z-30 flex h-16 shrink-0 items-center gap-x-4 border-b border-l border-white/5 bg-gray-900 px-4 shadow-sm sm:gap-x-6 sm:px-6 lg:px-8"
        >
            <button
                type="button"
                class="-m-2.5 p-2.5 lg:hidden"
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

            <!-- Separator -->
            <div class="h-6 w-px bg-white/10 lg:hidden" aria-hidden="true" />

            <div class="flex flex-1 gap-x-4 self-stretch lg:gap-x-6">
                <form class="relative flex flex-1" action="#" method="GET">
                    <label for="search-field" class="sr-only">Search</label>
                    <svg
                        class="pointer-events-none absolute inset-y-0 left-0 h-full w-5 text-gray-400"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                        aria-hidden="true"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z"
                            clip-rule="evenodd"
                        />
                    </svg>
                    <input
                        id="search-field"
                        class="block h-full w-full border-0 bg-gray-900 py-0 pl-8 pr-0 placeholder:text-gray-400 focus:ring-0 sm:text-sm"
                        placeholder="Search..."
                        type="search"
                        name="search"
                    />
                </form>
                <div class="flex items-center gap-x-4 lg:gap-x-6">
                    <button
                        type="button"
                        class="-m-2.5 p-2.5 text-gray-400 hover:text-gray-400"
                    >
                        <span class="sr-only">View notifications</span>
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
                                d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0"
                            />
                        </svg>
                    </button>

                    <!-- Separator -->
                    <div
                        class="hidden lg:block lg:h-6 lg:w-px lg:bg-white/10"
                        aria-hidden="true"
                    />

                    <!-- Profile dropdown -->
                    <Avatar email={data.email} avatar={data.avatar} first_name={data.first_name} last_name={data.last_name} />
                </div>
            </div>
        </div>

        <div class="container mx-auto p-6 sm:p-8 lg:p-10">
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
