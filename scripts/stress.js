import http from "k6/http";

export const options = {
    vus: 500, // users
    duration: "10s",
    // rps: 10500,
};

const randomOffset = () => Math.floor(Math.random() * 100);
const randomLimit = () => Math.floor(Math.random() * 100) + 1;
const randomAge = () => Math.floor(Math.random() * 101);
const randomGender = () => ["M", "F", ""][Math.floor(Math.random() * 3)];
const randomPlatform = () => ["android", "ios", "web", ""][Math.floor(Math.random() * 4)];
const randomCountry = () => ["US", "CA", "MX", "BR", "JP", "KR", "CN", "RU", "AU", "NZ", "TW", ""][Math.floor(Math.random() * 12)];

export default function () {
    let url = `http://localhost:8080/api/v1/ad?offset=${randomOffset()}&limit=${randomLimit()}`;

    const age = randomAge();
    if (age) url += `&age=${age}`;
    const gender = randomGender();
    if (gender) url += `&gender=${gender}`;
    const platform = randomPlatform();
    if (platform) url += `&platform=${platform}`;
    const country = randomCountry();
    if (country) url += `&country=${country}`;

    http.get(url.toString(), { tags: { name: "ad" } });
}
