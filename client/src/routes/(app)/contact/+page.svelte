<script>
    import { enhance } from "$app/forms";
    import { toast } from "$lib/ui/toast.store";
    import Input from "$lib/form/Input.svelte";
    import Button from "$lib/form/Button.svelte";

    const contact = {
        first_name: "",
        last_name: "",
        email: "",
        phone: "",
        message: "",
    };
</script>

<form
    class="max-w-2xl"
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
    <div class="grid grid-cols-1 gap-x-8 sm:grid-cols-2">
        <Input
            bind:value={contact.first_name}
            label="First name"
            name="first_name"
            type="text"
            autocomplete="given-name"
        />
        <Input
            bind:value={contact.last_name}
            label="Last name"
            name="last_name"
            type="text"
            autocomplete="family-name"
        />
        <div class="sm:col-span-2">
            <Input
                bind:value={contact.email}
                label="Email"
                name="email"
                type="email"
                autocomplete="email"
            />
        </div>
        <div class="sm:col-span-2">
            <Input
                bind:value={contact.phone}
                label="Phone number"
                name="phone"
                type="tel"
                autocomplete="tel"
            />
        </div>
        <div class="sm:col-span-2">
            <Input
                bind:value={contact.message}
                label="Message"
                name="message"
                type="text"
                rows={4}
            />
        </div>
    </div>
    <div class="mt-8 flex justify-end">
        <Button>Send message</Button>
    </div>
</form>
