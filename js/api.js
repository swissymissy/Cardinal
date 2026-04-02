const API_BASE = "";

async function apiRequest(endpoint, method, body = null, token = null) {
  const headers = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE}${endpoint}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : null,
  });

  return res;
}

// authenticated request with silent token refresh on 401
async function apiRequestAuth(endpoint, method, body = null) {
  let token = getAccessToken();
  let res = await apiRequest(endpoint, method, body, token);

  if (res.status === 401) {
    // try silent refresh
    const refreshToken = getRefreshToken();
    if (!refreshToken) {
      logout();
      return null;
    }

    const refreshRes = await apiRequest("/api/refresh", "POST", null, refreshToken);
    if (!refreshRes.ok) {
      logout();
      return null;
    }

    const refreshData = await refreshRes.json();
    saveAccessToken(refreshData.token);

    // retry original request with new token
    res = await apiRequest(endpoint, method, body, refreshData.token);
  }

  return res;
}
