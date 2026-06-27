import { useState, useEffect } from "react";
import { profileAPI } from "../api";
import { useAuth } from "../context/AuthContext";
import Loading from "../components/Loading";
import {
  User,
  Mail,
  Phone,
  Shield,
  Building2,
  Edit2,
  Save,
} from "lucide-react";
import toast from "react-hot-toast";

const ROLE_LABELS = {
  buyer: "Buyer / Pencari Properti",
  salesman: "Salesman / Agen",
  tenant_admin: "Admin Agency / Pemilik",
  platform_admin: "Admin Platform",
};

export default function ProfilePage() {
  const { user: authUser, refreshUser } = useAuth();
  const [profile, setProfile] = useState(null);
  const [editing, setEditing] = useState(false);
  const [saving, setSaving] = useState(false);
  const [form, setForm] = useState({ name: "", phone: "" });

  useEffect(() => {
    if (authUser) {
      setProfile(authUser);
      setForm({ name: authUser.name || "", phone: authUser.phone || "" });
    }
  }, [authUser]);

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    try {
      await profileAPI.update({ name: form.name, phone: form.phone });
      toast.success("Profil berhasil diperbarui!");
      await refreshUser();
      setEditing(false);
    } catch (err) {
      toast.error(err.response?.data?.error?.message || "Gagal menyimpan.");
    } finally {
      setSaving(false);
    }
  };

  if (!profile) return <Loading fullScreen />;

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold flex items-center gap-2 mb-6">
        <User className="w-6 h-6 text-primary-600" /> Profil Saya
      </h1>

      <div className="space-y-5">
        {/* Info card */}
        <div className="card p-5">
          <div className="flex items-center gap-4 mb-5">
            {profile.photo_url ? (
              <img
                src={profile.photo_url}
                alt=""
                className="w-16 h-16 rounded-full object-cover"
              />
            ) : (
              <div className="w-16 h-16 rounded-full bg-primary-100 flex items-center justify-center text-primary-600 font-bold text-xl">
                {profile.name?.charAt(0) || "U"}
              </div>
            )}
            <div>
              <h2 className="font-semibold text-lg">{profile.name}</h2>
              <p className="text-sm text-gray-500">{profile.email}</p>
            </div>
          </div>

          {/* Info grid */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
            <div className="flex items-center gap-2 text-gray-600">
              <Mail className="w-4 h-4 text-gray-400" />
              <span>{profile.email}</span>
            </div>
            <div className="flex items-center gap-2 text-gray-600">
              <Phone className="w-4 h-4 text-gray-400" />
              <span>{profile.phone || "—"}</span>
            </div>
            <div className="flex items-center gap-2 text-gray-600">
              <Shield className="w-4 h-4 text-gray-400" />
              <span>{ROLE_LABELS[profile.role] || profile.role}</span>
            </div>
            {profile.tenant && (
              <div className="flex items-center gap-2 text-gray-600">
                <Building2 className="w-4 h-4 text-gray-400" />
                <span>{profile.tenant.name}</span>
              </div>
            )}
          </div>

          {!editing && (
            <button
              onClick={() => setEditing(true)}
              className="btn-secondary flex items-center gap-1 mt-5 text-sm"
            >
              <Edit2 className="w-4 h-4" /> Edit Profil
            </button>
          )}
        </div>

        {/* Edit form */}
        {editing && (
          <form onSubmit={handleSave} className="card p-5 space-y-4">
            <h3 className="font-semibold">Edit Profil</h3>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Nama Lengkap
              </label>
              <input
                type="text"
                required
                className="input-field"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Nomor Telepon
              </label>
              <input
                type="text"
                className="input-field"
                value={form.phone}
                onChange={(e) => setForm({ ...form, phone: e.target.value })}
                placeholder="0812XXXXXXXX"
              />
            </div>
            <div className="flex gap-2">
              <button
                type="submit"
                disabled={saving}
                className="btn-primary flex items-center gap-1"
              >
                <Save className="w-4 h-4" />
                {saving ? "Menyimpan..." : "Simpan"}
              </button>
              <button
                type="button"
                onClick={() => setEditing(false)}
                className="btn-secondary"
              >
                Batal
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
