/**
 * Document table row component
 */

import { Download, Eye, Edit, Trash2, Archive, Send, CheckCircle, XCircle, MoreVertical } from "lucide-react";
import { useState } from "react";
import { Document } from "@/lib/api/documents";
import { STATUS_LABELS, STATUS_COLORS } from "@/lib/documents/constants";
import { formatDate, formatDocumentId } from "@/lib/documents/formatting";

interface DocumentRowProps {
  document: Document;
  onView?: (doc: Document) => void;
  onDownload?: (doc: Document) => void;
  onEdit?: (doc: Document) => void;
  onDelete?: (doc: Document) => void;
  onArchive?: (doc: Document) => void;
  onSend?: (doc: Document) => void;
  onApprove?: (doc: Document) => void;
  onReject?: (doc: Document) => void;
}

export function DocumentRow({ 
  document, 
  onView, 
  onDownload,
  onEdit,
  onDelete,
  onArchive,
  onSend,
  onApprove,
  onReject
}: DocumentRowProps) {
  const [showMore, setShowMore] = useState(false);
  
  const statusColor = STATUS_COLORS[document.status as keyof typeof STATUS_COLORS] || STATUS_COLORS.draft;
  const statusLabel = STATUS_LABELS[document.status as keyof typeof STATUS_LABELS] || document.status;

  // Determine which actions are available based on status
  const canEdit = document.status === "draft";
  const canSend = document.status === "draft";
  const canApprove = document.status === "sent";
  const canReject = document.status === "sent";
  const canArchive = document.status !== "archived";
  const canDelete = document.status === "draft";

  return (
    <tr className="hover:bg-gray-800/50 transition-colors">
      <td className="px-6 py-4">
        <div>
          <p className="text-white font-medium">{document.title}</p>
          <p className="text-xs text-gray-500">{formatDocumentId(document.id)}</p>
        </div>
      </td>
      <td className="px-6 py-4">
        <span className="text-sm text-gray-300">{document.organization_id}</span>
      </td>
      <td className="px-6 py-4">
        <span className={`inline-flex px-3 py-1 rounded-full text-xs font-semibold ${statusColor}`}>
          {statusLabel}
        </span>
      </td>
      <td className="px-6 py-4">
        <span className="text-sm text-gray-300">{document.created_by || "N/A"}</span>
      </td>
      <td className="px-6 py-4">
        <span className="text-sm text-gray-400">{formatDate(document.created_at)}</span>
      </td>
      <td className="px-6 py-4">
        <div className="flex items-center justify-end gap-1 relative">
          {/* Primary actions */}
          <button
            onClick={() => onView?.(document)}
            className="p-2 hover:bg-gray-800 rounded text-gray-400 hover:text-white transition"
            title="Просмотр"
          >
            <Eye className="w-4 h-4" />
          </button>
          
          {canSend && (
            <button
              onClick={() => onSend?.(document)}
              className="p-2 hover:bg-gray-800 rounded text-blue-400 hover:text-blue-300 transition"
              title="Отправить на одобрение"
            >
              <Send className="w-4 h-4" />
            </button>
          )}

          {canApprove && (
            <button
              onClick={() => onApprove?.(document)}
              className="p-2 hover:bg-gray-800 rounded text-green-400 hover:text-green-300 transition"
              title="Одобрить"
            >
              <CheckCircle className="w-4 h-4" />
            </button>
          )}

          {canReject && (
            <button
              onClick={() => onReject?.(document)}
              className="p-2 hover:bg-gray-800 rounded text-red-400 hover:text-red-300 transition"
              title="Отклонить"
            >
              <XCircle className="w-4 h-4" />
            </button>
          )}

          {/* More menu */}
          <div className="relative">
            <button
              onClick={() => setShowMore(!showMore)}
              className="p-2 hover:bg-gray-800 rounded text-gray-400 hover:text-white transition"
              title="Дополнительно"
            >
              <MoreVertical className="w-4 h-4" />
            </button>

            {showMore && (
              <div className="absolute right-0 mt-2 w-48 bg-gray-800 border border-gray-700 rounded-lg shadow-lg z-10">
                <button
                  onClick={() => {
                    onDownload?.(document);
                    setShowMore(false);
                  }}
                  className="w-full flex items-center gap-2 px-4 py-2 text-gray-300 hover:bg-gray-700 hover:text-white transition rounded-t-lg"
                >
                  <Download className="w-4 h-4" />
                  Скачать
                </button>

                {canEdit && (
                  <button
                    onClick={() => {
                      onEdit?.(document);
                      setShowMore(false);
                    }}
                    className="w-full flex items-center gap-2 px-4 py-2 text-gray-300 hover:bg-gray-700 hover:text-white transition"
                  >
                    <Edit className="w-4 h-4" />
                    Редактировать
                  </button>
                )}

                {canArchive && (
                  <button
                    onClick={() => {
                      onArchive?.(document);
                      setShowMore(false);
                    }}
                    className="w-full flex items-center gap-2 px-4 py-2 text-gray-300 hover:bg-gray-700 hover:text-white transition"
                  >
                    <Archive className="w-4 h-4" />
                    В архив
                  </button>
                )}

                {canDelete && (
                  <button
                    onClick={() => {
                      onDelete?.(document);
                      setShowMore(false);
                    }}
                    className="w-full flex items-center gap-2 px-4 py-2 text-red-400 hover:bg-gray-700 hover:text-red-300 transition rounded-b-lg"
                  >
                    <Trash2 className="w-4 h-4" />
                    Удалить
                  </button>
                )}
              </div>
            )}
          </div>
        </div>
      </td>
    </tr>
  );
}
