async function fetchNotifications() {
  const res = await apiRequestAuth("/api/notifications", "GET");
  if (!res || !res.ok) return [];
  return await res.json();
}

function getUnreadCount(notifications) {
  return notifications.filter(n => !n.is_read).length;
}

function formatTime(isoString) {
  const date = new Date(isoString);
  const now = new Date();
  const diff = Math.floor((now - date) / 1000);
  if (diff < 60) return `${diff}s ago`;
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  return date.toLocaleDateString();
}

function renderDropdown(notifications) {
  const dropdown = document.getElementById("notifDropdown");
  if (!dropdown) return;

  const unread = getUnreadCount(notifications);

  let html = `
    <div class="notif-dropdown-header">
      <span>Notifications</span>
      ${unread > 0 ? `<button class="btn btn-sm btn-outline" onclick="markAllRead()">✓ Mark all read</button>` : ""}
    </div>
  `;

  if (notifications.length === 0) {
    html += `<div class="notif-empty">No notifications yet</div>`;
  } else {
    for (const n of notifications) {
      html += `
        <div class="notif-item ${n.is_read ? "" : "unread"}" style="cursor:pointer" onclick="goToChirp('${n.chirp_id}', '${n.notif_id}', ${n.is_read})">
          <div class="notif-item-body">
            <div>${n.body}</div>
            <div class="notif-item-time">${formatTime(n.created_at)}</div>
          </div>
          ${!n.is_read ? `<button class="btn btn-sm btn-ghost" onclick="event.stopPropagation(); markOneRead('${n.notif_id}')">✓</button>` : ""}
        </div>
      `;
    }
  }

  dropdown.innerHTML = html;
}

function updateBadge(count) {
  const badge = document.getElementById("notifBadge");
  if (!badge) return;
  if (count > 0) {
    badge.textContent = count > 99 ? "99+" : count;
    badge.style.display = "block";
  } else {
    badge.style.display = "none";
  }
}

async function refreshNotifications() {
  const notifications = await fetchNotifications();
  updateBadge(getUnreadCount(notifications));
  renderDropdown(notifications);
}

async function markAllRead() {
  await apiRequestAuth("/api/notifications", "PUT");
  await refreshNotifications();
}

async function markOneRead(notifID) {
  await apiRequestAuth(`/api/notifications/${notifID}`, "PUT");
  await refreshNotifications();
}

async function goToChirp(chirpID, notifID, isRead) {
  if (!isRead) {
    await apiRequestAuth(`/api/notifications/${notifID}`, "PUT");
  }
  window.location.href = `/dashboard.html#chirp-${chirpID}`;
}

function initNotificationBell() {
  const bell = document.getElementById("notifBell");
  const dropdown = document.getElementById("notifDropdown");
  const wrapper = document.getElementById("notifBellWrapper");
  if (!bell || !dropdown || !wrapper) return;

  // toggle dropdown on bell click
  bell.addEventListener("click", async (e) => {
    e.stopPropagation();
    const isOpen = dropdown.classList.contains("open");
    if (!isOpen) {
      dropdown.classList.add("open");
      await refreshNotifications();
    } else {
      dropdown.classList.remove("open");
    }
  });

  // close when clicking outside
  document.addEventListener("click", (e) => {
    if (!wrapper.contains(e.target)) {
      dropdown.classList.remove("open");
    }
  });

  // poll badge count every 60 seconds
  refreshNotifications();
  setInterval(refreshNotifications, 60000);
}
