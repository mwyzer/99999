import axios from "axios";

const api = axios.create({
  baseURL: "/api/v1",
  headers: { "Content-Type": "application/json" },
  timeout: 30000,
});

// Request interceptor — attach JWT
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor — handle 401, format errors
api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      if (window.location.pathname !== "/login") {
        window.location.href = "/login";
      }
    }
    return Promise.reject(err);
  },
);

// ==================== AUTH ====================
export const authAPI = {
  register: (data) => api.post("/auth/register", data),
  login: (data) => api.post("/auth/login", data),
};

// ==================== PROFILE ====================
export const profileAPI = {
  get: () => api.get("/me/profile"),
  update: (data) => api.put("/me/profile", data),
};

// ==================== PUBLIC ====================
export const publicAPI = {
  listProperties: (params) => api.get("/properties", { params }),
  getProperty: (id) => api.get(`/properties/${id}`),
  featured: (params) => api.get("/properties/featured", { params }),
  nearby: (params) => api.get("/properties/nearby", { params }),
  cities: () => api.get("/locations/cities"),
  propertyTypes: () => api.get("/property-types"),
  facilities: () => api.get("/facilities"),
  locations: (params) => api.get("/locations", { params }),
  logView: (id) => api.post(`/properties/${id}/view`),
};

// ==================== BUYER ====================
export const buyerAPI = {
  listSaved: (params) => api.get("/buyer/favorites", { params }),
  save: (propertyId) => api.post(`/buyer/favorites/${propertyId}`),
  unsave: (propertyId) => api.delete(`/buyer/favorites/${propertyId}`),
  // Inquiries
  createInquiry: (data) => api.post("/buyer/inquiries", data),
  listInquiries: (params) => api.get("/buyer/inquiries", { params }),
};

// ==================== SALESMAN ====================
export const salesmanAPI = {
  dashboard: () => api.get("/salesman/dashboard"),
  listListings: (params) => api.get("/salesman/listings", { params }),
  createListing: (data) => api.post("/salesman/listings", data),
  getListing: (id) => api.get(`/salesman/listings/${id}`),
  updateListing: (id, data) => api.put(`/salesman/listings/${id}`, data),
  deleteListing: (id) => api.delete(`/salesman/listings/${id}`),
  submitListing: (id) => api.post(`/salesman/listings/${id}/submit`),
  deactivateListing: (id) => api.post(`/salesman/listings/${id}/deactivate`),
  markSold: (id) => api.post(`/salesman/listings/${id}/mark-sold`),
  markRented: (id) => api.post(`/salesman/listings/${id}/mark-rented`),
  uploadPhotos: (id, formData) =>
    api.post(`/salesman/listings/${id}/photos`, formData),
  deletePhoto: (listingId, photoId) =>
    api.delete(`/salesman/listings/${listingId}/photos/${photoId}`),
  reorderPhotos: (listingId, data) =>
    api.put(`/salesman/listings/${listingId}/photos/reorder`, data),
  getQuota: () => api.get("/salesman/quota"),
  // Inquiries
  listInquiries: (params) => api.get("/salesman/inquiries", { params }),
  updateInquiry: (id, data) => api.put(`/salesman/inquiries/${id}`, data),
};

// ==================== TENANT ADMIN ====================
export const tenantAPI = {
  dashboard: () => api.get("/tenant/dashboard"),
  getProfile: () => api.get("/tenant/profile"),
  updateProfile: (data) => api.put("/tenant/profile", data),
  listSalesmen: (params) => api.get("/tenant/salesmen", { params }),
  addSalesman: (data) => api.post("/tenant/salesmen", data),
  removeSalesman: (id) => api.delete(`/tenant/salesmen/${id}`),
  listListings: (params) => api.get("/tenant/listings", { params }),
  getSubscription: () => api.get("/tenant/subscription"),
  requestUpgrade: () =>
    api.post("/tenant/subscription/upgrade", { plan_type: "premium" }),
  // Inquiries
  listInquiries: (params) => api.get("/tenant/inquiries", { params }),
};

// ==================== PLATFORM ADMIN ====================
export const adminAPI = {
  dashboard: () => api.get("/admin/dashboard"),
  listTenants: (params) => api.get("/admin/tenants", { params }),
  createTenant: (data) => api.post("/admin/tenants", data),
  getTenant: (id) => api.get(`/admin/tenants/${id}`),
  suspendTenant: (id) => api.post(`/admin/tenants/${id}/suspend`),
  activateTenant: (id) => api.post(`/admin/tenants/${id}/activate`),
  changePlan: (id, data) => api.put(`/admin/tenants/${id}/plan`, data),
  listPending: (params) => api.get("/admin/listings/pending", { params }),
  approveListing: (id) => api.post(`/admin/listings/${id}/approve`),
  rejectListing: (id, data) => api.post(`/admin/listings/${id}/reject`, data),
  auditLogs: (params) => api.get("/admin/audit-logs", { params }),
  // All listings
  listAllListings: (params) => api.get("/admin/listings", { params }),
  // Master data: subscription plans
  listSubscriptionPlans: () => api.get("/admin/subscription-plans"),
  createSubscriptionPlan: (data) => api.post("/admin/subscription-plans", data),
  updateSubscriptionPlan: (id, data) =>
    api.put(`/admin/subscription-plans/${id}`, data),
  // Master data: property types
  listPropertyTypes: () => api.get("/admin/property-types"),
  createPropertyType: (data) => api.post("/admin/property-types", data),
  updatePropertyType: (id, data) =>
    api.put(`/admin/property-types/${id}`, data),
  // Master data: facilities
  listFacilities: () => api.get("/admin/facilities"),
  createFacility: (data) => api.post("/admin/facilities", data),
  updateFacility: (id, data) => api.put(`/admin/facilities/${id}`, data),
  // Master data: locations
  listLocations: (params) => api.get("/admin/locations", { params }),
  createLocation: (data) => api.post("/admin/locations", data),
  updateLocation: (id, data) => api.put(`/admin/locations/${id}`, data),
  // Change subscription by tenant ID
  changePlanByID: (id, data) =>
    api.put(`/admin/tenants/${id}/subscription`, data),
};

export default api;
