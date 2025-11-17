package product

import "db-less-store/pkg/db"

type ProductRepository struct {
	Database *db.Db
}

func NewProductRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{Database: database}
}

func (r *ProductRepository) CreateProduct(product *Product) error {
	result := r.Database.DB.Create(product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProductRepository) GetProductById(id string) (*Product, error) {
	var product Product
	result := r.Database.DB.First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (r *ProductRepository) GetAllProducts() ([]*Product, error) {
	var products []*Product
	result := r.Database.DB.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (r *ProductRepository) UpdateProduct(product *Product) error {
	result := r.Database.DB.Save(product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(id string) error {
	result := r.Database.DB.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProductRepository) ProductExists(id string) (bool, error) {
	var count int64
	result := r.Database.DB.Model(&Product{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
