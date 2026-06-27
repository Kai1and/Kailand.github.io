import { request } from "./httpClient.js";

export const bookingApi = {
  list() {
    return request("/bookings");
  },

  ownerList() {
    return request("/bookings/owner");
  },

  create(payload) {
    return request("/bookings", {
      method: "POST",
      body: JSON.stringify(payload)
    });
  },

  cancel(id) {
    return request(`/bookings/${id}/cancel`, {
      method: "PATCH"
    });
  },

  updateStatus(id, status) {
    return request(`/bookings/${id}/status`, {
      method: "PATCH",
      body: JSON.stringify({ status })
    });
  }
};
