/** @type {import('./$types').LayoutServerLoad} */
export function load({ locals }) {
    return {
        email: locals.user.email,
        first_name: locals.user.first_name,
        last_name: locals.user.last_name,
        avatar: locals.user.avatar,
    };
}
