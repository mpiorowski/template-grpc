<script>
    import Button from "$lib/form/Button.svelte";
    import Input from "$lib/form/Input.svelte";
    import Checkbox from "$lib/form/Checkbox.svelte";
    import Radio from "$lib/form/Radio.svelte";
    import Select from "$lib/form/Select.svelte";
    import Switch from "$lib/form/Switch.svelte";
    import Dropzone from "$lib/form/Dropzone.svelte";
    import FileInput from "$lib/form/FileInput.svelte";
    import { showToast, toast } from "$lib/ui/toast";
    import Modal from "$lib/ui/Modal.svelte";
    import Drawer from "$lib/ui/Drawer.svelte";
    import SelectMultiple from "$lib/form/SelectMultiple.svelte";
    import Tooltip from "$lib/ui/Tooltip.svelte";
    import { generateId } from "$lib/helpers";
    import SelectNative from "$lib/form/SelectNative.svelte";

    const select_native = /** @type {const} */ ([
        "Native Option 1",
        "Native Option 2",
        "Native Option 3",
    ]);
    const select_custom = /** @type {const} */ ([
        "Custom Option 1",
        "Custom Option 2",
        "Custom Option 3",
    ]);

    const radio = /** @type {const} */ ([
        { value: "radio_1", label: "Radio Option 1" },
        { value: "radio_2", label: "Radio Option 2" },
        { value: "radio_3", label: "Radio Option 3" },
    ]);

    const checkbox = /** @type {const} */ ([
        { value: "checkbox_1", label: "Checkbox Option 1" },
        { value: "checkbox_2", label: "Checkbox Option 2" },
        { value: "checkbox_3", label: "Checkbox Option 3" },
    ]);
    const multi_select = /** @type {const} */ ([
        "Multi Option 1",
        "Multi Option 2",
        "Multi Option 3",
        "Multi Option 4",
        "Multi Option 5",
    ]);

    /** @type {boolean} */
    let openModal = false;

    /** @type {boolean} */
    let openDrawer = false;

    /** @typedef {{
     * input: string;
     * textarea: string;
     * select_native: string;
     * select_custom: string;
     * radio: string;
     * checkbox: string[];
     * switch: boolean;
     * multi_select: string;
     * resume: File;
     * coverPhoto: File;
     * }} Form */

    /**
     * @type Form
     */
    const form = {
        input: "",
        textarea: "",
        select_native: "",
        select_custom: select_custom[0],
        switch: false,
        radio: radio[0].value,
        checkbox: [],
        multi_select: "",
        resume: new File([""], ""),
        coverPhoto: new File([""], ""),
    };

    /** @type {Record<keyof Form, string>} */
    let fields = {
        input: "",
        textarea: "",
        select_native: "",
        select_custom: "",
        switch: "",
        radio: "",
        checkbox: "",
        multi_select: "",
        resume: "",
        coverPhoto: "",
    };
    /**
     * @param {SubmitEvent & {currentTarget: HTMLFormElement}} event
     * @returns {void}
     */
    function handleSubmit(event) {
        console.info(event);
        fields = {
            input: "This field is required",
            textarea: "This field is required",
            select_native: "This field is required",
            select_custom: "This field is required",
            radio: "This field is required",
            checkbox: "This field is required",
            switch: "This field is required",
            multi_select: "This field is required",
            resume: "This field is required",
            coverPhoto: "This field is required",
        };
        const firstError = Object.keys(fields)[0];
        showToast({
            id: generateId(),
            title: "Validation failed",
            description: `Found ${Object.keys(fields).length} errors.`,
            type: "error",
            duration: 6000,
            action: {
                label: "Go to first error",
                onClick: () => {
                    /** @type {HTMLInputElement | null} */
                    const input = document.querySelector(
                        `[name="${firstError}"]`,
                    );
                    input?.focus();
                },
            },
        });
        const form_data_string = JSON.stringify(form, null, 2);
        toast.success("Saved", form_data_string);
    }
</script>

