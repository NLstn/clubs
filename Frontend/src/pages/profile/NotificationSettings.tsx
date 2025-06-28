import Layout from "../../components/layout/Layout";
import ProfileSidebar from "./ProfileSidebar";
import { useNotificationSettings } from "../../hooks/useNotificationSettings";
import { useState } from "react";

type NotificationSetting = {
  clubId?: string;
  email: boolean;
  inApp: boolean;
  id?: string;
};

const NotificationSettings = () => {
  const { settings, setSettings, saveSettings } = useNotificationSettings<NotificationSetting[]>();
  const [message, setMessage] = useState("");

  const handleToggle = (index: number, field: "email" | "inApp") => {
    const copy = [...settings];
    copy[index][field] = !copy[index][field];
    setSettings(copy);
  };

  const handleSave = async () => {
    try {
      await saveSettings(settings);
      setMessage("Saved");
      setTimeout(() => setMessage(""), 2000);
    } catch {
      setMessage("Failed to save");
    }
  };

  return (
    <Layout title="Notification Settings">
      <div style={{ display: "flex", minHeight: "calc(100vh - 90px)" }}>
        <ProfileSidebar />
        <div
          style={{
            flex: "1 1 auto",
            padding: "20px",
            maxWidth: "calc(100% - 200px)",
          }}
        >
          <h2>Notification Settings</h2>
          {message && <div>{message}</div>}
          <table>
            <thead>
              <tr>
                <th>Club</th>
                <th>Email</th>
                <th>In App</th>
              </tr>
            </thead>
            <tbody>
              {settings.map((s, idx) => (
                <tr key={s.id || idx}>
                  <td>{s.clubId || "All Clubs"}</td>
                  <td>
                    <input
                      type="checkbox"
                      checked={s.email}
                      onChange={() => handleToggle(idx, "email")}
                    />
                  </td>
                  <td>
                    <input
                      type="checkbox"
                      checked={s.inApp}
                      onChange={() => handleToggle(idx, "inApp")}
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <button onClick={handleSave} style={{ marginTop: "10px" }}>
            Save
          </button>
        </div>
      </div>
    </Layout>
  );
};

export default NotificationSettings;
