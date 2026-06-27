import { CalendarDays, MessageCircle, MapPin, ShieldCheck, X } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { bookingApi } from "../api/bookingApi.js";
import { chatApi } from "../api/chatApi.js";
import { equipmentApi } from "../api/equipmentApi.js";
import { demoEquipment } from "../data/demoEquipment.js";
import { useAuthStore } from "../store/authStore.js";

export default function EquipmentDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();
  const [item, setItem] = useState(null);
  const [showRent, setShowRent] = useState(false);
  const [rentForm, setRentForm] = useState({ start_at: "", end_at: "", comment: "" });
  const [notice, setNotice] = useState("");

  useEffect(() => {
    equipmentApi
      .get(id)
      .then(setItem)
      .catch(() => setItem(demoEquipment.find((entry) => String(entry.id) === String(id)) ?? demoEquipment[0]));
  }, [id]);

  const total = useMemo(() => {
    if (!rentForm.start_at || !rentForm.end_at || !item) return 0;
    const start = new Date(rentForm.start_at);
    const end = new Date(rentForm.end_at);
    const days = Math.max(1, Math.ceil((end - start) / 86400000));
    return days * Number(item.price_per_day || 0);
  }, [rentForm, item]);

  if (!item) {
    return <div className="empty-state">Загрузка объявления...</div>;
  }

  const openChat = async () => {
    if (!isAuthenticated()) {
      navigate("/login");
      return;
    }
    try {
      const conversation = await chatApi.start({
        equipment_id: item.id,
        message: `Здравствуйте! Интересует аренда: ${item.name}`
      });
      navigate(`/chats/${conversation.id}`);
    } catch (err) {
      setNotice(err.message);
    }
  };

  const rent = async (event) => {
    event.preventDefault();
    if (!isAuthenticated()) {
      navigate("/login");
      return;
    }
    try {
      await bookingApi.create({
        equipment_id: item.id,
        start_at: new Date(rentForm.start_at).toISOString(),
        end_at: new Date(rentForm.end_at).toISOString(),
        comment: rentForm.comment
      });
      setItem({ ...item, available: false });
      setNotice("Заявка отправлена. Объявление ушло в стоп-лист до ответа владельца.");
      setShowRent(false);
    } catch (err) {
      setNotice(err.message);
    }
  };

  return (
    <div className="detail-layout">
      <section className="detail-main">
        <img className="detail-image" src={item.image_url || demoEquipment[0].image_url} alt={item.name} />
        <h2>{item.name}</h2>
        <div className="meta-line">
          <span>
            <MapPin size={16} />
            {item.location || "Город не указан"}
          </span>
          <span>{item.category?.name}</span>
        </div>
        <p className="description">{item.description}</p>
      </section>

      <aside className="seller-panel">
        <strong className="price">{Number(item.price_per_day || 0).toLocaleString("ru-RU")} ₽/сутки</strong>
        <em className={item.available ? "badge ok" : "badge busy"}>{item.available ? "Доступно для брони" : "В стоп-листе"}</em>
        <button className="primary-button" onClick={openChat}>
          <MessageCircle size={18} />
          Написать продавцу
        </button>
        <button className="ghost-button" onClick={() => setShowRent(true)} disabled={!item.available}>
          <CalendarDays size={18} />
          {item.available ? "Арендовать" : "Уже забронировано"}
        </button>
        {notice && <div className="notice compact">{notice}</div>}

        <Link className="seller-card" to={`/profiles/${item.owner?.id ?? item.owner_id}`}>
          <img src={item.owner?.avatar_url || "https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?auto=format&fit=crop&w=300&q=80"} alt="" />
          <div>
            <strong>{item.owner?.name ?? "Владелец"}</strong>
            <span>{item.owner?.city ?? item.location}</span>
          </div>
        </Link>

        <div className="safe-note">
          <ShieldCheck size={18} />
          Данные защищены JWT-авторизацией, а переписка доступна только участникам сделки.
        </div>
      </aside>

      {showRent && (
        <div className="modal-backdrop">
          <form className="modal" onSubmit={rent}>
            <div className="modal-head">
              <h2>Выберите срок аренды</h2>
              <button type="button" className="icon-button" aria-label="Закрыть" onClick={() => setShowRent(false)}>
                <X size={18} />
              </button>
            </div>
            <label>
              Начало
              <input
                type="datetime-local"
                value={rentForm.start_at}
                onChange={(e) => setRentForm({ ...rentForm, start_at: e.target.value })}
                required
              />
            </label>
            <label>
              Окончание
              <input
                type="datetime-local"
                value={rentForm.end_at}
                onChange={(e) => setRentForm({ ...rentForm, end_at: e.target.value })}
                required
              />
            </label>
            <label>
              Комментарий
              <input
                value={rentForm.comment}
                onChange={(e) => setRentForm({ ...rentForm, comment: e.target.value })}
                placeholder="Например: заберу вечером"
              />
            </label>
            <strong>Итого: {total.toLocaleString("ru-RU")} ₽</strong>
            <div className="modal-actions">
              <button type="button" className="ghost-button" onClick={() => setShowRent(false)}>
                Отмена
              </button>
              <button className="primary-button">Отправить заявку</button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}
