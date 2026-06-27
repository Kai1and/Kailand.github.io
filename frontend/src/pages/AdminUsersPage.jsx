import { Ban, CheckCircle2, ShieldAlert, X } from "lucide-react";
import { useEffect, useState } from "react";
import { userApi } from "../api/userApi.js";
import { useAuthStore } from "../store/authStore.js";

const banReasons = [
  "Подозрение на мошенничество",
  "Нарушение правил размещения",
  "Жалобы пользователей",
  "Передача запрещенных контактов",
  "Нецелевое использование сервиса"
];

export default function AdminUsersPage() {
  const { user } = useAuthStore();
  const [users, setUsers] = useState([]);
  const [notice, setNotice] = useState("");
  const [banTarget, setBanTarget] = useState(null);
  const [banForm, setBanForm] = useState({ reason: banReasons[0], evidence: "" });

  const load = () => userApi.list().then(setUsers).catch((error) => setNotice(error.message));

  useEffect(() => {
    load();
  }, []);

  const block = async (event) => {
    event.preventDefault();
    if (!banTarget) return;
    try {
      await userApi.setBlocked(banTarget.id, { blocked: true, ...banForm });
      setBanTarget(null);
      setBanForm({ reason: banReasons[0], evidence: "" });
      load();
    } catch (error) {
      setNotice(error.message);
    }
  };

  const unblock = async (item) => {
    try {
      await userApi.setBlocked(item.id, { blocked: false });
      load();
    } catch (error) {
      setNotice(error.message);
    }
  };

  return (
    <section className="panel">
      <div className="panel-head">
        <h2>Пользователи</h2>
      </div>
      {notice && <div className="notice">{notice}</div>}
      <div className="table">
        <div className="row header">
          <span>Пользователь</span>
          <span>Роль</span>
          <span>Город</span>
          <span>Доступ</span>
        </div>
        {users.map((item) => (
          <div className="row" key={item.id}>
            <span>
              {item.name}
              <small>{item.email}</small>
              {item.blocked && item.ban_reason && <small>Причина: {item.ban_reason}</small>}
            </span>
            <span>{item.role}</span>
            <span>{item.city || "-"}</span>
            <span>
              {item.id === user?.id ? (
                "Текущий администратор"
              ) : item.blocked ? (
                <button className="table-button" onClick={() => unblock(item)}>
                  <CheckCircle2 size={16} />
                  Разблокировать
                </button>
              ) : (
                <button className="danger-button" onClick={() => setBanTarget(item)}>
                  <Ban size={16} />
                  Заблокировать
                </button>
              )}
            </span>
          </div>
        ))}
      </div>

      {banTarget && (
        <div className="modal-backdrop">
          <form className="modal" onSubmit={block}>
            <div className="modal-head">
              <h2>Бан пользователя</h2>
              <button type="button" className="icon-button" aria-label="Закрыть" onClick={() => setBanTarget(null)}>
                <X size={18} />
              </button>
            </div>
            <div className="safe-note">
              <ShieldAlert size={18} />
              {banTarget.name} потеряет доступ к аккаунту и защищенным разделам.
            </div>
            <label>
              Причина
              <select value={banForm.reason} onChange={(event) => setBanForm({ ...banForm, reason: event.target.value })}>
                {banReasons.map((reason) => (
                  <option key={reason}>{reason}</option>
                ))}
              </select>
            </label>
            <label>
              Доказательство
              <textarea
                value={banForm.evidence}
                onChange={(event) => setBanForm({ ...banForm, evidence: event.target.value })}
                placeholder="Например: ссылка на жалобу, описание нарушения, номер заявки"
              />
            </label>
            <div className="modal-actions">
              <button type="button" className="ghost-button" onClick={() => setBanTarget(null)}>
                Отмена
              </button>
              <button className="danger-button">
                <Ban size={16} />
                Заблокировать
              </button>
            </div>
          </form>
        </div>
      )}
    </section>
  );
}
