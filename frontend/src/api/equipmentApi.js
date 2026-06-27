import { request } from "./httpClient.js";

export const equipmentApi = {
  summary() {
    return request("/summary");
  },

  list() {
    return request("/equipment");
  },

  mine() {
    return request("/equipment/mine");
  },

  get(id) {
    return request(`/equipment/${id}`);
  },

  create(payload) {
    return request("/equipment", {
      method: "POST",
      body: JSON.stringify(payload)
    });
  },

  update(id, payload) {
    return request(`/equipment/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload)
    });
  },

  remove(id) {
    return request(`/equipment/${id}`, {
      method: "DELETE"
    });
  }
  ,
  setHidden(id, hidden) {
    return request(`/equipment/${id}/visibility`, {
      method: "PATCH",
      body: JSON.stringify({ hidden })
    });
  },

  moderationList() {
    return request("/moderation/equipment");
  },

  moderate(id, payload) {
    return request(`/moderation/equipment/${id}`, {
      method: "PATCH",
      body: JSON.stringify(payload)
    });
  }
};
