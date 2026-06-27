import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { authApi } from "../api/authApi.js";
import { useAuthStore } from "../store/authStore.js";

export default function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { setSession } = useAuthStore();
  const [form, setForm] = useState({ email: "", password: "" });
  const [error, setError] = useState("");

  const submit = async (event) => {
    event.preventDefault();
    setError("");
    try {
      const session = await authApi.login(form);
      setSession(session);
      navigate(location.state?.from?.pathname ?? "/");
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <main className="auth-page">
      <form className="auth-card auth-card-pro" onSubmit={submit}>
        <div className="auth-orbit" />
        <h1>Вход</h1>
        <p>Вернитесь к своим броням, объявлениям и сообщениям.</p>
        {error && <div className="notice">{error}</div>}
        <input placeholder="Email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
        <input
          placeholder="Пароль"
          type="password"
          value={form.password}
          onChange={(e) => setForm({ ...form, password: e.target.value })}
        />
        <button className="primary-button">Войти</button>
        <Link to="/register">Создать аккаунт</Link>
      </form>
    </main>
  );
}
