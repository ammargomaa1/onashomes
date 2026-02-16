package orders

import (
	"fmt"
	"time"

	"github.com/onas/ecommerce-api/internal/api/inventory"
	"github.com/onas/ecommerce-api/internal/api/orders/requests"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repo       *Repository
	invService *inventory.Service
}

func NewService(db *gorm.DB, repo *Repository, invService *inventory.Service) *Service {
	return &Service{
		db:         db,
		repo:       repo,
		invService: invService,
	}
}

// ... (imports and Service struct)

// CreateOrder handles order creation with strict inventory reservation
func (s *Service) CreateOrder(req requests.CreateOrderRequest, adminID int64) utils.IResource {
	var order *models.Order

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Validate StoreFront
		var storeFront models.StoreFront
		if err := tx.First(&storeFront, req.StoreFrontID).Error; err != nil {
			return fmt.Errorf("invalid store_front_id: %w", err)
		}

		repoTx := &Repository{db: tx}

		// Get Default Statuses
		pendingStatus, err := repoTx.GetOrderStatusBySlug("pending_payment")
		if err != nil {
			return fmt.Errorf("failed to get pending status: %w", err)
		}

		unpaidStatus, err := repoTx.GetPaymentStatusBySlug("unpaid")
		if err != nil {
			return fmt.Errorf("failed to get unpaid status: %w", err)
		}

		unfulfilledStatus, err := repoTx.GetFulfillmentStatusBySlug("unfulfilled")
		if err != nil {
			return fmt.Errorf("failed to get unfulfilled status: %w", err)
		}

		// Get Currency ID (Assuming storefront currency matches code in DB or fallback)
		currencyCode := storeFront.Currency
		if currencyCode == "" {
			currencyCode = "SAR"
		} // Default fallback
		currency, err := repoTx.GetCurrencyByCode(currencyCode)
		if err != nil {
			// Try default SAR if storefront currency is invalid/missing
			currency, err = repoTx.GetCurrencyByCode("SAR")
			if err != nil {
				return fmt.Errorf("failed to get currency: %w", err)
			}
		}

		// Prepare order items
		var orderItems []models.OrderItem
		var subtotal float64
		// ... (Item processing logic remains same, just verify imports if needed)

		// 2. Process Items
		for _, itemReq := range req.Items {
			// Fetch Variant & Product for Snapshot
			var variant models.ProductVariant
			if err := tx.Preload("Inventory").Joins("JOIN products ON products.id = product_variants.product_id").
				Preload("Product").
				First(&variant, itemReq.ProductVariantID).Error; err != nil {
				return fmt.Errorf("variant not found: %d", itemReq.ProductVariantID)
			}

			// ... (Inventory Reservation and Pricing Logic - same as before)
			// Determine Price
			unitPrice := variant.Product.Price
			if variant.Price != nil {
				unitPrice = *variant.Price
			}
			costPrice := 0.0
			if variant.CostPrice != nil {
				costPrice = *variant.CostPrice
			}

			// 3. Reserve Inventory
			if err := s.invService.ReserveStockWithTx(tx, variant.ID, req.StoreFrontID, itemReq.Quantity); err != nil {
				return fmt.Errorf("inventory reservation failed for SKU %s: %w", variant.SKU, err)
			}

			// Build Item
			totalPrice := unitPrice * float64(itemReq.Quantity)
			subtotal += totalPrice

			// Determine Product Name
			nameEn := variant.Product.NameEn
			nameAr := variant.Product.NameAr
			if nameEn == "" {
				nameEn = variant.Product.Name
			}

			orderItems = append(orderItems, models.OrderItem{
				ProductID:             variant.ProductID,
				ProductVariantID:      variant.ID,
				SKU:                   variant.SKU,
				ProductNameSnapshotEn: nameEn,
				ProductNameSnapshotAr: nameAr,
				UnitPrice:             unitPrice,
				CostPrice:             costPrice,
				Quantity:              itemReq.Quantity,
				TotalPrice:            totalPrice,
			})
		}

		// 4. Create Order Header
		orderNumber := fmt.Sprintf("ORD-%d", time.Now().UnixNano())

		newOrder := &models.Order{
			StoreFrontID: req.StoreFrontID,
			OrderNumber:  orderNumber,

			OrderStatusID:       pendingStatus.ID,
			PaymentStatusID:     unpaidStatus.ID,
			FulfillmentStatusID: unfulfilledStatus.ID,
			CurrencyID:          currency.ID,

			CustomerName:  req.CustomerName,
			CustomerEmail: req.CustomerEmail,
			CustomerPhone: req.CustomerPhone,

			Subtotal:       subtotal,
			ShippingAmount: req.ShippingAmount,
			TaxAmount:      req.TaxAmount,
			DiscountAmount: req.DiscountAmount,
			TotalAmount:    subtotal + req.ShippingAmount + req.TaxAmount - req.DiscountAmount,
			Notes:          req.Notes,
			CreatedByID:    adminID,
		}

		if err := repoTx.CreateOrder(tx, newOrder); err != nil {
			return err
		}

		// 5. Bulk Insert Items
		for i := range orderItems {
			orderItems[i].OrderID = newOrder.ID
		}
		if err := repoTx.CreateOrderItems(tx, orderItems); err != nil {
			return err
		}

		order = newOrder
		return nil
	})

	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	// Reload order to get full associations for response
	fullOrder, _ := s.repo.GetOrderByID(order.ID)
	return utils.NewCreatedResource("Order created successfully", fullOrder)
}

