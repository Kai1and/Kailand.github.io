import { request } from "./httpClient.js";

export const userApi = {
  profile(id) {
    return request(`/profiles/${id}`);
  },

  updateProfile(payload) {
    return request("/profile", {
      method: "PUT",
      body: JSON.stringify(payload)
    });
  }
  ,
  list() {
    return request("/users");
  },
  setBlocked(id, payload) {
    return request(`/users/${id}/blocked`, {
      method: "PATCH",
      body: JSON.stringify(typeof payload === "boolean" ? { blocked: payload } : payload)
    });
  }
};
