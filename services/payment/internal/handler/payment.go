package handler

import (
	"payment-system/pkg/logger"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/walker-16/payment-system/services/payment/internal/domain"
	"github.com/walker-16/payment-system/services/payment/internal/order"
	"github.com/walker-16/payment-system/services/payment/internal/repository"
)

// PaymentHandler handles payment HTTP requests.
type PaymentHandler struct {
	repository   repository.PaymentRepo
	orderService order.Service
	logger       logger.Logger
}

// NewPaymentHandler creates a new instance of PaymentHandler.
func NewPaymentHandler(orderService order.Service,
	repository repository.PaymentRepo,
	logger logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		orderService: orderService,
		repository:   repository,
		logger:       logger,
	}
}

// PaymentRequest represents the payload for creating a new payment.
type PaymentRequest struct {
	ExternalOrderID uuid.UUID `json:"external_order_id"`
}

// PaymentResponse represents the payload returned after a successful payment creation.
type PaymentResponse struct {
	PaymentID uuid.UUID `json:"payment_id"`
}

// CreatePayment handles POST /v1/payments requests.
// It validates headers, parses the request body, and returns a confirmation response.
// Headers required:
//   - idempotency-key: a unique key to ensure idempotent requests.
//   - x-user-id: the ID of the user making the request.
func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// get idempotency-key field from header.
	strIdempotencyKey := c.Get("idempotency-key")
	if strIdempotencyKey == "" {
		return fiber.NewError(fiber.StatusBadRequest,
			"idempotency-key header is required")
	}
	idempotencyKet, err := uuid.Parse(strIdempotencyKey)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest,
			"idempotency-key invalid")
	}

	// get x-user-id fiel from header.
	strUserID := c.Get("x-user-id")
	if strUserID == "" {
		return fiber.NewError(fiber.StatusBadRequest,
			"x-user-id header is required")
	}
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest,
			"x-user-id invalid")
	}

	// request body parse.
	var request PaymentRequest
	if err := c.BodyParser(&request); err != nil {
		h.logger.Error("failed to parse payment request",
			logger.Error(err))
		return fiber.NewError(fiber.StatusBadRequest,
			"invalid JSON body")
	}

	// check extenal orderID.
	if request.ExternalOrderID.String() == "" {
		return fiber.NewError(fiber.StatusBadRequest,
			"external_order_id is required")
	}

	// retrieve the order details by external order ID and ensure it
	// belongs to the specified user ID.
	order, err := h.orderService.GetOrderByExternalIDForUser(
		ctx, request.ExternalOrderID, uint32(userID))
	if err != nil {
		h.logger.Error("order validation failed", logger.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// create payment.
	paymentID := uuid.New()
	payment := &domain.Payment{
		PaymentID:       uuid.New(),
		ExternalOrderID: request.ExternalOrderID,
		UserID:          uint32(userID),
		IdempotencyKey:  idempotencyKet,
		Amount:          order.Amount,
		Currency:        order.Currency,
		Status:          "PENDING",
	}

	// TODO: check idempotency-id.

	// insert payment.
	err = h.repository.InsertPayment(ctx, payment)
	if err != nil {
		h.logger.Error("failed to insert payment", logger.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError,
			"failed to create payment")
	}

	// create payment response.
	response := &PaymentResponse{
		PaymentID: paymentID,
	}
	return c.Status(fiber.StatusAccepted).JSON(response)
}
