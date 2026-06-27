import { useState, useEffect, useMemo } from "react";
import { useLocation } from "react-router-dom";
import { tenantAPI, profileAPI } from "../api";
import { useAuth } from "../context/AuthContext";
import Loading from "../components/Loading";
import { DashboardLayout } from "../components/DashboardSidebar";
import {
  Users,
  Home,
  Crown,
  UserPlus,
  Trash2,
  Edit2,
  FileText,
  MapPin,
  MessageSquare,
  MessageCircle,
  Clock,
  Archive,
  CheckCircle,
  User,
  Mail,
  Phone,
  Save,
} from "lucide-react";
import toast from "react-hot-toast";

export default function TenantDashboard() {
  const location = useLocation();
  const { user: authUser, refreshUser } = useAuth();
  const tab = useMemo(
    () => location.hash?.replace("#", "") || "overview",
    [location.hash],
  );
  const [dashboard, setDashboard] = useState(null);
  const [salesmen, setSalesmen] = useState([]);
  const [subscription, setSubscription] = useState(null);
  const [loading, setLoading] = useState(true);
  const [newSalesman, setNewSalesman] = useState({
    name: "",
    email: "",
    phone: "",
    password: "",
  });
  // Listings tab
  const [listings, setListings] = useState([]);
  const [listingsMeta, setListingsMeta] = useState({ page: 1, total_pages: 1 });
  const [listingStatusFilter, setListingStatusFilter] = useState("");
  // Profile tab
  const [profile, setProfile] = useState(null);
  const [profileForm, setProfileForm] = useState({
    organization_name: "",
    description: "",
    phone: "",
    address: "",
  });
  const [savingProfile, setSavingProfile] = useState(false);
  // My Profile (me)
  const [meForm, setMeForm] = useState({ name: "", phone: "" });
  const [savingMe, setSavingMe] = useState(false);
  // Upgrade
  const [upgrading, setUpgrading] = useState(false);
  // Inquiries
  const [inquiries, setInquiries] = useState([]);
  const [inqLoading, setInqLoading] = useState(false);
  const [inqMeta, setInqMeta] = useState({ page: 1, total_pages: 1 });

  const fetchAll = async () => {
    try {
      const [d, s, sub] = await Promise.all([
        tenantAPI.dashboard(),
        tenantAPI.listSalesmen({ per_page: 50 }),
        tenantAPI.getSubscription(),
      ]);
      setDashboard(d.data.data);
      setSalesmen(s.data.data);
      setSubscription(sub.data.data);
    } catch {
      toast.error("Gagal memuat data.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAll();
  }, []);

  const handleAddSalesman = async (e) => {
    e.preventDefault();
    try {
      await tenantAPI.addSalesman(newSalesman);
      toast.success("Salesman berhasil ditambahkan!");
      setNewSalesman({ name: "", email: "", phone: "", password: "" });
      fetchAll();
    } catch (err) {
      toast.error(
        err.response?.data?.error?.message || "Gagal menambah salesman.",
      );
    }
  };

  const handleRemoveSalesman = async (id) => {
    if (!confirm("Yakin nonaktifkan salesman ini?")) return;
    try {
      await tenantAPI.removeSalesman(id);
      toast.success("Salesman dinonaktifkan.");
      fetchAll();
    } catch {
      toast.error("Gagal.");
    }
  };

  const handleUpgrade = async () => {
    setUpgrading(true);
    try {
      await tenantAPI.requestUpgrade();
      toast.success("Permintaan upgrade telah dikirim!");
      setSubscription((prev) => ({ ...prev, plan_type: "pending_upgrade" }));
    } catch {
      toast.error("Gagal mengirim permintaan upgrade.");
    } finally {
      setUpgrading(false);
    }
  };

  const fetchListings = async (page = 1) => {
    try {
      const params = { page, per_page: 12 };
      if (listingStatusFilter) params.status = listingStatusFilter;
      const { data } = await tenantAPI.listListings(params);
      setListings(data.data);
      setListingsMeta(data.meta);
    } catch {
      toast.error("Gagal memuat listing.");
    }
  };

  const fetchProfile = async () => {
    try {
      const { data } = await tenantAPI.getProfile();
      const p = data.data;
      setProfile(p);
      setProfileForm({
        organization_name: p.organization_name || "",
        description: p.description || "",
        phone: p.phone || "",
        address: p.address || "",
      });
    } catch {
      toast.error("Gagal memuat profil.");
    }
  };

  const handleSaveProfile = async (e) => {
    e.preventDefault();
    setSavingProfile(true);
    try {
      await tenantAPI.updateProfile(profileForm);
      toast.success("Profil berhasil diperbarui.");
      fetchProfile();
    } catch {
      toast.error("Gagal menyimpan profil.");
    } finally {
      setSavingProfile(false);
    }
  };

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

  useEffect(() => {
    if (tab === "listings") fetchListings();
    if (tab === "profile") fetchProfile();
    if (tab === "inquiries") fetchInquiries();
  }, [tab]);

  const fetchInquiries = async (page = 1) => {
    setInqLoading(true);
    try {
      const { data } = await tenantAPI.listInquiries({ page, per_page: 20 });
      setInquiries(data.data);
      setInqMeta(data.meta);
    } catch {
      toast.error("Gagal memuat pertanyaan.");
    } finally {
      setInqLoading(false);
    }
  };

  if (loading) return <Loading fullScreen />;

  return (
    <DashboardLayout role="tenant_admin">
      <div className="p-4 sm:p-6 lg:p-8">
        {/* Stats */}
        {dashboard && (
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-6">
            <StatCard
              icon={<Home className="w-5 h-5" />}
              label="Total Listing"
              value={dashboard.total_listings}
            />
            <StatCard
              icon={<Home className="w-5 h-5" />}
              label="Aktif"
              value={dashboard.active_listings}
            />
            <StatCard
              icon={<Users className="w-5 h-5" />}
              label="Salesman"
              value={`${dashboard.total_salesmen}/${dashboard.max_salesmen}`}
            />
            <StatCard
              icon={<Crown className="w-5 h-5" />}
              label="Paket"
              value={
                <span
                  className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-bold ${
                    dashboard.plan?.type === "premium"
                      ? "bg-amber-100 text-amber-700"
                      : "bg-gray-100 text-gray-600"
                  }`}
                >
                  {dashboard.plan?.type === "premium" ? "Premium" : "Free"}
                </span>
              }
            />
          </div>
        )}

        {/* Overview */}
        {tab === "overview" && dashboard && (
          <>
            <div className="card p-5 mb-6">
              <h3 className="font-semibold mb-3">Status Listing</h3>
              <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
                {Object.entries(dashboard.status_breakdown || {}).map(
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

            {/* Subscription status card */}
            {subscription && (
              <div className="card p-5 mb-6">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="font-semibold text-lg flex items-center gap-2">
                    <Crown className="w-5 h-5 text-amber-500" />
                    Status Langganan
                  </h3>
                  <span
                    className={`inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-bold ${
                      subscription.plan_type === "premium"
                        ? "bg-amber-100 text-amber-700"
                        : subscription.plan_type === "pending_upgrade"
                          ? "bg-blue-100 text-blue-700"
                          : "bg-gray-100 text-gray-600"
                    }`}
                  >
                    {subscription.plan_type === "premium"
                      ? "Premium"
                      : subscription.plan_type === "pending_upgrade"
                        ? "Menunggu Upgrade"
                        : "Free"}
                  </span>
                </div>
                <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-4">
                  <div className="text-center p-3 bg-gray-50 rounded-lg">
                    <p className="text-xl font-bold">
                      {subscription.usage?.salesmen_used || 0}/
                      {subscription.max_salesmen}
                    </p>
                    <p className="text-xs text-gray-500">Salesman</p>
                  </div>
                  <div className="text-center p-3 bg-gray-50 rounded-lg">
                    <p className="text-xl font-bold">
                      {subscription.max_listings_per_salesman === 999999
                        ? "∞"
                        : subscription.max_listings_per_salesman}
                    </p>
                    <p className="text-xs text-gray-500">
                      Max Listing/Salesman
                    </p>
                  </div>
                  <div className="text-center p-3 bg-gray-50 rounded-lg">
                    <p className="text-xl font-bold">
                      {subscription.usage?.total_active_listings || 0}
                    </p>
                    <p className="text-xs text-gray-500">Listing Aktif</p>
                  </div>
                  <div className="text-center p-3 bg-gray-50 rounded-lg">
                    <p className="text-xl font-bold">
                      {subscription.plan_type === "premium"
                        ? "Unlimited"
                        : subscription.max_listings_per_salesman === 999999
                          ? "Unlimited"
                          : subscription.max_salesmen *
                            subscription.max_listings_per_salesman}
                    </p>
                    <p className="text-xs text-gray-500">Total Kuota Listing</p>
                  </div>
                </div>
                {subscription.plan_type === "free" && (
                  <button
                    onClick={handleUpgrade}
                    disabled={upgrading}
                    className="btn-primary flex items-center gap-2"
                  >
                    <Crown className="w-4 h-4" />
                    {upgrading ? "Mengirim..." : "Upgrade ke Premium"}
                  </button>
                )}
                {subscription.plan_type === "pending_upgrade" && (
                  <p className="text-sm text-blue-600 bg-blue-50 px-3 py-2 rounded-lg">
                    Permintaan upgrade Anda sedang diproses. Tim kami akan
                    segera menghubungi.
                  </p>
                )}
                {subscription.plan_type === "premium" && (
                  <p className="text-sm text-amber-600 bg-amber-50 px-3 py-2 rounded-lg">
                    Anda sudah menggunakan paket Premium. Nikmati semua fitur
                    tanpa batasan!
                  </p>
                )}
              </div>
            )}
          </>
        )}

        {/* Salesmen */}
        {tab === "salesmen" && (
          <div className="space-y-6">
            <div className="card p-5">
              <h3 className="font-semibold mb-3">Tambah Salesman Baru</h3>
              <form
                onSubmit={handleAddSalesman}
                className="grid grid-cols-1 sm:grid-cols-4 gap-3"
              >
                <input
                  type="text"
                  required
                  className="input-field text-sm"
                  placeholder="Nama"
                  value={newSalesman.name}
                  onChange={(e) =>
                    setNewSalesman({ ...newSalesman, name: e.target.value })
                  }
                />
                <input
                  type="email"
                  required
                  className="input-field text-sm"
                  placeholder="Email"
                  value={newSalesman.email}
                  onChange={(e) =>
                    setNewSalesman({ ...newSalesman, email: e.target.value })
                  }
                />
                <input
                  type="text"
                  required
                  className="input-field text-sm"
                  placeholder="Telepon"
                  value={newSalesman.phone}
                  onChange={(e) =>
                    setNewSalesman({ ...newSalesman, phone: e.target.value })
                  }
                />
                <div className="flex gap-2">
                  <input
                    type="password"
                    required
                    className="input-field text-sm flex-1"
                    placeholder="Password"
                    value={newSalesman.password}
                    onChange={(e) =>
                      setNewSalesman({
                        ...newSalesman,
                        password: e.target.value,
                      })
                    }
                  />
                  <button
                    type="submit"
                    className="btn-primary text-sm flex items-center gap-1"
                  >
                    <UserPlus className="w-4 h-4" /> Tambah
                  </button>
                </div>
              </form>
            </div>

            <div className="space-y-3">
              {salesmen.map((s) => (
                <div
                  key={s.id}
                  className="card p-4 flex items-center justify-between"
                >
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-full bg-primary-100 flex items-center justify-center text-primary-600 font-bold">
                      {s.name.charAt(0)}
                    </div>
                    <div>
                      <p className="font-medium text-sm">{s.name}</p>
                      <p className="text-xs text-gray-500">
                        {s.email} · {s.phone}
                      </p>
                      <p className="text-xs text-gray-400">
                        {s.listing_count?.active || 0} listing aktif
                      </p>
                    </div>
                  </div>
                  <button
                    onClick={() => handleRemoveSalesman(s.id)}
                    className="text-red-500 hover:text-red-700 p-2"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Subscription */}
        {tab === "subscription" && subscription && (
          <div className="card p-5">
            <div className="flex items-center justify-between mb-4">
              <h3 className="font-semibold text-lg">Paket Langganan</h3>
              <span
                className={`inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-bold ${
                  subscription.plan_type === "premium"
                    ? "bg-amber-100 text-amber-700"
                    : subscription.plan_type === "pending_upgrade"
                      ? "bg-blue-100 text-blue-700"
                      : "bg-gray-100 text-gray-600"
                }`}
              >
                {subscription.plan_type === "premium"
                  ? "Premium"
                  : subscription.plan_type === "pending_upgrade"
                    ? "Menunggu Upgrade"
                    : "Free"}
              </span>
            </div>
            <div className="grid grid-cols-2 gap-4 mb-4">
              <div className="p-3 bg-gray-50 rounded-lg">
                <p className="text-xs text-gray-500">Paket</p>
                <p className="font-bold text-lg capitalize">
                  {subscription.plan_type === "pending_upgrade"
                    ? "Free (Menunggu Upgrade)"
                    : subscription.plan_type}
                </p>
              </div>
              <div className="p-3 bg-gray-50 rounded-lg">
                <p className="text-xs text-gray-500">Salesman</p>
                <p className="font-bold text-lg">
                  {subscription.usage?.salesmen_used}/
                  {subscription.max_salesmen}
                </p>
              </div>
              <div className="p-3 bg-gray-50 rounded-lg">
                <p className="text-xs text-gray-500">Max Listing/Salesman</p>
                <p className="font-bold text-lg">
                  {subscription.max_listings_per_salesman === 999999
                    ? "Unlimited"
                    : subscription.max_listings_per_salesman}
                </p>
              </div>
              <div className="p-3 bg-gray-50 rounded-lg">
                <p className="text-xs text-gray-500">Total Listing Aktif</p>
                <p className="font-bold text-lg">
                  {subscription.usage?.total_active_listings}
                </p>
              </div>
            </div>
            {subscription.plan_type === "free" && (
              <button
                onClick={handleUpgrade}
                disabled={upgrading}
                className="btn-primary flex items-center gap-2"
              >
                <Crown className="w-4 h-4" />
                {upgrading ? "Mengirim..." : "Upgrade ke Premium"}
              </button>
            )}
            {subscription.plan_type === "pending_upgrade" && (
              <p className="text-sm text-blue-600 bg-blue-50 px-3 py-2 rounded-lg">
                Permintaan upgrade Anda sedang diproses. Tim kami akan segera
                menghubungi.
              </p>
            )}
            {subscription.plan_type === "premium" && (
              <p className="text-sm text-amber-600 bg-amber-50 px-3 py-2 rounded-lg">
                Anda sudah menggunakan paket Premium. Nikmati semua fitur tanpa
                batasan!
              </p>
            )}
          </div>
        )}

        {/* Listings */}
        {tab === "listings" && (
          <div className="space-y-4">
            <div className="flex gap-2 flex-wrap">
              {[
                "",
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
                  onClick={() => {
                    setListingStatusFilter(s);
                    setTimeout(() => fetchListings(1), 0);
                  }}
                  className={`px-3 py-1 text-sm rounded-full border transition-colors ${
                    listingStatusFilter === s
                      ? "bg-primary-600 text-white border-primary-600"
                      : "border-gray-300 text-gray-600 hover:border-gray-400"
                  }`}
                >
                  {s || "Semua"}
                </button>
              ))}
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {listings.map((l) => (
                <div key={l.id} className="card p-4">
                  <div className="flex items-center justify-between mb-2">
                    <span
                      className={`badge text-xs ${
                        l.status === "approved"
                          ? "badge-success"
                          : l.status === "pending"
                            ? "badge-warning"
                            : l.status === "rejected"
                              ? "badge-error"
                              : "badge"
                      }`}
                    >
                      {l.status}
                    </span>
                    <span className="text-xs text-gray-400">
                      {new Date(l.created_at).toLocaleDateString("id-ID")}
                    </span>
                  </div>
                  <h4 className="font-medium text-sm mb-1 line-clamp-1">
                    {l.title}
                  </h4>
                  <p className="text-primary-600 font-bold text-sm mb-2">
                    {l.price}
                  </p>
                  <div className="flex items-center gap-1 text-xs text-gray-500">
                    <MapPin className="w-3 h-3" />
                    {l.salesman?.name || "—"}
                  </div>
                </div>
              ))}
            </div>
            {listings.length === 0 && (
              <p className="text-center text-gray-400 py-8">
                Belum ada listing.
              </p>
            )}
            {listingsMeta.total_pages > 1 && (
              <div className="flex justify-center gap-3 mt-4">
                <button
                  onClick={() => fetchListings(listingsMeta.page - 1)}
                  disabled={listingsMeta.page <= 1}
                  className="btn-secondary text-sm"
                >
                  ← Sebelumnya
                </button>
                <span className="text-sm text-gray-500 py-2">
                  {listingsMeta.page}/{listingsMeta.total_pages}
                </span>
                <button
                  onClick={() => fetchListings(listingsMeta.page + 1)}
                  disabled={listingsMeta.page >= listingsMeta.total_pages}
                  className="btn-secondary text-sm"
                >
                  Selanjutnya →
                </button>
              </div>
            )}
          </div>
        )}

        {/* Profile */}
        {tab === "profile" && (
          <div className="card p-5 max-w-lg">
            <h3 className="font-semibold text-lg mb-4">Profil Agency</h3>
            {profile && (
              <form onSubmit={handleSaveProfile} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Nama Organisasi
                  </label>
                  <input
                    type="text"
                    required
                    className="input-field"
                    value={profileForm.organization_name}
                    onChange={(e) =>
                      setProfileForm({
                        ...profileForm,
                        organization_name: e.target.value,
                      })
                    }
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Deskripsi
                  </label>
                  <textarea
                    className="input-field"
                    rows={3}
                    value={profileForm.description}
                    onChange={(e) =>
                      setProfileForm({
                        ...profileForm,
                        description: e.target.value,
                      })
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
                    value={profileForm.phone}
                    onChange={(e) =>
                      setProfileForm({ ...profileForm, phone: e.target.value })
                    }
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Alamat
                  </label>
                  <textarea
                    className="input-field"
                    rows={2}
                    value={profileForm.address}
                    onChange={(e) =>
                      setProfileForm({
                        ...profileForm,
                        address: e.target.value,
                      })
                    }
                  />
                </div>
                <button
                  type="submit"
                  disabled={savingProfile}
                  className="btn-primary flex items-center gap-2"
                >
                  <Edit2 className="w-4 h-4" />
                  {savingProfile ? "Menyimpan..." : "Simpan Perubahan"}
                </button>
              </form>
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
                {authUser.role && (
                  <span className="inline-block mt-1 text-xs px-2 py-0.5 rounded-full bg-primary-50 text-primary-700">
                    {authUser.role === "tenant_admin"
                      ? "Admin Agency"
                      : authUser.role}
                  </span>
                )}
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

        {/* Inquiries */}
        {tab === "inquiries" && (
          <div className="space-y-4">
            {inqLoading ? (
              <Loading />
            ) : inquiries.length === 0 ? (
              <p className="text-center text-gray-400 py-8">
                Belum ada pertanyaan.
              </p>
            ) : (
              inquiries.map((inq) => {
                const StatusIcon =
                  inq.status === "unread"
                    ? Clock
                    : inq.status === "read"
                      ? MessageCircle
                      : inq.status === "replied"
                        ? CheckCircle
                        : Archive;
                const statusCls =
                  inq.status === "unread"
                    ? "bg-yellow-100 text-yellow-700"
                    : inq.status === "read"
                      ? "bg-blue-100 text-blue-700"
                      : inq.status === "replied"
                        ? "bg-green-100 text-green-700"
                        : "bg-gray-100 text-gray-600";
                const statusLabel =
                  inq.status === "unread"
                    ? "Belum Dibaca"
                    : inq.status === "read"
                      ? "Dibaca"
                      : inq.status === "replied"
                        ? "Dibalas"
                        : "Ditutup";
                return (
                  <div key={inq.id} className="card p-4">
                    <div className="flex items-start justify-between mb-2">
                      <div>
                        <p className="font-medium text-sm">
                          {inq.buyer?.name || "Buyer"}
                        </p>
                        <p className="text-xs text-gray-500">
                          {inq.buyer?.email}
                        </p>
                      </div>
                      <span
                        className={`inline-flex items-center gap-1 px-2 py-0.5 text-xs rounded-full ${statusCls}`}
                      >
                        <StatusIcon className="w-3 h-3" />
                        {statusLabel}
                      </span>
                    </div>
                    {inq.property && (
                      <p className="text-xs text-primary-600 mb-1">
                        Listing: {inq.property.title}
                      </p>
                    )}
                    <p className="text-sm text-gray-700">
                      {inq.message || "—"}
                    </p>
                    <p className="text-xs text-gray-400 mt-2">
                      {new Date(inq.created_at).toLocaleDateString("id-ID", {
                        day: "numeric",
                        month: "long",
                        year: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </p>
                  </div>
                );
              })
            )}
            {inqMeta.total_pages > 1 && (
              <div className="flex justify-center gap-3 mt-4">
                <button
                  onClick={() => fetchInquiries(inqMeta.page - 1)}
                  disabled={inqMeta.page <= 1}
                  className="btn-secondary text-sm"
                >
                  ← Sebelumnya
                </button>
                <span className="text-sm text-gray-500 py-2">
                  {inqMeta.page}/{inqMeta.total_pages}
                </span>
                <button
                  onClick={() => fetchInquiries(inqMeta.page + 1)}
                  disabled={inqMeta.page >= inqMeta.total_pages}
                  className="btn-secondary text-sm"
                >
                  Selanjutnya →
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}

function StatCard({ icon, label, value }) {
  return (
    <div className="card p-4 flex items-center gap-3">
      <div className="text-primary-600">{icon}</div>
      <div>
        <p className="text-xs text-gray-500">{label}</p>
        <p className="text-xl font-bold">{value}</p>
      </div>
    </div>
  );
}
