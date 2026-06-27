import { useState, useEffect, useCallback } from "react";
import { publicAPI, buyerAPI } from "../api";
import { useAuth } from "../context/AuthContext";
import PropertyCard from "../components/PropertyCard";
import PropertyFilter from "../components/PropertyFilter";
import Loading from "../components/Loading";
import { MapPin, TrendingUp, Navigation, ArrowUp } from "lucide-react";
import toast from "react-hot-toast";

export default function HomePage() {
  const { isAuthenticated, role } = useAuth();
  const [properties, setProperties] = useState([]);
  const [featured, setFeatured] = useState([]);
  const [cities, setCities] = useState([]);
  const [loading, setLoading] = useState(true);
  const [meta, setMeta] = useState({ page: 1, total_pages: 1 });
  const [savedIds, setSavedIds] = useState(new Set());
  const [filters, setFilters] = useState({});
  const [nearby, setNearby] = useState([]);
  const [nearbyLoading, setNearbyLoading] = useState(false);
  const [geoError, setGeoError] = useState(null);
  const [showTopBtn, setShowTopBtn] = useState(false);

  useEffect(() => {
    const onScroll = () => setShowTopBtn(window.scrollY > 500);
    window.addEventListener("scroll", onScroll);
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  const fetchProperties = useCallback(async (params = {}, page = 1) => {
    setLoading(true);
    try {
      const { data } = await publicAPI.listProperties({
        ...params,
        page,
        per_page: 12,
      });
      setProperties(data.data);
      setMeta(data.meta);
    } catch {
      toast.error("Gagal memuat data properti.");
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchFeatured = async () => {
    try {
      const { data } = await publicAPI.featured({ limit: 6 });
      setFeatured(data.data.properties || []);
    } catch {
      /* silent */
    }
  };

  const fetchCities = async () => {
    try {
      const { data } = await publicAPI.cities();
      setCities(data.data || []);
    } catch {
      /* silent */
    }
  };

  const fetchNearby = async () => {
    if (!navigator.geolocation) {
      setGeoError("Browser Anda tidak mendukung geolokasi.");
      return;
    }
    setNearbyLoading(true);
    navigator.geolocation.getCurrentPosition(
      async (pos) => {
        try {
          const { data } = await publicAPI.nearby({
            latitude: pos.coords.latitude,
            longitude: pos.coords.longitude,
            radius_km: 25,
            limit: 6,
          });
          setNearby(data.data || []);
        } catch {
          /* silent */
        } finally {
          setNearbyLoading(false);
        }
      },
      () => {
        setGeoError(
          "Tidak dapat mengakses lokasi. Izinkan akses lokasi di browser.",
        );
        setNearbyLoading(false);
      },
      { timeout: 5000 },
    );
  };

  const fetchSaved = async () => {
    if (!isAuthenticated || role !== "buyer") return;
    try {
      const { data } = await buyerAPI.listSaved({ per_page: 100 });
      const ids = new Set(data.data.map((s) => s.id));
      setSavedIds(ids);
    } catch {
      /* silent */
    }
  };

  useEffect(() => {
    fetchProperties();
    fetchFeatured();
    fetchCities();
    fetchNearby();
  }, []);
  useEffect(() => {
    fetchSaved();
  }, [isAuthenticated]);

  const handleFilter = (f) => {
    setFilters(f);
    fetchProperties(f, 1);
  };

  const handleSave = async (id) => {
    if (!isAuthenticated) {
      toast.error("Silakan login terlebih dahulu.");
      return;
    }
    try {
      await buyerAPI.save(id);
      setSavedIds((prev) => new Set([...prev, id]));
      toast.success("Disimpan ke favorit!");
    } catch {
      toast.error("Gagal menyimpan.");
    }
  };

  const handleUnsave = async (id) => {
    try {
      await buyerAPI.unsave(id);
      setSavedIds((prev) => {
        const n = new Set(prev);
        n.delete(id);
        return n;
      });
      toast.success("Dihapus dari favorit.");
    } catch {
      toast.error("Gagal menghapus.");
    }
  };

  const handlePageChange = (page) => {
    fetchProperties(filters, page);
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  return (
    <div>
      {/* Hero */}
      <section className="bg-gradient-to-br from-primary-700 to-primary-900 text-white py-12 sm:py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h1 className="text-3xl sm:text-5xl font-extrabold mb-4">
            Temukan Properti Impian Anda
          </h1>
          <p className="text-primary-100 text-lg mb-8 max-w-2xl mx-auto">
            Jelajahi ribuan listing properti dari agensi, bank, dan perusahaan
            terpercaya — rumah, apartemen, tanah, dan ruko di seluruh Indonesia.
          </p>
          <div className="flex justify-center gap-4">
            <span className="flex items-center gap-1 text-primary-200 text-sm">
              <MapPin className="w-4 h-4" /> Jabodetabek • Bandung • Surabaya
            </span>
          </div>
        </div>
      </section>

      {/* Featured */}
      {featured.length > 0 && (
        <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
          <h2 className="text-xl font-bold flex items-center gap-2 mb-6">
            <TrendingUp className="w-5 h-5 text-primary-600" />
            Properti Pilihan
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {featured.map((p) => (
              <PropertyCard
                key={p.id}
                property={p}
                onSave={role === "buyer" ? handleSave : null}
                onUnsave={role === "buyer" ? handleUnsave : null}
                saved={savedIds.has(p.id)}
              />
            ))}
          </div>
        </section>
      )}

      {/* Nearby Properties */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 bg-gray-50">
        <h2 className="text-xl font-bold flex items-center gap-2 mb-6">
          <Navigation className="w-5 h-5 text-primary-600" />
          Properti di Sekitar Anda
        </h2>
        {nearbyLoading ? (
          <Loading />
        ) : geoError ? (
          <div className="text-center py-8">
            <p className="text-sm text-gray-400 mb-2">{geoError}</p>
            <button
              onClick={fetchNearby}
              className="text-sm text-primary-600 hover:text-primary-700"
            >
              Coba Lagi
            </button>
          </div>
        ) : nearby.length === 0 ? (
          <p className="text-sm text-gray-400 text-center py-8">
            Tidak ada properti di sekitar lokasi Anda.
          </p>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {nearby.map((p) => (
              <PropertyCard
                key={p.id}
                property={p}
                onSave={role === "buyer" ? handleSave : null}
                onUnsave={role === "buyer" ? handleUnsave : null}
                saved={savedIds.has(p.id)}
              />
            ))}
          </div>
        )}
      </section>

      {/* All listings */}
      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        <h2 className="text-xl font-bold mb-6">Semua Properti</h2>
        <PropertyFilter onFilter={handleFilter} cities={cities} />

        {loading ? (
          <Loading />
        ) : (
          <>
            {properties.length === 0 ? (
              <div className="text-center py-16 text-gray-400">
                <p className="text-lg">Tidak ada properti ditemukan.</p>
                <p className="text-sm mt-1">Coba ubah filter pencarian Anda.</p>
              </div>
            ) : (
              <>
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                  {properties.map((p) => (
                    <PropertyCard
                      key={p.id}
                      property={p}
                      onSave={role === "buyer" ? handleSave : null}
                      onUnsave={role === "buyer" ? handleUnsave : null}
                      saved={savedIds.has(p.id)}
                    />
                  ))}
                </div>

                {/* Pagination */}
                {meta.total_pages > 1 && (
                  <div className="flex justify-center gap-2 mt-8">
                    {Array.from(
                      { length: meta.total_pages },
                      (_, i) => i + 1,
                    ).map((page) => (
                      <button
                        key={page}
                        onClick={() => handlePageChange(page)}
                        className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                          page === meta.page
                            ? "bg-primary-600 text-white"
                            : "bg-gray-100 text-gray-600 hover:bg-gray-200"
                        }`}
                      >
                        {page}
                      </button>
                    ))}
                  </div>
                )}
              </>
            )}
          </>
        )}
      </section>

      {/* Back to top */}
      {showTopBtn && (
        <button
          onClick={() => window.scrollTo({ top: 0, behavior: "smooth" })}
          className="fixed bottom-6 right-6 z-40 bg-primary-600 text-white p-3 rounded-full shadow-lg hover:bg-primary-700 transition-colors"
        >
          <ArrowUp className="w-5 h-5" />
        </button>
      )}
    </div>
  );
}
