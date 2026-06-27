import { NavLink, useLocation, useNavigate } from "react-router-dom";
import {
  LayoutDashboard,
  Building2,
  CheckCircle,
  Users,
  Home,
  Plus,
  Crown,
  FileText,
  Shield,
  Heart,
  User,
  X,
  Menu,
  MessageSquare,
  PanelLeftClose,
  PanelLeftOpen,
  LogOut,
} from "lucide-react";
import { useState } from "react";
import { useAuth } from "../context/AuthContext";

const menuConfigs = {
  platform_admin: [
    {
      to: "/admin/dashboard",
      icon: LayoutDashboard,
      label: "Overview",
      exact: true,
    },
    {
      to: "/admin/dashboard",
      icon: Building2,
      label: "Tenant",
      hash: "tenants",
    },
    {
      to: "/admin/dashboard",
      icon: CheckCircle,
      label: "Review Listing",
      hash: "pending",
    },
    {
      to: "/admin/dashboard",
      icon: FileText,
      label: "Audit Logs",
      hash: "audit-logs",
    },
    { to: "/admin/dashboard", icon: User, label: "Profil Saya", hash: "me" },
  ],
  tenant_admin: [
    {
      to: "/tenant/dashboard",
      icon: LayoutDashboard,
      label: "Overview",
      exact: true,
    },
    {
      to: "/tenant/dashboard",
      icon: Users,
      label: "Salesman",
      hash: "salesmen",
    },
    {
      to: "/tenant/dashboard",
      icon: Home,
      label: "Listings",
      hash: "listings",
    },
    {
      to: "/tenant/dashboard",
      icon: Crown,
      label: "Paket",
      hash: "subscription",
    },
    {
      to: "/tenant/dashboard",
      icon: Building2,
      label: "Profil Tenant",
      hash: "profile",
    },
    {
      to: "/tenant/dashboard",
      icon: MessageSquare,
      label: "Pertanyaan",
      hash: "inquiries",
    },
    { to: "/tenant/dashboard", icon: User, label: "Profil Saya", hash: "me" },
  ],
  salesman: [
    {
      to: "/salesman/dashboard",
      icon: LayoutDashboard,
      label: "Dashboard",
      exact: true,
    },
    { to: "/salesman/listings/new", icon: Plus, label: "Listing Baru" },
    {
      to: "/salesman/dashboard",
      icon: Home,
      label: "Listing Saya",
      hash: "all",
    },
    { to: "/salesman/dashboard", icon: User, label: "Profil Saya", hash: "me" },
  ],
  buyer: [
    { to: "/saved", icon: Heart, label: "Favorit Saya" },
    { to: "/", icon: Home, label: "Cari Properti" },
    { to: "/profile", icon: User, label: "Profil Saya" },
  ],
};

