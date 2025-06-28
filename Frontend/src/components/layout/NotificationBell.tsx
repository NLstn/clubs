import { useState, useRef, useEffect } from "react";
import { useNotifications } from "../../hooks/useNotifications";
import "./NotificationBell.css";

const NotificationBell = () => {
  const { notifications, markRead, deleteNotification } = useNotifications();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className="notification-section" ref={ref}>
      <div className="bell" onClick={() => setOpen(!open)}>
        🔔
      </div>
      {open && (
        <div className="notification-dropdown">
          {notifications.length === 0 && (
            <div className="empty">No notifications</div>
          )}
          {notifications.map((n) => (
            <div
              key={n.id}
              className={"notification-item" + (n.read ? " read" : "")}
            >
              <span>{n.message}</span>
              <div className="actions">
                {!n.read && (
                  <button onClick={() => markRead(n.id, true)}>
                    Mark read
                  </button>
                )}
                <button onClick={() => deleteNotification(n.id)}>Delete</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default NotificationBell;
