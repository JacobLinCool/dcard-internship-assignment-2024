import http from "k6/http";
import { URL } from "https://jslib.k6.io/url/1.0.0/index.js";

export const options = {
    vus: 300, // users
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
    const url = new URL("http://localhost:8080/api/v1/ad");

    url.searchParams.append("offset", randomOffset());
    url.searchParams.append("limit", randomLimit());

    const age = randomAge();
    if (age) url.searchParams.append("age", age);
    const gender = randomGender();
    if (gender) url.searchParams.append("gender", gender);
    const platform = randomPlatform();
    if (platform) url.searchParams.append("platform", platform);
    const country = randomCountry();
    if (country) url.searchParams.append("country", country);

    http.get(url.toString(), { tags: { name: "ad" } });
}
