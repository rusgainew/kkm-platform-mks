export function getRedirectPathByRole(role?: string): string {
  if (role === "store_manager") {
    return "/manager-dashboard";
  } else if (role === "cashier") {
    return "/pos";
  } else if (role === "network_admin" || role === "admin") {
    return "/dashboard";
  }
  return "/dashboard";
}
