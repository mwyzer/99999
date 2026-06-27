import { Link, useNavigate, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import {
  Building2,
  Menu,
  X,
  LogOut,
  User,
  Heart,
  LayoutDashboard,
  Settings,
} from "lucide-react";
import { useState } from "react";

export default function Navbar() {
  const { isAuthenticated, user, role, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [open, setOpen] = useState(false);

  const isActive = (path) => location.pathname === path;

  const handleLogout = () => {
    logout();
    navigate("/");
  };

  const dashboardLink = () => {
    switch (role) {
      case "platform_admin":
        return "/admin/dashboard";
      case "tenant_admin":
        return "/tenant/dashboard";
      case "salesman":
        return "/salesman/dashboard";
      default:
        return "/";
    }
  };

  return (
    <nav className="bg-white shadow-sm border-b border-gray-100 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          {/* Logo */}
          <div className="flex items-center">
            <Link
              to="/"
              className="flex items-center gap-2 text-xl font-bold text-primary-700"
            >
              <Building2 className="w-6 h-6" />
              <span>PropertyHub</span>
            </Link>
          </div>

          {/* Desktop nav */}
          <div className="hidden md:flex items-center gap-4">
            <Link
              to="/"
              className={`transition-colors text-sm font-medium ${
                isActive("/")
                  ? "text-primary-600"
                  : "text-gray-600 hover:text-primary-600"
              }`}
            >
              Cari Properti
            </Link>

            {isAuthenticated ? (
              <>
                {role === "buyer" && (
                  <Link
                    to="/saved"
                    className={`flex items-center gap-1 text-sm ${
                      isActive("/saved")
                        ? "text-primary-600"
                        : "text-gray-600 hover:text-primary-600"
                    }`}
                  >
                    <Heart className="w-4 h-4" /> Favorit
                  </Link>
                )}
                {(role === "salesman" || role === "tenant_admin") && (
                  <Link
                    to={dashboardLink()}
                    className={`flex items-center gap-1 text-sm ${
                      isActive(dashboardLink())
                        ? "text-primary-600"
                        : "text-gray-600 hover:text-primary-600"
                    }`}
                  >
                    <LayoutDashboard className="w-4 h-4" /> Dashboard
                  </Link>
                )}
                {role === "platform_admin" && (
                  <Link
                    to={dashboardLink()}
                    className={`flex items-center gap-1 text-sm ${
                      isActive("/admin/dashboard")
                        ? "text-primary-600"
                        : "text-gray-600 hover:text-primary-600"
                    }`}
                  >
                    <Settings className="w-4 h-4" /> Admin Panel
                  </Link>
                )}
                <div className="flex items-center gap-3 ml-2 pl-4 border-l">
                  <Link
                    to="/profile"
                    className={`flex items-center gap-1 text-sm ${
                      isActive("/profile")
                        ? "text-primary-600"
                        : "text-gray-600 hover:text-primary-600"
                    }`}
                  >
                    <User className="w-4 h-4" />
                  </Link>
                  {user?.photo_url ? (
                    <img
                      src={user.photo_url}
                      alt=""
                      className="w-8 h-8 rounded-full object-cover"
                    />
                  ) : (
                    <div className="w-8 h-8 rounded-full bg-primary-100 flex items-center justify-center text-primary-600 font-bold text-sm">
                      {user?.name?.charAt(0) || "U"}
                    </div>
                  )}
                  <span className="text-sm text-gray-500">{user?.name}</span>
                  <button
                    onClick={handleLogout}
                    className="text-sm text-red-500 hover:text-red-700 flex items-center gap-1"
                  >
                    <LogOut className="w-4 h-4" /> Keluar
                  </button>
                </div>
              </>
            ) : (
              <>
                <Link
                  to="/login"
                  className="text-sm text-gray-600 hover:text-primary-600 font-medium"
                >
                  Masuk
                </Link>
                <Link to="/register" className="btn-primary text-sm">
                  Daftar
                </Link>
              </>
            )}
          </div>

          {/* Mobile menu button */}
          <div className="md:hidden flex items-center">
            <button onClick={() => setOpen(!open)} className="text-gray-600">
              {open ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile menu */}
      {open && (
        <div className="md:hidden border-t bg-white px-4 py-3 space-y-2">
          <Link
            to="/"
            onClick={() => setOpen(false)}
            className="block text-gray-600 py-2"
          >
            Cari Properti
          </Link>
          {isAuthenticated ? (
            <>
              {role === "buyer" && (
                <Link
                  to="/saved"
                  onClick={() => setOpen(false)}
                  className="block text-gray-600 py-2"
                >
                  Favorit
                </Link>
              )}
              {(role === "salesman" || role === "tenant_admin") && (
                <Link
                  to={dashboardLink()}
                  onClick={() => setOpen(false)}
                  className="block text-gray-600 py-2"
                >
                  Dashboard
                </Link>
              )}
              {role === "platform_admin" && (
                <Link
                  to={dashboardLink()}
                  onClick={() => setOpen(false)}
                  className="block text-gray-600 py-2"
                >
                  Admin Panel
                </Link>
              )}
              <Link
                to="/profile"
                onClick={() => setOpen(false)}
                className="block text-gray-600 py-2"
              >
                Profil Saya
              </Link>
              <button
                onClick={() => {
                  handleLogout();
                  setOpen(false);
                }}
                className="block text-red-500 py-2 w-full text-left"
              >
                Keluar
              </button>
            </>
          ) : (
            <>
              <Link
                to="/login"
                onClick={() => setOpen(false)}
                className="block text-gray-600 py-2"
              >
                Masuk
              </Link>
              <Link
                to="/register"
                onClick={() => setOpen(false)}
                className="block btn-primary text-center"
              >
                Daftar
              </Link>
            </>
          )}
        </div>
      )}
    </nav>
  );
}
