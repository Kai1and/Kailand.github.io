import { useEffect, useState } from "react";

let state = {
  token: localStorage.getItem("token"),
  user: JSON.parse(localStorage.getItem("user") ?? "null")
};

const listeners = new Set();

function notify() {
  listeners.forEach((listener) => listener(state));
}

function setSession({ token, user }) {
  state = { token, user };
  localStorage.setItem("token", token);
  localStorage.setItem("user", JSON.stringify(user));
  notify();
}

function logout() {
  state = { token: null, user: null };
  localStorage.removeItem("token");
  localStorage.removeItem("user");
  notify();
}

export function useAuthStore() {
  const [snapshot, setSnapshot] = useState(state);

  useEffect(() => {
    listeners.add(setSnapshot);
    return () => listeners.delete(setSnapshot);
  }, []);

  return {
    ...snapshot,
    isAuthenticated: () => Boolean(state.token),
    setSession,
    logout
  };
}
