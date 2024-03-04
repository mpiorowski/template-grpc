<script>
    import { page } from "$app/stores";

    /** @type {() => void} */
    export let close = () => {};

    $: current = $page.url.pathname.split("/")[1];
</script>

<div
    class="flex grow flex-col gap-y-5 overflow-y-auto bg-black/10 px-6 pb-4 ring-1 ring-white/5"
>
    <div class="flex h-16 shrink-0 items-center">
        <a on:click={close} href="/" class="text-2xl font-bold text-indigo-600">
            Template
        </a>
    </div>
    <nav class="flex flex-1 flex-col">
        <ul role="list" class="-mx-2 space-y-1">
            <!-- Current: "bg-gray-800 text-white", Default: "text-gray-400 hover:text-white hover:bg-gray-800" -->
            <li>
                <a
                    on:click={close}
                    href="/contact"
                    class="group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6
                    {current === 'contact'
                        ? 'bg-gray-800 text-white'
                        : 'text-gray-400 hover:bg-gray-800 hover:text-white'}"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        class="feather feather-phone"
                    >
                        <path
                            d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"
                        ></path>
                    </svg>

                    Contact
                </a>
            </li>
            <li>
                <a
                    on:click={close}
                    href="/profile"
                    class="group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6
                    {current === 'profile'
                        ? 'bg-gray-800 text-white'
                        : 'text-gray-400 hover:bg-gray-800 hover:text-white'}"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        class="feather feather-user"
                    >
                        <path
                            d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"
                        ></path>
                        <circle cx="12" cy="7" r="4"></circle>
                    </svg>
                    Profile
                </a>
            </li>
            <li>
                <a
                    on:click={close}
                    href="/articles"
                    class="group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6
                    {current === 'articles'
                        ? 'bg-gray-800 text-white'
                        : 'text-gray-400 hover:bg-gray-800 hover:text-white'}"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="24"
                        height="24"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        class="feather feather-file"
                    >
                        <path
                            d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"
                        ></path>
                        <polyline points="13 2 13 9 20 9"></polyline>
                    </svg>
                    Articles
                </a>
            </li>
        </ul>
    </nav>
</div>
