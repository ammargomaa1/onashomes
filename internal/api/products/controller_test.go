package products

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/api/products/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

type fakeService struct {
	publicListFn func(ctx context.Context, pg *utils.Pagination, q string) ([]AdminProductListItem, int64, error)
	publicGetFn  func(ctx context.Context, id int64) (*AdminProductDetail, error)
}

// Implement Service interface for fakeService
func (f *fakeService) CreateProduct(ctx context.Context, req requests.AdminProductRequest) (*AdminProductDetail, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeService) UpdateProduct(ctx context.Context, id int64, req requests.AdminProductRequest) (*AdminProductDetail, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeService) DeleteProduct(ctx context.Context, id int64) error {
	return errors.New("not implemented")
}

func (f *fakeService) ListProducts(ctx context.Context, pagination *utils.Pagination) ([]AdminProductListItem, int64, error) {
	return nil, 0, errors.New("not implemented")
}

func (f *fakeService) GetProductByID(ctx context.Context, id int64) (*AdminProductDetail, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeService) CreateVariant(ctx context.Context, productID int64, req requests.AdminVariantRequest) (*AdminVariant, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeService) UpdateVariant(ctx context.Context, productID, variantID int64, req requests.AdminVariantRequest) (*AdminVariant, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeService) ListVariantAddOns(ctx context.Context, productID, variantID int64) ([]AdminVariantAddOn, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeService) CreateVariantAddOn(ctx context.Context, productID, variantID, addOnProductID int64) (*AdminVariantAddOn, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeService) UpdateVariantAddOn(ctx context.Context, productID, variantID, addOnID, addOnProductID int64) (*AdminVariantAddOn, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeService) DeleteVariantAddOn(ctx context.Context, productID, variantID, addOnID int64) error {
	return errors.New("not implemented")
}

func (f *fakeService) PublicListProducts(ctx context.Context, pg *utils.Pagination, q string) ([]AdminProductListItem, int64, error) {
	if f.publicListFn != nil {
		return f.publicListFn(ctx, pg, q)
	}
	return nil, 0, nil
}

func (f *fakeService) PublicGetProduct(ctx context.Context, id int64) (*AdminProductDetail, error) {
	if f.publicGetFn != nil {
		return f.publicGetFn(ctx, id)
	}
	return nil, nil
}

func TestPublicListProducts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fs := &fakeService{
		publicListFn: func(ctx context.Context, pg *utils.Pagination, q string) ([]AdminProductListItem, int64, error) {
			items := []AdminProductListItem{
				{ID: 1, Name: "T-Shirt", Price: 199.99, IsActive: true},
			}
			return items, int64(len(items)), nil
		},
	}
	ctrl := NewController(fs)

	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/products", ctrl.PublicListProducts)

	req, _ := http.NewRequest(http.MethodGet, "/products?q=shirt&page=1&limit=10", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), "T-Shirt") {
		t.Fatalf("expected response body to contain product name, got %s", w.Body.String())
	}
}

func TestPublicGetProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fs := &fakeService{
		publicGetFn: func(ctx context.Context, id int64) (*AdminProductDetail, error) {
			return nil, errors.New("not found")
		},
	}
	ctrl := NewController(fs)

	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/products/:id", ctrl.PublicGetProduct)

	req, _ := http.NewRequest(http.MethodGet, "/products/123", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