// ConfirmOrder confirms stock deduction (Paid -> Confirmed)
func (s *Service) ConfirmOrder(id int64) utils.IResource {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repoTx := &Repository{db: tx}

		order, err := repoTx.GetOrderByID(id)
		if err != nil {
			return err
		}

		// Check Status
		if order.OrderStatus.Slug == "confirmed" || order.OrderStatus.Slug == "completed" {
			return fmt.Errorf("order already confirmed")
		}
		if order.OrderStatus.Slug == "cancelled" {
			return fmt.Errorf("cannot confirm cancelled order")
		}

		// Deduct Stock
		for _, item := range order.Items {
			if err := s.invService.ConfirmStockDeductionWithTx(tx, item.ProductVariantID, order.StoreFrontID, item.Quantity); err != nil {
				return fmt.Errorf("failed to confirm stock deduction for item %s: %w", item.SKU, err)
			}
		}

		// Update Status
		confirmedStatus, err := repoTx.GetOrderStatusBySlug("confirmed")
		if err != nil {
			return err
		}

		if err := repoTx.UpdateStatus(tx, order.ID, confirmedStatus.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}
	return utils.NewOKResource("Order confirmed", nil)
}

// CancelOrder releases reserved stock
func (s *Service) CancelOrder(id int64) utils.IResource {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		repoTx := &Repository{db: tx}

		order, err := repoTx.GetOrderByID(id)
		if err != nil {
			return err
		}

		if order.OrderStatus.Slug == "cancelled" {
			return nil // Already cancelled
		}
		if order.OrderStatus.Slug == "completed" || order.OrderStatus.Slug == "fulfilled" {
			return fmt.Errorf("cannot cancel fulfilled/completed order")
		}

		shouldRelease := order.OrderStatus.Slug == "draft" ||
			order.OrderStatus.Slug == "pending_payment" ||
			order.OrderStatus.Slug == "paid"

		if shouldRelease {
			for _, item := range order.Items {
				if err := s.invService.ReleaseReservedStockWithTx(tx, item.ProductVariantID, order.StoreFrontID, item.Quantity); err != nil {
					return fmt.Errorf("failed to release stock for item %s: %w", item.SKU, err)
				}
			}
		} else if order.OrderStatus.Slug == "confirmed" {
			return fmt.Errorf("cancellation of confirmed orders requires restocking (not implemented in P1)")
		}

		cancelledStatus, err := repoTx.GetOrderStatusBySlug("cancelled")
		if err != nil {
			return err
		}

		if err := repoTx.UpdateStatus(tx, order.ID, cancelledStatus.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}
	return utils.NewOKResource("Order cancelled", nil)
}

func (s *Service) ListOrders(filter requests.OrderFilterRequest, pagination *utils.Pagination) utils.IResource {
	orders, total, err := s.repo.ListOrders(filter, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to list orders", err)
	}
	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Orders retrieved", orders, pagination.GetMeta())
}

