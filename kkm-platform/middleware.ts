import { NextRequest, NextResponse } from "next/server";
import { locales, defaultLocale, Locale } from "@/i18n/config";

/**
 * Middleware для определения и обработки языка
 * - Проверяет cookies, headers, pathname
 * - Редиректит на правильный locale
 * - Устанавливает NEXT_LOCALE cookie
 */
export function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;

  // Проверяем, начинается ли путь с одного из поддерживаемых локалей
  const pathnameHasLocale = locales.some(
    (locale) => pathname.startsWith(`/${locale}/`) || pathname === `/${locale}`
  );

  if (pathnameHasLocale) {
    return NextResponse.next();
  }

  // Определяем locale из различных источников
  let locale: Locale = defaultLocale;

  // 1. Проверяем cookies
  const localeCookie = request.cookies.get("NEXT_LOCALE")?.value as Locale;
  if (localeCookie && locales.includes(localeCookie)) {
    locale = localeCookie;
  } else {
    // 2. Проверяем Accept-Language header
    const acceptLanguage = request.headers.get("accept-language");
    if (acceptLanguage) {
      const preferredLocale = acceptLanguage
        .split(",")[0]
        .split("-")[0]
        .toLowerCase() as Locale;

      if (locales.includes(preferredLocale)) {
        locale = preferredLocale;
      }
    }
  }

  // Редиректим на URL с locale
  const newPathname = `/${locale}${pathname}`;
  const response = NextResponse.redirect(new URL(newPathname, request.url));

  // Устанавливаем cookie
  response.cookies.set("NEXT_LOCALE", locale, {
    maxAge: 60 * 60 * 24 * 365, // 1 год
    path: "/",
  });

  return response;
}

export const config = {
  matcher: [
    // Пропускаем системные пути
    "/((?!_next|api|favicon.ico|robots.txt|sitemap.xml).*)",
  ],
};
