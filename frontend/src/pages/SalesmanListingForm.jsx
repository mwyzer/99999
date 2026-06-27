import { useState, useEffect, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { salesmanAPI } from "../api";
import {
  ArrowLeft,
  Save,
  Send,
  Upload,
  X,
  ImagePlus,
  GripVertical,
} from "lucide-react";
import toast from "react-hot-toast";
import { DashboardLayout } from "../components/DashboardSidebar";

const PROPERTY_TYPES = [
  { value: "house", label: "Rumah" },
  { value: "land", label: "Tanah" },
  { value: "apartment", label: "Apartemen" },
  { value: "shophouse", label: "Ruko" },
  { value: "warehouse", label: "Gudang" },
  { value: "office", label: "Kantor" },
  { value: "villa", label: "Villa" },
];

const CERT_TYPES = ["SHM", "SHGB", "Girik", "Lainnya"];

const FACILITY_OPTIONS = [
  { key: "carport", label: "Carport" },
  { key: "garage", label: "Garasi" },
  { key: "garden", label: "Taman" },
  { key: "security_24h", label: "Keamanan 24 Jam" },
  { key: "swimming_pool", label: "Kolam Renang" },
  { key: "furnished", label: "Furnished" },
  { key: "ac", label: "AC" },
  { key: "balcony", label: "Balkon" },
  { key: "gym", label: "Gym" },
  { key: "pet_friendly", label: "Ramah Hewan" },
  { key: "wifi", label: "WiFi" },
  { key: "water_heater", label: "Water Heater" },
];

export default function SalesmanListingForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isEdit = !!id;
  const fileInputRef = useRef(null);
  const [submitting, setSubmitting] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [existingPhotos, setExistingPhotos] = useState([]);
  const [newPhotoFiles, setNewPhotoFiles] = useState([]);
  const [form, setForm] = useState({
    title: "",
    description: "",
    price: "",
    listing_type: "sale",
    property_type: "house",
    source_type: "regular",
    address: "",
    city: "",
    province: "",
    land_area: "",
    building_area: "",
    bedrooms: "",
    bathrooms: "",
    floors: "",
    certificate_type: "",
    facilities: {},
  });

  useEffect(() => {
    if (isEdit) {
      (async () => {
        try {
          const { data } = await salesmanAPI.getListing(id);
          const l = data.data;
          setForm({
            title: l.title || "",
            description: l.description || "",
            price: l.price || "",
            listing_type: l.listing_type || "sale",
            property_type: l.property_type || "house",
            source_type: l.source_type || "regular",
            address: l.address || "",
            city: l.city || "",
            province: l.province || "",
            land_area: l.land_area || "",
            building_area: l.building_area || "",
            bedrooms: l.bedrooms ?? "",
            bathrooms: l.bathrooms ?? "",
            floors: l.floors ?? "",
            certificate_type: l.certificate_type || "",
            facilities: l.facilities || {},
          });
          if (l.photos) {
            setExistingPhotos(l.photos);
          }
        } catch {
          toast.error("Gagal memuat listing.");
          navigate("/salesman/dashboard");
        }
      })();
    }
  }, [id, isEdit, navigate]);

  const update = (key, value) => setForm((prev) => ({ ...prev, [key]: value }));

  const toggleFacility = (key) => {
    setForm((prev) => ({
      ...prev,
      facilities: {
        ...prev.facilities,
        [key]: !prev.facilities[key],
      },
    }));
  };

  // ── Photo management ──
  const handleFileSelect = (e) => {
    const files = Array.from(e.target.files);
    if (files.length === 0) return;

    const totalAfter =
      existingPhotos.length + newPhotoFiles.length + files.length;
    if (totalAfter > 10) {
      toast.error(
        `Maksimal 10 foto. Saat ini sudah ada ${existingPhotos.length + newPhotoFiles.length} foto.`,
      );
      return;
    }

    // Validate each file
    const allowedTypes = ["image/jpeg", "image/png", "image/webp"];
    const maxSize = 5 * 1024 * 1024; // 5MB
    for (const f of files) {
      if (!allowedTypes.includes(f.type)) {
        toast.error(
          `File "${f.name}" tidak didukung. Gunakan JPG, PNG, atau WebP.`,
        );
        return;
      }
      if (f.size > maxSize) {
        toast.error(`File "${f.name}" terlalu besar. Maksimal 5MB.`);
        return;
      }
    }

    setNewPhotoFiles((prev) => [
      ...prev,
      ...files.map((f) => ({
        file: f,
        preview: URL.createObjectURL(f),
        id: crypto.randomUUID(),
      })),
    ]);
    // Reset input so same file can be re-selected
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const removeNewPhoto = (tempId) => {
    setNewPhotoFiles((prev) => {
      const item = prev.find((p) => p.id === tempId);
      if (item) URL.revokeObjectURL(item.preview);
      return prev.filter((p) => p.id !== tempId);
    });
  };

  const deleteExistingPhoto = async (photoId) => {
    if (!window.confirm("Hapus foto ini?")) return;
    try {
      await salesmanAPI.deletePhoto(id, photoId);
      setExistingPhotos((prev) => prev.filter((p) => p.id !== photoId));
      toast.success("Foto dihapus.");
    } catch {
      toast.error("Gagal menghapus foto.");
    }
  };

  const uploadNewPhotos = async () => {
    if (newPhotoFiles.length === 0) return [];
    setUploading(true);
    try {
      const formData = new FormData();
      newPhotoFiles.forEach((p) => formData.append("photos", p.file));
      const { data } = await salesmanAPI.uploadPhotos(id, formData);
      toast.success(`${data.data.photos_count} foto berhasil diunggah.`);
      const uploaded = data.data.photos || [];
      setExistingPhotos((prev) => [...prev, ...uploaded]);
      // Clean up previews
      newPhotoFiles.forEach((p) => URL.revokeObjectURL(p.preview));
      setNewPhotoFiles([]);
      return uploaded;
    } catch (err) {
      toast.error(
        err.response?.data?.error?.message || "Gagal mengunggah foto.",
      );
      return [];
    } finally {
      setUploading(false);
    }
  };

  const movePhoto = async (photoId, direction) => {
    const idx = existingPhotos.findIndex((p) => p.id === photoId);
    if (idx < 0) return;
    const newIdx = idx + direction;
    if (newIdx < 0 || newIdx >= existingPhotos.length) return;

    const reordered = [...existingPhotos];
    const [moved] = reordered.splice(idx, 1);
    reordered.splice(newIdx, 0, moved);
    setExistingPhotos(reordered);

    try {
      await salesmanAPI.reorderPhotos(id, {
        photo_ids: reordered.map((p) => p.id),
      });
    } catch {
      toast.error("Gagal mengubah urutan foto.");
      // Revert
      setExistingPhotos(existingPhotos);
    }
  };

  const handleSubmit = async (e, submitAfter = false) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      const payload = {
        ...form,
        price: parseFloat(form.price),
        land_area: form.land_area ? parseFloat(form.land_area) : null,
        building_area: form.building_area
          ? parseFloat(form.building_area)
          : null,
        bedrooms: form.bedrooms !== "" ? parseInt(form.bedrooms) : null,
        bathrooms: form.bathrooms !== "" ? parseInt(form.bathrooms) : null,
        floors: form.floors !== "" ? parseInt(form.floors) : null,
        certificate_type: form.certificate_type || null,
        facilities: form.facilities || {},
      };

      let listing;
      if (isEdit) {
        const { data } = await salesmanAPI.updateListing(id, payload);
        listing = data.data;
      } else {
        const { data } = await salesmanAPI.createListing(payload);
        listing = data.data;
      }

      if (submitAfter) {
        await salesmanAPI.submitListing(listing.id);
        // Upload photos if any pending
        if (newPhotoFiles.length > 0) {
          const uploadListingId = listing.id;
          const formData = new FormData();
          newPhotoFiles.forEach((p) => formData.append("photos", p.file));
          try {
            await salesmanAPI.uploadPhotos(uploadListingId, formData);
          } catch {
            // Non-fatal — listing is already submitted
          }
          newPhotoFiles.forEach((p) => URL.revokeObjectURL(p.preview));
          setNewPhotoFiles([]);
        }
        toast.success("Listing dibuat & diajukan untuk review!");
      } else {
        // Navigate to edit page so user can add photos
        if (!isEdit && newPhotoFiles.length === 0) {
          toast.success(
            "Listing disimpan sebagai draft. Tambahkan foto untuk melengkapi.",
          );
        }
        toast.success(
          isEdit ? "Listing diperbarui." : "Listing disimpan sebagai draft.",
        );
      }

      if (!isEdit && !submitAfter) {
        // New listing: navigate to edit page so photos can be uploaded
        navigate(`/salesman/listings/${listing.id}/edit`);
      } else {
        navigate("/salesman/dashboard");
      }
    } catch (err) {
      toast.error(err.response?.data?.error?.message || "Gagal menyimpan.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <DashboardLayout role="salesman">
      <div className="p-4 sm:p-6 lg:p-8 max-w-3xl">
        <button
          onClick={() => navigate("/salesman/dashboard")}
          className="flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 mb-4"
        >
          <ArrowLeft className="w-4 h-4" /> Kembali ke Dashboard
        </button>

        <h1 className="text-2xl font-bold mb-6">
          {isEdit ? "Edit Listing" : "Listing Baru"}
        </h1>

        <form onSubmit={(e) => handleSubmit(e, false)} className="space-y-5">
          {/* Basic Info */}
          <div className="card p-5 space-y-4">
            <h2 className="font-semibold text-lg">Informasi Dasar</h2>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Judul Listing *
              </label>
              <input
                type="text"
                required
                className="input-field"
                value={form.title}
                onChange={(e) => update("title", e.target.value)}
                placeholder="Contoh: Rumah Minimalis di Jakarta Selatan"
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Harga (IDR) *
                </label>
                <input
                  type="number"
                  required
                  className="input-field"
                  value={form.price}
                  onChange={(e) => update("price", e.target.value)}
                  placeholder="2500000000"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Tipe Listing *
                </label>
                <select
                  className="input-field"
                  value={form.listing_type}
                  onChange={(e) => update("listing_type", e.target.value)}
                >
                  <option value="sale">Dijual</option>
                  <option value="rent">Disewa</option>
                </select>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Tipe Properti *
                </label>
                <select
                  className="input-field"
                  value={form.property_type}
                  onChange={(e) => update("property_type", e.target.value)}
                >
                  {PROPERTY_TYPES.map((t) => (
                    <option key={t.value} value={t.value}>
                      {t.label}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Sumber
                </label>
                <select
                  className="input-field"
                  value={form.source_type}
                  onChange={(e) => update("source_type", e.target.value)}
                >
                  <option value="regular">Reguler</option>
                  <option value="bank_auction">Lelang Bank</option>
                  <option value="company_asset">Aset Perusahaan</option>
                </select>
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Deskripsi
              </label>
              <textarea
                className="input-field"
                rows={4}
                value={form.description}
                onChange={(e) => update("description", e.target.value)}
                placeholder="Deskripsikan properti Anda..."
              />
            </div>
          </div>

          {/* Location */}
          <div className="card p-5 space-y-4">
            <h2 className="font-semibold text-lg">Lokasi</h2>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Alamat
              </label>
              <input
                type="text"
                className="input-field"
                value={form.address}
                onChange={(e) => update("address", e.target.value)}
                placeholder="Alamat lengkap"
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Kota
                </label>
                <input
                  type="text"
                  className="input-field"
                  value={form.city}
                  onChange={(e) => update("city", e.target.value)}
                  placeholder="Jakarta Selatan"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Provinsi
                </label>
                <input
                  type="text"
                  className="input-field"
                  value={form.province}
                  onChange={(e) => update("province", e.target.value)}
                  placeholder="DKI Jakarta"
                />
              </div>
            </div>
          </div>

          {/* Specs */}
          <div className="card p-5 space-y-4">
            <h2 className="font-semibold text-lg">Spesifikasi</h2>
            <div className="grid grid-cols-3 gap-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Luas Tanah (m²)
                </label>
                <input
                  type="number"
                  className="input-field"
                  value={form.land_area}
                  onChange={(e) => update("land_area", e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Luas Bangunan (m²)
                </label>
                <input
                  type="number"
                  className="input-field"
                  value={form.building_area}
                  onChange={(e) => update("building_area", e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Sertifikat
                </label>
                <select
                  className="input-field"
                  value={form.certificate_type}
                  onChange={(e) => update("certificate_type", e.target.value)}
                >
                  <option value="">Pilih...</option>
                  {CERT_TYPES.map((c) => (
                    <option key={c} value={c}>
                      {c}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <div className="grid grid-cols-3 gap-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Kamar Tidur
                </label>
                <input
                  type="number"
                  className="input-field"
                  value={form.bedrooms}
                  onChange={(e) => update("bedrooms", e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Kamar Mandi
                </label>
                <input
                  type="number"
                  className="input-field"
                  value={form.bathrooms}
                  onChange={(e) => update("bathrooms", e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Lantai
                </label>
                <input
                  type="number"
                  className="input-field"
                  value={form.floors}
                  onChange={(e) => update("floors", e.target.value)}
                />
              </div>
            </div>
          </div>

          {/* Facilities */}
          <div className="card p-5 space-y-4">
            <h2 className="font-semibold text-lg">Fasilitas</h2>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
              {FACILITY_OPTIONS.map((f) => (
                <label
                  key={f.key}
                  className={`flex items-center gap-2 p-2 rounded-lg border cursor-pointer transition-colors ${
                    form.facilities[f.key]
                      ? "border-primary-500 bg-primary-50 text-primary-700"
                      : "border-gray-200 hover:border-gray-300 text-gray-600"
                  }`}
                >
                  <input
                    type="checkbox"
                    checked={!!form.facilities[f.key]}
                    onChange={() => toggleFacility(f.key)}
                    className="w-4 h-4 text-primary-600 rounded"
                  />
                  <span className="text-sm">{f.label}</span>
                </label>
              ))}
            </div>
          </div>

          {/* Photos — only in edit mode */}
          {isEdit && (
            <div className="card p-5 space-y-4">
              <h2 className="font-semibold text-lg">Foto Properti</h2>
              <p className="text-sm text-gray-500">
                Maksimal 10 foto. Format: JPG, PNG, WebP. Ukuran: max 5MB per
                foto.
              </p>

              {/* Existing photos */}
              {existingPhotos.length > 0 && (
                <div>
                  <p className="text-sm font-medium text-gray-600 mb-2">
                    Foto Terunggah ({existingPhotos.length}/10)
                  </p>
                  <div className="grid grid-cols-3 sm:grid-cols-5 gap-3">
                    {existingPhotos.map((photo, idx) => (
                      <div key={photo.id} className="relative group">
                        <img
                          src={photo.medium_url || photo.watermarked_url}
                          alt={photo.file_name || `Foto ${idx + 1}`}
                          className="w-full h-28 object-cover rounded-lg border"
                        />
                        <span className="absolute top-1 left-1 bg-black/60 text-white text-xs px-1.5 py-0.5 rounded">
                          {idx + 1}
                        </span>
                        <div className="absolute inset-0 bg-black/0 group-hover:bg-black/30 rounded-lg transition-colors flex items-center justify-center gap-1 opacity-0 group-hover:opacity-100">
                          {idx > 0 && (
                            <button
                              type="button"
                              onClick={() => movePhoto(photo.id, -1)}
                              className="bg-white rounded-full p-1 hover:bg-gray-100"
                              title="Geser kiri"
                            >
                              <GripVertical className="w-4 h-4" />
                            </button>
                          )}
                          {idx < existingPhotos.length - 1 && (
                            <button
                              type="button"
                              onClick={() => movePhoto(photo.id, 1)}
                              className="bg-white rounded-full p-1 hover:bg-gray-100"
                              title="Geser kanan"
                            >
                              <GripVertical className="w-4 h-4 rotate-180" />
                            </button>
                          )}
                          <button
                            type="button"
                            onClick={() => deleteExistingPhoto(photo.id)}
                            className="bg-red-500 text-white rounded-full p-1 hover:bg-red-600"
                            title="Hapus foto"
                          >
                            <X className="w-4 h-4" />
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* New photo previews */}
              {newPhotoFiles.length > 0 && (
                <div>
                  <p className="text-sm font-medium text-gray-600 mb-2">
                    Foto Baru ({newPhotoFiles.length})
                  </p>
                  <div className="grid grid-cols-3 sm:grid-cols-5 gap-3">
                    {newPhotoFiles.map((item) => (
                      <div key={item.id} className="relative">
                        <img
                          src={item.preview}
                          alt="Preview"
                          className="w-full h-28 object-cover rounded-lg border border-blue-300 ring-2 ring-blue-200"
                        />
                        <button
                          type="button"
                          onClick={() => removeNewPhoto(item.id)}
                          className="absolute -top-1.5 -right-1.5 bg-red-500 text-white rounded-full p-0.5 hover:bg-red-600"
                        >
                          <X className="w-3.5 h-3.5" />
                        </button>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Upload button */}
              <div className="flex items-center gap-3">
                <input
                  ref={fileInputRef}
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  multiple
                  className="hidden"
                  onChange={handleFileSelect}
                />
                <button
                  type="button"
                  onClick={() => fileInputRef.current?.click()}
                  disabled={
                    existingPhotos.length + newPhotoFiles.length >= 10 ||
                    uploading
                  }
                  className="btn-secondary flex items-center gap-1 text-sm"
                >
                  <ImagePlus className="w-4 h-4" /> Pilih Foto
                </button>
                {newPhotoFiles.length > 0 && (
                  <button
                    type="button"
                    onClick={uploadNewPhotos}
                    disabled={uploading}
                    className="btn-primary flex items-center gap-1 text-sm"
                  >
                    <Upload className="w-4 h-4" />
                    {uploading
                      ? "Mengunggah..."
                      : `Unggah ${newPhotoFiles.length} Foto`}
                  </button>
                )}
              </div>
            </div>
          )}

          {!isEdit && (
            <div className="card p-5 bg-blue-50 border border-blue-200">
              <p className="text-sm text-blue-700">
                💡 <strong>Tip:</strong> Simpan listing sebagai draft terlebih
                dahulu, lalu Anda dapat menambahkan foto properti di halaman
                edit.
              </p>
            </div>
          )}

          {/* Actions */}
          <div className="flex gap-3 justify-end">
            <button
              type="submit"
              disabled={submitting}
              className="btn-secondary flex items-center gap-1"
            >
              <Save className="w-4 h-4" /> Simpan Draft
            </button>
            <button
              type="button"
              disabled={submitting}
              onClick={(e) => handleSubmit(e, true)}
              className="btn-primary flex items-center gap-1"
            >
              <Send className="w-4 h-4" /> Simpan & Ajukan Review
            </button>
          </div>
        </form>
      </div>
    </DashboardLayout>
  );
}
