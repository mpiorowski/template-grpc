import api from '$lib/api';

/** @type {import('./$types').PageServerLoad} */
export function load({ locals }) {
    const data = await api(
    return {
        locals,
    };
}
