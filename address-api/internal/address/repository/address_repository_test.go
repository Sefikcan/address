package repository

import (
	"context"
	"github.com/sefikcan/address-api/internal/address/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func SetupTestDB() (*gorm.DB, func()) {
	// Create a new in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&entity.Address{})
	if err != nil {
		return nil, nil
	}

	// Return a cleanup function to close the DB connection
	return db, func() {
		db.Exec(`DROP TABLE addresses`)
	}
}

func TestAddressRepository_Create(t *testing.T) {
	db, teardown := SetupTestDB()
	defer teardown()

	repo := NewAddressRepository(db)

	address := entity.Address{
		Id:          1,
		Country:     "Test St",
		City:        "Test City",
		FullAddress: "test test",
		UserId:      "1",
	}
	result, err := repo.Create(context.Background(), address)

	assert.NoError(t, err)
	assert.Equal(t, address.Id, result.Id)
	assert.Equal(t, address.City, result.City)
	assert.Equal(t, address.Country, result.Country)
	assert.Equal(t, address.UserId, result.UserId)
	assert.Equal(t, address.FullAddress, result.FullAddress)
}

func TestAddressRepository_GetById(t *testing.T) {
	db, teardown := SetupTestDB()
	defer teardown()

	repo := NewAddressRepository(db)
	address := entity.Address{
		Id:          1,
		Country:     "Test St",
		City:        "Test City",
		FullAddress: "test test",
		UserId:      "1",
	}
	_, _ = repo.Create(context.Background(), address)

	result, err := repo.GetById(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, address.Id, result.Id)
	assert.Equal(t, address.City, result.City)
	assert.Equal(t, address.Country, result.Country)
	assert.Equal(t, address.UserId, result.UserId)
	assert.Equal(t, address.FullAddress, result.FullAddress)
}

func TestAddressRepository_Delete(t *testing.T) {
	db, teardown := SetupTestDB()
	defer teardown()

	repo := NewAddressRepository(db)
	address := entity.Address{
		Id:          1,
		Country:     "Test St",
		City:        "Test City",
		FullAddress: "test test",
		UserId:      "1",
	}
	_, _ = repo.Create(context.Background(), address)

	err := repo.Delete(context.Background(), 1)

	assert.NoError(t, err)
}

func TestAddressRepository_GetAll(t *testing.T) {
	db, teardown := SetupTestDB()
	defer teardown()

	repo := NewAddressRepository(db)

	// Create multiple addresses
	addresses := []entity.Address{
		{
			Id:          1,
			Country:     "Test St",
			City:        "Test City",
			FullAddress: "test test",
			UserId:      "1",
		},
		{
			Id:          2,
			Country:     "Test 2 St",
			City:        "Test 2 City",
			FullAddress: "test 2 test",
			UserId:      "2",
		},
	}
	for _, addr := range addresses {
		_, _ = repo.Create(context.Background(), addr)
	}

	pageable, err := repo.GetAll(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(len(addresses)), pageable.TotalItems)
	assert.Len(t, pageable.Items, 2)
}

func TestAddressRepository_Update(t *testing.T) {
	db, teardown := SetupTestDB()
	defer teardown()

	repo := NewAddressRepository(db)

	// Create an address to update
	address := entity.Address{
		Id:          1,
		Country:     "Test St",
		City:        "Test City",
		FullAddress: "test test",
		UserId:      "1",
	}
	_, _ = repo.Create(context.Background(), address)

	// Update the address
	address.City = "Updated St 2"
	result, err := repo.Update(context.Background(), address)

	assert.NoError(t, err)
	assert.Equal(t, address.Id, result.Id)
	assert.Equal(t, address.City, result.City)
	assert.Equal(t, address.Country, result.Country)
	assert.Equal(t, address.UserId, result.UserId)
	assert.Equal(t, address.FullAddress, result.FullAddress)
}
