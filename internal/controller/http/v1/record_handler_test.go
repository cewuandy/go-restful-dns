package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"github.com/samber/do"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cewuandy/go-restful-dns/internal/controller/http/middleware"
	"github.com/cewuandy/go-restful-dns/internal/domain"
	"github.com/cewuandy/go-restful-dns/internal/domain/mocks"
	"github.com/cewuandy/go-restful-dns/pkg/gin/routes"
)

type recordHandlerTestSuite struct {
	suite.Suite

	handler domain.RecordHandler

	recordUsecase *mocks.RecordUseCase

	r *gin.Engine

	exampleA domain.A
}

func TestRecordHandler(t *testing.T) {
	suite.Run(t, &recordHandlerTestSuite{})
}

func (t *recordHandlerTestSuite) SetupSuite() {
	injector := do.New()
	t.recordUsecase = &mocks.RecordUseCase{}
	do.ProvideValue[domain.RecordUseCase](injector, t.recordUsecase)
	do.Provide[domain.RecordHandler](injector, NewRecordHandler)
	do.Provide[domain.ErrorHandler](injector, middleware.NewErrorHandler)

	t.r = gin.New()
	t.r.Use(do.MustInvoke[domain.ErrorHandler](injector).HandleError)

	routes.RegisterRecordRoutes(t.r, do.MustInvoke[domain.RecordHandler](injector))

	t.exampleA = domain.A{
		Hdr: domain.RR_Header{
			Name:   "test.com.",
			Rrtype: domain.TypeA,
			Class:  domain.ClassINET,
			Ttl:    1440,
		},
		Address: net.ParseIP("1.1.1.1"),
	}
}

func (t *recordHandlerTestSuite) SetupTest() {
	var (
		anyContext  = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyRR       = mock.MatchedBy(func(rr dns.RR) bool { return true })
		anyQuestion = mock.AnythingOfType("domain.Question")
		a           = dns.A{
			Hdr: dns.RR_Header{
				Name:   "test.com.",
				Rrtype: 1,
				Class:  1,
				Ttl:    1440,
			}, A: net.ParseIP("1.1.1.1"),
		}
	)

	rr, _ := dns.NewRR(a.String())

	t.recordUsecase.
		On("CreateRecord", anyContext, anyRR).
		Return(nil)
	t.recordUsecase.
		On("GetRecord", anyContext, anyQuestion).
		Return(rr, nil)
	t.recordUsecase.
		On("ListRecords", anyContext).
		Return([]dns.RR{rr}, nil)
	t.recordUsecase.
		On("UpdateRecord", anyContext, anyRR).
		Return(nil)
	t.recordUsecase.
		On("DeleteRecord", anyContext, anyQuestion).
		Return(nil)
}

func (t *recordHandlerTestSuite) SetupErrorTest() {
	t.SetupTest()
	t.recordUsecase.ExpectedCalls = nil
}

func (t *recordHandlerTestSuite) TestCreateRecordAPI() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyRR      = mock.MatchedBy(func(rr dns.RR) bool { return true })
	)

	t.Run(
		"success", func() {
			recorder := httptest.NewRecorder()
			raw, _ := json.Marshal(t.exampleA)
			request, err := http.NewRequest(
				http.MethodPost, "/api/v1/record/a", bytes.NewBuffer(raw),
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusCreated, recorder.Code)
			t.Contains(recorder.Body.String(), "test.com.")
		},
	)

	t.Run(
		"support_error", func() {
			recorder := httptest.NewRecorder()
			raw, _ := json.Marshal(t.exampleA)
			request, err := http.NewRequest(
				http.MethodPost, "/api/v1/record/unsupportRecord", bytes.NewBuffer(raw),
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusBadRequest, recorder.Code)
			t.Contains(recorder.Body.String(), "This type doesn't support currently")
		},
	)

	t.Run(
		"bind_json_error", func() {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(
				http.MethodPost, "/api/v1/record/a", nil,
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusBadRequest, recorder.Code)
			t.Contains(recorder.Body.String(), "Bind JSON error:")
		},
	)

	t.Run(
		"CreateRecord_error", func() {
			t.SetupErrorTest()
			t.recordUsecase.
				On("CreateRecord", anyContext, anyRR).
				Return(&domain.Error{Message: "test-error", StatusCode: http.StatusBadRequest})
			recorder := httptest.NewRecorder()
			raw, _ := json.Marshal(t.exampleA)
			request, err := http.NewRequest(
				http.MethodPost, "/api/v1/record/a", bytes.NewBuffer(raw),
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusBadRequest, recorder.Code)
			t.Contains(recorder.Body.String(), "test-error")
		},
	)
}

