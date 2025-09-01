import http from "k6/http";
import { check } from "k6";

export const options = {
  vus: 50,
  duration: "5m",
};

export default function () {
  const res = http.get("http://localhost:8002/api/extra");

  check(res, {
    "mensaje enviado: ": (r) => r.status === 200,
  });

  if (res.status !== 200) {
    console.log("📣", JSON.stringify(res.body));
  }
}
