package database

import (
	"log"
	"math"
	"strconv"
	"time"

	"property-hub-backend/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// SeedAllData seeds comprehensive dummy data for development/testing.
// Minimum 5 entries per category: tenants, users (per role), listings, photos, saved, audit logs.
func SeedAllData() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("[DB] Seed skipped — data already exists")
		return
	}

	log.Println("[DB] ═══════════════════════════════════════════")
	log.Println("[DB]  Seeding comprehensive dummy data...")
	log.Println("[DB] ═══════════════════════════════════════════")

	// ================================================================
	// 1. TENANTS (5)
	// ================================================================
	log.Println("[DB] • Seeding 5 tenants...")

	tenants := []models.Tenant{
		{
			ID:               uuid.MustParse("b0000000-0000-0000-0000-000000000001"),
			OrganizationName: "PropertiJaya Agency",
			SubdomainSlug:    "propertijaya",
			Description:      sp("Agensi properti terpercaya sejak 2010. Melayani jual-beli dan sewa rumah, apartemen, dan ruko di Jabodetabek."),
			Phone:            sp("0215551234"),
			Address:          sp("Jl. Jend. Sudirman No. 123, Jakarta Pusat 10220"),
			Status:           models.TenantStatusActive,
		},
		{
			ID:               uuid.MustParse("b0000000-0000-0000-0000-000000000002"),
			OrganizationName: "BankMaju Asset Management",
			SubdomainSlug:    "bankmaju",
			Description:      sp("Divisi pengelolaan aset dan lelang BankMaju. Fokus pada properti lelang dan aset bermasalah."),
			Phone:            sp("0215555678"),
			Address:          sp("Jl. MH Thamrin No. 45, Jakarta Pusat 10340"),
			Status:           models.TenantStatusActive,
		},
		{
			ID:               uuid.MustParse("b0000000-0000-0000-0000-000000000003"),
			OrganizationName: "GriyaSentosa Properti",
			SubdomainSlug:    "griyasentosa",
			Description:      sp("Pengembang perumahan premium di kawasan Bandung dan sekitarnya. Spesialis rumah tapak dan villa."),
			Phone:            sp("0228881234"),
			Address:          sp("Jl. Dago No. 88, Bandung 40135"),
			Status:           models.TenantStatusActive,
		},
		{
			ID:               uuid.MustParse("b0000000-0000-0000-0000-000000000004"),
			OrganizationName: "MegaRaya Developments",
			SubdomainSlug:    "megaraya",
			Description:      sp("Developer properti komersial. Spesialis ruko, gudang, dan ruang kantor di Surabaya dan Jawa Timur."),
			Phone:            sp("0317775678"),
			Address:          sp("Jl. Basuki Rahmat No. 55, Surabaya 60271"),
			Status:           models.TenantStatusActive,
		},
		{
			ID:               uuid.MustParse("b0000000-0000-0000-0000-000000000005"),
			OrganizationName: "CiptaGraha Corporindo",
			SubdomainSlug:    "ciptagraha",
			Description:      sp("Perusahaan pengelola aset korporasi. Menjual dan menyewakan aset properti perusahaan."),
			Phone:            sp("0219994321"),
			Address:          sp("Jl. HR Rasuna Said, Kuningan, Jakarta Selatan 12940"),
			Status:           models.TenantStatusActive,
		},
	}
	for i := range tenants {
		DB.Create(&tenants[i])
	}

	// ================================================================
	// 2. SUBSCRIPTIONS (5 — one per tenant)
	// ================================================================
	log.Println("[DB] • Seeding 5 subscriptions...")

	subscriptions := []models.Subscription{
		{ID: uuid.MustParse("c0000000-0000-0000-0000-000000000001"), TenantID: tenants[0].ID, PlanType: models.PlanFree, MaxSalesmen: 5, MaxListingsPerSalesman: 5},
		{ID: uuid.MustParse("c0000000-0000-0000-0000-000000000002"), TenantID: tenants[1].ID, PlanType: models.PlanPremium, MaxSalesmen: 999999, MaxListingsPerSalesman: 999999},
		{ID: uuid.MustParse("c0000000-0000-0000-0000-000000000003"), TenantID: tenants[2].ID, PlanType: models.PlanPremium, MaxSalesmen: 20, MaxListingsPerSalesman: 50},
		{ID: uuid.MustParse("c0000000-0000-0000-0000-000000000004"), TenantID: tenants[3].ID, PlanType: models.PlanFree, MaxSalesmen: 5, MaxListingsPerSalesman: 5},
		{ID: uuid.MustParse("c0000000-0000-0000-0000-000000000005"), TenantID: tenants[4].ID, PlanType: models.PlanPremium, MaxSalesmen: 15, MaxListingsPerSalesman: 30},
	}
	for i := range subscriptions {
		DB.Create(&subscriptions[i])
	}

	// ================================================================
	// 3. USERS — by role (min 5 per role)
	// ================================================================
	log.Println("[DB] • Seeding users...")

	// --- 3a. Platform Admin (1) ---
	platformAdmin := models.User{
		ID:           uuid.MustParse("a0000000-0000-0000-0000-000000000001"),
		Email:        "admin@propertyhub.id",
		PasswordHash: hp("Admin@123"),
		Name:         "Super Admin",
		Phone:        sp("081100000001"),
		Role:         models.RolePlatformAdmin,
		Status:       models.UserStatusActive,
	}
	DB.Create(&platformAdmin)

	// --- 3b. Tenant Admins (5 — one per tenant) ---
	tenantAdmins := []models.User{
		{ID: uuid.MustParse("d0000000-0000-0000-0000-000000000001"), TenantID: up(tenants[0].ID), Email: "budi@propertijaya.id", PasswordHash: hp("Budi@123"), Name: "Budi Santoso", Phone: sp("081200000001"), Role: models.RoleTenantAdmin, Status: models.UserStatusActive},
		{ID: uuid.MustParse("d0000000-0000-0000-0000-000000000002"), TenantID: up(tenants[1].ID), Email: "admin@bankmaju.id", PasswordHash: hp("Bank@123"), Name: "Hendra Gunawan", Phone: sp("081200000002"), Role: models.RoleTenantAdmin, Status: models.UserStatusActive},
		{ID: uuid.MustParse("d0000000-0000-0000-0000-000000000003"), TenantID: up(tenants[2].ID), Email: "dewi@griyasentosa.id", PasswordHash: hp("Dewi@123"), Name: "Dewi Lestari", Phone: sp("081200000003"), Role: models.RoleTenantAdmin, Status: models.UserStatusActive},
		{ID: uuid.MustParse("d0000000-0000-0000-0000-000000000004"), TenantID: up(tenants[3].ID), Email: "admin@megaraya.id", PasswordHash: hp("Mega@123"), Name: "Agus Priyanto", Phone: sp("081200000004"), Role: models.RoleTenantAdmin, Status: models.UserStatusActive},
		{ID: uuid.MustParse("d0000000-0000-0000-0000-000000000005"), TenantID: up(tenants[4].ID), Email: "admin@ciptagraha.id", PasswordHash: hp("Cipta@123"), Name: "Ratna Sari Dewi", Phone: sp("081200000005"), Role: models.RoleTenantAdmin, Status: models.UserStatusActive},
	}
	for i := range tenantAdmins {
		DB.Create(&tenantAdmins[i])
	}

	// --- 3c. Salesmen (10 — 2 per tenant) ---
	salesmen := []models.User{
		// PropertiJaya
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000001"), TenantID: up(tenants[0].ID), Email: "andi@propertijaya.id", PasswordHash: hp("Andi@123"), Name: "Andi Pratama", Phone: sp("081300000001"), Role: models.RoleSalesman, Status: models.UserStatusActive},
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000002"), TenantID: up(tenants[0].ID), Email: "siti@propertijaya.id", PasswordHash: hp("Siti@123"), Name: "Siti Nurhaliza", Phone: sp("081300000002"), Role: models.RoleSalesman, Status: models.UserStatusActive},
		// BankMaju
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000003"), TenantID: up(tenants[1].ID), Email: "rudi@bankmaju.id", PasswordHash: hp("Rudi@123"), Name: "Rudi Hermawan", Phone: sp("081300000003"), Role: models.RoleSalesman, Status: models.UserStatusActive},
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000004"), TenantID: up(tenants[1].ID), Email: "maya@bankmaju.id", PasswordHash: hp("Maya@123"), Name: "Maya Anggraini", Phone: sp("081300000004"), Role: models.RoleSalesman, Status: models.UserStatusActive},
		// GriyaSentosa
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000005"), TenantID: up(tenants[2].ID), Email: "eko@griyasentosa.id", PasswordHash: hp("Eko@123"), Name: "Eko Prasetyo", Phone: sp("081300000005"), Role: models.RoleSalesman, Status: models.UserStatusActive},
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000006"), TenantID: up(tenants[2].ID), Email: "nina@griyasentosa.id", PasswordHash: hp("Nina@123"), Name: "Nina Kusumawati", Phone: sp("081300000006"), Role: models.RoleSalesman, Status: models.UserStatusActive},
		// MegaRaya
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000007"), TenantID: up(tenants[3].ID), Email: "bagus@megaraya.id", PasswordHash: hp("Bagus@123"), Name: "Bagus Wijaya", Phone: sp("081300000007"), Role: models.RoleSalesman, Status: models.UserStatusActive},
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000008"), TenantID: up(tenants[3].ID), Email: "lina@megaraya.id", PasswordHash: hp("Lina@123"), Name: "Lina Marlina", Phone: sp("081300000008"), Role: models.RoleSalesman, Status: models.UserStatusActive},
		// CiptaGraha
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000009"), TenantID: up(tenants[4].ID), Email: "tony@ciptagraha.id", PasswordHash: hp("Tony@123"), Name: "Tony Kusuma", Phone: sp("081300000009"), Role: models.RoleSalesman, Status: models.UserStatusActive},
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000010"), TenantID: up(tenants[4].ID), Email: "dian@ciptagraha.id", PasswordHash: hp("Dian@123"), Name: "Dian Permata", Phone: sp("081300000010"), Role: models.RoleSalesman, Status: models.UserStatusActive},
	}
	for i := range salesmen {
		DB.Create(&salesmen[i])
	}

	// --- 3d. Buyers (7) ---
	buyers := []models.User{
		{ID: uuid.MustParse("f0000000-0000-0000-0000-000000000001"), Email: "rina@email.com", PasswordHash: hp("Rina@123"), Name: "Rina Wijaya", Phone: sp("081400000001"), Role: models.RoleBuyer, Status: models.UserStatusActive},
		{ID: uuid.MustParse("f0000000-0000-0000-0000-000000000002"), Email: "doni@email.com", PasswordHash: hp("Doni@123"), Name: "Doni Kusuma", Phone: sp("081400000002"), Role: models.RoleBuyer, Status: models.UserStatusActive},
		{ID: uuid.MustParse("f0000000-0000-0000-0000-000000000003"), Email: "sari@email.com", PasswordHash: hp("Sari@123"), Name: "Sari Indah", Phone: sp("081400000003"), Role: models.RoleBuyer, Status: models.UserStatusActive},
		{ID: uuid.MustParse("f0000000-0000-0000-0000-000000000004"), Email: "bambang@email.com", PasswordHash: hp("Bambang@123"), Name: "Bambang Supriyadi", Phone: sp("081400000004"), Role: models.RoleBuyer, Status: models.UserStatusActive},
		{ID: uuid.MustParse("f0000000-0000-0000-0000-000000000005"), Email: "yanti@email.com", PasswordHash: hp("Yanti@123"), Name: "Yanti Susanti", Phone: sp("081400000005"), Role: models.RoleBuyer, Status: models.UserStatusActive},
		{ID: uuid.MustParse("f0000000-0000-0000-0000-000000000006"), Email: "ahmad@email.com", PasswordHash: hp("Ahmad@123"), Name: "Ahmad Fauzi", Phone: sp("081400000006"), Role: models.RoleBuyer, Status: models.UserStatusActive},
		{ID: uuid.MustParse("f0000000-0000-0000-0000-000000000007"), Email: "fitri@email.com", PasswordHash: hp("Fitri@123"), Name: "Fitri Handayani", Phone: sp("081400000007"), Role: models.RoleBuyer, Status: models.UserStatusActive},
	}
	for i := range buyers {
		DB.Create(&buyers[i])
	}

	// ================================================================
	// 4. PROPERTY LISTINGS (20 — min 5 per source type, mixed status)
	// ================================================================
	log.Println("[DB] • Seeding 20 property listings...")

	now := time.Now()
	lat := func(v float64) *float64 { return &v }
	fl := func(v float64) *float64 { return &v }
	iu := func(v int) *int { return &v }

	facilitiesHouse := models.Facilities{
		"listrik":     "2200 Watt",
		"air":         "PDAM + Sumur",
		"carport":     "2 mobil",
		"keamanan":    "24 jam",
		"hadap":       "Timur",
		"lantai":      "keramik",
		"dapur":       "full furnished",
	}
	facilitiesApartment := models.Facilities{
		"listrik":   "1300 Watt",
		"air":       "PDAM",
		"parkir":    "1 mobil",
		"keamanan":  "24 jam + CCTV",
		"lift":      "ya",
		"gym":       "ya",
		"kolam":     "ya",
		"furnished": "semi furnished",
	}
	facilitiesLand := models.Facilities{
		"listrik":    "tersedia",
		"air":        "tersedia",
		"bentuk":     "kotak",
		"akses_jalan": "aspal, 6 meter",
	}
	facilitiesVilla := models.Facilities{
		"listrik":    "4400 Watt",
		"air":        "Sumur bor",
		"carport":    "3 mobil",
		"keamanan":   "security + CCTV",
		"pemandangan": "pegunungan",
		"kolam_renang": "pribadi",
		"taman":      "luas",
		"hadap":      "Selatan",
	}
	facilitiesRuko := models.Facilities{
		"listrik":    "5500 Watt",
		"air":        "PDAM",
		"parkir":     "4 mobil",
		"akses_jalan": "jalan utama",
		"hadap":      "Utara",
		"lantai":     "granit",
		"toilet":     "2",
	}

	listings := []models.PropertyListing{
		// ── PropertiJaya (Tenant 1) — Regular, 6 listings ──
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000001"), TenantID: tenants[0].ID, SalesmanID: salesmen[0].ID,
			Title: "Rumah Minimalis 2 Lantai di Bintaro Sektor 7", Description: sp("Rumah modern minimalis dengan desain terbuka, lingkungan asri dan tenang. Dekat dengan Bintaro Plaza, sekolah internasional, dan akses tol."),
			Price: 1850000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeHouse, SourceType: models.SourceTypeRegular,
			Address: sp("Jl. Bintaro Utama 7 No. 15, Bintaro"), City: sp("Tangerang Selatan"), Province: sp("Banten"),
			Latitude: lat(-6.2765), Longitude: lat(106.7183), LandArea: fl(150), BuildingArea: fl(200), Bedrooms: iu(4), Bathrooms: iu(3), Floors: iu(2),
			CertificateType: sp(models.CertSHM), Facilities: facilitiesHouse, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000002"), TenantID: tenants[0].ID, SalesmanID: salesmen[0].ID,
			Title: "Apartemen Studio Fully Furnished — Gandaria City", Description: sp("Studio apartemen fully furnished di kawasan premium. Fasilitas lengkap: gym, kolam renang, 24h security. Langsung connected ke mall."),
			Price: 650000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeApartment, SourceType: models.SourceTypeRegular,
			Address: sp("Tower A Lt. 15, Gandaria City"), City: sp("Jakarta Selatan"), Province: sp("DKI Jakarta"),
			Latitude: lat(-6.2436), Longitude: lat(106.7834), LandArea: nil, BuildingArea: fl(32), Bedrooms: iu(1), Bathrooms: iu(1), Floors: iu(1),
			CertificateType: sp(models.CertSHGB), Facilities: facilitiesApartment, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000003"), TenantID: tenants[0].ID, SalesmanID: salesmen[1].ID,
			Title: "Rumah Hook Taman Minimalis BSD City", Description: sp("Posisi hook dengan taman luas. Cluster premium dengan one-gate system, dekat BSD Junction dan AEON Mall."),
			Price: 2200000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeHouse, SourceType: models.SourceTypeRegular,
			Address: sp("Cluster Savana, BSD City"), City: sp("Tangerang Selatan"), Province: sp("Banten"),
			Latitude: lat(-6.3016), Longitude: lat(106.6717), LandArea: fl(180), BuildingArea: fl(250), Bedrooms: iu(5), Bathrooms: iu(4), Floors: iu(2),
			CertificateType: sp(models.CertSHM), Facilities: facilitiesHouse, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000004"), TenantID: tenants[0].ID, SalesmanID: salesmen[1].ID,
			Title: "Sewa Apartemen 2BR — Kuningan City", Description: sp("Apartemen 2 bedroom furnished, view city. Termasuk IPL. Akses langsung ke mall dan MRT."),
			Price: 75000000, ListingType: models.ListingTypeRent, PropertyType: models.PropertyTypeApartment, SourceType: models.SourceTypeRegular,
			Address: sp("Kuningan City Tower B Lt. 22"), City: sp("Jakarta Selatan"), Province: sp("DKI Jakarta"),
			Latitude: lat(-6.2285), Longitude: lat(106.8242), LandArea: nil, BuildingArea: fl(48), Bedrooms: iu(2), Bathrooms: iu(1), Floors: iu(1),
			CertificateType: nil, Facilities: facilitiesApartment, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000005"), TenantID: tenants[0].ID, SalesmanID: salesmen[0].ID,
			Title: "Rumah Minimalis 3 Lantai di Kelapa Gading", Description: sp("Rumah baru, arsitektur modern. Lokasi strategis di pusat Kelapa Gading, dekat mall dan rumah sakit."),
			Price: 4500000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeHouse, SourceType: models.SourceTypeRegular,
			Address: sp("Jl. Boulevard Raya, Kelapa Gading"), City: sp("Jakarta Utara"), Province: sp("DKI Jakarta"),
			Latitude: lat(-6.1629), Longitude: lat(106.9065), LandArea: fl(200), BuildingArea: fl(350), Bedrooms: iu(5), Bathrooms: iu(4), Floors: iu(3),
			CertificateType: sp(models.CertSHM), Facilities: facilitiesHouse, Status: models.ListingStatusPending,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000006"), TenantID: tenants[0].ID, SalesmanID: salesmen[1].ID,
			Title: "Sewa Rumah Tahunan di Kebayoran Baru", Description: sp("Rumah klasik dengan halaman luas. Cocok untuk keluarga atau kantor representatif."),
			Price: 120000000, ListingType: models.ListingTypeRent, PropertyType: models.PropertyTypeHouse, SourceType: models.SourceTypeRegular,
			Address: sp("Jl. Wijaya No. 34, Kebayoran Baru"), City: sp("Jakarta Selatan"), Province: sp("DKI Jakarta"),
			Latitude: lat(-6.2441), Longitude: lat(106.7935), LandArea: fl(300), BuildingArea: fl(400), Bedrooms: iu(6), Bathrooms: iu(5), Floors: iu(2),
			CertificateType: sp(models.CertSHM), Facilities: facilitiesHouse, Status: models.ListingStatusSold,
		},

		// ── BankMaju (Tenant 2) — Bank Auction, 5 listings ──
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000007"), TenantID: tenants[1].ID, SalesmanID: salesmen[2].ID,
			Title: "Rumah Lelang Murah di Depok — SHM", Description: sp("Rumah lelang bank, kondisi butuh renovasi ringan. Harga di bawah pasar, cocok untuk investasi."),
			Price: 450000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeHouse, SourceType: models.SourceTypeBankAuction,
			Address: sp("Jl. Margonda Raya No. 200, Depok"), City: sp("Depok"), Province: sp("Jawa Barat"),
			Latitude: lat(-6.4001), Longitude: lat(106.8186), LandArea: fl(100), BuildingArea: fl(120), Bedrooms: iu(3), Bathrooms: iu(2), Floors: iu(1),
			CertificateType: sp(models.CertSHM), Facilities: models.Facilities{"listrik": "900 Watt", "air": "Sumur"}, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000008"), TenantID: tenants[1].ID, SalesmanID: salesmen[2].ID,
			Title: "Apartemen Lelang — Kalibata City Tower Sakura", Description: sp("Apartemen lelang bank, fully furnished. Lokasi strategis dekat stasiun Kalibata."),
			Price: 380000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeApartment, SourceType: models.SourceTypeBankAuction,
			Address: sp("Tower Sakura Lt. 10, Kalibata City"), City: sp("Jakarta Selatan"), Province: sp("DKI Jakarta"),
			Latitude: lat(-6.2587), Longitude: lat(106.8360), LandArea: nil, BuildingArea: fl(28), Bedrooms: iu(1), Bathrooms: iu(1), Floors: iu(1),
			CertificateType: sp(models.CertSHGB), Facilities: models.Facilities{"listrik": "1300 Watt", "parkir": "1"}, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000009"), TenantID: tenants[1].ID, SalesmanID: salesmen[3].ID,
			Title: "Tanah Kavling Lelang di Bogor Selatan", Description: sp("Tanah kavling luas, cocok untuk villa atau investasi. Akses jalan lebar, view gunung."),
			Price: 280000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeLand, SourceType: models.SourceTypeBankAuction,
			Address: sp("Jl. Raya Puncak KM 15, Bogor"), City: sp("Bogor"), Province: sp("Jawa Barat"),
			Latitude: lat(-6.6716), Longitude: lat(106.8800), LandArea: fl(500), BuildingArea: nil, Bedrooms: nil, Bathrooms: nil, Floors: nil,
			CertificateType: sp(models.CertSHM), Facilities: facilitiesLand, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000010"), TenantID: tenants[1].ID, SalesmanID: salesmen[3].ID,
			Title: "Ruko Lelang di Pulo Gadung — 3 Lantai", Description: sp("Ruko lelang di kawasan bisnis Pulo Gadung. Lantai 1 showroom, lt 2-3 gudang/kantor."),
			Price: 1500000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeShophouse, SourceType: models.SourceTypeBankAuction,
			Address: sp("Kawasan Industri Pulo Gadung Blok C-12"), City: sp("Jakarta Timur"), Province: sp("DKI Jakarta"),
			Latitude: lat(-6.1881), Longitude: lat(106.9088), LandArea: fl(120), BuildingArea: fl(360), Bedrooms: nil, Bathrooms: iu(3), Floors: iu(3),
			CertificateType: sp(models.CertSHGB), Facilities: facilitiesRuko, Status: models.ListingStatusPending,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000011"), TenantID: tenants[1].ID, SalesmanID: salesmen[2].ID,
			Title: "Rumah Lelang di Cibubur — Cluster Premium", Description: sp("Rumah lelang di cluster premium, kondisi bagus. Harga di bawah NJOP."),
			Price: 890000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeHouse, SourceType: models.SourceTypeBankAuction,
			Address: sp("Cluster Legenda Wisata, Cibubur"), City: sp("Bogor"), Province: sp("Jawa Barat"),
			Latitude: lat(-6.3680), Longitude: lat(106.9280), LandArea: fl(135), BuildingArea: fl(180), Bedrooms: iu(4), Bathrooms: iu(3), Floors: iu(2),
			CertificateType: sp(models.CertSHM), Facilities: facilitiesHouse, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},

		// ── GriyaSentosa (Tenant 3) — Regular & Villa, 4 listings ──
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000012"), TenantID: tenants[2].ID, SalesmanID: salesmen[4].ID,
			Title: "Villa Mewah Pemandangan Gunung — Lembang Asri", Description: sp("Villa premium dengan pemandangan pegunungan yang memukau. Cocok untuk retreat atau investasi properti wisata."),
			Price: 3500000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeVilla, SourceType: models.SourceTypeRegular,
			Address: sp("Jl. Maribaya No. 88, Lembang"), City: sp("Bandung Barat"), Province: sp("Jawa Barat"),
			Latitude: lat(-6.8148), Longitude: lat(107.6186), LandArea: fl(800), BuildingArea: fl(300), Bedrooms: iu(5), Bathrooms: iu(4), Floors: iu(2),
			CertificateType: sp(models.CertSHM), Facilities: facilitiesVilla, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000013"), TenantID: tenants[2].ID, SalesmanID: salesmen[4].ID,
			Title: "Rumah Modern di Dago Pakar — Full City View", Description: sp("Rumah modern dengan view kota Bandung. Desain arsitek, material premium. Very instagrammable!"),
			Price: 5200000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeHouse, SourceType: models.SourceTypeRegular,
			Address: sp("Dago Pakar Residence, Bandung"), City: sp("Bandung"), Province: sp("Jawa Barat"),
			Latitude: lat(-6.8711), Longitude: lat(107.6373), LandArea: fl(400), BuildingArea: fl(350), Bedrooms: iu(4), Bathrooms: iu(3), Floors: iu(2),
			CertificateType: sp(models.CertSHM), Facilities: facilitiesHouse, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000014"), TenantID: tenants[2].ID, SalesmanID: salesmen[5].ID,
			Title: "Sewa Villa Bulanan — Punclut Bandung", Description: sp("Villa nyaman untuk staycation. Tersedia bulanan, termasuk staff dan maintenance."),
			Price: 25000000, ListingType: models.ListingTypeRent, PropertyType: models.PropertyTypeVilla, SourceType: models.SourceTypeRegular,
			Address: sp("Jl. Punclut, Lembang"), City: sp("Bandung Barat"), Province: sp("Jawa Barat"),
			Latitude: lat(-6.8100), Longitude: lat(107.6250), LandArea: fl(500), BuildingArea: fl(200), Bedrooms: iu(3), Bathrooms: iu(2), Floors: iu(1),
			CertificateType: nil, Facilities: facilitiesVilla, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000015"), TenantID: tenants[2].ID, SalesmanID: salesmen[5].ID,
			Title: "Tanah Strategis di Pusat Kota Bandung", Description: sp("Tanah datar di lokasi premium. Cocok untuk apartemen, hotel, atau mixed-use development."),
			Price: 15000000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeLand, SourceType: models.SourceTypeRegular,
			Address: sp("Jl. Merdeka No. 100, Bandung"), City: sp("Bandung"), Province: sp("Jawa Barat"),
			Latitude: lat(-6.9148), Longitude: lat(107.6098), LandArea: fl(1500), BuildingArea: nil, Bedrooms: nil, Bathrooms: nil, Floors: nil,
			CertificateType: sp(models.CertSHM), Facilities: facilitiesLand, Status: models.ListingStatusDraft,
		},

		// ── MegaRaya (Tenant 4) — Commercial, 3 listings ──
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000016"), TenantID: tenants[3].ID, SalesmanID: salesmen[6].ID,
			Title: "Ruko 3 Lantai di Rungkut — Surabaya Timur", Description: sp("Ruko strategis di kawasan bisnis Rungkut. Ramai dan mudah diakses. Cocok untuk toko, gudang, atau kantor."),
			Price: 2750000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeShophouse, SourceType: models.SourceTypeRegular,
			Address: sp("Jl. Rungkut Industri No. 45, Surabaya"), City: sp("Surabaya"), Province: sp("Jawa Timur"),
			Latitude: lat(-7.3223), Longitude: lat(112.7688), LandArea: fl(100), BuildingArea: fl(300), Bedrooms: nil, Bathrooms: iu(2), Floors: iu(3),
			CertificateType: sp(models.CertSHGB), Facilities: facilitiesRuko, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000017"), TenantID: tenants[3].ID, SalesmanID: salesmen[7].ID,
			Title: "Gudang Modern 500m² — SIER Surabaya", Description: sp("Gudang dengan akses kontainer, plafon tinggi 8m, lantai beton. Dalam kawasan industri SIER."),
			Price: 4800000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeWarehouse, SourceType: models.SourceTypeRegular,
			Address: sp("Kawasan SIER Blok B-5, Surabaya"), City: sp("Surabaya"), Province: sp("Jawa Timur"),
			Latitude: lat(-7.3309), Longitude: lat(112.7425), LandArea: fl(500), BuildingArea: fl(480), Bedrooms: nil, Bathrooms: iu(1), Floors: iu(1),
			CertificateType: sp(models.CertSHGB), Facilities: models.Facilities{"listrik": "22000 Watt", "air": "PDAM", "akses_kontainer": "ya", "plafon": "8 meter", "lantai": "beton"}, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000018"), TenantID: tenants[3].ID, SalesmanID: salesmen[6].ID,
			Title: "Sewa Ruang Kantor — CBD Surabaya", Description: sp("Ruang kantor fully furnished di pusat bisnis Surabaya. Termasuk meeting room dan pantry."),
			Price: 180000000, ListingType: models.ListingTypeRent, PropertyType: models.PropertyTypeOffice, SourceType: models.SourceTypeRegular,
			Address: sp("CBD Surabaya Tower Lt. 12"), City: sp("Surabaya"), Province: sp("Jawa Timur"),
			Latitude: lat(-7.2741), Longitude: lat(112.7448), LandArea: nil, BuildingArea: fl(200), Bedrooms: nil, Bathrooms: iu(2), Floors: iu(1),
			CertificateType: nil, Facilities: models.Facilities{"listrik": "5500 Watt", "ac": "central", "lift": "ya", "parkir": "basement", "pantry": "ya", "meeting_room": "ya"}, Status: models.ListingStatusPending,
		},

		// ── CiptaGraha (Tenant 5) — Company Asset, 2 listings ──
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000019"), TenantID: tenants[4].ID, SalesmanID: salesmen[8].ID,
			Title: "Rumah Dinas Eks Perusahaan — Menteng Jakarta", Description: sp("Rumah bersejarah di kawasan Menteng. Bekas rumah dinas perusahaan multinasional. Arsitektur kolonial, lahan luas."),
			Price: 18000000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeHouse, SourceType: models.SourceTypeCompanyAsset,
			Address: sp("Jl. Imam Bonjol No. 55, Menteng"), City: sp("Jakarta Pusat"), Province: sp("DKI Jakarta"),
			Latitude: lat(-6.1918), Longitude: lat(106.8306), LandArea: fl(1200), BuildingArea: fl(600), Bedrooms: iu(7), Bathrooms: iu(6), Floors: iu(2),
			CertificateType: sp(models.CertSHM), Facilities: models.Facilities{"listrik": "11000 Watt", "air": "PDAM", "carport": "6 mobil", "keamanan": "security 24 jam", "taman": "luas"}, Status: models.ListingStatusApproved, ApprovedBy: up(platformAdmin.ID), ApprovedAt: &now,
		},
		{
			ID: uuid.MustParse("10000000-0000-0000-0000-000000000020"), TenantID: tenants[4].ID, SalesmanID: salesmen[9].ID,
			Title: "Ruang Kantor Eks Perusahaan — SCBD Lot 10", Description: sp("Ruang kantor full floor di kawasan SCBD. View kota, fully furnished, siap pakai."),
			Price: 8500000000, ListingType: models.ListingTypeSale, PropertyType: models.PropertyTypeOffice, SourceType: models.SourceTypeCompanyAsset,
			Address: sp("SCBD Lot 10 Lt. 25, Jakarta"), City: sp("Jakarta Selatan"), Province: sp("DKI Jakarta"),
			Latitude: lat(-6.2243), Longitude: lat(106.8102), LandArea: nil, BuildingArea: fl(350), Bedrooms: nil, Bathrooms: iu(3), Floors: iu(1),
			CertificateType: sp(models.CertSHGB), Facilities: models.Facilities{"listrik": "33000 Watt", "ac": "central", "lift": "private", "parkir": "6 mobil", "view": "city"}, Status: models.ListingStatusDraft,
		},
	}
	for i := range listings {
		DB.Create(&listings[i])
	}

	// ================================================================
	// 5. PROPERTY PHOTOS (1–3 per listing)
	// ================================================================
	log.Println("[DB] • Seeding property photos...")
	// Using picsum.photos as placeholder — real app would use actual uploads
	baseURL := "https://picsum.photos"

	type photoDef struct {
		id        string
		listingID uuid.UUID
		sortOrder int
	}
	photoDefs := []photoDef{
		// Listing 1 — 3 photos
		{"20000000-0000-0000-0000-000000000001", listings[0].ID, 0},
		{"20000000-0000-0000-0000-000000000002", listings[0].ID, 1},
		{"20000000-0000-0000-0000-000000000003", listings[0].ID, 2},
		// Listing 2 — 2 photos
		{"20000000-0000-0000-0000-000000000004", listings[1].ID, 0},
		{"20000000-0000-0000-0000-000000000005", listings[1].ID, 1},
		// Listing 3 — 2 photos
		{"20000000-0000-0000-0000-000000000006", listings[2].ID, 0},
		{"20000000-0000-0000-0000-000000000007", listings[2].ID, 1},
		// Listing 4 — 1 photo
		{"20000000-0000-0000-0000-000000000008", listings[3].ID, 0},
		// Listing 5 — 2 photos
		{"20000000-0000-0000-0000-000000000009", listings[4].ID, 0},
		{"20000000-0000-0000-0000-000000000010", listings[4].ID, 1},
		// Listing 6 — 1 photo
		{"20000000-0000-0000-0000-000000000011", listings[5].ID, 0},
		// Listing 7 — 2 photos
		{"20000000-0000-0000-0000-000000000012", listings[6].ID, 0},
		{"20000000-0000-0000-0000-000000000013", listings[6].ID, 1},
		// Listing 8 — 2 photos
		{"20000000-0000-0000-0000-000000000014", listings[7].ID, 0},
		{"20000000-0000-0000-0000-000000000015", listings[7].ID, 1},
		// Listing 12 (Villa) — 3 photos
		{"20000000-0000-0000-0000-000000000016", listings[11].ID, 0},
		{"20000000-0000-0000-0000-000000000017", listings[11].ID, 1},
		{"20000000-0000-0000-0000-000000000018", listings[11].ID, 2},
		// Listing 13 — 1 photo
		{"20000000-0000-0000-0000-000000000019", listings[12].ID, 0},
		// Listing 16 — 2 photos
		{"20000000-0000-0000-0000-000000000020", listings[15].ID, 0},
		{"20000000-0000-0000-0000-000000000021", listings[15].ID, 1},
		// Listing 19 — 2 photos
		{"20000000-0000-0000-0000-000000000022", listings[18].ID, 0},
		{"20000000-0000-0000-0000-000000000023", listings[18].ID, 1},
	}

	for _, pd := range photoDefs {
		seed := int(math.Abs(float64(int(pd.sortOrder)*73 + 447)))
		photo := models.PropertyPhoto{
			ID:             uuid.MustParse(pd.id),
			ListingID:      pd.listingID,
			OriginalURL:    baseURL + "/seed/" + itoa(seed) + "/800/600",
			WatermarkedURL: baseURL + "/seed/" + itoa(seed) + "/800/600",
			MediumURL:      sp(baseURL + "/seed/" + itoa(seed) + "/400/300"),
			ThumbnailURL:   sp(baseURL + "/seed/" + itoa(seed) + "/150/150"),
			SortOrder:      pd.sortOrder,
		}
		DB.Create(&photo)
	}

	// ================================================================
	// 6. SAVED PROPERTIES (buyers bookmarking listings)
	// ================================================================
	log.Println("[DB] • Seeding saved properties...")

	saved := []models.SavedProperty{
		{ID: uuid.MustParse("30000000-0000-0000-0000-000000000001"), BuyerID: buyers[0].ID, ListingID: listings[0].ID},   // Rina → Rumah Bintaro
		{ID: uuid.MustParse("30000000-0000-0000-0000-000000000002"), BuyerID: buyers[0].ID, ListingID: listings[11].ID},  // Rina → Villa Lembang
		{ID: uuid.MustParse("30000000-0000-0000-0000-000000000003"), BuyerID: buyers[0].ID, ListingID: listings[3].ID},   // Rina → Apt Kuningan sewa
		{ID: uuid.MustParse("30000000-0000-0000-0000-000000000004"), BuyerID: buyers[1].ID, ListingID: listings[6].ID},   // Doni → Lelang Depok
		{ID: uuid.MustParse("30000000-0000-0000-0000-000000000005"), BuyerID: buyers[1].ID, ListingID: listings[7].ID},   // Doni → Lelang Kalibata
		{ID: uuid.MustParse("30000000-0000-0000-0000-000000000006"), BuyerID: buyers[2].ID, ListingID: listings[12].ID},  // Sari → Dago Pakar
		{ID: uuid.MustParse("30000000-0000-0000-0000-000000000007"), BuyerID: buyers[3].ID, ListingID: listings[15].ID},  // Bambang → Ruko Surabaya
		{ID: uuid.MustParse("30000000-0000-0000-0000-000000000008"), BuyerID: buyers[4].ID, ListingID: listings[0].ID},   // Yanti → Rumah Bintaro
		{ID: uuid.MustParse("30000000-0000-0000-0000-000000000009"), BuyerID: buyers[5].ID, ListingID: listings[18].ID},  // Ahmad → Menteng
		{ID: uuid.MustParse("30000000-0000-0000-0000-000000000010"), BuyerID: buyers[6].ID, ListingID: listings[1].ID},   // Fitri → Apt Gandaria
	}
	for i := range saved {
		DB.Create(&saved[i])
	}

	// ================================================================
	// 7. AUDIT LOGS (activity trail)
	// ================================================================
	log.Println("[DB] • Seeding audit logs...")

	auditLogs := []models.AuditLog{
		{ID: uuid.MustParse("40000000-0000-0000-0000-000000000001"), UserID: up(platformAdmin.ID), UserRole: models.RolePlatformAdmin, Action: models.AuditActionApprove, EntityType: "property_listing", EntityID: listings[0].ID.String(), NewValues: mp("status", models.ListingStatusApproved)},
		{ID: uuid.MustParse("40000000-0000-0000-0000-000000000002"), UserID: up(platformAdmin.ID), UserRole: models.RolePlatformAdmin, Action: models.AuditActionApprove, EntityType: "property_listing", EntityID: listings[1].ID.String(), NewValues: mp("status", models.ListingStatusApproved)},
		{ID: uuid.MustParse("40000000-0000-0000-0000-000000000003"), UserID: up(platformAdmin.ID), UserRole: models.RolePlatformAdmin, Action: models.AuditActionApprove, EntityType: "property_listing", EntityID: listings[3].ID.String(), NewValues: mp("status", models.ListingStatusApproved)},
		{ID: uuid.MustParse("40000000-0000-0000-0000-000000000004"), UserID: up(salesmen[0].ID), UserRole: models.RoleSalesman, Action: models.AuditActionCreate, EntityType: "property_listing", EntityID: listings[0].ID.String(), NewValues: mp("title", "Rumah Minimalis 2 Lantai di Bintaro Sektor 7")},
		{ID: uuid.MustParse("40000000-0000-0000-0000-000000000005"), UserID: up(salesmen[1].ID), UserRole: models.RoleSalesman, Action: models.AuditActionCreate, EntityType: "property_listing", EntityID: listings[5].ID.String(), NewValues: mp("title", "Sewa Rumah Tahunan di Kebayoran Baru")},
		{ID: uuid.MustParse("40000000-0000-0000-0000-000000000006"), UserID: up(tenantAdmins[0].ID), UserRole: models.RoleTenantAdmin, Action: models.AuditActionCreate, EntityType: "user", EntityID: salesmen[0].ID.String(), NewValues: mp("name", "Andi Pratama")},
		{ID: uuid.MustParse("40000000-0000-0000-0000-000000000007"), UserID: up(platformAdmin.ID), UserRole: models.RolePlatformAdmin, Action: models.AuditActionActivate, EntityType: "tenant", EntityID: tenants[0].ID.String(), NewValues: mp("status", "active")},
		{ID: uuid.MustParse("40000000-0000-0000-0000-000000000008"), UserID: up(platformAdmin.ID), UserRole: models.RolePlatformAdmin, Action: models.AuditActionActivate, EntityType: "tenant", EntityID: tenants[2].ID.String(), NewValues: mp("status", "active")},
		{ID: uuid.MustParse("40000000-0000-0000-0000-000000000009"), UserID: up(salesmen[2].ID), UserRole: models.RoleSalesman, Action: models.AuditActionUpdate, EntityType: "property_listing", EntityID: listings[6].ID.String(), OldValues: mp("price", "500,000,000"), NewValues: mp("price", "450,000,000")},
		{ID: uuid.MustParse("40000000-0000-0000-0000-000000000010"), UserID: up(platformAdmin.ID), UserRole: models.RolePlatformAdmin, Action: models.AuditActionApprove, EntityType: "property_listing", EntityID: listings[17].ID.String(), NewValues: mp("status", models.ListingStatusApproved)},
	}

	for i := range auditLogs {
		DB.Create(&auditLogs[i])
	}

	// ================================================================
	// SUMMARY
	// ================================================================
	log.Println("[DB] ═══════════════════════════════════════════")
	log.Println("[DB]  ✅ Seed data complete!")
	log.Printf("[DB]  📊 %d tenants", len(tenants))
	log.Printf("[DB]  📊 %d subscriptions", len(subscriptions))
	log.Printf("[DB]  📊 1 platform admin")
	log.Printf("[DB]  📊 %d tenant admins", len(tenantAdmins))
	log.Printf("[DB]  📊 %d salesmen", len(salesmen))
	log.Printf("[DB]  📊 %d buyers", len(buyers))
	log.Printf("[DB]  📊 %d property listings", len(listings))
	log.Printf("[DB]  📊 %d property photos", len(photoDefs))
	log.Printf("[DB]  📊 %d saved properties", len(saved))
	log.Printf("[DB]  📊 %d audit logs", len(auditLogs))
	log.Println("[DB] ═══════════════════════════════════════════")
	log.Println("[DB]  🔑 Demo accounts (password: <Name>@123)")
	log.Println("[DB]     Platform Admin: admin@propertyhub.id / Admin@123")
	log.Println("[DB]     Tenant Admins:  budi@propertijaya.id / Budi@123")
	log.Println("[DB]                     admin@bankmaju.id / Bank@123")
	log.Println("[DB]     Salesmen:       andi@propertijaya.id / Andi@123")
	log.Println("[DB]     Buyers:         rina@email.com / Rina@123")
	log.Println("[DB] ═══════════════════════════════════════════")
}

// ── Helper functions for seeding ──

func sp(s string) *string             { return &s }
func hp(pw string) string             { bytes, _ := bcrypt.GenerateFromPassword([]byte(pw), 12); return string(bytes) }
func up(u uuid.UUID) *uuid.UUID       { return &u }
func itoa(i int) string               { return strconv.Itoa(i) }
func mp(k, v string) *models.JSONMap  { m := models.JSONMap{k: v}; return &m }
