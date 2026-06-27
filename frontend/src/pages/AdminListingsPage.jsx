import { useEffect, useState } from "react";
import { equipmentApi } from "../api/equipmentApi.js";

export default function AdminListingsPage() {
  const [items, setItems] = useState([]);
  const [notice, setNotice] = useState("");
  const load = () => equipmentApi.list().then(setItems).catch((error) => setNotice(error.message));
  useEffect(() => { load(); }, []);
  const hide = async (item) => { try { await equipmentApi.setHidden(item.id, true); load(); } catch (error) { setNotice(error.message); } };
  const remove = async (item) => { if (!window.confirm(`Удалить «${item.name}»?`)) return; try { await equipmentApi.remove(item.id); load(); } catch (error) { setNotice(error.message); } };
  return <section className="panel"><div className="panel-head"><h2>Управление объявлениями</h2></div>{notice && <div className="notice">{notice}</div>}
    <div className="admin-listings">{items.map((item) => <article key={item.id}><img src={item.image_url} alt="" /><div><strong>{item.name}</strong><span>{item.owner?.name} · {item.price_per_day} ₽/сутки</span></div><div className="row-actions"><button className="table-button" onClick={() => hide(item)}>Скрыть</button><button className="danger-button" onClick={() => remove(item)}>Удалить</button></div></article>)}</div>
  </section>;
}
