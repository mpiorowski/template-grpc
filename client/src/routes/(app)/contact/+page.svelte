<script>
    import { enhance } from "$app/forms";
    import { toast } from "$lib/ui/toast.store";
</script>

<div class="relative isolate bg-gray-900">
    <div class="mx-auto grid max-w-7xl grid-cols-1 lg:grid-cols-2">
        <form
            action="?/contact"
            method="post"
            use:enhance={() => {
                return async ({ result, update }) => {
                    if (result.type === "success") {
                        toast.success(
                            "Send",
                            "Your message has been sent successfully.",
                        );
                    } else {
                        toast.error(
                            "Send",
                            "An error occurred while sending your message.",
                        );
                    }
                    await update();
                };
            }}
        >
            <div class="mx-auto max-w-xl lg:mr-0 lg:max-w-lg">
                <div class="grid grid-cols-1 gap-x-8 gap-y-6 sm:grid-cols-2">
                    <div>
                        <label
                            for="first-name"
                            class="block text-sm font-semibold leading-6 text-white"
                        >
                            First name
                        </label>
                        <div class="mt-2.5">
                            <input
                                type="text"
                                name="first_name"
                                id="first-name"
                                autocomplete="given-name"
                                class="block w-full rounded-md border-0 bg-white/5 px-3.5 py-2 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                            />
                        </div>
                    </div>
                    <div>
                        <label
                            for="last-name"
                            class="block text-sm font-semibold leading-6 text-white"
                        >
                            Last name
                        </label>
                        <div class="mt-2.5">
                            <input
                                type="text"
                                name="last_name"
                                id="last-name"
                                autocomplete="family-name"
                                class="block w-full rounded-md border-0 bg-white/5 px-3.5 py-2 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                            />
                        </div>
                    </div>
                    <div class="sm:col-span-2">
                        <label
                            for="email"
                            class="block text-sm font-semibold leading-6 text-white"
                        >
                            Email
                        </label>
                        <div class="mt-2.5">
                            <input
                                type="email"
                                name="email"
                                id="email"
                                autocomplete="email"
                                class="block w-full rounded-md border-0 bg-white/5 px-3.5 py-2 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                            />
                        </div>
                    </div>
                    <div class="sm:col-span-2">
                        <label
                            for="phone-number"
                            class="block text-sm font-semibold leading-6 text-white"
                        >
                            Phone number
                        </label>
                        <div class="mt-2.5">
                            <input
                                type="tel"
                                name="phone"
                                id="phone-number"
                                autocomplete="tel"
                                class="block w-full rounded-md border-0 bg-white/5 px-3.5 py-2 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                            />
                        </div>
                    </div>
                    <div class="sm:col-span-2">
                        <label
                            for="message"
                            class="block text-sm font-semibold leading-6 text-white"
                        >
                            Message
                        </label>
                        <div class="mt-2.5">
                            <textarea
                                name="message"
                                id="message"
                                rows="4"
                                class="block w-full rounded-md border-0 bg-white/5 px-3.5 py-2 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                            ></textarea>
                        </div>
                    </div>
                </div>
                <div class="mt-8 flex justify-end">
                    <button
                        type="submit"
                        class="rounded-md bg-indigo-500 px-3.5 py-2.5 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"
                    >
                        Send message
                    </button>
                </div>
            </div>
        </form>
    </div>
</div>
