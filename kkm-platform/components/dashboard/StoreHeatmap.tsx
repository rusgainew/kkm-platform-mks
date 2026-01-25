import { MapPin } from 'lucide-react';

interface StoreData {
  id: string;
  name: string;
  revenue: number;
  coordinates: { x: number; y: number };
}

interface StoreHeatmapProps {
  stores: StoreData[];
  title?: string;
}

export default function StoreHeatmap({ stores, title = 'География продаж' }: StoreHeatmapProps) {
  const maxRevenue = Math.max(...stores.map((s) => s.revenue));

  const getHeatColor = (revenue: number) => {
    const intensity = revenue / maxRevenue;
    if (intensity > 0.7) return 'from-emerald-500/80 to-emerald-400/60 shadow-emerald-500/50';
    if (intensity > 0.4) return 'from-blue-500/80 to-blue-400/60 shadow-blue-500/50';
    return 'from-yellow-500/80 to-yellow-400/60 shadow-yellow-500/50';
  };

  const getSize = (revenue: number) => {
    const intensity = revenue / maxRevenue;
    if (intensity > 0.7) return 'w-6 h-6';
    if (intensity > 0.4) return 'w-5 h-5';
    return 'w-4 h-4';
  };

  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 p-6 shadow-xl">
      <h3 className="text-lg font-semibold text-white mb-6 flex items-center gap-2">
        <span className="w-1 h-6 bg-linear-to-b from-blue-400 to-emerald-400 rounded-full"></span>
        {title}
      </h3>

      <div className="relative bg-gray-800/50 rounded-lg border border-gray-700 h-100 overflow-hidden backdrop-blur-sm">
        {/* Grid Pattern */}
        <div className="absolute inset-0 opacity-10">
          <div className="grid grid-cols-10 grid-rows-10 h-full">
            {Array.from({ length: 100 }).map((_, i) => (
              <div key={i} className="border border-blue-500/20"></div>
            ))}
          </div>
        </div>

        {/* Store Markers */}
        {stores.map((store) => (
          <div
            key={store.id}
            className="absolute group cursor-pointer"
            style={{
              left: `${store.coordinates.x}%`,
              top: `${store.coordinates.y}%`,
              transform: 'translate(-50%, -50%)',
            }}
          >
            {/* Glow Effect */}
            <div
              className={`absolute inset-0 rounded-full blur-xl bg-linear-to-br ${getHeatColor(
                store.revenue
              )} animate-pulse`}
            ></div>

            {/* Marker */}
            <div
              className={`relative ${getSize(
                store.revenue
              )} bg-linear-to-br ${getHeatColor(
                store.revenue
              )} rounded-full flex items-center justify-center shadow-lg transition-transform group-hover:scale-125`}
            >
              <MapPin className="w-3 h-3 text-white" />
            </div>

            {/* Tooltip */}
            <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none z-10">
              <div className="bg-gray-950 border border-gray-700 rounded-lg px-3 py-2 shadow-2xl whitespace-nowrap">
                <p className="text-sm font-semibold text-white">{store.name}</p>
                <p className="text-xs text-emerald-400">{store.revenue.toLocaleString('ru-RU')} ₽</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Legend */}
      <div className="mt-4 flex items-center justify-between text-xs text-gray-400">
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 rounded-full bg-linear-to-br from-yellow-500 to-yellow-400"></div>
          <span>Низкая выручка</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 rounded-full bg-linear-to-br from-blue-500 to-blue-400"></div>
          <span>Средняя выручка</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 rounded-full bg-linear-to-br from-emerald-500 to-emerald-400"></div>
          <span>Высокая выручка</span>
        </div>
      </div>
    </div>
  );
}
