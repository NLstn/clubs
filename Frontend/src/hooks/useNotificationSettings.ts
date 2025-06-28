import { useEffect, useState } from "react";
import api from "../utils/api";

export interface NotificationSetting {
  id?: string;
  userId?: string;
  clubId?: string;
  email: boolean;
  inApp: boolean;
}

export const useNotificationSettings = () => {
  const [settings, setSettings] = useState<NotificationSetting[]>([]);

  const fetchSettings = async () => {
    try {
      const res = await api.get("/api/v1/notificationSettings");
      setSettings(res.data);
    } catch (err) {
      console.error("Failed to fetch notification settings", err);
    }
  };

  const saveSettings = async (s: NotificationSetting[]) => {
    await api.post("/api/v1/notificationSettings", s);
    setSettings(s);
  };

  useEffect(() => {
    fetchSettings();
  }, []);

  return { settings, setSettings, fetchSettings, saveSettings };
};
