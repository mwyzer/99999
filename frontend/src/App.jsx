import { Routes, Route, Navigate, useLocation } from "react-router-dom";
import { useAuth } from "./context/AuthContext";
import Navbar from "./components/Navbar";
import Footer from "./components/Footer";
import Loading from "./components/Loading";

// Pages
import HomePage from "./pages/HomePage";
import PropertyDetailPage from "./pages/PropertyDetailPage";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import BuyerSavedPage from "./pages/BuyerSavedPage";
import BuyerInquiriesPage from "./pages/BuyerInquiriesPage";
import SalesmanDashboard from "./pages/SalesmanDashboard";
import SalesmanListingForm from "./pages/SalesmanListingForm";
import TenantDashboard from "./pages/TenantDashboard";
import AdminDashboard from "./pages/AdminDashboard";
import ProfilePage from "./pages/ProfilePage";
import NotFoundPage from "./pages/NotFoundPage";

function ProtectedRoute({ children, roles }) {
  const { isAuthenticated, role, loading } = useAuth();
  if (loading) return <Loading fullScreen />;
  if (!isAuthenticated) return <Navigate to="/login" />;
  if (roles && !roles.includes(role)) return <Navigate to="/" />;
  return children;
}

function GuestRoute({ children }) {
  const { isAuthenticated, loading } = useAuth();
  if (loading) return <Loading fullScreen />;
  if (isAuthenticated) return <Navigate to="/" />;
  return children;
}

export default function App() {
  const { loading } = useAuth();
  const location = useLocation();
  if (loading) return <Loading fullScreen />;

  // Dashboard routes have their own sidebar — hide global Navbar/Footer
  const isDashboardRoute =
    location.pathname.startsWith("/admin/") ||
    location.pathname.startsWith("/tenant/") ||
    location.pathname.startsWith("/salesman/") ||
    location.pathname === "/saved" ||
    location.pathname === "/inquiries";

  return (
    <div className="min-h-screen flex flex-col">
      {!isDashboardRoute && <Navbar />}
      <main className="flex-1">
        <Routes>
          {/* Public */}
          <Route path="/" element={<HomePage />} />
          <Route path="/properties/:id" element={<PropertyDetailPage />} />
          <Route
            path="/login"
            element={
              <GuestRoute>
                <LoginPage />
              </GuestRoute>
            }
          />
          <Route
            path="/register"
            element={
              <GuestRoute>
                <RegisterPage />
              </GuestRoute>
            }
          />

          {/* Buyer */}
          <Route
            path="/saved"
            element={
              <ProtectedRoute roles={["buyer"]}>
                <BuyerSavedPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/inquiries"
            element={
              <ProtectedRoute roles={["buyer"]}>
                <BuyerInquiriesPage />
              </ProtectedRoute>
            }
          />

          {/* Salesman */}
          <Route
            path="/salesman/dashboard"
            element={
              <ProtectedRoute roles={["salesman", "tenant_admin"]}>
                <SalesmanDashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/salesman/listings/new"
            element={
              <ProtectedRoute roles={["salesman", "tenant_admin"]}>
                <SalesmanListingForm />
              </ProtectedRoute>
            }
          />
          <Route
            path="/salesman/listings/:id/edit"
            element={
              <ProtectedRoute roles={["salesman", "tenant_admin"]}>
                <SalesmanListingForm />
              </ProtectedRoute>
            }
          />

          {/* Tenant Admin */}
          <Route
            path="/tenant/dashboard"
            element={
              <ProtectedRoute roles={["tenant_admin"]}>
                <TenantDashboard />
              </ProtectedRoute>
            }
          />

          {/* Platform Admin */}
          <Route
            path="/admin/dashboard"
            element={
              <ProtectedRoute roles={["platform_admin"]}>
                <AdminDashboard />
              </ProtectedRoute>
            }
          />

          {/* Profile — all authenticated roles */}
          <Route
            path="/profile"
            element={
              <ProtectedRoute
                roles={["buyer", "salesman", "tenant_admin", "platform_admin"]}
              >
                <ProfilePage />
              </ProtectedRoute>
            }
          />

          {/* 404 */}
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </main>
      {!isDashboardRoute && <Footer />}
    </div>
  );
}
