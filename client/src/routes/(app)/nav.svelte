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
        <a
            on:click={close}
            href="/contact"
            class="text-2xl font-bold text-indigo-600"
        >
            Template
        </a>
    </div>
    <nav class="flex flex-1 flex-col">
        <ul role="list" class="-mx-2 space-y-1">
            <!-- Current: "bg-gray-800 text-white", Default: "text-gray-400 hover:text-white hover:bg-gray-800" -->
            <li>
                <a
                    on:click={close}
                    href="/profile"
                    aria-current={current === "profile" ? "page" : undefined}
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
                    href="/notes"
                    aria-current={current === "notes" ? "page" : undefined}
                    class="group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6
                    {current === 'notes'
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
                        class="feather feather-pen-tool"
                    >
                        <path d="M12 19l7-7 3 3-7 7-3-3z"></path>
                        <path
                            d="M18 13l-1.5-7.5L2 2l3.5 14.5L13 18l5-5z"
                        ></path>
                        <path d="M2 2l7.586 7.586"></path>
                        <circle cx="11" cy="11" r="2"></circle>
                    </svg>
                    Notes
                </a>
            </li>
            <li>
                <a
                    on:click={close}
                    href="/billing"
                    aria-current={current === "billing" ? "page" : undefined}
                    class="group flex gap-x-3 rounded-md p-2 text-sm font-semibold leading-6
                    {current === 'billing'
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
                        class="feather feather-credit-card"
                    >
                        <rect
                            x="1"
                            y="4"
                            width="22"
                            height="16"
                            rx="2"
                            ry="2"
                        ></rect>
                        <line x1="1" y1="10" x2="23" y2="10"></line>
                    </svg>
                    Billing
                </a>
            </li>
        </ul>
    </nav>
</div>
