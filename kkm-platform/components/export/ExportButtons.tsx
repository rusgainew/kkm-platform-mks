'use client';

import { Download, FileText, Table2 } from 'lucide-react';
import { exportToPDF, exportToCSV, exportToJSON } from '@/lib/export';

export interface ExportButtonsProps {
  data: any[];
  filename?: string;
  title?: string;
}

export function ExportButtons({ data, filename = 'export', title = 'Отчет' }: ExportButtonsProps) {
  const handlePDF = () => {
    exportToPDF(data, { filename: `${filename}.pdf`, title });
  };

  const handleCSV = () => {
    exportToCSV(data, { filename: `${filename}.csv` });
  };

  const handleJSON = () => {
    exportToJSON(data, { filename: `${filename}.json` });
  };

  return (
    <div className="flex gap-2">
      <button
        onClick={handlePDF}
        className="flex items-center gap-2 px-3 py-2 text-sm bg-red-600 text-white rounded hover:bg-red-700 transition"
        title="Скачать PDF"
      >
        <FileText className="w-4 h-4" />
        PDF
      </button>

      <button
        onClick={handleCSV}
        className="flex items-center gap-2 px-3 py-2 text-sm bg-green-600 text-white rounded hover:bg-green-700 transition"
        title="Скачать CSV"
      >
        <Table2 className="w-4 h-4" />
        CSV
      </button>

      <button
        onClick={handleJSON}
        className="flex items-center gap-2 px-3 py-2 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 transition"
        title="Скачать JSON"
      >
        <Download className="w-4 h-4" />
        JSON
      </button>
    </div>
  );
}
