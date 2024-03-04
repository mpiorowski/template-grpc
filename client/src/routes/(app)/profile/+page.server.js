import { getAllValues, getValue } from "$lib/helpers";
import { profileService } from "$lib/server/grpc";
import { perf } from "$lib/server/logger";
import { createMetadata } from "$lib/server/metadata";
import { grpcSafe } from "$lib/server/safe";
import { Status } from "@grpc/grpc-js/build/src/constants";
import { error, fail } from "@sveltejs/kit";

/** @type {import('./$types').PageServerLoad} */
export async function load({ locals }) {
    const end = perf("load_profile");
    const metadata = createMetadata(locals.user.id);
    const profile = await new Promise((r) =>
        profileService.GetProfile({}, metadata, grpcSafe(r)),
    );
    if (!profile.success) {
        if (profile.code === Status.NOT_FOUND) {
            const newProfile = await new Promise((r) =>
                profileService.InsertProfile({}, metadata, grpcSafe(r)),
            );
            if (!newProfile.success) {
                throw error(500, newProfile.error);
            }
        }
        throw error(500, profile.error);
    }
    end();
    return {
        profile: profile.data,
    };
}

/** @type {import('./$types').Actions} */
export const actions = {
    create_profile: async ({ request, locals }) => {
        const end = perf("create_profile");
        const form = await request.formData();

        /** @type {import("$lib/proto/proto/Profile").Profile} */
        const data = {
            active: getValue(form, "active") === "on",
            username: getValue(form, "username"),
            about: getValue(form, "about"),
            first_name: getValue(form, "first_name"),
            last_name: getValue(form, "last_name"),
            email: getValue(form, "email"),
            country: getValue(form, "country"),
            street_address: getValue(form, "street_address"),
            city: getValue(form, "city"),
            state: getValue(form, "state"),
            zip: getValue(form, "zip"),
            email_notifications: getAllValues(form, "email_notifications"),
            push_notification: getValue(form, "push_notification"),
            resume: "",
            cover: "",
            position: getValue(form, "position"),
            skills: getValue(form, "skills"),
        };

        const metadata = createMetadata(locals.user.id);
        /** @type {import("$lib/server/safe.types").GrpcSafe<import("$lib/proto/proto/Profile").Profile__Output>} */
        const profile = await new Promise((r) =>
            profileService.UpdateProfile(data, metadata, grpcSafe(r)),
        );
        if (!profile.success) {
            return fail(500, { error: profile.error });
        }
        end();
        return { profile: profile.data };
    },
};
