import { MapPin, ShieldCheck, UserRound } from "lucide-react";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { userApi } from "../api/userApi.js";
import { demoEquipment } from "../data/demoEquipment.js";

export default function ProfilePage() {
  const { id } = useParams();
  const [profile, setProfile] = useState(null);

  useEffect(() => {
    userApi
      .profile(id)
      .then(setProfile)
      .catch(() => setProfile(demoEquipment.find((item) => String(item.owner.id) === String(id))?.owner ?? demoEquipment[0].owner));
  }, [id]);

  if (!profile) {
    return <div className="empty-state">Загрузка профиля...</div>;
  }

  return (
    <section className="profile-page">
      <div className="profile-card">
        <img src={profile.avatar_url || demoEquipment[0].owner.avatar_url} alt="" />
        <div>
          <h2>{profile.name}</h2>
          <p>
            <MapPin size={16} />
            {profile.city || "Город не указан"}
          </p>
          <p>
            <UserRound size={16} />
            {profile.role === "admin" ? "Администратор" : "Пользователь"}
          </p>
        </div>
      </div>
      <section className="panel">
        <div className="panel-head">
          <h2>О продавце</h2>
        </div>
        <p className="description">{profile.bio || "Пользователь пока не добавил описание профиля."}</p>
        <div className="safe-note inline">
          <ShieldCheck size={18} />
          В публичном профиле не отображается пароль и служебные данные учетной записи.
        </div>
      </section>
    </section>
  );
}
