package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubReportService is a minimal stub implementing service.ReportService for handler tests.
type stubReportService struct {
	lastListQuery dto.ListReportsQuery
}

var _ service.ReportService = (*stubReportService)(nil)

func (s *stubReportService) CreateProcessingReport(_ context.Context, _, _, _ string, _ json.RawMessage) (*model.Report, error) {
	return &model.Report{}, nil
}
func (s *stubReportService) SucceedReport(_ context.Context, _ string, _ json.RawMessage, _ string, _ *int, _ *string, _ []model.GeneratedFile) (*dto.ReportDTO, error) {
	return &dto.ReportDTO{}, nil
}
func (s *stubReportService) FailReport(_ context.Context, _ string, _ string) error { return nil }
func (s *stubReportService) FallbackReport(_ context.Context, _ string, _ json.RawMessage, _ string) error {
	return nil
}
func (s *stubReportService) GetReport(_ context.Context, _ string) (*dto.ReportDTO, error) {
	return &dto.ReportDTO{}, nil
}
func (s *stubReportService) ListReports(_ context.Context, query dto.ListReportsQuery) (*dto.PaginatedResponse[dto.ReportDTO], error) {
	s.lastListQuery = query
	return &dto.PaginatedResponse[dto.ReportDTO]{Items: []dto.ReportDTO{}, Total: 0, Page: query.Page, PageSize: query.PageSize}, nil
}
func (s *stubReportService) DeleteReport(_ context.Context, _ string) error { return nil }
func (s *stubReportService) GetDashboardStats(_ context.Context) (*dto.DashboardStatsDTO, error) {
	return &dto.DashboardStatsDTO{}, nil
}

func TestListReports_PassesStatusParameter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &stubReportService{}
	router := gin.New()
	RegisterReportRoutes(&router.RouterGroup, stub)

	req := httptest.NewRequest(http.MethodGet, "/reports?status=succeeded&page=1&page_size=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "succeeded", stub.lastListQuery.Status)
	assert.Equal(t, 1, stub.lastListQuery.Page)
	assert.Equal(t, 5, stub.lastListQuery.PageSize)
}

func TestListReports_CombinedFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &stubReportService{}
	router := gin.New()
	RegisterReportRoutes(&router.RouterGroup, stub)

	req := httptest.NewRequest(http.MethodGet, "/reports?tool_type=ui_review&status=failed&sort=oldest&page=2&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ui_review", stub.lastListQuery.ToolType)
	assert.Equal(t, "failed", stub.lastListQuery.Status)
	assert.Equal(t, "oldest", stub.lastListQuery.Sort)
	assert.Equal(t, 2, stub.lastListQuery.Page)
	assert.Equal(t, 20, stub.lastListQuery.PageSize)
}

func TestListReports_NoStatusParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &stubReportService{}
	router := gin.New()
	RegisterReportRoutes(&router.RouterGroup, stub)

	req := httptest.NewRequest(http.MethodGet, "/reports", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "", stub.lastListQuery.Status, "status should be empty when not provided")
	assert.Equal(t, 0, stub.lastListQuery.Page, "raw page=0, SetDefaults called in service layer")
	assert.Equal(t, 0, stub.lastListQuery.PageSize, "raw page_size=0, SetDefaults called in service layer")
}

func TestListReports_ResponseStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &stubReportService{}
	router := gin.New()
	RegisterReportRoutes(&router.RouterGroup, stub)

	req := httptest.NewRequest(http.MethodGet, "/reports?status=succeeded", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp util.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, util.CodeSuccess, resp.Code)

	dataBytes, _ := json.Marshal(resp.Data)
	var paginated dto.PaginatedResponse[dto.ReportDTO]
	err = json.Unmarshal(dataBytes, &paginated)
	require.NoError(t, err)
	assert.Equal(t, int64(0), paginated.Total)
	assert.NotNil(t, paginated.Items)
}

func TestListReports_InvalidStatusPassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &stubReportService{}
	router := gin.New()
	RegisterReportRoutes(&router.RouterGroup, stub)

	req := httptest.NewRequest(http.MethodGet, "/reports?status=invalid_xss<script>", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
