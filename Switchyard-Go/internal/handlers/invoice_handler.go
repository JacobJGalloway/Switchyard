package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
	"github.com/JacobJGalloway/switchyard-go/internal/repository"
)

// InvoiceHandler serves read-only invoice and delivery confirmation lookups.
type InvoiceHandler struct {
	invoiceRepo repository.InvoiceRepository
}

func NewInvoiceHandler(invoiceRepo repository.InvoiceRepository) *InvoiceHandler {
	return &InvoiceHandler{invoiceRepo: invoiceRepo}
}

// GetInvoice handles GET /api/invoice/:id
func (h *InvoiceHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}
	inv, err := h.invoiceRepo.GetInvoice(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "invoice not found")
		return
	}
	writeJSON(w, http.StatusOK, inv)
}

// GetInvoicesByStore handles GET /api/invoice/store/:storeId
func (h *InvoiceHandler) GetInvoicesByStore(w http.ResponseWriter, r *http.Request) {
	storeID := chi.URLParam(r, "storeId")
	if storeID == "" {
		writeError(w, http.StatusBadRequest, "storeId is required")
		return
	}
	invoices, err := h.invoiceRepo.GetInvoicesByStore(r.Context(), storeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch invoices")
		return
	}
	if invoices == nil {
		invoices = []*models.InternalInvoice{}
	}
	writeJSON(w, http.StatusOK, invoices)
}
