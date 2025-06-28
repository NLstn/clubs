import { useEffect, useState } from "react";
import api from "../utils/api";

export interface Notification {
  id: string;
  userId: string;
  clubId?: string;
  message: string;
  read: boolean;
  createdAt: string;
}

export const useNotifications = () => {
  const [notifications, setNotifications] = useState<Notification[]>([]);

  const fetchNotifications = async () => {
    try {
      const res = await api.get("/api/v1/notifications");
      setNotifications(res.data);
    } catch (err) {
      console.error("Failed to fetch notifications", err);
    }
  };

  const markRead = async (id: string, read: boolean) => {
    await api.put(`/api/v1/notifications/${id}`, { read });
    setNotifications((n) =>
      n.map((not) => (not.id === id ? { ...not, read } : not)),
    );
  };

  const deleteNotification = async (id: string) => {
    await api.delete(`/api/v1/notifications/${id}`);
    setNotifications((n) => n.filter((not) => not.id !== id));
  };

  useEffect(() => {
    fetchNotifications();
  }, []);

  return { notifications, fetchNotifications, markRead, deleteNotification };
};
