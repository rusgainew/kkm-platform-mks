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
    // Try direct token first (kkm-platform stores it as localStorage.token)
    let token = localStorage.getItem("token");

    if (token) {
      console.log(
        "[bearerAuth] Found token in localStorage.token, length:",
        token.length
      );
      return {
        Authorization: `Bearer ${token}`,
      };
    }

    // Fallback: try auth-storage (from users.ts or other stores)
    const authStorage = localStorage.getItem("auth-storage");
    if (authStorage) {
      const auth = JSON.parse(authStorage);
      token =
        auth?.state?.tokens?.access_token ||
        auth?.state?.tokens?.accessToken ||
        auth?.state?.token ||
        auth?.token;

      if (token) {
        console.log(
          "[bearerAuth] Found token in auth-storage, length:",
          token.length
        );
        return {
          Authorization: `Bearer ${token}`,
        };
      }
    }

    console.warn("[bearerAuth] No token found in localStorage");
    console.log("[bearerAuth] Available keys:", Object.keys(localStorage));
    return {};
  } catch (error) {
    console.error("[bearerAuth] Error getting auth token:", error);
    return {};
  }
}
