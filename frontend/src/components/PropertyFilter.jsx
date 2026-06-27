import { Search, SlidersHorizontal, ArrowUpDown } from "lucide-react";
import { useState, useMemo } from "react";

const PROPERTY_TYPES = [
  { value: "", label: "Semua Tipe" },
  { value: "house", label: "Rumah" },
  { value: "apartment", label: "Apartemen" },
  { value: "land", label: "Tanah" },
  { value: "shophouse", label: "Ruko" },
  { value: "warehouse", label: "Gudang" },
  { value: "office", label: "Kantor" },
  { value: "villa", label: "Villa" },
];

const LISTING_TYPES = [
  { value: "", label: "Semua" },
  { value: "sale", label: "Dijual" },
  { value: "rent", label: "Disewa" },
];

const SOURCE_TYPES = [
  { value: "", label: "Semua Sumber" },
  { value: "regular", label: "Reguler" },
  { value: "bank_auction", label: "Lelang Bank" },
  { value: "company_asset", label: "Aset Perusahaan" },
];

const SORT_OPTIONS = [
  { value: "created_at-desc", label: "Terbaru" },
  { value: "created_at-asc", label: "Terlama" },
  { value: "price-asc", label: "Harga Terendah" },
  { value: "price-desc", label: "Harga Tertinggi" },
];

export default function PropertyFilter({ onFilter, cities }) {
  const [showMore, setShowMore] = useState(false);
  const [filters, setFilters] = useState({
    q: "",
    property_type: "",
    listing_type: "",
    source_type: "",
    city: "",
    price_min: "",
    price_max: "",
    latitude: "",
    longitude: "",
    radius_km: "10",
  });
  const [sort, setSort] = useState("created_at-desc");

  const activeCount = useMemo(() => {
    let c = 0;
    if (filters.property_type) c++;
    if (filters.listing_type) c++;
    if (filters.source_type) c++;
    if (filters.city) c++;
    if (filters.price_min || filters.price_max) c++;
    if (filters.latitude && filters.longitude) c++;
    return c;
  }, [filters]);

  const handleChange = (key, value) => {
    const updated = { ...filters, [key]: value };
    setFilters(updated);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const clean = {};
    Object.entries(filters).forEach(([k, v]) => {
      if (v !== "" && v !== null) clean[k] = v;
    });
    const [sortKey, sortOrder] = sort.split("-");
    clean.sort = sortKey;
    clean.order = sortOrder;
    onFilter(clean);
  };

  const handleReset = () => {
    setFilters({
      q: "",
      property_type: "",
      listing_type: "",
      source_type: "",
      city: "",
      price_min: "",
      price_max: "",
      latitude: "",
      longitude: "",
      radius_km: "10",
    });
    setSort("created_at-desc");
    onFilter({ sort: "created_at", order: "desc" });
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="bg-white rounded-xl shadow-sm border border-gray-100 p-4 mb-6"
    >
      {/* Search bar */}
      <div className="flex gap-2">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            type="text"
            placeholder="Cari properti..."
            value={filters.q}
            onChange={(e) => handleChange("q", e.target.value)}
            className="input-field pl-9"
          />
        </div>
        <button type="submit" className="btn-primary flex items-center gap-1">
          <Search className="w-4 h-4" /> Cari
        </button>
        <button
          type="button"
          onClick={() => setShowMore(!showMore)}
          className="btn-secondary flex items-center gap-1"
        >
          <SlidersHorizontal className="w-4 h-4" />
          <span className="hidden sm:inline">Filter</span>
        </button>
      </div>

      {/* Advanced filters */}
      {showMore && (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3 mt-3 pt-3 border-t">
          <select
            value={filters.property_type}
            onChange={(e) => handleChange("property_type", e.target.value)}
            className="input-field text-sm"
          >
            {PROPERTY_TYPES.map((t) => (
              <option key={t.value} value={t.value}>
                {t.label}
              </option>
            ))}
          </select>
          <select
            value={filters.listing_type}
            onChange={(e) => handleChange("listing_type", e.target.value)}
            className="input-field text-sm"
          >
            {LISTING_TYPES.map((t) => (
              <option key={t.value} value={t.value}>
                {t.label}
              </option>
            ))}
          </select>
          <select
            value={filters.source_type}
            onChange={(e) => handleChange("source_type", e.target.value)}
            className="input-field text-sm"
          >
            {SOURCE_TYPES.map((t) => (
              <option key={t.value} value={t.value}>
                {t.label}
              </option>
            ))}
          </select>
          <select
            value={filters.city}
            onChange={(e) => handleChange("city", e.target.value)}
            className="input-field text-sm"
          >
            <option value="">Semua Kota</option>
            {cities?.map((c) => (
              <option key={c.city} value={c.city}>
                {c.city}, {c.province}
              </option>
            ))}
          </select>
          <input
            type="number"
            placeholder="Harga Min (Rp)"
            value={filters.price_min}
            onChange={(e) => handleChange("price_min", e.target.value)}
            className="input-field text-sm"
          />
          <input
            type="number"
            placeholder="Harga Max (Rp)"
            value={filters.price_max}
            onChange={(e) => handleChange("price_max", e.target.value)}
            className="input-field text-sm"
          />
          <select
            value={sort}
            onChange={(e) => setSort(e.target.value)}
            className="input-field text-sm"
          >
            <ArrowUpDown className="w-3 h-3 inline mr-1" />
            {SORT_OPTIONS.map((o) => (
              <option key={o.value} value={o.value}>
                {o.label}
              </option>
            ))}
          </select>
          {/* Radius slider */}
          <div className="col-span-full sm:col-span-2 space-y-1">
            <label className="text-xs text-gray-500 flex items-center justify-between">
              <span>📍 Radius Pencarian</span>
              <span className="font-medium text-primary-600">
                {filters.radius_km} km
              </span>
            </label>
            <input
              type="range"
              min="1"
              max="50"
              step="1"
              value={filters.radius_km}
              onChange={(e) => handleChange("radius_km", e.target.value)}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-primary-600"
            />
            <div className="flex justify-between text-[10px] text-gray-400">
              <span>1 km</span>
              <span>25 km</span>
              <span>50 km</span>
            </div>
          </div>
          <div className="col-span-full flex items-center justify-between">
            {activeCount > 0 && (
              <span className="text-xs text-primary-600 bg-primary-50 px-2 py-1 rounded-full">
                {activeCount} filter aktif
              </span>
            )}
            <button
              type="button"
              onClick={handleReset}
              className="text-sm text-gray-500 hover:text-red-500 transition-colors"
            >
              Reset Filter
            </button>
          </div>
        </div>
      )}
    </form>
  );
}
