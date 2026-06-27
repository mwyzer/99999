import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { Building2, Eye, EyeOff } from "lucide-react";
import toast from "react-hot-toast";

export default function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [form, setForm] = useState({ email: "", password: "" });
  const [submitting, setSubmitting] = useState(false);
  const [showPw, setShowPw] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      const user = await login(form.email, form.password);
      toast.success(`Selamat datang, ${user.user.name}!`);
      // Redirect based on role
      switch (user.user.role) {
        case "platform_admin":
          navigate("/admin/dashboard");
          break;
        case "tenant_admin":
          navigate("/tenant/dashboard");
          break;
        case "salesman":
          navigate("/salesman/dashboard");
          break;
        default:
          navigate("/");
      }
    } catch (err) {
      toast.error(err.response?.data?.error?.message || "Login gagal.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <Building2 className="w-12 h-12 text-primary-600 mx-auto mb-3" />
          <h1 className="text-2xl font-bold text-gray-900">
            Masuk ke PropertyHub
          </h1>
          <p className="text-gray-500 mt-1">
            Masuk untuk mengelola properti Anda
          </p>
        </div>

        <form onSubmit={handleSubmit} className="card p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              type="email"
              required
              className="input-field"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              placeholder="nama@email.com"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Password
            </label>
            <div className="relative">
              <input
                type={showPw ? "text" : "password"}
                required
                className="input-field pr-10"
                value={form.password}
                onChange={(e) => setForm({ ...form, password: e.target.value })}
                placeholder="Minimal 8 karakter"
              />
              <button
                type="button"
                onClick={() => setShowPw(!showPw)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
              >
                {showPw ? (
                  <EyeOff className="w-4 h-4" />
                ) : (
                  <Eye className="w-4 h-4" />
                )}
              </button>
            </div>
          </div>
          <button
            type="submit"
            disabled={submitting}
            className="btn-primary w-full"
          >
            {submitting ? "Memproses..." : "Masuk"}
          </button>
        </form>

        <p className="text-center mt-4 text-sm text-gray-500">
          Belum punya akun?{" "}
          <Link
            to="/register"
            className="text-primary-600 hover:text-primary-700 font-medium"
          >
            Daftar di sini
          </Link>
        </p>

        {/* Demo credentials */}
        <div className="mt-6 p-3 bg-blue-50 rounded-lg text-xs text-blue-800">
          <p className="font-semibold mb-1">Akun Demo:</p>
          <p>Admin Platform: admin@propertyhub.id / Admin@123</p>
          <p>Tenant Admin: budi@propertijaya.id / Budi@123</p>
          <p>Salesman: andi@propertijaya.id / Andi@123</p>
          <p>Buyer: rina@email.com / Rina@123</p>
        </div>
      </div>
    </div>
  );
}
