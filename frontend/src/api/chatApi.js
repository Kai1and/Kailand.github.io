import { request } from "./httpClient.js";

export const chatApi = {
  list() {
    return request("/chats");
  },

  start(payload) {
    return request("/chats", {
      method: "POST",
      body: JSON.stringify(payload)
    });
  },

  messages(id) {
    return request(`/chats/${id}/messages`);
  },

  send(id, payload) {
    return request(`/chats/${id}/messages`, {
      method: "POST",
      body: JSON.stringify(typeof payload === "string" ? { body: payload } : payload)
    });
  }
};
