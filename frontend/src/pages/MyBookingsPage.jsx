import { CalendarClock, XCircle } from "lucide-react";
import { useEffect, useState } from "react";
import { bookingApi } from "../api/bookingApi.js";

const statusText = {
  pending: "Ожидает ответа",
  approved: "Подтверждена",
  rejected: "Отклонена",
  cancelled: "Отменена",
  returned: "Завершена"
};

export default function MyBookingsPage() {
  const [bookings, setBookings] = useState([]);
  const [error, setError] = useState("");

  const load = () => bookingApi.list().then(setBookings).catch((err) => setError(err.message));

  useEffect(() => {
    load();
  }, []);

  const cancel = async (id) => {
    try {
      await bookingApi.cancel(id);
      load();
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <section className="panel">
      <div className="panel-head">
        <h2>Мои брони</h2>
      </div>
      {error && <div className="notice">{error}</div>}
      <div className="booking-list">
        {bookings.map((booking) => (
          <article className="booking-card" key={booking.id}>
            <img src={booking.equipment?.image_url} alt="" />
            <div>
              <strong>{booking.equipment?.name ?? `Заявка #${booking.id}`}</strong>
              <span>
                <CalendarClock size={16} />
                {formatDate(booking.start_at)} - {formatDate(booking.end_at)}
              </span>
              {booking.comment && <small>{booking.comment}</small>}
            </div>
            <div className="booking-side">
              <em className={`badge ${booking.status === "pending" ? "wait" : booking.status === "approved" ? "ok" : "busy"}`}>
                {statusText[booking.status] ?? booking.status}
              </em>
              {(booking.status === "pending" || booking.status === "approved") && (
                <button className="table-button" onClick={() => cancel(booking.id)}>
                  <XCircle size={16} />
                  Отменить
                </button>
              )}
            </div>
          </article>
        ))}
      </div>
      {!bookings.length && !error && <div className="empty-state">Брони появятся здесь.</div>}
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
