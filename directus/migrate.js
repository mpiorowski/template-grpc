const BASE_DIRECTUS_URL = "http://localhost:8055";
const BASE_ACCESS_TOKEN = "AmSoehfgbTQ0YDdjsnbwP3b8XZ_6TSpg";

const TARGET_DIRECTUS_URL = "https://template-directus.fly.dev";
const TARGET_ACCESS_TOKEN = "9lqFQfmy8l1Y4Q-ylIjj_1hYO5E4eZNH";

async function main() {
    const snapshot = await getSnapshot();
    const diff = await getDiff(snapshot);
    await applyDiff(diff);
}

main();

async function getSnapshot() {
    const URL = `${BASE_DIRECTUS_URL}/schema/snapshot?access_token=${BASE_ACCESS_TOKEN}`;
    const { data } = await fetch(URL)
        .then((r) => r.json())
        .catch((e) => console.error(e));
    return data;
}

async function getDiff(snapshot) {
    const URL = `${TARGET_DIRECTUS_URL}/schema/diff?access_token=${TARGET_ACCESS_TOKEN}&force=true`;
    const { data } = await fetch(URL, {
        method: "POST",
        body: JSON.stringify(snapshot),
        headers: {
            "Content-Type": "application/json",
        },
    })
        .then((r) => r.json())
        .catch((e) => console.error(e));
    return data;
}

async function applyDiff(diff) {
    const URL = `${TARGET_DIRECTUS_URL}/schema/apply?access_token=${TARGET_ACCESS_TOKEN}&force=true`;

    await fetch(URL, {
        method: "POST",
        body: JSON.stringify(diff),
        headers: {
            "Content-Type": "application/json",
        },
    });
}
