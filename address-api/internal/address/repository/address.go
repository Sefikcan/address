package repository

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sefikcan/address/address-api/internal/address/entity"
	"github.com/sefikcan/address/address-api/internal/common"
	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(ctx context.Context, address entity.Address) (entity.Address, error)
	Update(ctx context.Context, address entity.Address) (entity.Address, error)
	GetById(ctx context.Context, id int) (entity.Address, error)
	GetAll(ctx context.Context, page, pageSize int) (*common.Pageable[entity.Address], error)
	Delete(ctx context.Context, id int) error
}

type addressRepository struct {
	db *gorm.DB
}

func (a addressRepository) Update(ctx context.Context, address entity.Address) (entity.Address, error) {
	if result := a.db.WithContext(ctx).Save(&address); result.Error != nil {
		return entity.Address{}, errors.Wrap(result.Error, "addressRepository.Update.DbError")
	}

	return address, nil
}

func (a addressRepository) GetAll(ctx context.Context, page, pageSize int) (*common.Pageable[entity.Address], error) {
	var addresses []entity.Address
	var totalItems int64

	query := a.db.Model(&entity.Address{})

	// TODO: Add filter for query

	if err := query.WithContext(ctx).Count(&totalItems).Error; err != nil {
		return nil, errors.Wrap(err, "addressRepository.GetAll.CountDbError")
	}

	offset := (page - 1) * pageSize

	if err := query.Limit(pageSize).Offset(offset).WithContext(ctx).Find(&addresses).Error; err != nil {
		return nil, errors.Wrap(err, "addressRepository.GetAll.DbError")
	}

	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))

	pageable := &common.Pageable[entity.Address]{
		Items:       addresses,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}

	return pageable, nil
}

func (a addressRepository) Create(ctx context.Context, address entity.Address) (entity.Address, error) {
	if result := a.db.WithContext(ctx).Create(&address); result.Error != nil {
		return entity.Address{}, errors.Wrap(result.Error, "addressRepository.Create.DbError")
	}

	return address, nil
}

func (a addressRepository) GetById(ctx context.Context, id int) (entity.Address, error) {
	currentAddress := entity.Address{}
	err := a.db.WithContext(ctx).Where(`id = ?`, id).First(&currentAddress).Error
	if err != nil {
		return entity.Address{}, errors.Wrap(err, "addressRepository.GetById.DbError")
	}

	return currentAddress, err
}

func (a addressRepository) Delete(ctx context.Context, id int) error {
	if result := a.db.WithContext(ctx).Delete(&entity.Address{Id: id}); result.Error != nil {
		return errors.Wrap(result.Error, "addressRepository.Delete.DbError")
	}

	return nil
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{
		db: db,
	}
}
