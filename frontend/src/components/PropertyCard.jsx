import { Link } from "react-router-dom";
import { MapPin, Bed, Bath, Maximize, Phone } from "lucide-react";

const PROPERTY_LABELS = {
  house: "Rumah",
  land: "Tanah",
  apartment: "Apartemen",
  shophouse: "Ruko",
  warehouse: "Gudang",
  office: "Kantor",
  villa: "Villa",
};

const SOURCE_BADGES = {
  bank_auction: { label: "Lelang Bank", cls: "badge-auction" },
  company_asset: { label: "Aset Perusahaan", cls: "badge-company" },
  regular: null,
};

function formatPrice(price) {
  const p = parseFloat(price);
  if (p >= 1_000_000_000) {
    return `Rp ${(p / 1_000_000_000).toFixed(1)}M`;
  }
  if (p >= 1_000_000) {
    return `Rp ${(p / 1_000_000).toFixed(0)}M`;
  }
  return `Rp ${p.toLocaleString("id-ID")}`;
}

export default function PropertyCard({ property, onSave, onUnsave, saved }) {
  const {
    id,
    title,
    price,
    listing_type,
    property_type,
    source_type,
    city,
    bedrooms,
    bathrooms,
    building_area,
    land_area,
    main_photo_url,
    salesman,
    tenant,
  } = property;

  const badge = SOURCE_BADGES[source_type];
  const typeLabel = PROPERTY_LABELS[property_type] || property_type;
  const waMessage = encodeURIComponent(
    `Halo, saya tertarik dengan properti ${title} yang saya lihat di PropertyHub`,
  );

  return (
    <div className="card group">
      {/* Image */}
      <Link
        to={`/properties/${id}`}
        className="block relative overflow-hidden aspect-[4/3] bg-gray-200"
      >
        <img
          src={main_photo_url || ""}
          alt={title}
          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          loading="lazy"
          onError={(e) => {
            e.target.style.display = "none";
            e.target.nextSibling.style.display = "flex";
          }}
        />
        {/* Watermark overlay */}
        <div className="absolute inset-0 pointer-events-none select-none overflow-hidden opacity-[0.12]">
          <div
            className="absolute -inset-full text-[10px] font-bold text-black whitespace-nowrap"
            style={{
              backgroundImage:
                "repeating-linear-gradient(135deg, transparent, transparent 40px, currentColor 40px, currentColor 42px)",
              WebkitBackgroundClip: "text",
              backgroundClip: "text",
              color: "transparent",
            }}
          >
            {Array.from({ length: 20 }, (_, i) => (
              <span key={i} className="inline-block px-8 py-6 -rotate-45">
                PropertyHub
              </span>
            ))}
          </div>
        </div>
        <div
          className="w-full h-full items-center justify-center text-gray-400"
          style={{ display: main_photo_url ? "none" : "flex" }}
        >
          <Maximize className="w-12 h-12" />
        </div>
        <div className="absolute top-2 left-2 flex gap-1">
          <span
            className={listing_type === "sale" ? "badge-sale" : "badge-rent"}
          >
            {listing_type === "sale" ? "Dijual" : "Disewa"}
          </span>
          {badge && <span className={badge.cls}>{badge.label}</span>}
        </div>
        <span className="absolute top-2 right-2 badge bg-white/90 text-gray-700 text-xs">
          {typeLabel}
        </span>
      </Link>

      {/* Content */}
      <div className="p-4">
        <Link to={`/properties/${id}`}>
          <h3 className="font-semibold text-gray-900 line-clamp-2 mb-1 hover:text-primary-600 transition-colors">
            {title}
          </h3>
        </Link>

        <p className="text-lg font-bold text-primary-700 mb-2">
          {formatPrice(price)}
        </p>

        {city && (
          <p className="flex items-center gap-1 text-sm text-gray-500 mb-2">
            <MapPin className="w-3.5 h-3.5" /> {city}
          </p>
        )}

        <div className="flex items-center gap-4 text-xs text-gray-500 mb-3">
          {bedrooms != null && (
            <span className="flex items-center gap-1">
              <Bed className="w-3.5 h-3.5" /> {bedrooms} KT
            </span>
          )}
          {bathrooms != null && (
            <span className="flex items-center gap-1">
              <Bath className="w-3.5 h-3.5" /> {bathrooms} KM
            </span>
          )}
          {(building_area || land_area) && (
            <span className="flex items-center gap-1">
              <Maximize className="w-3.5 h-3.5" /> {building_area || land_area}{" "}
              m²
            </span>
          )}
        </div>

        {/* Agent + Tenant */}
        <div className="flex items-center justify-between pt-3 border-t border-gray-100">
          <div className="flex items-center gap-2 min-w-0">
            {tenant?.logo_url && (
              <img
                src={tenant.logo_url}
                alt={tenant.name}
                className="w-6 h-6 rounded object-cover"
              />
            )}
            <div className="text-xs text-gray-500 truncate">
              {salesman?.name && (
                <p className="font-medium text-gray-700 truncate">
                  {salesman.name}
                </p>
              )}
              {tenant?.name && <p className="truncate">{tenant.name}</p>}
            </div>
          </div>
          {salesman?.phone && (
            <a
              href={`https://wa.me/${salesman.phone.replace(/[^0-9]/g, "")}?text=${waMessage}`}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-1 text-white bg-green-500 hover:bg-green-600 px-2.5 py-1.5 rounded-full text-xs font-medium transition-colors flex-shrink-0"
            >
              <Phone className="w-3 h-3" /> WA
            </a>
          )}
        </div>

        {/* Save button (buyer only) */}
        {onSave && (
          <button
            onClick={() => (saved ? onUnsave(id) : onSave(id))}
            className={`mt-2 w-full text-xs py-1.5 rounded-lg font-medium transition-colors ${
              saved
                ? "bg-red-50 text-red-600 hover:bg-red-100"
                : "bg-gray-50 text-gray-600 hover:bg-gray-100"
            }`}
          >
            {saved ? "❤️ Hapus dari Favorit" : "🤍 Simpan ke Favorit"}
          </button>
        )}
      </div>
    </div>
  );
}
