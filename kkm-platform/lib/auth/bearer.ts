/**
 * Get authorization header with bearer token
 * @returns Object with Authorization header or empty object if no token
 */
export function bearerAuth(): Record<string, string> {
  if (typeof window === "undefined") {
    console.log("[bearerAuth] SSR environment, no localStorage");
    return {};
  }

  try {
    // Try direct token first
    let token =
      localStorage.getItem("accessToken") ||
      localStorage.getItem("authToken") ||
      localStorage.getItem("token");

    if (token) {
      console.log(
        "[bearerAuth] Found token in localStorage.token, length:",
        token.length,
      );
      return {
        Authorization: `Bearer ${token}`,
      };
    }

    const readTokenFromPersistedState = (raw: string | null) => {
      if (!raw) return null;
      const auth = JSON.parse(raw);
      return (
        auth?.state?.tokens?.access_token ||
        auth?.state?.tokens?.accessToken ||
        auth?.state?.token ||
        auth?.token ||
        null
      );
    };

    // Fallback: try auth-storage/auth-store (Zustand persist)
    token =
      readTokenFromPersistedState(localStorage.getItem("auth-storage")) ||
      readTokenFromPersistedState(localStorage.getItem("auth-store"));

    if (token) {
      console.log(
        "[bearerAuth] Found token in persisted store, length:",
        token.length,
      );
      return {
        Authorization: `Bearer ${token}`,
      };
    }

    console.warn("[bearerAuth] No token found in localStorage");
    console.log("[bearerAuth] Available keys:", Object.keys(localStorage));
    return {};
  } catch (error) {
    console.error("[bearerAuth] Error getting auth token:", error);
    return {};
  }
}
