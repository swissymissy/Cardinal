function saveTokens(access, refresh) {
  localStorage.setItem("access_token", access);
  localStorage.setItem("refresh_token", refresh);
}

function saveAccessToken(token) {
  localStorage.setItem("access_token", token);
}

function getAccessToken() {
  return localStorage.getItem("access_token");
}

function getRefreshToken() {
  return localStorage.getItem("refresh_token");
}

function saveUserInfo(id, username) {
  localStorage.setItem("current_user_id", id);
  localStorage.setItem("current_username", username);
}

function getCurrentUserID() {
  return localStorage.getItem("current_user_id");
}

function getCurrentUsername() {
  return localStorage.getItem("current_username");
}

async function logout() {
  const refreshToken = getRefreshToken();
  if (refreshToken) {
    // revoke refresh token on server before clearing
    await apiRequest("/api/revoke", "POST", null, refreshToken);
  }
  localStorage.clear();
  window.location.href = "/";
}

function requireAuth() {
  if (!getAccessToken()) {
    window.location.href = "/";
  }
}
