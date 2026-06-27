import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { equipmentApi } from "../api/equipmentApi.js";
import { demoEquipment } from "../data/demoEquipment.js";

export default function EquipmentPage() {
	const [equipment, setEquipment] = useState([]);
	const [error, setError] = useState("");
	const [search, setSearch] = useState("");

  useEffect(() => {
    equipmentApi
      .list()
      .then((items) => setEquipment(items.length ? items : demoEquipment))
      .catch((err) => {
        setError(err.message);
        setEquipment(demoEquipment);
      });
  }, []);

	const found = equipment.filter((item) => item.name.toLowerCase().includes(search.trim().toLowerCase()));

	return (
    <section className="panel">
      <div className="panel-head">
        <h2>Оборудование</h2>
		<input placeholder="Поиск по названию" aria-label="Поиск оборудования" value={search} onChange={(event) => setSearch(event.target.value)} />
      </div>
      {error && <div className="notice">{error}</div>}
      <div className="card-grid">
		{found.map((item) => (
          <Link className="listing-card" to={`/equipment/${item.id}`} key={item.id}>
            <img src={item.image_url || demoEquipment[0].image_url} alt={item.name} />
            <strong>{item.name}</strong>
            <span>{Number(item.price_per_day || 0).toLocaleString("ru-RU")} ₽/сутки</span>
            <small>{item.location || item.category?.name || "Город не указан"}</small>
            <em className={item.available ? "badge ok" : "badge busy"}>{item.available ? "Доступно" : "Занято"}</em>
          </Link>
        ))}
      </div>
		{!found.length && <div className="empty-state">По вашему запросу ничего не найдено.</div>}
    </section>
  );
}
