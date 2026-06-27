import { CheckCircle2, Mail, UserRound } from "lucide-react";
import { useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { authApi } from "../api/authApi.js";
import { useAuthStore } from "../store/authStore.js";

const namePattern = /^[A-Za-zА-Яа-яЁё][A-Za-zА-Яа-яЁё\s'-]{1,58}[A-Za-zА-Яа-яЁё]$/;
const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export default function RegisterPage() {
  const navigate = useNavigate();
  const { setSession } = useAuthStore();
  const [form, setForm] = useState({ name: "", email: "", password: "" });
  const [error, setError] = useState("");

  const checks = useMemo(() => ({
    name: namePattern.test(form.name.trim()),
    email: emailPattern.test(form.email.trim()),
    password: form.password.length >= 6
  }), [form]);

  const canSubmit = checks.name && checks.email && checks.password;

  const submit = async (event) => {
    event.preventDefault();
    setError("");
    if (!canSubmit) {
      setError("Проверьте имя, email и пароль перед регистрацией.");
      return;
    }
    try {
      const session = await authApi.register({ ...form, name: form.name.trim(), email: form.email.trim().toLowerCase() });
      setSession(session);
      navigate("/");
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <main className="auth-page">
      <form className="auth-card auth-card-pro" onSubmit={submit}>
        <div className="auth-orbit" />
        <h1>Регистрация</h1>
        <p>Создайте профиль, чтобы бронировать технику и публиковать объявления.</p>
        {error && <div className="notice">{error}</div>}
        <label className={checks.name || !form.name ? "field" : "field invalid"}>
          <span><UserRound size={16} /> Имя</span>
          <input placeholder="Например: Алексей Власов" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
          <small>Только буквы, пробел, апостроф или дефис.</small>
        </label>
        <label className={checks.email || !form.email ? "field" : "field invalid"}>
          <span><Mail size={16} /> Email</span>
          <input placeholder="name@example.com" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
          <small>Введите корректный адрес электронной почты.</small>
        </label>
        <label className={checks.password || !form.password ? "field" : "field invalid"}>
          <span><CheckCircle2 size={16} /> Пароль</span>
          <input
            placeholder="Минимум 6 символов"
            type="password"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
          />
          <small>Не короче 6 символов.</small>
        </label>
        <button className="primary-button" disabled={!canSubmit}>Зарегистрироваться</button>
        <Link to="/login">Уже есть аккаунт</Link>
      </form>
    </main>
  );
}
