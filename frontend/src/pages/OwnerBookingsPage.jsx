import { CheckCircle2, RotateCcw, XCircle } from "lucide-react";
import { useEffect, useState } from "react";
import { bookingApi } from "../api/bookingApi.js";

const labels = {
  pending: "Ожидает ответа",
  approved: "Подтверждена",
  rejected: "Отклонена",
  cancelled: "Отменена клиентом",
  returned: "Возвращено"
};

export default function OwnerBookingsPage() {
  const [bookings, setBookings] = useState([]);
  const [notice, setNotice] = useState("");

  const load = () => bookingApi.ownerList().then(setBookings).catch((error) => setNotice(error.message));

  useEffect(() => {
    load();
  }, []);

  const setStatus = async (id, status) => {
    try {
      await bookingApi.updateStatus(id, status);
      load();
    } catch (error) {
      setNotice(error.message);
    }
  };

  return (
    <section className="panel">
      <div className="panel-head">
        <h2>Заявки владельца</h2>
      </div>
      {notice && <div className="notice">{notice}</div>}
      <div className="booking-list">
        {bookings.map((booking) => (
          <article className="booking-card" key={booking.id}>
            <img src={booking.equipment?.image_url} alt="" />
            <div>
              <strong>{booking.equipment?.name}</strong>
              <span>{booking.user?.name} · {booking.user?.email}</span>
              <small>{formatDate(booking.start_at)} - {formatDate(booking.end_at)}</small>
              {booking.comment && <small>{booking.comment}</small>}
            </div>
            <div className="booking-side">
              <em className={`badge ${booking.status === "pending" ? "wait" : booking.status === "approved" ? "ok" : "busy"}`}>
                {labels[booking.status] ?? booking.status}
              </em>
              {booking.status === "pending" && (
                <div className="row-actions">
                  <button className="table-button" onClick={() => setStatus(booking.id, "approved")}>
                    <CheckCircle2 size={16} />
                    Принять
                  </button>
                  <button className="danger-button" onClick={() => setStatus(booking.id, "rejected")}>
                    <XCircle size={16} />
                    Отклонить
                  </button>
                </div>
              )}
              {booking.status === "approved" && (
                <button className="table-button" onClick={() => setStatus(booking.id, "returned")}>
                  <RotateCcw size={16} />
                  Возврат
                </button>
              )}
            </div>
          </article>
        ))}
      </div>
      {!bookings.length && !notice && <div className="empty-state">Заявки на ваше оборудование появятся здесь.</div>}
    </section>
  );
}

function formatDate(value) {
  return new Date(value).toLocaleString("ru-RU", {
    day: "2-digit",
    month: "2-digit",
    hour: "2-digit",
    minute: "2-digit"
  });
}
