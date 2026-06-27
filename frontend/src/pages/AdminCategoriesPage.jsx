import { useEffect, useState } from "react";
import { categoryApi } from "../api/categoryApi.js";

export default function AdminCategoriesPage() {
  const [categories, setCategories] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    categoryApi.list().then(setCategories).catch((err) => setError(err.message));
  }, []);

  return (
    <section className="panel">
      <div className="panel-head">
        <h2>Категории</h2>
      </div>
      {error && <div className="notice">{error}</div>}
      {categories.map((category) => (
        <div className="list-item" key={category.id}>
          <strong>{category.name}</strong>
          <span>{category.description}</span>
        </div>
      ))}
      {!categories.length && !error && <div className="empty-state">Категории пока не созданы.</div>}
    </section>
  );
}
