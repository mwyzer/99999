import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { buyerAPI } from "../api";
import Loading from "../components/Loading";
import { DashboardLayout } from "../components/DashboardSidebar";
import {
  MessageSquare,
  ChevronLeft,
  ChevronRight,
  Clock,
  CheckCircle,
  MessageCircle,
  Archive,
} from "lucide-react";
import toast from "react-hot-toast";

const STATUS_CONFIG = {
  unread: {
    label: "Belum Dibaca",
    icon: Clock,
    cls: "bg-yellow-100 text-yellow-700",
  },
  read: {
    label: "Dibaca",
    icon: CheckCircle,
    cls: "bg-blue-100 text-blue-700",
  },
  replied: {
    label: "Dibalas",
    icon: MessageCircle,
    cls: "bg-green-100 text-green-700",
  },
  closed: { label: "Ditutup", icon: Archive, cls: "bg-gray-100 text-gray-600" },
};

export default function BuyerInquiriesPage() {
  const [inquiries, setInquiries] = useState([]);
  const [loading, setLoading] = useState(true);
  const [meta, setMeta] = useState({ page: 1, total_pages: 1, total: 0 });

  const fetchInquiries = async (page = 1) => {
    setLoading(true);
    try {
      const { data } = await buyerAPI.listInquiries({ page, per_page: 12 });
      setInquiries(data.data);
      setMeta(data.meta);
    } catch {
      toast.error("Gagal memuat pertanyaan.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchInquiries();
  }, []);

  return (
    <DashboardLayout role="buyer">
      <div className="p-4 sm:p-6 lg:p-8">
        <div className="flex items-center gap-2 mb-6">
          <MessageSquare className="w-6 h-6 text-primary-600" />
          <h1 className="text-2xl font-bold">Pertanyaan Saya</h1>
        </div>

        {loading ? (
          <Loading />
        ) : inquiries.length === 0 ? (
          <div className="text-center py-16 text-gray-400">
            <MessageSquare className="w-12 h-12 mx-auto mb-3 text-gray-300" />
            <p className="text-lg">Belum ada pertanyaan.</p>
            <p className="text-sm mt-1">
              Kirim pertanyaan dari halaman detail properti yang Anda minati.
            </p>
          </div>
        ) : (
          <>
            <div className="space-y-3">
              {inquiries.map((inq) => {
                const cfg = STATUS_CONFIG[inq.status] || STATUS_CONFIG.unread;
                const Icon = cfg.icon;
                return (
                  <div key={inq.id} className="card p-4">
                    <div className="flex items-start justify-between gap-3">
                      <div className="min-w-0 flex-1">
                        <Link
                          to={`/properties/${inq.property_id}`}
                          className="font-medium text-gray-900 hover:text-primary-600 transition-colors line-clamp-1"
                        >
                          {inq.property_title || "Properti"}
                        </Link>
                        {inq.message && (
                          <p className="text-sm text-gray-600 mt-1 line-clamp-2">
                            {inq.message}
                          </p>
                        )}
                        <p className="text-xs text-gray-400 mt-1.5">
                          {new Date(inq.created_at).toLocaleDateString(
                            "id-ID",
                            {
                              day: "numeric",
                              month: "long",
                              year: "numeric",
                              hour: "2-digit",
                              minute: "2-digit",
                            },
                          )}
                        </p>
                      </div>
                      <span
                        className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium flex-shrink-0 ${cfg.cls}`}
                      >
                        <Icon className="w-3 h-3" />
                        {cfg.label}
                      </span>
                    </div>
                  </div>
                );
              })}
            </div>

            {/* Pagination */}
            {meta.total_pages > 1 && (
              <div className="flex items-center justify-center gap-3 mt-8">
                <button
                  onClick={() => fetchInquiries(meta.page - 1)}
                  disabled={meta.page <= 1}
                  className="btn-secondary flex items-center gap-1 text-sm"
                >
                  <ChevronLeft className="w-4 h-4" /> Sebelumnya
                </button>
                <span className="text-sm text-gray-500">
                  Halaman {meta.page} dari {meta.total_pages}
                </span>
                <button
                  onClick={() => fetchInquiries(meta.page + 1)}
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