func (t *recordHandlerTestSuite) TestGetRecordAPI() {
	var (
		anyContext  = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyQuestion = mock.AnythingOfType("domain.Question")
	)
	t.Run(
		"success", func() {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(
				http.MethodGet, "/api/v1/record?name=test.com.&qtype=A&qclass=INET", nil,
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusOK, recorder.Code)
			t.Contains(recorder.Body.String(), "test.com.")
		},
	)

	t.Run(
		"bind_query_error", func() {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(
				http.MethodGet, "/api/v1/record?name=test.com.&qtype=A", nil,
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusBadRequest, recorder.Code)
			t.Contains(recorder.Body.String(), "Bind query error:")
		},
	)

	t.Run(
		"GetRecord_error", func() {
			t.SetupErrorTest()
			t.recordUsecase.
				On("GetRecord", anyContext, anyQuestion).
				Return(nil, &domain.Error{Message: "test-error", StatusCode: http.StatusBadRequest})
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(
				http.MethodGet, "/api/v1/record?name=test.com.&qtype=A&qclass=INET", nil,
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusBadRequest, recorder.Code)
			t.Contains(recorder.Body.String(), "test-error")
		},
	)
}

func (t *recordHandlerTestSuite) TestListRecordsAPI() {
	var anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })

	t.Run(
		"success", func() {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/api/v1/records", nil)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusOK, recorder.Code)
			t.Contains(recorder.Body.String(), "test.com.")

		},
	)

	t.Run(
		"ListRecords_error", func() {
			t.SetupErrorTest()
			t.recordUsecase.
				On("ListRecords", anyContext).
				Return(nil, &domain.Error{Message: "test-error", StatusCode: http.StatusBadRequest})
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/api/v1/records", nil)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusBadRequest, recorder.Code)
			t.Contains(recorder.Body.String(), "test-error")
		},
	)
}

func (t *recordHandlerTestSuite) TestUpdateRecordAPI() {
	var (
		anyContext = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyRR      = mock.MatchedBy(func(rr dns.RR) bool { return true })
	)

	t.Run(
		"success", func() {
			recorder := httptest.NewRecorder()
			raw, _ := json.Marshal(t.exampleA)
			request, err := http.NewRequest(
				http.MethodPut, "/api/v1/record/a", bytes.NewBuffer(raw),
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusOK, recorder.Code)
			t.Contains(recorder.Body.String(), "test.com.")
		},
	)

	t.Run(
		"support_error", func() {
			recorder := httptest.NewRecorder()
			raw, _ := json.Marshal(t.exampleA)
			request, err := http.NewRequest(
				http.MethodPut, "/api/v1/record/unsupportRecord", bytes.NewBuffer(raw),
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusBadRequest, recorder.Code)
			t.Contains(recorder.Body.String(), "This type doesn't support currently")
		},
	)

	t.Run(
		"bind_json_error", func() {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(
				http.MethodPut, "/api/v1/record/a", nil,
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusBadRequest, recorder.Code)
			t.Contains(recorder.Body.String(), "Bind JSON error:")
		},
	)

	t.Run(
		"UpdateRecord_error", func() {
			t.SetupErrorTest()
			t.recordUsecase.
				On("UpdateRecord", anyContext, anyRR).
				Return(&domain.Error{Message: "test-error", StatusCode: http.StatusBadRequest})
			recorder := httptest.NewRecorder()
			raw, _ := json.Marshal(t.exampleA)
			request, err := http.NewRequest(
				http.MethodPut, "/api/v1/record/a", bytes.NewBuffer(raw),
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusBadRequest, recorder.Code)
			t.Contains(recorder.Body.String(), "test-error")
		},
	)
}

func (t *recordHandlerTestSuite) TestDeleteRecordAPI() {
	var (
		anyContext  = mock.MatchedBy(func(ctx context.Context) bool { return true })
		anyQuestion = mock.AnythingOfType("domain.Question")
		q           = domain.Question{
			Name:   "test.com.",
			Qtype:  domain.TypeA,
			Qclass: domain.ClassINET,
		}
	)

	t.Run(
		"success", func() {
			recorder := httptest.NewRecorder()
			raw, _ := json.Marshal(q)
			request, err := http.NewRequest(
				http.MethodDelete, "/api/v1/record", bytes.NewBuffer(raw),
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusNoContent, recorder.Code)
			t.Equal(0, len(recorder.Body.String()))
		},
	)

	t.Run(
		"bind_json_error", func() {
			recorder := httptest.NewRecorder()
			errorQ := domain.Question{}
			raw, _ := json.Marshal(errorQ)
			request, err := http.NewRequest(
				http.MethodDelete, "/api/v1/record", bytes.NewBuffer(raw),
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusBadRequest, recorder.Code)
			t.Contains(recorder.Body.String(), "Bind JSON error:")
		},
	)

	t.Run(
		"DeleteRecord_error", func() {
			t.SetupErrorTest()
			t.recordUsecase.
				On("DeleteRecord", anyContext, anyQuestion).
				Return(&domain.Error{Message: "test-error", StatusCode: http.StatusBadRequest})
			recorder := httptest.NewRecorder()
			raw, _ := json.Marshal(q)
			request, err := http.NewRequest(
				http.MethodDelete, "/api/v1/record", bytes.NewBuffer(raw),
			)
			t.Nil(err)

			t.r.ServeHTTP(recorder, request)

			t.Equal(http.StatusBadRequest, recorder.Code)
			t.Contains(recorder.Body.String(), "test-error")
		},
	)
}
