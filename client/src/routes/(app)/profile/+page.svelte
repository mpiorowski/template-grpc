<script>
    import { page } from "$app/stores";
    import Button from "$lib/form/Button.svelte";
    import { toast } from "$lib/ui/toast.store";

    /** @type {import("./$types").PageData} */
    export let data;

    $: if ($page.url.searchParams.get("success") === "created") {
        toast.success("Created", "User has been created");
    }
</script>

<div class="flex max-w-4xl place-content-center justify-between">
    <h1>Users</h1>
    <Button href="/users/-1">Add new user</Button>
</div>

<div class="mt-10 flex max-w-4xl flex-col gap-4">
    {#each data.users as user}
        <div class="flex flex-col gap-2 rounded-xl border border-gray-600 p-4">
            <h2>{user.username}</h2>
            <p>{user.email}</p>
            <p>{user.first_name} {user.last_name}</p>
            <p>{user.active ? "Active" : "Inactive"}</p>
            <Button href="/users/{user.id}">Edit</Button>
        </div>
    {/each}
</div>
