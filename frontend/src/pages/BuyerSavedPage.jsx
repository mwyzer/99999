import { useState, useEffect } from "react";
import { buyerAPI } from "../api";
import PropertyCard from "../components/PropertyCard";
import Loading from "../components/Loading";
import { DashboardLayout } from "../components/DashboardSidebar";
import { Heart, ChevronLeft, ChevronRight } from "lucide-react";
import toast from "react-hot-toast";

export default function BuyerSavedPage() {
  const [saved, setSaved] = useState([]);
  const [loading, setLoading] = useState(true);
  const [meta, setMeta] = useState({
    page: 1,
    total_pages: 1,
    total: 0,
    per_page: 12,
  });

  const fetchSaved = async (page = 1) => {
    setLoading(true);
    try {
      const { data } = await buyerAPI.listSaved({ page, per_page: 12 });
      setSaved(data.data);
      setMeta(data.meta);
    } catch {
      toast.error("Gagal memuat favorit.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSaved();
  }, []);

  const handleUnsave = async (id) => {
    try {
      await buyerAPI.unsave(id);
      setSaved((prev) => prev.filter((s) => s.id !== id));
      toast.success("Dihapus dari favorit.");
    } catch {
      toast.error("Gagal menghapus.");
    }
  };

  return (
    <DashboardLayout role="buyer">
      <div className="p-4 sm:p-6 lg:p-8">
        <div className="flex items-center gap-2 mb-6">
          <Heart className="w-6 h-6 text-red-500" />
          <h1 className="text-2xl font-bold">Properti Favorit Saya</h1>
        </div>

        {loading ? (
          <Loading />
        ) : saved.length === 0 ? (
          <div className="text-center py-16 text-gray-400">
            <Heart className="w-12 h-12 mx-auto mb-3 text-gray-300" />
            <p className="text-lg">Belum ada properti favorit.</p>
            <p className="text-sm mt-1">
              Jelajahi listing dan simpan properti yang Anda sukai.
            </p>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {saved.map((p) => (
                <PropertyCard
                  key={p.id}
                  property={p}
                  onUnsave={handleUnsave}
                  saved
                />
              ))}
            </div>

            {/* Pagination */}
            {meta.total_pages > 1 && (
              <div className="flex items-center justify-center gap-3 mt-8">
                <button
                  onClick={() => fetchSaved(meta.page - 1)}
                  disabled={meta.page <= 1}
                  className="btn-secondary flex items-center gap-1 text-sm"
                >
                  <ChevronLeft className="w-4 h-4" /> Sebelumnya
                </button>
                <span className="text-sm text-gray-500">
                  Halaman {meta.page} dari {meta.total_pages}
                </span>
                <button
                  onClick={() => fetchSaved(meta.page + 1)}
                  disabled={meta.page >= meta.total_pages}
                  className="btn-secondary flex items-center gap-1 text-sm"
                >
                  Selanjutnya <ChevronRight className="w-4 h-4" />
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </DashboardLayout>
  );
}
