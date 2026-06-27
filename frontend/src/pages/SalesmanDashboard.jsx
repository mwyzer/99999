import { useState, useEffect, useMemo } from "react";
import { Link, useLocation } from "react-router-dom";
import { salesmanAPI, profileAPI } from "../api";
import { useAuth } from "../context/AuthContext";
import Loading from "../components/Loading";
import { DashboardLayout } from "../components/DashboardSidebar";
import {
  Plus,
  Edit,
  Send,
  CheckCircle,
  XCircle,
  Ban,
  DollarSign,
  ExternalLink,
  MessageSquare,
  MessageCircle,
  Clock,
  Archive,
  User,
  Mail,
  Save,
} from "lucide-react";
import toast from "react-hot-toast";

const STATUS_COLORS = {
  draft: "bg-gray-100 text-gray-700",
  pending: "bg-yellow-100 text-yellow-700",
  approved: "bg-green-100 text-green-700",
  rejected: "bg-red-100 text-red-700",
  sold: "bg-blue-100 text-blue-700",
  rented: "bg-purple-100 text-purple-700",
  inactive: "bg-gray-200 text-gray-500",
};

const STATUS_LABELS = {
  draft: "Draft",
  pending: "Pending",
  approved: "Disetujui",
  rejected: "Ditolak",
  sold: "Terjual",
  rented: "Tersewa",
  inactive: "Nonaktif",
};

