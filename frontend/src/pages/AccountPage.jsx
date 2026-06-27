import { ImagePlus, Save } from "lucide-react";
import { useState } from "react";
import { userApi } from "../api/userApi.js";
import { useAuthStore } from "../store/authStore.js";

export default function AccountPage() {
  const { user, setSession, token } = useAuthStore();
  const [form, setForm] = useState({ name: user?.name ?? "", phone: user?.phone ?? "", city: user?.city ?? "", avatar_url: user?.avatar_url ?? "", bio: user?.bio ?? "" });
  const [notice, setNotice] = useState("");

  const submit = async (event) => {
    event.preventDefault();
    try {
      const updated = await userApi.updateProfile(form);
      setSession({ token, user: updated });
      setNotice("Профиль сохранен.");
    } catch (error) {
      setNotice(error.message);
    }
  };

  const pickAvatar = (event) => {
    const file = event.target.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = () => setForm((current) => ({ ...current, avatar_url: reader.result }));
    reader.readAsDataURL(file);
  };

  return (
    <div className="management-grid">
      <section className="profile-card account-summary animated-card">
        <img src={form.avatar_url || "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?auto=format&fit=crop&w=300&q=80"} alt="Аватар" />
        <div>
          <h2>{user?.name}</h2>
          <p>{user?.email}</p>
          <p>{roleName(user?.role)}</p>
        </div>
      </section>
      <form className="panel form-panel animated-form" onSubmit={submit}>
        <div className="panel-head">
          <h2>Данные аккаунта</h2>
        </div>
        {notice && <div className="notice">{notice}</div>}
        <label className="file-picker">
          <ImagePlus size={18} />
          Выбрать аватар с устройства
          <input type="file" accept="image/*" onChange={pickAvatar} />
        </label>
        <input placeholder="Имя" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
        <input placeholder="Телефон" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
        <input placeholder="Город" value={form.city} onChange={(e) => setForm({ ...form, city: e.target.value })} />
        <input placeholder="Ссылка на аватар" value={form.avatar_url} onChange={(e) => setForm({ ...form, avatar_url: e.target.value })} />
        <textarea placeholder="О себе" value={form.bio} onChange={(e) => setForm({ ...form, bio: e.target.value })} />
        <button className="primary-button">
          <Save size={18} />
          Сохранить изменения
        </button>
      </form>
    </div>
  );
}

function roleName(role) {
  if (role === "admin") return "Администратор";
  if (role === "moderator") return "Модератор";
  return "Пользователь";
}
