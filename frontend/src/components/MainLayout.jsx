import { CalendarDays, ClipboardCheck, Crown, FolderTree, Home, LogOut, MessageCircle, Moon, Package, Shield, Sparkles, Sun, UserRound, Users } from "lucide-react";
import { useEffect, useState } from "react";
import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { chatApi } from "../api/chatApi.js";
import { useAuthStore } from "../store/authStore.js";

export default function MainLayout() {
  const navigate = useNavigate();
  const { user, logout, isAuthenticated } = useAuthStore();
  const signedIn = isAuthenticated();
  const [theme, setTheme] = useState(() => localStorage.getItem("theme") || "light");
  const [unreadChats, setUnreadChats] = useState(0);

  useEffect(() => {
    document.documentElement.dataset.theme = theme;
    localStorage.setItem("theme", theme);
  }, [theme]);

  useEffect(() => {
    if (!signedIn) {
      setUnreadChats(0);
      return;
    }
    const loadUnread = () =>
      chatApi
        .list()
        .then((items) => setUnreadChats(items.filter((item) => item.unread_count > 0).length))
        .catch(() => setUnreadChats(0));
    loadUnread();
    const timer = window.setInterval(loadUnread, 15000);
    return () => window.clearInterval(timer);
  }, [signedIn]);

  const handleLogout = () => {
    logout();
    navigate("/");
  };

  const role = roleMeta(user?.role);

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-icon"><Package size={26} /></span>
          <span className="brand-text">Equipment Booking</span>
        </div>

        <nav className="nav-list">
          <NavLink to="/">
            <Home size={18} />
            Главная
          </NavLink>
          {signedIn && (
            <>
              <NavLink to="/bookings">
                <CalendarDays size={18} />
                Мои брони
              </NavLink>
              <NavLink to="/chats">
                <MessageCircle size={18} />
                Сообщения
                {unreadChats > 0 && <span className="nav-badge">{unreadChats}</span>}
              </NavLink>
              <NavLink to="/listings">
                <FolderTree size={18} />
                Мои объявления
              </NavLink>
              <NavLink to="/owner/bookings">
                <Shield size={18} />
                Заявки владельца
              </NavLink>
			  <NavLink to="/account">
				<UserRound size={18} />
				Мой профиль
			  </NavLink>
            </>
          )}
          {user?.role === "admin" && (
            <>
              <NavLink to="/admin/users">
                <Users size={18} />
                Пользователи
              </NavLink>
              <NavLink to="/admin/categories">
                <FolderTree size={18} />
                Категории
              </NavLink>
			  <NavLink to="/admin/listings">
				<Package size={18} />
				Управление объявлениями
			  </NavLink>
            </>
          )}
          {(user?.role === "admin" || user?.role === "moderator") && (
            <NavLink to="/moderation">
              <ClipboardCheck size={18} />
              Модерация
            </NavLink>
          )}
        </nav>
      </aside>

      <section className="workspace">
        <header className="topbar">
          <div>
            <p className="eyebrow">Прокат и учет</p>
            <h1>Бронирование оборудования</h1>
          </div>
          {signedIn ? (
            <div className="auth-actions">
              <div className={`role-chip ${role.className}`}>
                {role.icon}
                <span>{role.label}</span>
              </div>
              <button className="icon-button" aria-label="Сменить тему" onClick={() => setTheme(theme === "dark" ? "light" : "dark")}>
                {theme === "dark" ? <Sun size={18} /> : <Moon size={18} />}
              </button>
              <button className="ghost-button" onClick={handleLogout}>
                <LogOut size={18} />
                Выйти
              </button>
            </div>
          ) : (
            <div className="auth-actions">
              <button className="icon-button" aria-label="Сменить тему" onClick={() => setTheme(theme === "dark" ? "light" : "dark")}>
                {theme === "dark" ? <Sun size={18} /> : <Moon size={18} />}
              </button>
              <button className="ghost-button" onClick={() => navigate("/login")}>
                Войти
              </button>
              <button className="primary-button" onClick={() => navigate("/register")}>
                Регистрация
              </button>
            </div>
          )}
        </header>
        <Outlet />
      </section>
    </div>
  );
}

function roleMeta(role) {
  if (role === "admin") {
    return { label: "Админ", className: "admin", icon: <Crown size={16} /> };
  }
  if (role === "moderator") {
    return { label: "Модератор", className: "moderator", icon: <Shield size={16} /> };
  }
  return { label: "Пользователь", className: "user", icon: <Sparkles size={16} /> };
}
