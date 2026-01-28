import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "../globals.css";
import { QueryProvider } from "@/lib/providers/QueryProvider";
import { ToastProvider } from "@/lib/providers/ToastProvider";
import AuthInitializer from "@/components/auth/AuthInitializer";
import ErrorBoundary from "@/components/error/ErrorBoundary";
import { RateLimitWarningBanner } from "@/components/rate-limit/RateLimitWarningBanner";
import { RateLimitProvider } from "@/lib/providers/RateLimitProvider";

const inter = Inter({
  subsets: ["latin", "cyrillic"],
  variable: "--font-inter",
});

export const metadata: Metadata = {
  title: "Касса POS - Централизованная система",
  description: "Централизованный кассовый терминал (POS-система)",
};

export default async function LocaleLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}>) {
  const { locale } = await params;

  return (
    <html lang={locale}>
      <body className={`${inter.variable} font-sans antialiased`}>
        <ErrorBoundary>
          <QueryProvider>
            <ToastProvider />
            <RateLimitProvider>
              <RateLimitWarningBanner />
              <AuthInitializer>
                {children}
              </AuthInitializer>
            </RateLimitProvider>
          </QueryProvider>
        </ErrorBoundary>
      </body>
    </html>
  );
}
