import { useState, useEffect, useMemo } from "react";
import { useLocation } from "react-router-dom";
import { adminAPI, profileAPI } from "../api";
import { useAuth } from "../context/AuthContext";
import Loading from "../components/Loading";
import { DashboardLayout } from "../components/DashboardSidebar";
import {
  Building2,
  CheckCircle,
  XCircle,
  Ban,
  RefreshCw,
  Users,
  Home,
  Plus,
  Eye,
  FileText,
  Clock,
  Search,
  ChevronLeft,
  ChevronRight,
  Crown,
  User,
  Mail,
  Save,
  ArrowUpCircle,
  Trash2,
} from "lucide-react";
import toast from "react-hot-toast";

export default function AdminDashboard() {
  const location = useLocation();
  const { user: authUser, refreshUser } = useAuth();
  const hash = useMemo(
    () => location.hash?.replace("#", "") || "",
    [location.hash],
  );
  const [dashboard, setDashboard] = useState(null);
  const [tenants, setTenants] = useState([]);
  const [pendingListings, setPendingListings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [tab, setTab] = useState(hash || "overview");

  // Sync tab with hash changes from sidebar navigation
  useEffect(() => {
    if (hash && hash !== "overview") setTab(hash);
  }, [hash]);

  const [rejectForm, setRejectForm] = useState({ id: null, reason: "" });
  // Create tenant
  const [createForm, setCreateForm] = useState({
    organization_name: "",
    subdomain_slug: "",
    admin_name: "",
    admin_email: "",
    admin_phone: "",
    admin_password: "",
    plan_type: "free",
  });
  const [creating, setCreating] = useState(false);
  // Audit logs
  const [auditLogs, setAuditLogs] = useState([]);
  const [auditMeta, setAuditMeta] = useState({ page: 1, total_pages: 1 });
  const [auditFilters, setAuditFilters] = useState({ action: "", user_id: "" });
  // Change plan
  const pendingUpgradeTenants = tenants.filter(
    (t) => t.plan_type?.startsWith("pending_") && t.status === "active",
  );

  // My Profile
  const [meForm, setMeForm] = useState({ name: "", phone: "" });
  const [savingMe, setSavingMe] = useState(false);

  const fetchAll = async () => {
    try {
      const [d, t, p] = await Promise.all([
        adminAPI.dashboard(),
        adminAPI.listTenants({ per_page: 50 }),
        adminAPI.listPending({ per_page: 50 }),
      ]);
      setDashboard(d.data.data);
      setTenants(t.data.data);
      setPendingListings(p.data.data);
    } catch {
      toast.error("Gagal memuat data admin.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAll();
  }, []);

  useEffect(() => {
    if (tab === "audit-logs") fetchAuditLogs();
  }, [tab]);

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

  const handleApprove = async (id) => {
    try {
      await adminAPI.approveListing(id);
      toast.success("Listing disetujui!");
      fetchAll();
    } catch (err) {
      toast.error(err.response?.data?.error?.message || "Gagal.");
    }
  };

  const handleReject = async () => {
    if (!rejectForm.reason || rejectForm.reason.length < 10) {
      toast.error("Alasan penolakan minimal 10 karakter.");
      return;
    }
    try {
      await adminAPI.rejectListing(rejectForm.id, {
        reason: rejectForm.reason,
      });
      toast.success("Listing ditolak.");
      setRejectForm({ id: null, reason: "" });
      fetchAll();
    } catch (err) {
      toast.error(err.response?.data?.error?.message || "Gagal.");
    }
  };

  const handleTenantAction = async (id, action) => {
    const actions = {
      suspend: adminAPI.suspendTenant,
      activate: adminAPI.activateTenant,
    };
    try {
      await actions[action](id);
      toast.success(
        `Tenant ${action === "suspend" ? "dinonaktifkan" : "diaktifkan"}!`,
      );
      fetchAll();
    } catch {
      toast.error("Gagal.");
    }
  };

  const handleDeleteTenant = async (id, name) => {
    if (
      !confirm(
        `Yakin hapus tenant "${name}"? Semua data (user, listing, subscription) akan dihapus permanen.`,
      )
    )
      return;
    try {
      await adminAPI.deleteTenant(id);
      toast.success("Tenant berhasil dihapus!");
      fetchAll();
    } catch (err) {
      toast.error(
        err.response?.data?.error?.message || "Gagal menghapus tenant.",
      );
    }
  };

  const handleCreateTenant = async (e) => {
    e.preventDefault();
    setCreating(true);
    try {
      await adminAPI.createTenant(createForm);
      toast.success("Tenant berhasil dibuat!");
      setCreateForm({
        organization_name: "",
        subdomain_slug: "",
        admin_name: "",
        admin_email: "",
        admin_phone: "",
        admin_password: "",
        plan_type: "free",
      });
      fetchAll();
    } catch (err) {
      const details = err.response?.data?.error?.details;
      if (details?.length) {
        const fieldErrors = details
          .map((d) => `• ${d.field}: ${d.message}`)
          .join("\n");
        toast.error(fieldErrors, { duration: 6000 });
      } else {
        toast.error(
          err.response?.data?.error?.message || "Gagal membuat tenant.",
        );
      }
    } finally {
      setCreating(false);
    }
  };

  const fetchAuditLogs = async (page = 1) => {
    try {
      const params = { page, per_page: 20, ...auditFilters };
      const { data } = await adminAPI.auditLogs(params);
      setAuditLogs(data.data);
      setAuditMeta(data.meta);
    } catch {
      toast.error("Gagal memuat audit logs.");
    }
  };

  const handleChangePlan = async (tenantId, planType) => {
    try {
      await adminAPI.changePlan(tenantId, { plan_type: planType });
      toast.success(`Paket tenant diubah ke ${planType}!`);
      fetchAll();
    } catch (err) {
      const msg =
        err.response?.data?.error?.message ||
        err.response?.data?.message ||
        "Gagal mengubah paket.";
      toast.error(msg);
    }
  };

  if (loading) return <Loading fullScreen />;

  return (
    <DashboardLayout role="platform_admin">
      <div className="p-4 sm:p-6 lg:p-8">
        {/* Stats */}
        {dashboard && (
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-6">
            <StatCard
              icon={<Building2 className="w-5 h-5" />}
              label="Tenant"
              value={dashboard.total_tenants}
            />
            <StatCard
              icon={<Users className="w-5 h-5" />}
              label="User"
              value={dashboard.total_users}
            />
            <StatCard
              icon={<Home className="w-5 h-5" />}
              label="Listing"
              value={dashboard.total_listings}
            />
            <StatCard
              icon={<Eye className="w-5 h-5" />}
              label="Pending Review"
              value={dashboard.pending_reviews}
              color="bg-yellow-50 text-yellow-700"
            />
            <StatCard
              icon={<ArrowUpCircle className="w-5 h-5" />}
              label="Upgrade Request"
              value={dashboard.pending_upgrades || 0}
              color="bg-blue-50 text-blue-700"
            />
          </div>
        )}

        {/* Tabs */}
        <div className="flex gap-2 mb-6 border-b">
          {[
            { id: "overview", label: "Ringkasan" },
            {
              id: "pending",
              label: `Review Pending (${pendingListings.length})`,
            },
            {
              id: "upgrades",
              label: `Upgrade (${pendingUpgradeTenants.length})`,
            },
            { id: "tenants", label: "Daftar Tenant" },
            { id: "create-tenant", label: "Tambah Tenant" },
            { id: "audit-logs", label: "Audit Log" },
            { id: "me", label: "Profil Saya" },
          ].map((t) => (
            <button
              key={t.id}
              onClick={() => setTab(t.id)}
              className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
                tab === t.id
                  ? "border-primary-600 text-primary-600"
                  : "border-transparent text-gray-500"
              }`}
            >
              {t.label}
            </button>
          ))}
        </div>

        {/* Pending Listings */}
        {tab === "pending" && (
          <div className="space-y-3">
            {pendingListings.length === 0 ? (
              <p className="text-center py-10 text-gray-400">
                Tidak ada listing pending.
              </p>
            ) : (
              pendingListings.map((l) => (
                <div key={l.id} className="card p-4">
                  <div className="flex justify-between items-start gap-4">
                    <div>
                      <h3 className="font-medium">{l.title}</h3>
                      <p className="text-sm text-gray-500">
                        {l.city} · {l.property_type}
                      </p>
                      <p className="text-xs text-gray-400 mt-1">
                        Tenant: {l.tenant?.name} · Salesman: {l.salesman?.name}
                      </p>
                      <p className="text-sm font-semibold text-primary-700 mt-1">
                        Rp {parseFloat(l.price).toLocaleString("id-ID")}
                      </p>
                    </div>
                    <div className="flex gap-2 flex-shrink-0">
                      <button
                        onClick={() => handleApprove(l.id)}
                        className="btn-success text-xs py-1 px-3 flex items-center gap-1"
                      >
                        <CheckCircle className="w-3 h-3" /> Setujui
                      </button>
                      <button
                        onClick={() => setRejectForm({ id: l.id, reason: "" })}
                        className="btn-danger text-xs py-1 px-3 flex items-center gap-1"
                      >
                        <XCircle className="w-3 h-3" /> Tolak
                      </button>
                    </div>
                  </div>
                  {/* Reject reason form */}
                  {rejectForm.id === l.id && (
                    <div className="mt-3 pt-3 border-t">
                      <textarea
                        className="input-field text-sm mb-2"
                        rows={2}
                        placeholder="Alasan penolakan (min. 10 karakter)..."
                        value={rejectForm.reason}
                        onChange={(e) =>
                          setRejectForm({
                            ...rejectForm,
                            reason: e.target.value,
                          })
                        }
                      />
                      <div className="flex gap-2">
                        <button
                          onClick={handleReject}
                          className="btn-danger text-xs py-1 px-3"
                        >
                          Konfirmasi Tolak
                        </button>
                        <button
                          onClick={() =>
                            setRejectForm({ id: null, reason: "" })
                          }
                          className="btn-secondary text-xs py-1 px-3"
                        >
                          Batal
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              ))
            )}
          </div>
        )}

        {/* Upgrade Requests */}
        {tab === "upgrades" && (
          <div className="space-y-3">
            {pendingUpgradeTenants.length === 0 ? (
              <div className="text-center py-16 text-gray-400">
                <ArrowUpCircle className="w-12 h-12 mx-auto mb-3 text-gray-300" />
                <p className="text-lg">Tidak ada permintaan.</p>
                <p className="text-sm mt-1">
                  Tenant yang mengajukan perubahan paket akan muncul di sini.
                </p>
              </div>
            ) : (
              pendingUpgradeTenants.map((t) => {
                const reqLabel =
                  t.plan_type === "pending_upgrade"
                    ? "Ingin Premium"
                    : t.plan_type === "pending_free"
                      ? "Ingin Free"
                      : "Ingin Nonaktif";
                const approvePlan =
                  t.plan_type === "pending_upgrade"
                    ? "premium"
                    : t.plan_type === "pending_free"
                      ? "free"
                      : null;
                const isDisableReq = t.plan_type === "pending_disable";

                return (
                  <div
                    key={t.id}
                    className="card p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-l-4 border-blue-400"
                  >
                    <div className="min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <p className="font-medium">{t.organization_name}</p>
                        <span className="inline-block px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-700">
                          {reqLabel}
                        </span>
                      </div>
                      <p className="text-xs text-gray-500">
                        {t.subdomain_slug} · {t.total_listings} listing ·{" "}
                        {t.total_users} user
                      </p>
                    </div>
                    <div className="flex gap-2 flex-shrink-0">
                      {isDisableReq ? (
                        <>
                          <button
                            onClick={() => handleTenantAction(t.id, "suspend")}
                            className="btn-danger text-xs py-1.5 px-3 flex items-center gap-1"
                          >
                            <Ban className="w-3.5 h-3.5" /> Nonaktifkan
                          </button>
                          <button
                            onClick={() => handleChangePlan(t.id, "free")}
                            className="btn-secondary text-xs py-1.5 px-3 flex items-center gap-1"
                          >
                            <XCircle className="w-3.5 h-3.5" /> Tolak
                          </button>
                        </>
                      ) : (
                        <>
                          <button
                            onClick={() => handleChangePlan(t.id, approvePlan)}
                            className="btn-success text-xs py-1.5 px-3 flex items-center gap-1"
                          >
                            <CheckCircle className="w-3.5 h-3.5" /> Setujui
                          </button>
                          <button
                            onClick={() =>
                              handleChangePlan(
                                t.id,
                                approvePlan === "premium" ? "free" : "premium",
                              )
                            }
                            className="btn-secondary text-xs py-1.5 px-3 flex items-center gap-1"
                          >
                            <XCircle className="w-3.5 h-3.5" /> Tolak
                          </button>
                        </>
                      )}
                    </div>
                  </div>
                );
              })
            )}
          </div>
        )}

        {/* Tenants */}
        {tab === "tenants" && (
          <div className="space-y-3">
            {tenants.map((t) => {
              const planLabel =
                t.plan_type === "premium"
                  ? "Premium"
                  : t.plan_type === "pending_upgrade"
                    ? "Menunggu Premium"
                    : t.plan_type === "pending_free"
                      ? "Menunggu Free"
                      : t.plan_type === "pending_disable"
                        ? "Menunggu Nonaktif"
                        : "Free";
              const planCls =
                t.plan_type === "premium"
                  ? "bg-amber-100 text-amber-700"
                  : t.plan_type?.startsWith("pending_")
                    ? "bg-blue-100 text-blue-700"
                    : "bg-gray-100 text-gray-600";
              const isPending = t.plan_type?.startsWith("pending_");
              const isSuspended = t.status !== "active";

              return (
                <div
                  key={t.id}
                  className="card p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
                >
                  <div className="min-w-0">
                    <p className="font-medium">{t.organization_name}</p>
                    <p className="text-xs text-gray-500">
                      {t.subdomain_slug} · {t.total_listings} listing ·{" "}
                      {t.total_users} user
                    </p>
                    <div className="flex gap-1.5 mt-1.5">
                      <span
                        className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium ${planCls}`}
                      >
                        {planLabel}
                      </span>
                      <span
                        className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium ${
                          t.status === "active"
                            ? "bg-green-100 text-green-700"
                            : "bg-red-100 text-red-700"
                        }`}
                      >
                        {t.status === "active" ? "Aktif" : "Nonaktif"}
                      </span>
                    </div>
                  </div>
                  <div className="flex gap-1.5 flex-shrink-0">
                    {isSuspended ? (
                      <button
                        onClick={() => handleTenantAction(t.id, "activate")}
                        className="btn-success text-xs py-1.5 px-3 flex items-center gap-1"
                        title="Aktifkan kembali"
                      >
                        <RefreshCw className="w-3.5 h-3.5" /> Aktifkan
                      </button>
                    ) : (
                      <>
                        <button
                          onClick={() => handleChangePlan(t.id, "free")}
                          className={`text-xs py-1.5 px-3 rounded font-medium transition-colors ${
                            t.plan_type === "free"
                              ? "bg-gray-200 text-gray-500 cursor-default"
                              : "bg-gray-100 hover:bg-gray-200 text-gray-700"
                          }`}
                          disabled={t.plan_type === "free"}
                        >
                          Free
                        </button>
                        <button
                          onClick={() => handleChangePlan(t.id, "premium")}
                          className={`text-xs py-1.5 px-3 rounded font-medium transition-colors ${
                            t.plan_type === "premium"
                              ? "bg-amber-200 text-amber-600 cursor-default"
                              : "bg-amber-100 hover:bg-amber-200 text-amber-700"
                          }`}
                          disabled={t.plan_type === "premium"}
                        >
                          Premium
                        </button>
                        <button
                          onClick={() => handleTenantAction(t.id, "suspend")}
                          className="text-xs py-1.5 px-3 rounded font-medium bg-red-50 hover:bg-red-100 text-red-600 transition-colors flex items-center gap-1"
                          title="Nonaktifkan tenant"
                        >
                          <Ban className="w-3 h-3" /> Disable
                        </button>
                        <button
                          onClick={() =>
                            handleDeleteTenant(t.id, t.organization_name)
                          }
                          className="text-xs py-1.5 px-2 rounded font-medium text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors"
                          title="Hapus tenant permanen"
                        >
                          <Trash2 className="w-3.5 h-3.5" />
                        </button>
                      </>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {/* Overview */}
        {tab === "overview" && dashboard && (
          <div className="card p-5">
            <h3 className="font-semibold mb-3">Distribusi Listing</h3>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
              {Object.entries(dashboard.listings_by_status || {}).map(
                ([k, v]) => (
                  <div
                    key={k}
                    className="text-center p-3 bg-gray-50 rounded-lg"
                  >
                    <p className="text-2xl font-bold">{v}</p>
                    <p className="text-xs text-gray-500 capitalize">{k}</p>
                  </div>
                ),
              )}
            </div>
          </div>
        )}

        {/* Create Tenant */}
        {tab === "create-tenant" && (
          <div className="card p-5 max-w-lg">
            <h3 className="font-semibold text-lg mb-4 flex items-center gap-2">
              <Plus className="w-5 h-5" /> Buat Tenant Baru
            </h3>
            <form onSubmit={handleCreateTenant} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Nama Organisasi *
                </label>
                <input
                  type="text"
                  required
                  className="input-field"
                  value={createForm.organization_name}
                  onChange={(e) =>
                    setCreateForm({
                      ...createForm,
                      organization_name: e.target.value,
                    })
                  }
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Subdomain *
                </label>
                <input
                  type="text"
                  required
                  className="input-field"
                  value={createForm.subdomain_slug}
                  onChange={(e) =>
                    setCreateForm({
                      ...createForm,
                      subdomain_slug: e.target.value,
                    })
                  }
                  placeholder="agencyname"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Nama Admin *
                  </label>
                  <input
                    type="text"
                    required
                    className="input-field"
                    value={createForm.admin_name}
                    onChange={(e) =>
                      setCreateForm({
                        ...createForm,
                        admin_name: e.target.value,
                      })
                    }
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Email Admin *
                  </label>
                  <input
                    type="email"
                    required
                    className="input-field"
                    value={createForm.admin_email}
                    onChange={(e) =>
                      setCreateForm({
                        ...createForm,
                        admin_email: e.target.value,
                      })
                    }
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Telepon Admin *
                  </label>
                  <input
                    type="text"
                    required
                    className="input-field"
                    value={createForm.admin_phone}
                    onChange={(e) =>
                      setCreateForm({
                        ...createForm,
                        admin_phone: e.target.value,
                      })
                    }
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Password Admin *
                  </label>
                  <input
                    type="password"
                    required
                    className="input-field"
                    value={createForm.admin_password}
                    onChange={(e) =>
                      setCreateForm({
                        ...createForm,
                        admin_password: e.target.value,
                      })
                    }
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Paket
                </label>
                <select
                  className="input-field"
                  value={createForm.plan_type}
                  onChange={(e) =>
                    setCreateForm({ ...createForm, plan_type: e.target.value })
                  }
                >
                  <option value="free">Free</option>
                  <option value="premium">Premium</option>
                </select>
              </div>
              <button
                type="submit"
                disabled={creating}
                className="btn-primary flex items-center gap-2"
              >
                <Plus className="w-4 h-4" />{" "}
                {creating ? "Membuat..." : "Buat Tenant"}
              </button>
            </form>
          </div>
        )}

        {/* Audit Logs */}
        {tab === "audit-logs" && (
          <div className="space-y-4">
            <div className="flex gap-3 flex-wrap">
              <input
                type="text"
                placeholder="User ID..."
                className="input-field text-sm max-w-[200px]"
                value={auditFilters.user_id}
                onChange={(e) =>
                  setAuditFilters({ ...auditFilters, user_id: e.target.value })
                }
              />
              <select
                className="input-field text-sm max-w-[180px]"
                value={auditFilters.action}
                onChange={(e) =>
                  setAuditFilters({ ...auditFilters, action: e.target.value })
                }
              >
                <option value="">Semua Aksi</option>
                <option value="create">Create</option>
                <option value="update">Update</option>
                <option value="delete">Delete</option>
                <option value="approve">Approve</option>
                <option value="reject">Reject</option>
              </select>
              <button
                onClick={() => fetchAuditLogs(1)}
                className="btn-secondary flex items-center gap-1 text-sm"
              >
                <Search className="w-4 h-4" /> Filter
              </button>
            </div>
            {auditLogs.length === 0 ? (
              <p className="text-center text-gray-400 py-8">
                {tab === "audit-logs"
                  ? "Klik Filter untuk memuat audit log."
                  : "Tidak ada data."}
              </p>
            ) : (
              <div className="space-y-2">
                {auditLogs.map((log) => (
                  <div key={log.id} className="card p-3 text-sm">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <span
                          className={`px-2 py-0.5 rounded text-xs font-medium ${
                            log.action === "approve"
                              ? "bg-green-100 text-green-700"
                              : log.action === "reject"
                                ? "bg-red-100 text-red-700"
                                : log.action === "delete"
                                  ? "bg-red-100 text-red-700"
                                  : log.action === "create"
                                    ? "bg-blue-100 text-blue-700"
                                    : "bg-gray-100 text-gray-600"
                          }`}
                        >
                          {log.action}
                        </span>
                        <span className="text-gray-500">
                          {log.entity_type}#{log.entity_id?.substring(0, 8)}
                        </span>
                      </div>
                      <span className="text-xs text-gray-400 flex items-center gap-1">
                        <Clock className="w-3 h-3" />{" "}
                        {new Date(log.created_at).toLocaleString("id-ID")}
                      </span>
                    </div>
                    {log.user && (
                      <p className="text-xs text-gray-400 mt-1">
                        By: {log.user.name} ({log.user.role})
                      </p>
                    )}
                  </div>
                ))}
              </div>
            )}
            {auditMeta.total_pages > 1 && (
              <div className="flex justify-center gap-3">
                <button
                  onClick={() => fetchAuditLogs(auditMeta.page - 1)}
                  disabled={auditMeta.page <= 1}
                  className="btn-secondary text-sm"
                >
                  <ChevronLeft className="w-4 h-4" />
                </button>
                <span className="text-sm text-gray-500 py-2">
                  {auditMeta.page}/{auditMeta.total_pages}
                </span>
                <button
                  onClick={() => fetchAuditLogs(auditMeta.page + 1)}
                  disabled={auditMeta.page >= auditMeta.total_pages}
                  className="btn-secondary text-sm"
                >
                  <ChevronRight className="w-4 h-4" />
                </button>
              </div>
            )}
          </div>
        )}

        {/* My Profile */}
        {tab === "me" && authUser && (
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

function StatCard({ icon, label, value, color }) {
  return (
    <div className={`card p-4 flex items-center gap-3 ${color || ""}`}>
      <div className="text-primary-600">{icon}</div>
      <div>
        <p className="text-xs text-gray-500">{label}</p>
        <p className="text-xl font-bold">{value}</p>
      </div>
    </div>
  );
}