export default function DashboardSidebar({
  role,
  onClose,
  collapsed,
  onToggleCollapse,
}) {
  const location = useLocation();
  const navigate = useNavigate();
  const { logout } = useAuth();
  const menu = menuConfigs[role] || [];

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  const isActive = (item) => {
    if (item.exact) return location.pathname === item.to && !location.hash;
    if (item.hash) return location.hash === `#${item.hash}`;
    return location.pathname === item.to;
  };

  const handleClick = (item) => {
    if (item.hash) {
      navigate(`${item.to}#${item.hash}`, { replace: true });
    }
    if (onClose) onClose();
  };

  return (
    <aside
      className={`${
        collapsed ? "w-16" : "w-64"
      } min-h-screen bg-white border-r border-gray-200 flex flex-col shrink-0 transition-all duration-200`}
    >
      {/* Brand */}
      <div className="flex items-center justify-between h-16 px-3 border-b border-gray-100">
        {!collapsed && (
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center shrink-0">
              <Building2 className="w-5 h-5 text-white" />
            </div>
            <span className="font-bold text-gray-900 text-lg">PropertyHub</span>
          </div>
        )}
        {collapsed && (
          <div className="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center mx-auto shrink-0">
            <Building2 className="w-5 h-5 text-white" />
          </div>
        )}
        {onClose && (
          <button
            onClick={onClose}
            className="lg:hidden p-1 text-gray-400 hover:text-gray-600"
          >
            <X className="w-5 h-5" />
          </button>
        )}
      </div>

      {/* Role badge */}
      {!collapsed && (
        <div className="px-5 py-3">
          <span className="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-medium rounded-full bg-primary-50 text-primary-700">
            <Shield className="w-3 h-3" />
            {role === "platform_admin" && "Platform Admin"}
            {role === "tenant_admin" && "Tenant Admin"}
            {role === "salesman" && "Salesman"}
            {role === "buyer" && "Buyer"}
          </span>
        </div>
      )}

      {/* Menu */}
      <nav className="flex-1 px-2 py-2 space-y-0.5 overflow-y-auto">
        {menu.map((item) => {
          const active = isActive(item);
          const Icon = item.icon;
          return (
            <NavLink
              key={item.label}
              to={item.hash ? `${item.to}#${item.hash}` : item.to}
              onClick={() => handleClick(item)}
              title={collapsed ? item.label : undefined}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors ${
                collapsed ? "justify-center" : ""
              } ${
                active
                  ? "bg-primary-50 text-primary-700"
                  : "text-gray-600 hover:bg-gray-50 hover:text-gray-900"
              }`}
            >
              <Icon
                className={`w-5 h-5 shrink-0 ${active ? "text-primary-600" : "text-gray-400"}`}
              />
              {!collapsed && item.label}
            </NavLink>
          );
        })}
      </nav>

      {/* Logout */}
      <div className="px-2 py-2 border-t border-gray-100">
        <button
          onClick={handleLogout}
          className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm text-red-500 hover:bg-red-50 hover:text-red-600 transition-colors"
          title="Logout"
        >
          <LogOut className="w-5 h-5 shrink-0" />
          {!collapsed && "Logout"}
        </button>
      </div>

      {/* Collapse toggle */}
      <div className="px-2 py-2 border-t border-gray-100">
        <button
          onClick={onToggleCollapse}
          className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm text-gray-400 hover:bg-gray-50 hover:text-gray-600 transition-colors"
          title={collapsed ? "Buka sidebar" : "Tutup sidebar"}
        >
          {collapsed ? (
            <PanelLeftOpen className="w-5 h-5 shrink-0" />
          ) : (
            <PanelLeftClose className="w-5 h-5 shrink-0" />
          )}
          {!collapsed && "Tutup"}
        </button>
      </div>

      {/* Footer */}
      {!collapsed && (
        <div className="px-5 py-3 border-t border-gray-100">
          <p className="text-xs text-gray-400">© 2026 PropertyHub</p>
        </div>
      )}
    </aside>
  );
}

export function DashboardLayout({ children, role }) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  return (
    <div className="flex min-h-[calc(100vh-64px)]">
      {/* Mobile overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar — hidden on mobile, shown on lg+ */}
      <div
        className={`fixed lg:sticky top-0 z-50 h-screen transition-transform lg:translate-x-0 ${
          sidebarOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <DashboardSidebar
          role={role}
          onClose={() => setSidebarOpen(false)}
          collapsed={sidebarCollapsed}
          onToggleCollapse={() => setSidebarCollapsed((prev) => !prev)}
        />
      </div>

      {/* Main content */}
      <div className="flex-1 min-w-0">
        {/* Mobile header bar */}
        <div className="lg:hidden flex items-center gap-3 px-4 py-3 border-b border-gray-100 bg-white sticky top-0 z-30">
          <button
            onClick={() => setSidebarOpen(true)}
            className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg"
          >
            <Menu className="w-5 h-5" />
          </button>
          <span className="font-semibold text-gray-900">PropertyHub</span>
        </div>
        {children}
      </div>
    </div>
  );
}
