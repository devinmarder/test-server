import http from "k6/http";
import {check, sleep} from "k6";

let url, opts, r;

export let options = {
    vus: 1,
    duration: "30s"
};

export default function (data) {
    url = "http://localhost:8989/codes/200";
    opts = {
        headers: {
            "Content-Type": "application/json",
        }
    };

    r = http.request("POST", url, 
`{}`, opts);

    check(r, {
        "Status was 200 (OK)": (r) => r.status === 200,
        "Server responded in 100ms or less": (r) => r.timings.duration <= 100,
        "Server responded in 200ms or less": (r) => r.timings.duration <= 200,
        "Server responded in 300ms or less": (r) => r.timings.duration <= 300,
        "Server responded in 400ms or less": (r) => r.timings.duration <= 400,
        "Server responded in 500ms or less": (r) => r.timings.duration <= 500
    });
}
