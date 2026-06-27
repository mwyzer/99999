import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { publicAPI, buyerAPI } from "../api";
import { useAuth } from "../context/AuthContext";
import Loading from "../components/Loading";
import {
  MapPin,
  Bed,
  Bath,
  Maximize,
  Building2,
  Phone,
  Calendar,
  Ruler,
  Heart,
  X,
  ChevronLeft,
  ChevronRight,
  MessageSquare,
  Send,
} from "lucide-react";
import toast from "react-hot-toast";

function formatPrice(p) {
  const price = parseFloat(p);
  return `Rp ${price.toLocaleString("id-ID")}`;
}

export default function PropertyDetailPage() {
  const { id } = useParams();
  const { isAuthenticated, role } = useAuth();
  const [property, setProperty] = useState(null);
  const [loading, setLoading] = useState(true);
  const [activePhoto, setActivePhoto] = useState(0);
  const [saved, setSaved] = useState(false);
  const [saving, setSaving] = useState(false);
  const [lightbox, setLightbox] = useState(false);
  const [inquiryOpen, setInquiryOpen] = useState(false);
  const [inquiryMsg, setInquiryMsg] = useState("");
  const [sendingInquiry, setSendingInquiry] = useState(false);

  useEffect(() => {
    (async () => {
      try {
        const { data } = await publicAPI.getProperty(id);
        setProperty(data.data);
        // Check if saved (for buyers)
        if (isAuthenticated && role === "buyer") {
          try {
            const savedRes = await buyerAPI.listSaved({ per_page: 100 });
            const ids = (savedRes.data.data || []).map((s) => s.id);
            setSaved(ids.includes(data.data.id));
          } catch {
            // Non-critical
          }
        }
      } catch {
        toast.error("Properti tidak ditemukan.");
      } finally {
        setLoading(false);
      }
    })();
  }, [id]);

  if (loading) return <Loading fullScreen />;
  if (!property)
    return (
      <div className="text-center py-20 text-gray-400">
        Properti tidak ditemukan.
      </div>
    );

  const handleToggleSave = async () => {
    if (!isAuthenticated || role !== "buyer") {
      toast.error("Silakan login sebagai buyer untuk menyimpan properti.");
      return;
    }
    setSaving(true);
    try {
      if (saved) {
        await buyerAPI.unsave(property.id);
        setSaved(false);
        toast.success("Dihapus dari favorit.");
      } else {
        await buyerAPI.save(property.id);
        setSaved(true);
        toast.success("Disimpan ke favorit!");
      }
    } catch (err) {
      toast.error(err.response?.data?.error?.message || "Gagal.");
    } finally {
      setSaving(false);
    }
  };

  const handleSendInquiry = async (e) => {
    e.preventDefault();
    if (!isAuthenticated || role !== "buyer") {
      toast.error("Silakan login sebagai buyer untuk mengirim pertanyaan.");
      return;
    }
    if (!inquiryMsg.trim() || inquiryMsg.trim().length < 3) {
      toast.error("Pertanyaan minimal 3 karakter.");
      return;
    }
    setSendingInquiry(true);
    try {
      await buyerAPI.createInquiry({
        property_id: property.id,
        message: inquiryMsg.trim(),
      });
      toast.success("Pertanyaan berhasil dikirim!");
      setInquiryMsg("");
      setInquiryOpen(false);
    } catch (err) {
      toast.error(err.response?.data?.error?.message || "Gagal mengirim.");
    } finally {
      setSendingInquiry(false);
    }
  };

  const waMessage = encodeURIComponent(
    `Halo, saya tertarik dengan properti ${property.title} yang saya lihat di PropertyHub`,
  );
  const waNumber = property.salesman?.phone?.replace(/[^0-9]/g, "") || "";

  const photos = property.photos || [];

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Photo Gallery */}
      {photos.length > 0 && (
        <div className="mb-8">
          <div
            className="aspect-[16/9] rounded-xl overflow-hidden bg-gray-200 mb-3 cursor-pointer relative group"
            onClick={() => setLightbox(true)}
          >
            <img
              src={
                photos[activePhoto]?.watermarked_url ||
                photos[activePhoto]?.medium_url ||
                ""
              }
              alt={property.title}
              className="w-full h-full object-cover"
              onError={(e) => {
                e.target.style.display = "none";
                const fb = document.getElementById("photo-fallback");
                if (fb) fb.style.display = "flex";
              }}
            />
            {/* Watermark overlay */}
            <div className="absolute inset-0 pointer-events-none select-none overflow-hidden opacity-[0.10]">
              <div
                className="absolute -inset-full text-[13px] font-bold text-white"
                style={{
                  backgroundImage:
                    "repeating-linear-gradient(135deg, transparent, transparent 50px, white 50px, white 52px)",
                  WebkitBackgroundClip: "text",
                  backgroundClip: "text",
                  color: "transparent",
                }}
              >
                {Array.from({ length: 30 }, (_, i) => (
                  <span key={i} className="inline-block px-10 py-8 -rotate-45">
                    PropertyHub
                  </span>
                ))}
              </div>
            </div>
            <div
              id="photo-fallback"
              className="w-full h-full items-center justify-center text-gray-400 bg-gray-200"
              style={{ display: "none" }}
            >
              <Maximize className="w-16 h-16" />
            </div>
            <div className="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-colors flex items-center justify-center">
              <span className="text-white text-sm bg-black/50 px-3 py-1 rounded-full opacity-0 group-hover:opacity-100 transition-opacity">
                Klik untuk memperbesar
              </span>
            </div>
          </div>
          {photos.length > 1 && (
            <div className="flex gap-2 overflow-x-auto pb-2">
              {photos.map((photo, idx) => (
                <button
                  key={photo.id}
                  onClick={() => setActivePhoto(idx)}
                  className={`flex-shrink-0 w-20 h-16 rounded-lg overflow-hidden border-2 transition-colors relative ${
                    idx === activePhoto
                      ? "border-primary-600"
                      : "border-transparent"
                  }`}
                >
                  <img
                    src={photo.thumbnail_url || photo.watermarked_url || ""}
                    alt=""
                    className="w-full h-full object-cover"
                    onError={(e) => {
                      e.target.style.display = "none";
                    }}
                  />
                  <div className="absolute inset-0 pointer-events-none select-none opacity-[0.12]">
                    <span className="absolute inset-0 flex items-center justify-center text-[6px] font-bold text-black -rotate-45">
                      PH
                    </span>
                  </div>
                </button>
              ))}
            </div>
          )}
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Content */}
        <div className="lg:col-span-2">
          <div className="flex flex-wrap gap-2 mb-3">
            <span
              className={
                property.listing_type === "sale" ? "badge-sale" : "badge-rent"
              }
            >
              {property.listing_type === "sale" ? "Dijual" : "Disewa"}
            </span>
            {property.source_type === "bank_auction" && (
              <span className="badge-auction">Lelang Bank</span>
            )}
            {property.source_type === "company_asset" && (
              <span className="badge-company">Aset Perusahaan</span>
            )}
          </div>

          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
            {property.title}
          </h1>
          <p className="text-2xl font-bold text-primary-700 mb-4">
            {formatPrice(property.price)}
          </p>

          {property.city && (
            <p className="flex items-center gap-1 text-gray-500 mb-4">
              <MapPin className="w-4 h-4" /> {property.city},{" "}
              {property.province}
            </p>
          )}

          {/* Specs */}
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-6 p-4 bg-gray-50 rounded-xl">
            {property.bedrooms != null && (
              <div className="flex items-center gap-2 text-sm">
                <Bed className="w-4 h-4 text-gray-400" />
                <span>{property.bedrooms} Kamar Tidur</span>
              </div>
            )}
            {property.bathrooms != null && (
              <div className="flex items-center gap-2 text-sm">
                <Bath className="w-4 h-4 text-gray-400" />
                <span>{property.bathrooms} Kamar Mandi</span>
              </div>
            )}
            {property.land_area && (
              <div className="flex items-center gap-2 text-sm">
                <Ruler className="w-4 h-4 text-gray-400" />
                <span>LT {property.land_area} m²</span>
              </div>
            )}
            {property.building_area && (
              <div className="flex items-center gap-2 text-sm">
                <Maximize className="w-4 h-4 text-gray-400" />
                <span>LB {property.building_area} m²</span>
              </div>
            )}
            {property.floors != null && (
              <div className="flex items-center gap-2 text-sm">
                <Building2 className="w-4 h-4 text-gray-400" />
                <span>{property.floors} Lantai</span>
              </div>
            )}
            {property.certificate_type && (
              <div className="flex items-center gap-2 text-sm">
                <Calendar className="w-4 h-4 text-gray-400" />
                <span>Sertifikat: {property.certificate_type}</span>
              </div>
            )}
          </div>

          {/* Description */}
          {property.description && (
            <div className="mb-6">
              <h3 className="font-semibold text-lg mb-2">Deskripsi</h3>
              <p className="text-gray-600 whitespace-pre-line">
                {property.description}
              </p>
            </div>
          )}

          {/* Facilities */}
          {property.facilities &&
            Object.keys(property.facilities).length > 0 && (
              <div>
                <h3 className="font-semibold text-lg mb-2">Fasilitas</h3>
                <div className="flex flex-wrap gap-2">
                  {Object.entries(property.facilities).map(([key, val]) => {
                    if (typeof val === "boolean" && val) {
                      return (
                        <span
                          key={key}
                          className="badge bg-green-100 text-green-700"
                        >
                          {key.replace(/_/g, " ")}
                        </span>
                      );
                    }
                    if (typeof val === "string") {
                      return (
                        <span
                          key={key}
                          className="badge bg-blue-100 text-blue-700"
                        >
                          {key.replace(/_/g, " ")}: {val}
                        </span>
                      );
                    }
                    return null;
                  })}
                </div>
              </div>
            )}

          {/* Lightbox */}
          {lightbox && photos.length > 0 && (
            <div
              className="fixed inset-0 z-50 bg-black/90 flex items-center justify-center"
              onClick={() => setLightbox(false)}
            >
              <button
                onClick={() => setLightbox(false)}
                className="absolute top-4 right-4 text-white bg-black/50 p-2 rounded-full hover:bg-black/70 z-10"
              >
                <X className="w-6 h-6" />
              </button>
              {photos.length > 1 && (
                <>
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      setActivePhoto((prev) =>
                        prev > 0 ? prev - 1 : photos.length - 1,
                      );
                    }}
                    className="absolute left-4 text-white bg-black/50 p-3 rounded-full hover:bg-black/70 z-10"
                  >
                    <ChevronLeft className="w-6 h-6" />
                  </button>
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      setActivePhoto((prev) =>
                        prev < photos.length - 1 ? prev + 1 : 0,
                      );
                    }}
                    className="absolute right-4 text-white bg-black/50 p-3 rounded-full hover:bg-black/70 z-10"
                  >
                    <ChevronRight className="w-6 h-6" />
                  </button>
                </>
              )}
              <img
                src={
                  photos[activePhoto]?.watermarked_url ||
                  photos[activePhoto]?.medium_url
                }
                alt={property.title}
                className="max-w-full max-h-[90vh] object-contain"
                onClick={(e) => e.stopPropagation()}
              />
              <p className="absolute bottom-4 text-white text-sm bg-black/50 px-3 py-1 rounded-full">
                {activePhoto + 1} / {photos.length}
              </p>
            </div>
          )}
        </div>
        <div className="lg:col-span-1">
          <div className="sticky top-20 space-y-4">
            {/* Save button for buyers */}
            {role === "buyer" && (
              <button
                onClick={handleToggleSave}
                disabled={saving}
                className={`w-full flex items-center justify-center gap-2 py-2.5 rounded-lg font-medium transition-colors ${
                  saved
                    ? "bg-red-50 text-red-600 border border-red-200 hover:bg-red-100"
                    : "bg-primary-50 text-primary-600 border border-primary-200 hover:bg-primary-100"
                }`}
              >
                <Heart
                  className={`w-5 h-5 ${saved ? "fill-red-500 text-red-500" : ""}`}
                />
                {saved ? "Hapus dari Favorit" : "Simpan ke Favorit"}
              </button>
            )}

            {/* Inquiry button + form for buyers */}
            {role === "buyer" && (
              <div className="card p-4">
                {!inquiryOpen ? (
                  <button
                    onClick={() => setInquiryOpen(true)}
                    className="w-full flex items-center justify-center gap-2 py-2.5 rounded-lg font-medium bg-blue-50 text-blue-600 border border-blue-200 hover:bg-blue-100 transition-colors"
                  >
                    <MessageSquare className="w-5 h-5" />
                    Tanyakan Properti Ini
                  </button>
                ) : (
                  <form onSubmit={handleSendInquiry} className="space-y-3">
                    <h4 className="font-semibold text-sm flex items-center gap-1">
                      <MessageSquare className="w-4 h-4 text-blue-600" />
                      Kirim Pertanyaan
                    </h4>
                    <textarea
                      rows={4}
                      placeholder="Tanyakan detail properti, jadwal survey, atau informasi lainnya..."
                      className="input-field w-full resize-none"
                      value={inquiryMsg}
                      onChange={(e) => setInquiryMsg(e.target.value)}
                      maxLength={2000}
                    />
                    <div className="flex gap-2">
                      <button
                        type="submit"
                        disabled={sendingInquiry}
                        className="btn-primary flex items-center gap-1 text-sm flex-1 justify-center"
                      >
                        <Send className="w-4 h-4" />
                        {sendingInquiry ? "Mengirim..." : "Kirim"}
                      </button>
                      <button
                        type="button"
                        onClick={() => {
                          setInquiryOpen(false);
                          setInquiryMsg("");
                        }}
                        className="btn-secondary text-sm"
                      >
                        Batal
                      </button>
                    </div>
                  </form>
                )}
              </div>
            )}

            {/* Agent card */}
            <div className="card p-4">
              <h3 className="font-semibold mb-3">Kontak Agen</h3>
              {property.salesman && (
                <div className="flex items-center gap-3 mb-3">
                  {property.salesman.photo_url ? (
                    <img
                      src={property.salesman.photo_url}
                      alt=""
                      className="w-12 h-12 rounded-full object-cover"
                    />
                  ) : (
                    <div className="w-12 h-12 rounded-full bg-primary-100 flex items-center justify-center text-primary-600 font-bold">
                      {property.salesman.name?.charAt(0)}
                    </div>
                  )}
                  <div>
                    <p className="font-medium text-sm">
                      {property.salesman.name}
                    </p>
                  </div>
                </div>
              )}
              {waNumber && (
                <a
                  href={`https://wa.me/${waNumber}?text=${waMessage}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="btn-success w-full flex items-center justify-center gap-2"
                >
                  <Phone className="w-4 h-4" /> Hubungi via WhatsApp
                </a>
              )}
            </div>

            {/* Tenant info */}
            {property.tenant && (
              <div className="card p-4">
                <h3 className="font-semibold text-sm mb-2">Dipasarkan oleh</h3>
                <div className="flex items-center gap-3">
                  {property.tenant.logo_url ? (
                    <img
                      src={property.tenant.logo_url}
                      alt=""
                      className="w-10 h-10 rounded-lg object-cover"
                    />
                  ) : (
                    <div className="w-10 h-10 rounded-lg bg-gray-100 flex items-center justify-center">
                      <Building2 className="w-5 h-5 text-gray-400" />
                    </div>
                  )}
                  <div>
                    <p className="font-medium text-sm">
                      {property.tenant.name}
                    </p>
                    {property.tenant.phone && (
                      <p className="text-xs text-gray-500">
                        {property.tenant.phone}
                      </p>
                    )}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
