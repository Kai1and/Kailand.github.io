import { Edit3, PlusCircle, Save } from "lucide-react";
import { useEffect, useState } from "react";
import { categoryApi } from "../api/categoryApi.js";
import { equipmentApi } from "../api/equipmentApi.js";

const emptyForm = { category_id: "", name: "", description: "", image_url: "", location: "", price_per_day: "", serial: "" };

const moderationLabel = {
  pending: "На модерации",
  approved: "Опубликовано",
  rejected: "Отклонено"
};

export default function MyListingsPage() {
  const [items, setItems] = useState([]);
  const [categories, setCategories] = useState([]);
  const [form, setForm] = useState(emptyForm);
  const [editing, setEditing] = useState(null);
  const [notice, setNotice] = useState("");

  const load = () => equipmentApi.mine().then(setItems).catch((error) => setNotice(error.message));

  useEffect(() => {
    load();
    categoryApi.list().then(setCategories);
  }, []);

  const submit = async (event) => {
    event.preventDefault();
    const payload = { ...form, category_id: Number(form.category_id), price_per_day: Number(form.price_per_day) };
    try {
      if (editing) {
        await equipmentApi.update(editing, payload);
        setNotice("Изменения отправлены на модерацию.");
      } else {
        await equipmentApi.create(payload);
        setNotice("Объявление отправлено на модерацию.");
      }
      setEditing(null);
      setForm(emptyForm);
      load();
    } catch (error) {
      setNotice(error.message);
    }
  };

  const edit = (item) => {
    setEditing(item.id);
    setForm({
      category_id: String(item.category_id),
      name: item.name,
      description: item.description,
      image_url: item.image_url,
      location: item.location,
      price_per_day: String(item.price_per_day),
      serial: item.serial
    });
  };

  return (
    <div className="management-grid">
      <section className="panel">
        <div className="panel-head">
          <h2>Мои объявления</h2>
        </div>
        <div className="admin-listings">
          {items.map((item) => (
            <article key={item.id}>
              <img src={item.image_url} alt="" />
              <div>
                <strong>{item.name}</strong>
                <span>{item.price_per_day} ₽/сутки · {moderationLabel[item.moderation_status] ?? item.moderation_status}</span>
                {item.reject_reason && <span>Причина: {item.reject_reason}</span>}
              </div>
              <button className="table-button" onClick={() => edit(item)}>
                <Edit3 size={16} />
                Редактировать
              </button>
            </article>
          ))}
        </div>
        {!items.length && <div className="empty-state">У вас пока нет объявлений.</div>}
      </section>
      <form className="panel form-panel animated-form" onSubmit={submit}>
        <div className="panel-head">
          <h2>{editing ? "Редактирование" : "Новое объявление"}</h2>
        </div>
        {notice && <div className="notice">{notice}</div>}
        <input placeholder="Название" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
        <select value={form.category_id} onChange={(e) => setForm({ ...form, category_id: e.target.value })} required>
          <option value="">Категория</option>
          {categories.map((item) => <option key={item.id} value={item.id}>{item.name}</option>)}
        </select>
        <input placeholder="Цена за сутки, ₽" type="number" min="1" value={form.price_per_day} onChange={(e) => setForm({ ...form, price_per_day: e.target.value })} required />
        <input placeholder="Город" value={form.location} onChange={(e) => setForm({ ...form, location: e.target.value })} required />
        <input placeholder="Ссылка на фото" value={form.image_url} onChange={(e) => setForm({ ...form, image_url: e.target.value })} />
        <input placeholder="Серийный номер" value={form.serial} onChange={(e) => setForm({ ...form, serial: e.target.value })} />
        <textarea placeholder="Описание" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} required />
        <button className="primary-button">
          {editing ? <Save size={18} /> : <PlusCircle size={18} />}
          {editing ? "Сохранить и отправить" : "Отправить на модерацию"}
        </button>
      </form>
    </div>
  );
}
