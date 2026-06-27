import { request } from "./httpClient.js";

export const authApi = {
  register(payload) {
    return request("/auth/register", {
      method: "POST",
      body: JSON.stringify(payload)
    });
  },

  login(payload) {
    return request("/auth/login", {
      method: "POST",
      body: JSON.stringify(payload)
    });
  },

  me() {
    return request("/auth/me");
  }
};
