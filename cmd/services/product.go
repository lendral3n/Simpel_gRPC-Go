package services

import (
	"context"
	paginationPb "l3ngrpc/pb/pagination"
	productPb "l3ngrpc/pb/product"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ProductService struct {
	productPb.UnimplementedProductServiceServer
	DB *gorm.DB
}

type Product struct {
	gorm.Model
	Name         string
	Price        float64
	Stock        uint32
	CategoryId   uint
	CategoryName string
}

func (p ProductService) GetProducts(ctx context.Context, page *productPb.Page) (*productPb.Products, error) {
	var products []Product
	var total int64

	// Dapatkan total produk
	p.DB.Model(&Product{}).Count(&total)

	// Jumlah data per halaman
	perPage := 3

	// Hitung jumlah halaman
	lastPage := total / int64(perPage)
	if total%int64(perPage) != 0 {
		lastPage++
	}

	// Hitung halaman saat ini
	currentPage := page.GetPage()
	if currentPage > uint64(lastPage) {
		return nil, status.Errorf(codes.InvalidArgument, "currentPage is out of range.")
	}

	// Hitung offset berdasarkan halaman dan jumlah per halaman
	offset := (currentPage - 1) * uint64(perPage)

	// Query produk dengan limit dan offset
	result := p.DB.Offset(int(offset)).Limit(perPage).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	var pbProducts []*productPb.Product
	for _, product := range products {
		pbProduct := &productPb.Product{
			Id:    uint64(product.ID),
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Stock,
			Category: &productPb.Category{
				Id:   uint32(product.CategoryId),
				Name: product.CategoryName,
			},
		}
		pbProducts = append(pbProducts, pbProduct)
	}

	pbPagination := &paginationPb.Pagination{
		Total:       uint64(total),
		PerPage:     uint32(perPage),
		CurrentPage: currentPage,
		LastPage:    uint64(lastPage),
	}

	return &productPb.Products{
		Pagination: pbPagination,
		Data:       pbProducts,
	}, nil
}

// CreateProduct implements product.ProductServiceServer.
// Subtle: this method shadows the method (UnimplementedProductServiceServer).CreateProduct of ProductService.UnimplementedProductServiceServer.
func (p ProductService) CreateProduct(ctx context.Context, pbProduct *productPb.Product) (*productPb.Id, error) {
	// Buat produk baru dengan kategori yang sudah ada atau yang baru dibuat
	product := Product{
		Name:         pbProduct.Name,
		Price:        pbProduct.Price,
		Stock:        pbProduct.Stock,
		CategoryId:   uint(pbProduct.Category.Id),
		CategoryName: pbProduct.Category.Name,
	}

	// Simpan produk ke database
	if result := p.DB.Create(&product); result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	// Mengembalikan ID produk yang baru saja dibuat
	return &productPb.Id{
		Id: uint64(product.ID),
	}, nil
}

// DeleteProduct implements product.ProductServiceServer.
// Subtle: this method shadows the method (UnimplementedProductServiceServer).DeleteProduct of ProductService.UnimplementedProductServiceServer.
func (p ProductService) DeleteProduct(ctx context.Context, id *productPb.Id) (*productPb.Status, error) {
	if result := p.DB.Delete(&Product{}, id.GetId()); result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &productPb.Status{Status: 1}, nil
}

// GetProduct implements product.ProductServiceServer.
// Subtle: this method shadows the method (UnimplementedProductServiceServer).GetProduct of ProductService.UnimplementedProductServiceServer.
func (p ProductService) GetProduct(ctx context.Context, id *productPb.Id) (*productPb.Product, error) {
	var product Product

	if result := p.DB.First(&product, id.GetId()); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "Product with Id %d not found", id.GetId())
		}
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &productPb.Product{
		Id:    uint64(product.ID),
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
		Category: &productPb.Category{
			Id:   uint32(product.CategoryId),
			Name: product.CategoryName,
		},
	}, nil

}

// UpdateProduct implements product.ProductServiceServer.
// Subtle: this method shadows the method (UnimplementedProductServiceServer).UpdateProduct of ProductService.UnimplementedProductServiceServer.
func (p ProductService) UpdateProduct(ctx context.Context, pbProduct *productPb.Product) (*productPb.Status, error) {
	var product Product

	if result := p.DB.First(&product, pbProduct.GetId()); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "Product with ID %d not found", pbProduct.GetId())
		}
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	product.Name = pbProduct.GetName()
    product.Price = pbProduct.GetPrice()
    product.Stock = pbProduct.GetStock()
    product.CategoryId = uint(pbProduct.GetCategory().GetId())
    product.CategoryName = pbProduct.GetCategory().GetName()

	if result := p.DB.Updates(&product); result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}


	return &productPb.Status{Status: 1}, nil
}