func (s *Service) GetOrder(id int64) utils.IResource {
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return utils.NewNotFoundResource("Order not found", nil)
	}
	return utils.NewOKResource("Order details", order)
}

func (s *Service) GetOrderMeta() utils.IResource {
	orderStatuses, err := s.repo.ListOrderStatuses()
	if err != nil {
		return utils.NewInternalErrorResource("Failed to fetch order statuses", err)
	}

	paymentStatuses, err := s.repo.ListPaymentStatuses()
	if err != nil {
		return utils.NewInternalErrorResource("Failed to fetch payment statuses", err)
	}

	fulfillmentStatuses, err := s.repo.ListFulfillmentStatuses()
	if err != nil {
		return utils.NewInternalErrorResource("Failed to fetch fulfillment statuses", err)
	}

	currencies, err := s.repo.ListCurrencies()
	if err != nil {
		return utils.NewInternalErrorResource("Failed to fetch currencies", err)
	}

	return utils.NewOKResource("Order metadata", map[string]interface{}{
		"order_statuses":       orderStatuses,
		"payment_statuses":     paymentStatuses,
		"fulfillment_statuses": fulfillmentStatuses,
		"currencies":           currencies,
	})
}

// UpdateOrder updates order details and items (add/remove/update qty)
func (s *Service) UpdateOrder(id int64, req requests.UpdateOrderRequest) utils.IResource {
	var updatedOrder *models.Order

	err := s.db.Transaction(func(tx *gorm.DB) error {
		repoTx := &Repository{db: tx}

		// 1. Fetch Existing Order
		order, err := repoTx.GetOrderByID(id)
		if err != nil {
			return err
		}

		// Only allow updates for Draft/Pending/Paid orders. Confirmed/Completed orders are locked for structured flows.
		if order.OrderStatus.Slug == "confirmed" || order.OrderStatus.Slug == "completed" || order.OrderStatus.Slug == "cancelled" {
			return fmt.Errorf("cannot update order with status %s", order.OrderStatus.Slug)
		}

		// 2. Update Basic Info
		order.CustomerName = req.CustomerName
		order.CustomerEmail = req.CustomerEmail
		order.CustomerPhone = req.CustomerPhone
		order.Notes = req.Notes
		order.ShippingAmount = req.ShippingAmount
		order.TaxAmount = req.TaxAmount
		order.DiscountAmount = req.DiscountAmount
		// Note: StoreFront cannot be changed easily as it affects currency/inventory context. Ignoring for now.

		// 3. Process Items Diff
		var currentItems = make(map[int64]models.OrderItem)
		for _, item := range order.Items {
			currentItems[item.ID] = item
		}

		var subtotal float64
		// We will reconstruct the order items list
		// Strategies:
		// - ID == 0: New Item -> Reserve Stock -> Create
		// - ID > 0 & IsRemoved: Existing Item -> Release Stock -> Delete
		// - ID > 0 & !IsRemoved: Update Qty -> Diff Stock -> Update

		for _, itemReq := range req.Items {
			if itemReq.ID == 0 {
				// --- NEW ITEM ---
				if itemReq.IsRemoved {
					continue // Ignore new items marked removed
				}

				// Fetch Variant
				var variant models.ProductVariant
				if err := tx.Preload("Product").First(&variant, itemReq.ProductVariantID).Error; err != nil {
					return fmt.Errorf("variant not found: %d", itemReq.ProductVariantID)
				}

				// Reserve Stock
				if err := s.invService.ReserveStockWithTx(tx, variant.ID, order.StoreFrontID, itemReq.Quantity); err != nil {
					return fmt.Errorf("stock reservation failed for new item %s: %w", variant.SKU, err)
				}

				// Pricing
				unitPrice := variant.Product.Price
				if variant.Price != nil {
					unitPrice = *variant.Price
				}
				costPrice := 0.0
				if variant.CostPrice != nil {
					costPrice = *variant.CostPrice
				}
				totalPrice := unitPrice * float64(itemReq.Quantity)
				subtotal += totalPrice

				// Name
				nameEn := variant.Product.NameEn
				if nameEn == "" {
					nameEn = variant.Product.Name
				}

				// Add to DB
				newItem := models.OrderItem{
					OrderID:               order.ID,
					ProductID:             variant.ProductID,
					ProductVariantID:      variant.ID,
					SKU:                   variant.SKU,
					ProductNameSnapshotEn: nameEn,
					ProductNameSnapshotAr: variant.Product.NameAr,
					UnitPrice:             unitPrice,
					CostPrice:             costPrice,
					Quantity:              itemReq.Quantity,
					TotalPrice:            totalPrice,
				}
				if err := tx.Create(&newItem).Error; err != nil {
					return err
				}

			} else {
				// --- EXISTING ITEM ---
				existingItem, exists := currentItems[itemReq.ID]
				if !exists {
					return fmt.Errorf("item id %d not found in order", itemReq.ID)
				}
				delete(currentItems, itemReq.ID) // Mark as processed

				if itemReq.IsRemoved {
					// REMOVE: Release Stock -> Delete
					if err := s.invService.ReleaseReservedStockWithTx(tx, existingItem.ProductVariantID, order.StoreFrontID, existingItem.Quantity); err != nil {
						return fmt.Errorf("stock release failed for removed item %s: %w", existingItem.SKU, err)
					}
					if err := tx.Delete(&existingItem).Error; err != nil {
						return err
					}
					// Do not add to subtotal
				} else {
					// UPDATE: Check Quantity Diff
					qtyDiff := itemReq.Quantity - existingItem.Quantity

					if qtyDiff > 0 {
						// Increase: Reserve more
						if err := s.invService.ReserveStockWithTx(tx, existingItem.ProductVariantID, order.StoreFrontID, qtyDiff); err != nil {
							return fmt.Errorf("stock reservation failed for update %s: %w", existingItem.SKU, err)
						}
					} else if qtyDiff < 0 {
						// Decrease: Release some
						if err := s.invService.ReleaseReservedStockWithTx(tx, existingItem.ProductVariantID, order.StoreFrontID, -qtyDiff); err != nil {
							return fmt.Errorf("stock release failed for update %s: %w", existingItem.SKU, err)
						}
					}

					// Update Item fields
					existingItem.Quantity = itemReq.Quantity
					existingItem.TotalPrice = existingItem.UnitPrice * float64(itemReq.Quantity)
					// Note: We keep original unit price snapshot even if product price changed, unless policy dictates otherwise.

					if err := tx.Save(&existingItem).Error; err != nil {
						return err
					}
					subtotal += existingItem.TotalPrice
				}
			}
		}

		// Any items in currentItems that weren't in the request?
		// If the frontend sends the full list, missing items could imply deletion, but let's strictly use IsRemoved flag or assume the request contains ALL items.
		// For safety in this "Update" implementation, we assume the Request contains the *Complete* desired state or at least all *Changed* items.
		// However, typical frontend logic sends the full "Cart".
		// Let's stick to: Request items are the ones to process. If an item is NOT in the request, it is untouched?
		// Or: Request replaces items?
		// Safer approach for now: Request contains specific updates/adds/removes.
		// If frontend sends the full list every time, we need to handle "missing" items.
		// But let's assume the frontend sends:
		// 1. All current items (some maybe marked IsRemoved, some with new Qty)
		// 2. New items (ID=0)
		// If an ID is missing from request but exists in DB, we leave it alone (safer) OR we delete it (RESTful).
		// Let's assume WE LEAVE IT ALONE unless explicitly flagged IsRemoved. This allows partial updates if needed.
		// BUT we need to calculate total. If we leave it alone, we need to fetch it to add to Subtotal.

		// To correctly calculate Subtotal, we must account for ALL items remaining in the order.
		// So we iterate valid keys in currentItems that were NOT deleted.
		for _, remainingItem := range currentItems {
			subtotal += remainingItem.TotalPrice
		}

		// 4. Update Totals
		order.Subtotal = subtotal
		order.TotalAmount = subtotal + order.ShippingAmount + order.TaxAmount - order.DiscountAmount

		if err := tx.Save(order).Error; err != nil {
			return err
		}

		updatedOrder = order
		return nil
	})

	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	fullOrder, _ := s.repo.GetOrderByID(updatedOrder.ID)
	return utils.NewOKResource("Order updated successfully", fullOrder)
}
