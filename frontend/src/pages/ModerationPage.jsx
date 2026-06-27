import { CheckCircle2, XCircle } from "lucide-react";
import { useEffect, useState } from "react";
import { equipmentApi } from "../api/equipmentApi.js";

export default function ModerationPage() {
  const [items, setItems] = useState([]);
  const [notice, setNotice] = useState("");
  const [reasons, setReasons] = useState({});

  const load = () => equipmentApi.moderationList().then(setItems).catch((error) => setNotice(error.message));

  useEffect(() => {
    load();
  }, []);

  const moderate = async (item, status) => {
    try {
      await equipmentApi.moderate(item.id, { status, reason: reasons[item.id] || "Не соответствует правилам публикации" });
      setNotice(status === "approved" ? "Объявление опубликовано." : "Объявление отклонено.");
      load();
    } catch (error) {
      setNotice(error.message);
    }
  };

  return (
    <section className="panel">
      <div className="panel-head">
        <div>
          <h2>Модерация объявлений</h2>
          <p className="panel-subtitle">Проверьте фото, описание, цену и категорию перед публикацией</p>
        </div>
      </div>
      {notice && <div className="notice">{notice}</div>}
      <div className="moderation-board">
        {items.map((item) => (
          <article className="moderation-card" key={item.id}>
            <img src={item.image_url} alt="" />
            <div>
              <strong>{item.name}</strong>
              <span>{item.owner?.name} · {item.location} · {item.price_per_day} ₽/сутки</span>
              <p>{item.description}</p>
              <textarea
                placeholder="Причина отклонения"
                value={reasons[item.id] || ""}
                onChange={(event) => setReasons({ ...reasons, [item.id]: event.target.value })}
              />
            </div>
            <div className="row-actions">
              <button className="table-button" onClick={() => moderate(item, "approved")}>
                <CheckCircle2 size={16} />
                Одобрить
              </button>
              <button className="danger-button" onClick={() => moderate(item, "rejected")}>
                <XCircle size={16} />
                Отклонить
              </button>
            </div>
          </article>
        ))}
      </div>
      {!items.length && !notice && <div className="empty-state">Новых объявлений на проверку нет.</div>}
    </section>
  );
}
