/**
 * Document formatting utilities
 */

export const formatDate = (timestamp: string | number): string => {
  try {
    const date =
      typeof timestamp === "string"
        ? new Date(timestamp)
        : new Date(timestamp * 1000);
    return new Intl.DateTimeFormat("ru-RU", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date);
  } catch {
    return "Неверная дата";
  }
};

export const formatDocumentId = (id: string): string => {
  return id.substring(0, 8).toUpperCase();
};

export const truncateText = (text: string, length: number = 50): string => {
  return text.length > length ? `${text.substring(0, length)}...` : text;
};