{#if openModal}
    <Modal
        alert
        bind:open={openModal}
        title="Deactivate your account"
        description="Are you sure you want to deactivate your account? All of your data will be permanently removed. This action cannot be undone."
    >
        <form method="post" action="?/deactive">
            <Button variant="danger">Deactivate</Button>
        </form>
    </Modal>
{/if}
{#if openDrawer}
    <Drawer title="Info" open={openDrawer} close={() => (openDrawer = false)}>
        {#each Array(20) as _}
            <p class="mb-10 text-sm text-gray-500">
                {_} Lorem ipsum dolor sit amet consectetur adipisicing elit. Quas
                cupiditate laboriosam fugiat.
            </p>
        {/each}
        <svelte:fragment slot="submit">
            <Button>Save</Button>
        </svelte:fragment>
    </Drawer>
{/if}

<form on:submit|preventDefault={handleSubmit} class="m-auto max-w-2xl p-10">
    <div class="space-y-12">
        <div class="border-b border-gray-900/10 pb-12">
            <div>
                <h2
                    class="flex items-center gap-2 text-base font-semibold leading-7 text-gray-900"
                >
                    Basic Form
                    <Tooltip
                        text="Basic form elements like input, textarea, select, etc."
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            class="feather feather-info h-4"
                        >
                            <circle cx="12" cy="12" r="10" />
                            <line x1="12" y1="16" x2="12" y2="12" />
                            <line x1="12" y1="8" x2="12.01" y2="8" />
                        </svg>
                    </Tooltip>
                </h2>
                <p class="mt-1 text-sm leading-6 text-gray-600">
                    Basic form elements like input, textarea, select, etc.
                </p>
            </div>

            <div class="mt-10 grid grid-cols-1 gap-x-6 gap-y-2 sm:grid-cols-6">
                <div class="sm:col-span-4">
                    <Input
                        name="input"
                        label="Input"
                        bind:value={form.input}
                        error={fields.input}
                    />
                </div>

                <div class="col-span-full">
                    <Input
                        name="textarea"
                        label="Textarea"
                        bind:value={form.textarea}
                        rows={3}
                        helper="Write a few sentences about yourself."
                        error={fields.textarea}
                    />
                </div>

                <div class="sm:col-span-4">
                    <SelectNative
                        name="select_native"
                        label="Select Native"
                        bind:value={form.select_native}
                        error={fields.select_native}
                    >
                        {#each select_native as val}
                            <option value={val}>
                                {val}
                            </option>
                        {/each}
                    </SelectNative>
                </div>

                <div class="sm:col-span-4">
                    <Select
                        name="select_custom"
                        label="Select Custom"
                        bind:value={form.select_custom}
                        values={select_custom}
                        options={select_custom}
                    />
                </div>
            </div>
        </div>

        <div class="border-b border-gray-900/10 pb-12">
            <h2
                class="flex items-center gap-2 text-base font-semibold leading-7 text-gray-900"
            >
                Checkboxes and Radio
                <Tooltip
                    text="Checkboxes and radio buttons for selecting options."
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        class="feather feather-info h-4"
                    >
                        <circle cx="12" cy="12" r="10" />
                        <line x1="12" y1="16" x2="12" y2="12" />
                        <line x1="12" y1="8" x2="12.01" y2="8" />
                    </svg>
                </Tooltip>
            </h2>
            <p class="mt-1 text-sm leading-6 text-gray-600">
                Checkboxes and radio buttons for selecting options.
            </p>

            <div class="mt-10 space-y-10">
                <Switch
                    name="switch"
                    label="Switch"
                    bind:checked={form.switch}
                />
                <fieldset>
                    <legend
                        class="text-sm font-semibold leading-6 text-gray-900"
                    >
                        Checkboxes
                    </legend>
                    <p class="mt-1 text-sm leading-6 text-gray-600">
                        Choose one option from the list.
                    </p>
                    <div class="mt-6 space-y-6">
                    {#each checkbox as c}
                        <Checkbox
                            id="checkbox-{c.value}"
                            name="checkbox"
                            label="{c.label}"
                            value={c.value}
                            group={form.checkbox}
                            description="Checkbox description"
                        />
                    {/each}
                    </div>
                </fieldset>
                <fieldset>
                    <legend
                        class="text-sm font-semibold leading-6 text-gray-900"
                    >
                        Radio
                    </legend>
                    <p class="mt-1 text-sm leading-6 text-gray-600">
                        Choose one option from the list.
                    </p>
                    <div class="mt-6 space-y-6">
                        {#each radio as r, i}
                            <Radio
                                id="radio-{r.value}"
                                name="radio_{i}"
                                label="{r.label}"
                                value={r.value}
                                bind:group={form.radio}
                                description="Radio {i + 1} description"
                            />
                        {/each}
                    </div>
                </fieldset>
            </div>
        </div>
        <div class="pb-12">
            <h2
                class="flex items-center gap-2 text-base font-semibold leading-7 text-gray-900"
            >
                Advanced Form
                <Tooltip
                    text="Advanced form elements like select multiple, file input, dropzone, etc."
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        class="feather feather-info h-4"
                    >
                        <circle cx="12" cy="12" r="10" />
                        <line x1="12" y1="16" x2="12" y2="12" />
                        <line x1="12" y1="8" x2="12.01" y2="8" />
                    </svg>
                </Tooltip>
            </h2>
            <p class="mt-1 text-sm leading-6 text-gray-600">
                Share your profesional details so others can find you.
            </p>
            <div class="mt-6">
                <SelectMultiple
                    name="select_multiple"
                    label="Select Multiple"
                    bind:value={form.multi_select}
                    options={multi_select}
                    error={fields.multi_select}
                />
            </div>

            <div class="col-span-full">
                <FileInput
                    label="Resume"
                    name="resume"
                    bind:file={form.resume}
                    helper="PDF up to 5MB"
                    error={fields.resume}
                />
            </div>

            <div class="col-span-full mt-6">
                <Dropzone
                    name="cover_photo"
                    label="Cover photo"
                    bind:file={form.coverPhoto}
                    description="SVG, PNG, JPG, GIF up to 10MB"
                    accept="image/*"
                    error={fields.coverPhoto}
                />
            </div>
        </div>
    </div>

    <div
        class="sticky bottom-0 flex justify-end border-t border-gray-900/10 bg-white p-4"
    >
        <div class="inline-flex items-center gap-x-4">
            <Button
                type="button"
                variant="link"
                on:click={() => (openModal = true)}
            >
                Deactivate
            </Button>
            <Button
                type="button"
                variant="secondary"
                on:click={() => (openDrawer = true)}
            >
                Info
            </Button>
            <Button>Save</Button>
        </div>
    </div>
</form>
