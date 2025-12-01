package v1

import (
	"io"
	"net/http"
	"strings"

	"github.com/andibalo/flip-test/internal/entity"
	"github.com/andibalo/flip-test/internal/service"
	"github.com/andibalo/flip-test/pkg/common"
	"github.com/andibalo/flip-test/pkg/httpresp"
	"github.com/andibalo/flip-test/pkg/pagination"
	pkgsort "github.com/andibalo/flip-test/pkg/sort"
	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
)

const (
	//transactionBasePath = "/api/v1/transaction/"
	transactionBasePath = ""
)

type TransactionController struct {
	transactionService service.TransactionService
}

func NewTransactionController(transactionService service.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

func (tc *TransactionController) AddRoutes(r *gin.Engine) {
	v1 := r.Group(transactionBasePath)
	{
		v1.POST("/upload", tc.UploadCSV)
		v1.GET("/balance", tc.GetBalance)
		v1.GET("/issues", tc.GetIssues)
	}
}

// UploadCSV handles POST /upload
// @Summary Upload CSV file with transaction data
// @Description Accepts CSV file upload, parses it, and stores transactions in memory
// @Tags transactions
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file"
// @Success 200 {object} entity.UploadCSVResponse
// @Failure 400 {object} httpresp.HTTPErrResp
// @Failure 500 {object} httpresp.HTTPErrResp
// @Router /upload [post]
func (tc *TransactionController) UploadCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		httpresp.HttpRespError(c, oops.
			Code(httpresp.BadRequest.AsString()).
			With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
			Errorf("file parameter is required"))
		return
	}

	if file.Header.Get("Content-Type") != "text/csv" && !common.IsCsvFile(file.Filename) {
		httpresp.HttpRespError(c, oops.
			Code(httpresp.BadRequest.AsString()).
			With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
			Errorf("file must be a CSV"))
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		httpresp.HttpRespError(c, oops.
			Code(httpresp.ServerError.AsString()).
			With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).
			Errorf("failed to open uploaded file"))
		return
	}
	defer openedFile.Close()

	fileContent, err := io.ReadAll(openedFile)
	if err != nil {
		httpresp.HttpRespError(c, oops.
			Code(httpresp.ServerError.AsString()).
			With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).
			Errorf("failed to read file content"))
		return
	}

	count, err := tc.transactionService.UploadCSVFile(c.Request.Context(), fileContent)
	if err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	response := &entity.UploadCSVResponse{
		TransactionsUploaded: count,
	}

	httpresp.HttpRespSuccess(c, response, nil)
}

// GetBalance handles GET /balance
// @Summary Get total balance
// @Description Returns total balance (credits - debits from successful transactions only)
// @Tags transactions
// @Produce json
// @Success 200 {object} entity.BalanceResponse
// @Failure 500 {object} httpresp.HTTPErrResp
// @Router /balance [get]
func (tc *TransactionController) GetBalance(c *gin.Context) {
	balance, err := tc.transactionService.GetTotalBalance(c.Request.Context())
	if err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	response := &entity.BalanceResponse{
		TotalBalance: balance,
	}

	httpresp.HttpRespSuccess(c, response, nil)
}

// GetIssues handles GET /issues
// @Summary Get unsuccessful transactions
// @Description Returns non-successful transactions (FAILED + PENDING) with pagination
// @Tags transactions
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} entity.IssuesResponse
// @Failure 500 {object} httpresp.HTTPErrResp
// @Router /issues [get]
func (tc *TransactionController) GetIssues(c *gin.Context) {
	var queryParams entity.GetIssuesQueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	var sorts pkgsort.Sorts
	if queryParams.Sorts != "" {
		sortValues := strings.Split(queryParams.Sorts, ",")
		sorts = pkgsort.ParseMultipleSorts(sortValues)

		if len(sorts.Data()) > 1 {
			httpresp.HttpRespError(c, oops.
				Code(httpresp.BadRequest.AsString()).
				With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
				Errorf("only single column sorting is supported"))
			return
		}

		allowedColumns := []string{"timestamp", "name", "type", "amount", "status", "description"}
		if err := sorts.Validate(allowedColumns); err != nil {
			httpresp.HttpRespError(c, oops.
				Code(httpresp.BadRequest.AsString()).
				With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).
				Wrapf(err, "invalid sort parameter"))
			return
		}
	}

	page := queryParams.GetPageWithDefault()
	pageSize := queryParams.GetPageSizeWithDefault()

	transactions, totalCount, err := tc.transactionService.GetUnsuccessfulTransactions(c.Request.Context(), entity.GetIssuesFilter{
		Sorts:             sorts,
		PaginationRequest: queryParams.PaginationRequest,
	})
	if err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	response := &entity.IssuesResponse{
		Transactions: transactions,
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)
	paginationResp := &pagination.Pagination{
		CurrentPage:     int64(page),
		CurrentElements: int64(len(transactions)),
		TotalPages:      totalPages,
		TotalElements:   totalCount,
		SortBy:          queryParams.Sorts,
	}

	httpresp.HttpRespSuccess(c, response, paginationResp)
}
