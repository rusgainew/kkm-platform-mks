import { QueryProvider } from "@/lib/providers/QueryProvider";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <QueryProvider>{children}</QueryProvider>;
}
