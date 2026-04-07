function initNav() {
  const username = getCurrentUsername() || "Profile";
  const userID = getCurrentUserID();

  const nav = document.createElement("nav");
  nav.className = "nav";
  nav.innerHTML = `
    <a href="/dashboard.html" class="nav-logo">Cardinal</a>
    <a href="/dashboard.html" class="nav-link">Feed</a>
    <div class="nav-right">
      <a href="/profile.html?user=${userID}" class="nav-link">@${escapeHTML(username)}</a>
      <div class="notif-bell-wrapper" id="notifBellWrapper">
        <button class="notif-bell" id="notifBell" title="Notifications">🔔</button>
        <span class="notif-badge" id="notifBadge" style="display:none">0</span>
        <div class="notif-dropdown" id="notifDropdown"></div>
      </div>
      <button class="btn btn-ghost btn-sm" onclick="logout()">Logout</button>
    </div>
  `;

  document.body.prepend(nav);
  initNotificationBell();
}