export default function SalesmanDashboard() {
  const location = useLocation();
  const { user: authUser, refreshUser } = useAuth();
  const hash = useMemo(
    () => location.hash?.replace("#", "") || "",
    [location.hash],
  );
  const [dashboard, setDashboard] = useState(null);
  const [listings, setListings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [tab, setTab] = useState("all");
  const [mainTab, setMainTab] = useState(
    hash === "inquiries" ? "inquiries" : hash === "me" ? "me" : "listings",
  );

  // Sync mainTab with hash changes from sidebar navigation
  useEffect(() => {
    if (hash === "inquiries") setMainTab("inquiries");
    else if (hash === "me") setMainTab("me");
    else if (hash === "all" || hash === "") setMainTab("listings");
  }, [hash]);

  // Inquiries
  const [inquiries, setInquiries] = useState([]);
  const [inqLoading, setInqLoading] = useState(false);
  const [inqMeta, setInqMeta] = useState({ page: 1, total_pages: 1 });
  // My Profile
  const [meForm, setMeForm] = useState({ name: "", phone: "" });
  const [savingMe, setSavingMe] = useState(false);
  const [replyForm, setReplyForm] = useState({ id: null, status: "" });

  const fetchData = async () => {
    try {
      const [dashRes, listRes] = await Promise.all([
        salesmanAPI.dashboard(),
        salesmanAPI.listListings({ per_page: 50 }),
      ]);
      setDashboard(dashRes.data.data);
      setListings(listRes.data.data);
    } catch {
      toast.error("Gagal memuat dashboard.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    if (mainTab === "inquiries") fetchInquiries();
  }, [mainTab]);

  // Init me form from auth user
  useEffect(() => {
    if (authUser) {
      setMeForm({ name: authUser.name || "", phone: authUser.phone || "" });
    }
  }, [authUser]);

  const handleSaveMe = async (e) => {
    e.preventDefault();
    setSavingMe(true);
    try {
      await profileAPI.update({ name: meForm.name, phone: meForm.phone });
      toast.success("Profil berhasil diperbarui!");
      await refreshUser();
    } catch (err) {
      toast.error(err.response?.data?.error?.message || "Gagal menyimpan.");
    } finally {
      setSavingMe(false);
    }
  };

  const fetchInquiries = async (page = 1) => {
    setInqLoading(true);
    try {
      const { data } = await salesmanAPI.listInquiries({ page, per_page: 20 });
      setInquiries(data.data);
      setInqMeta(data.meta);
    } catch {
      toast.error("Gagal memuat pertanyaan.");
    } finally {
      setInqLoading(false);
    }
  };

  const handleUpdateInquiry = async (id, status) => {
    try {
      await salesmanAPI.updateInquiry(id, { status });
      toast.success("Status pertanyaan diperbarui.");
      setReplyForm({ id: null, status: "" });
      fetchInquiries();
    } catch {
      toast.error("Gagal memperbarui status.");
    }
  };

  const handleAction = async (action, id) => {
    try {
      const actions = {
        submit: salesmanAPI.submitListing,
        deactivate: salesmanAPI.deactivateListing,
        sold: salesmanAPI.markSold,
        rented: salesmanAPI.markRented,
        delete: salesmanAPI.deleteListing,
      };
      if (action === "delete" && !confirm("Yakin hapus listing ini?")) return;
      await actions[action](id);
      toast.success("Berhasil!");
      fetchData();
    } catch (err) {
      toast.error(err.response?.data?.error?.message || "Gagal.");
    }
  };

  if (loading) return <Loading fullScreen />;

  const filteredListings =
    tab === "all" ? listings : listings.filter((l) => l.status === tab);

  return (
    <DashboardLayout role="salesman">
      <div className="p-4 sm:p-6 lg:p-8">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold">Dashboard Salesman</h1>
          <Link
            to="/salesman/listings/new"
            className="btn-primary flex items-center gap-1"
          >
            <Plus className="w-4 h-4" /> Listing Baru
          </Link>
        </div>

        {/* Stats */}
        {dashboard && (
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-6">
            <StatCard
              label="Total Listing"
              value={dashboard.total_listings}
              color="bg-blue-50 text-blue-700"
            />
            <StatCard
              label="Aktif"
              value={dashboard.active_count}
              color="bg-green-50 text-green-700"
            />
            <StatCard
              label="Pending"
              value={dashboard.status_breakdown?.pending || 0}
              color="bg-yellow-50 text-yellow-700"
            />
            <StatCard
              label="Kuota"
              value={`${dashboard.quota?.used || 0}/${dashboard.quota?.max || 5}`}
              color="bg-purple-50 text-purple-700"
            />
          </div>
        )}

        {/* Main Tabs */}
        <div className="flex gap-2 mb-4 border-b">
          {[
            { id: "listings", label: "Listing" },
            { id: "inquiries", label: "Pertanyaan" },
            { id: "me", label: "Profil Saya" },
          ].map((t) => (
            <button
              key={t.id}
              onClick={() => setMainTab(t.id)}
              className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
                mainTab === t.id
                  ? "border-primary-600 text-primary-600"
                  : "border-transparent text-gray-500"
              }`}
            >
              {t.label}
            </button>
          ))}
        </div>

        {mainTab === "listings" && (
          <>
            {/* Status Tabs */}
            <div className="flex flex-wrap gap-2 mb-4">
              {[
                "all",
                "draft",
                "pending",
                "approved",
                "rejected",
                "sold",
                "rented",
                "inactive",
              ].map((s) => (
                <button
                  key={s}
                  onClick={() => setTab(s)}
                  className={`px-3 py-1.5 rounded-full text-xs font-medium transition-colors ${
                    tab === s
                      ? "bg-primary-600 text-white"
                      : "bg-gray-100 text-gray-600 hover:bg-gray-200"
                  }`}
                >
                  {s === "all" ? "Semua" : STATUS_LABELS[s]}
                </button>
              ))}
            </div>

            {/* Listings */}
            <div className="space-y-3">
              {filteredListings.length === 0 ? (
                <p className="text-center py-10 text-gray-400">
                  Tidak ada listing.
                </p>
              ) : (
                filteredListings.map((l) => (
                  <div
                    key={l.id}
                    className="card p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
                  >
                    <div className="flex items-center gap-3 min-w-0">
                      {l.main_photo_url ? (
                        <img
                          src={l.main_photo_url}
                          alt=""
                          className="w-16 h-12 rounded-lg object-cover flex-shrink-0"
                        />
                      ) : (
                        <div className="w-16 h-12 rounded-lg bg-gray-200 flex-shrink-0" />
                      )}
                      <div className="min-w-0">
                        <h3 className="font-medium text-sm truncate">
                          {l.title}
                        </h3>
                        <p className="text-xs text-gray-500">
                          {l.city} · {l.property_type}
                        </p>
                        <span
                          className={`inline-block mt-1 px-2 py-0.5 rounded-full text-xs font-medium ${STATUS_COLORS[l.status] || ""}`}
                        >
                          {STATUS_LABELS[l.status] || l.status}
                        </span>
                      </div>
                    </div>
                    <div className="flex gap-1 flex-shrink-0">
                      <a
                        href={`/properties/${l.id}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="btn-secondary text-xs py-1 px-2 flex items-center gap-1"
                        title="Lihat publik"
                      >
                        <ExternalLink className="w-3 h-3" />
                      </a>
                      {(l.status === "draft" || l.status === "rejected") && (
                        <>
                          <Link
                            to={`/salesman/listings/${l.id}/edit`}
                            className="btn-secondary text-xs py-1 px-2 flex items-center gap-1"
                          >
                            <Edit className="w-3 h-3" /> Edit
                          </Link>
                          <button
                            onClick={() => handleAction("submit", l.id)}
                            className="btn-primary text-xs py-1 px-2 flex items-center gap-1"
                          >
                            <Send className="w-3 h-3" /> Ajukan
                          </button>
                          <button
                            onClick={() => handleAction("delete", l.id)}
                            className="btn-danger text-xs py-1 px-2"
                          >
                            Hapus
                          </button>
                        </>
                      )}
                      {l.status === "approved" && (
                        <>
                          <button
                            onClick={() => handleAction("deactivate", l.id)}
                            className="btn-secondary text-xs py-1 px-2 flex items-center gap-1"
                          >
                            <Ban className="w-3 h-3" /> Nonaktifkan
                          </button>
                          <button
                            onClick={() => handleAction("sold", l.id)}
                            className="btn-success text-xs py-1 px-2 flex items-center gap-1"
                          >
                            <CheckCircle className="w-3 h-3" /> Terjual
                          </button>
                          <button
                            onClick={() => handleAction("rented", l.id)}
                            className="bg-purple-600 text-white text-xs py-1 px-2 rounded-lg flex items-center gap-1"
                          >
                            <DollarSign className="w-3 h-3" /> Tersewa
                          </button>
                        </>
                      )}
                      {l.status === "rejected" && (
                        <Link
                          to={`/salesman/listings/${l.id}/edit`}
                          className="btn-secondary text-xs py-1 px-2 flex items-center gap-1"
                        >
                          <Edit className="w-3 h-3" /> Perbaiki
                        </Link>
                      )}
                    </div>
                  </div>
                ))
              )}
            </div>
          </>
        )}

        {/* Inquiries Tab */}
        {mainTab === "inquiries" && (
          <>
            {inqLoading ? (
              <Loading />
            ) : inquiries.length === 0 ? (
              <div className="text-center py-16 text-gray-400">
                <MessageSquare className="w-12 h-12 mx-auto mb-3 text-gray-300" />
                <p className="text-lg">Belum ada pertanyaan.</p>
                <p className="text-sm mt-1">
                  Pertanyaan dari buyer akan muncul di sini.
                </p>
              </div>
            ) : (
              <div className="space-y-3">
                {inquiries.map((inq) => {
                  const statusCfg = {
                    unread: {
                      label: "Belum Dibaca",
                      cls: "bg-yellow-100 text-yellow-700",
                    },
                    read: { label: "Dibaca", cls: "bg-blue-100 text-blue-700" },
                    replied: {
                      label: "Dibalas",
                      cls: "bg-green-100 text-green-700",
                    },
                    closed: {
                      label: "Ditutup",
                      cls: "bg-gray-100 text-gray-600",
                    },
                  };
                  const cfg = statusCfg[inq.status] || statusCfg.unread;
                  return (
                    <div key={inq.id} className="card p-4">
                      <div className="flex items-start justify-between gap-3">
                        <div className="min-w-0 flex-1">
                          <div className="flex items-center gap-2 mb-1">
                            <span
                              className={`px-2 py-0.5 rounded-full text-xs font-medium ${cfg.cls}`}
                            >
                              {cfg.label}
                            </span>
                            {inq.property_title && (
                              <span className="text-sm font-medium text-gray-700 truncate">
                                {inq.property_title}
                              </span>
                            )}
                          </div>
                          {inq.buyer_name && (
                            <p className="text-xs text-gray-500 mb-1">
                              Dari: {inq.buyer_name} ({inq.buyer_email})
                            </p>
                          )}
                          {inq.message && (
                            <p className="text-sm text-gray-600 whitespace-pre-line">
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
                        <div className="flex gap-1 flex-shrink-0">
                          {inq.status === "unread" && (
                            <button
                              onClick={() =>
                                handleUpdateInquiry(inq.id, "read")
                              }
                              className="btn-secondary text-xs py-1 px-2"
                            >
                              Tandai Dibaca
                            </button>
                          )}
                          {(inq.status === "unread" ||
                            inq.status === "read") && (
                            <>
                              <button
                                onClick={() =>
                                  handleUpdateInquiry(inq.id, "replied")
                                }
                                className="btn-success text-xs py-1 px-2 flex items-center gap-1"
                              >
                                <MessageCircle className="w-3 h-3" /> Dibalas
                              </button>
                              <button
                                onClick={() =>
                                  handleUpdateInquiry(inq.id, "closed")
                                }
                                className="btn-secondary text-xs py-1 px-2 flex items-center gap-1"
                              >
                                <Archive className="w-3 h-3" />
                              </button>
                            </>
                          )}
                          {inq.status === "replied" && (
                            <button
                              onClick={() =>
                                handleUpdateInquiry(inq.id, "closed")
                              }
                              className="btn-secondary text-xs py-1 px-2 flex items-center gap-1"
                            >
                              <Archive className="w-3 h-3" /> Tutup
                            </button>
                          )}
                        </div>
                      </div>
                    </div>
                  );
                })}
                {inqMeta.total_pages > 1 && (
                  <div className="flex items-center justify-center gap-3 mt-4">
                    <button
                      onClick={() => fetchInquiries(inqMeta.page - 1)}
                      disabled={inqMeta.page <= 1}
                      className="btn-secondary text-sm"
                    >
                      Sebelumnya
                    </button>
                    <span className="text-sm text-gray-500">
                      Halaman {inqMeta.page} dari {inqMeta.total_pages}
                    </span>
                    <button
                      onClick={() => fetchInquiries(inqMeta.page + 1)}
                      disabled={inqMeta.page >= inqMeta.total_pages}
                      className="btn-secondary text-sm"
                    >
                      Selanjutnya
                    </button>
                  </div>
                )}
              </div>
            )}
          </>
        )}

        {mainTab === "me" && authUser && (
          <div className="card p-5 max-w-lg">
            <h3 className="font-semibold text-lg mb-4 flex items-center gap-2">
              <User className="w-5 h-5 text-primary-600" /> Profil Saya
            </h3>
            <div className="flex items-center gap-4 mb-5">
              {authUser.photo_url ? (
                <img
                  src={authUser.photo_url}
                  alt=""
                  className="w-16 h-16 rounded-full object-cover"
                />
              ) : (
                <div className="w-16 h-16 rounded-full bg-primary-100 flex items-center justify-center text-primary-600 font-bold text-xl">
                  {authUser.name?.charAt(0) || "U"}
                </div>
              )}
              <div>
                <h2 className="font-semibold text-lg">{authUser.name}</h2>
                <p className="text-sm text-gray-500 flex items-center gap-1">
                  <Mail className="w-3.5 h-3.5" /> {authUser.email}
                </p>
              </div>
            </div>
            <form onSubmit={handleSaveMe} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Nama
                </label>
                <input
                  type="text"
                  required
                  className="input-field"
                  value={meForm.name}
                  onChange={(e) =>
                    setMeForm({ ...meForm, name: e.target.value })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Telepon
                </label>
                <input
                  type="text"
                  className="input-field"
                  value={meForm.phone}
                  onChange={(e) =>
                    setMeForm({ ...meForm, phone: e.target.value })
                  }
                />
              </div>
              <button
                type="submit"
                disabled={savingMe}
                className="btn-primary flex items-center gap-2"
              >
                <Save className="w-4 h-4" />
                {savingMe ? "Menyimpan..." : "Simpan Perubahan"}
              </button>
            </form>
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}

function StatCard({ label, value, color }) {
  return (
    <div className={`rounded-xl p-4 ${color}`}>
      <p className="text-xs opacity-75">{label}</p>
      <p className="text-2xl font-bold mt-1">{value}</p>
    </div>
  );
}
