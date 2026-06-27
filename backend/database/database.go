package database

import (
	"log"
	"time"

	"property-hub-backend/config"
	"property-hub-backend/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	var err error

	logLevel := logger.Warn
	if cfg.AppEnv == "development" {
		logLevel = logger.Info
	}

	// Retry connection up to 30 times (waiting 2s between attempts = 60s total)
	// This handles slow PostgreSQL startup in Docker
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
			NowFunc: func() time.Time {
				loc, _ := time.LoadLocation(cfg.DBTimezone)
				return time.Now().In(loc)
			},
		})
		if err == nil {
			sqlDB, pingErr := DB.DB()
			if pingErr == nil {
				if pingErr = sqlDB.Ping(); pingErr == nil {
					break
				}
			}
		}
		if i < maxRetries-1 {
			log.Printf("[DB] Waiting for PostgreSQL... (attempt %d/%d)", i+1, maxRetries)
			time.Sleep(2 * time.Second)
		} else {
			log.Fatalf("❌ Failed to connect to database after %d attempts: %v", maxRetries, err)
		}
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetime) * time.Minute)

	// Ping
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}

	log.Println("✅ Database connected successfully")

	// Auto-migrate models
	if err := autoMigrate(); err != nil {
		log.Fatalf("❌ Auto-migration failed: %v", err)
	}
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.TenantUser{},
		&models.SubscriptionPlan{},
		&models.TenantSubscription{},
		&models.PropertyType{},
		&models.Location{},
		&models.Facility{},
		&models.PropertyListing{},
		&models.PropertyPhoto{},
		&models.PropertyFacility{},
		&models.BankAuctionDetail{},
		&models.CompanyAssetDetail{},
		&models.SavedProperty{},
		&models.Inquiry{},
		&models.PropertyView{},
		&models.AuditLog{},
	)
}

// SeedDefaultData seeds platform admin and sample tenants
func SeedDefaultData() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("[DB] Seed skipped — data already exists")
		return
	}

	log.Println("[DB] Seeding default data...")

	// Platform Admin
	platformAdmin := models.User{
		ID:           uuid.MustParse("a0000000-0000-0000-0000-000000000001"),
		Email:        "admin@propertyhub.id",
		PasswordHash: hashPassword("Admin@123"),
		Name:         "Super Admin",
		Phone:        strPtr("081100000001"),
		Role:         models.RolePlatformAdmin,
		Status:       models.UserStatusActive,
	}
	DB.Create(&platformAdmin)

	// Tenant 1: PropertiJaya (Free)
	t1ID := uuid.MustParse("b0000000-0000-0000-0000-000000000001")
	tenant1 := models.Tenant{
		ID:               t1ID,
		OrganizationName: "PropertiJaya Agency",
		SubdomainSlug:    "propertijaya",
		Description:      strPtr("Agensi properti terpercaya sejak 2010. Melayani jual-beli dan sewa rumah, apartemen, dan ruko di Jabodetabek."),
		Phone:            strPtr("0215551234"),
		Address:          strPtr("Jl. Sudirman No. 123, Jakarta Pusat"),
		Status:           models.TenantStatusActive,
	}
	DB.Create(&tenant1)

	DB.Create(&models.Subscription{
		ID:                     uuid.MustParse("c0000000-0000-0000-0000-000000000001"),
		TenantID:               t1ID,
		PlanType:               models.PlanFree,
		MaxSalesmen:            5,
		MaxListingsPerSalesman: 5,
	})

	admin1 := models.User{
		ID:           uuid.MustParse("d0000000-0000-0000-0000-000000000001"),
		TenantID:     uuidPtr(t1ID),
		Email:        "budi@propertijaya.id",
		PasswordHash: hashPassword("Budi@123"),
		Name:         "Budi Santoso",
		Phone:        strPtr("081200000001"),
		Role:         models.RoleTenantAdmin,
		Status:       models.UserStatusActive,
	}
	DB.Create(&admin1)

	salesmen := []models.User{
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000001"), TenantID: uuidPtr(t1ID), Email: "andi@propertijaya.id", PasswordHash: hashPassword("Andi@123"), Name: "Andi Pratama", Phone: strPtr("081200000002"), Role: models.RoleSalesman, Status: models.UserStatusActive},
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000002"), TenantID: uuidPtr(t1ID), Email: "siti@propertijaya.id", PasswordHash: hashPassword("Siti@123"), Name: "Siti Nurhaliza", Phone: strPtr("081200000003"), Role: models.RoleSalesman, Status: models.UserStatusActive},
		{ID: uuid.MustParse("e0000000-0000-0000-0000-000000000003"), TenantID: uuidPtr(t1ID), Email: "rudi@propertijaya.id", PasswordHash: hashPassword("Rudi@123"), Name: "Rudi Hermawan", Phone: strPtr("081200000004"), Role: models.RoleSalesman, Status: models.UserStatusActive},
	}
	for _, s := range salesmen {
		DB.Create(&s)
	}

	// Tenant 2: BankMaju (Premium)
	t2ID := uuid.MustParse("b0000000-0000-0000-0000-000000000002")
	tenant2 := models.Tenant{
		ID:               t2ID,
		OrganizationName: "BankMaju Asset Management",
		SubdomainSlug:    "bankmaju",
		Description:      strPtr("Divisi pengelolaan aset dan lelang BankMaju."),
		Phone:            strPtr("0215555678"),
		Address:          strPtr("Jl. Thamrin No. 45, Jakarta Pusat"),
		Status:           models.TenantStatusActive,
	}
	DB.Create(&tenant2)

	DB.Create(&models.Subscription{
		ID:                     uuid.MustParse("c0000000-0000-0000-0000-000000000002"),
		TenantID:               t2ID,
		PlanType:               models.PlanPremium,
		MaxSalesmen:            999999,
		MaxListingsPerSalesman: 999999,
	})

	admin2 := models.User{
		ID:           uuid.MustParse("d0000000-0000-0000-0000-000000000002"),
		TenantID:     uuidPtr(t2ID),
		Email:        "admin@bankmaju.id",
		PasswordHash: hashPassword("Bank@123"),
		Name:         "Hendra Gunawan",
		Phone:        strPtr("081100000002"),
		Role:         models.RoleTenantAdmin,
		Status:       models.UserStatusActive,
	}
	DB.Create(&admin2)

	buyers := []models.User{
		{ID: uuid.MustParse("f0000000-0000-0000-0000-000000000001"), Email: "rina@email.com", PasswordHash: hashPassword("Rina@123"), Name: "Rina Wijaya", Phone: strPtr("081300000001"), Role: models.RoleBuyer, Status: models.UserStatusActive},
		{ID: uuid.MustParse("f0000000-0000-0000-0000-000000000002"), Email: "doni@email.com", PasswordHash: hashPassword("Doni@123"), Name: "Doni Kusuma", Phone: strPtr("081300000002"), Role: models.RoleBuyer, Status: models.UserStatusActive},
	}
	for _, b := range buyers {
		DB.Create(&b)
	}

	log.Println("[DB] Seed data created (17 users, 2 tenants)")
}

func strPtr(s string) *string       { return &s }
func uuidPtr(u uuid.UUID) *uuid.UUID { return &u }

// hashPassword is a helper for seeding only
func hashPassword(pw string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pw), 12)
	return string(bytes)
}
