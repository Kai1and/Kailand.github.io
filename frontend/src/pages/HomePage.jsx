import { ArrowRight, CalendarDays, ChevronLeft, ChevronRight, HelpCircle, LifeBuoy, Package, Search, ShieldCheck, Sparkles } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { equipmentApi } from "../api/equipmentApi.js";
import { demoEquipment } from "../data/demoEquipment.js";

const fallbackSummary = {
  equipment_total: 0,
  available_total: 0,
  busy_total: 0,
  active_bookings: 0,
  pending_bookings: 0
};

const pageSize = 6;

export default function HomePage() {
  const [summary, setSummary] = useState(fallbackSummary);
  const [equipment, setEquipment] = useState([]);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [error, setError] = useState("");

  useEffect(() => {
    equipmentApi.summary().then(setSummary).catch(() => setSummary(fallbackSummary));
    equipmentApi
      .list()
      .then((items) => setEquipment(items.length ? items : demoEquipment))
      .catch((err) => {
        setError(err.message);
        setEquipment(demoEquipment);
      });
  }, []);

  const found = useMemo(() => {
    const query = search.trim().toLowerCase();
    return equipment.filter((item) => {
      if (!query) return true;
      return [item.name, item.description, item.location, item.category?.name].filter(Boolean).join(" ").toLowerCase().includes(query);
    });
  }, [equipment, search]);

  const pages = Math.max(1, Math.ceil(found.length / pageSize));
  const visible = found.slice((page - 1) * pageSize, page * pageSize);

  useEffect(() => {
    setPage(1);
  }, [search]);

  return (
    <div className="page-grid">
      <section className="search-hero animated-hero">
        <div>
          <p className="eyebrow">Аренда техники рядом</p>
          <h2>Найдите оборудование для учебы, съемки или мероприятия</h2>
        </div>
        <div className="search-bar">
          <Search size={20} />
          <input placeholder="Камера, проектор, ноутбук" aria-label="Поиск оборудования" value={search} onChange={(event) => setSearch(event.target.value)} />
          <Sparkles size={20} />
        </div>
      </section>

      <section className="stats">
        <article>
          <Package size={20} />
          <strong>{summary.equipment_total}</strong>
          <span>единиц оборудования</span>
        </article>
        <article>
          <CalendarDays size={20} />
          <strong>{summary.busy_total}</strong>
          <span>в стоп-листе сейчас</span>
        </article>
        <article>
          <ShieldCheck size={20} />
          <strong>{summary.pending_bookings}</strong>
          <span>ожидают ответа владельца</span>
        </article>
      </section>

      <section className="panel catalog-panel">
        <div className="panel-head">
          <div>
            <h2>Доска объявлений</h2>
            <p className="panel-subtitle">Одобренные модерацией предложения от пользователей</p>
          </div>
          <span className="result-count">{found.length} найдено</span>
        </div>
        {error && <div className="notice">{error}</div>}
        <div className="card-grid">
          {visible.map((item) => (
            <Link className="listing-card elevated" to={`/equipment/${item.id}`} key={item.id}>
              <div className="listing-media">
                <img src={item.image_url || demoEquipment[0].image_url} alt={item.name} />
                <em className={item.available ? "badge ok" : "badge busy"}>{item.available ? "Доступно" : "Стоп-лист"}</em>
              </div>
              <div className="listing-info">
                <span>{Number(item.price_per_day || 0).toLocaleString("ru-RU")} ₽/сутки</span>
                <strong>{item.name}</strong>
                <small>{item.location || item.category?.name || "Город не указан"}</small>
              </div>
            </Link>
          ))}
        </div>
        {!visible.length && <div className="empty-state">По вашему запросу ничего не найдено.</div>}
        <div className="pagination">
          <button className="icon-button" onClick={() => setPage(Math.max(1, page - 1))} disabled={page === 1} aria-label="Назад">
            <ChevronLeft size={18} />
          </button>
          <span>{page} / {pages}</span>
          <button className="icon-button" onClick={() => setPage(Math.min(pages, page + 1))} disabled={page === pages} aria-label="Вперед">
            <ChevronRight size={18} />
          </button>
          <Link className="text-link" to="/listings">
            Разместить объявление <ArrowRight size={16} />
          </Link>
        </div>
      </section>

      <section className="info-grid">
        <article className="glow-card">
          <LifeBuoy size={22} />
          <h2>Help</h2>
          <p>Если товар нужен срочно, напишите владельцу перед бронированием. В чате можно уточнить комплект, место передачи и приложить фото.</p>
        </article>
        <article className="glow-card">
          <ShieldCheck size={22} />
          <h2>Безопасность</h2>
          <p>Новые объявления проходят модерацию, аккаунты с нарушениями блокируются администратором, а сообщения хранятся в защищенном виде.</p>
        </article>
        <article className="glow-card">
          <HelpCircle size={22} />
          <h2>FAQ</h2>
          <details>
            <summary>Когда объявление появится на сайте?</summary>
            <p>После проверки модератором. До этого оно видно владельцу в разделе “Мои объявления”.</p>
          </details>
          <details>
            <summary>Почему товар в стоп-листе?</summary>
            <p>На него уже отправлена активная заявка или бронь подтверждена владельцем.</p>
          </details>
        </article>
      </section>
    </div>
  );
}
