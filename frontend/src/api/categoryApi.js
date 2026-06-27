import { request } from "./httpClient.js";

export const categoryApi = {
  list() {
    return request("/categories");
  },

  create(payload) {
    return request("/categories", {
      method: "POST",
      body: JSON.stringify(payload)
    });
  },

  update(id, payload) {
    return request(`/categories/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload)
    });
  },

  remove(id) {
    return request(`/categories/${id}`, {
      method: "DELETE"
    });
  }
};
