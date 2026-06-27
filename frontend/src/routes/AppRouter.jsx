import { Navigate, Route, Routes } from "react-router-dom";
import MainLayout from "../components/MainLayout.jsx";
import AdminCategoriesPage from "../pages/AdminCategoriesPage.jsx";
import AdminUsersPage from "../pages/AdminUsersPage.jsx";
import EquipmentDetailPage from "../pages/EquipmentDetailPage.jsx";
import HomePage from "../pages/HomePage.jsx";
import LoginPage from "../pages/LoginPage.jsx";
import MyBookingsPage from "../pages/MyBookingsPage.jsx";
import MyListingsPage from "../pages/MyListingsPage.jsx";
import OwnerBookingsPage from "../pages/OwnerBookingsPage.jsx";
import ChatPage from "../pages/ChatPage.jsx";
import ProfilePage from "../pages/ProfilePage.jsx";
import AccountPage from "../pages/AccountPage.jsx";
import AdminListingsPage from "../pages/AdminListingsPage.jsx";
import ModerationPage from "../pages/ModerationPage.jsx";
import RegisterPage from "../pages/RegisterPage.jsx";
import ProtectedRoute from "./ProtectedRoute.jsx";

export default function AppRouter() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route element={<MainLayout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/equipment" element={<Navigate to="/" replace />} />
        <Route path="/equipment/:id" element={<EquipmentDetailPage />} />
        <Route path="/profiles/:id" element={<ProfilePage />} />
        <Route element={<ProtectedRoute />}>
          <Route path="/bookings" element={<MyBookingsPage />} />
          <Route path="/listings" element={<MyListingsPage />} />
          <Route path="/owner/bookings" element={<OwnerBookingsPage />} />
          <Route path="/chats" element={<ChatPage />} />
          <Route path="/chats/:id" element={<ChatPage />} />
		  <Route path="/account" element={<AccountPage />} />
        </Route>
        <Route element={<ProtectedRoute roles={["admin"]} />}>
          <Route path="/admin/users" element={<AdminUsersPage />} />
          <Route path="/admin/categories" element={<AdminCategoriesPage />} />
		  <Route path="/admin/listings" element={<AdminListingsPage />} />
        </Route>
        <Route element={<ProtectedRoute roles={["admin", "moderator"]} />}>
          <Route path="/moderation" element={<ModerationPage />} />
        </Route>
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
