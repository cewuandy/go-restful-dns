package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"github.com/samber/do"
	"net/http"
	"reflect"

	"github.com/cewuandy/go-restful-dns/internal/domain"

	"github.com/pkg/errors"
)

type recordHandler struct {
	recordUseCase domain.RecordUseCase
}

// CreateRecordAPI ...
// @title CreateRecordAPI
// @description Create a new dns record
// @tags Record
// @accept json
// @param body body domain.A true "The example of A record request body"
// @success 201 {object} domain.A
// @failure 400 {object} domain.Error
// @router /record/{recordType} [POST]
func (r *recordHandler) CreateRecordAPI(ctx *gin.Context) {
	recordType := ctx.Param("recordType")
	v, ok := domain.RecordTypeMap[recordType]
	if !ok {
		err := domain.Error{
			Message:    "This type doesn't support currently",
			StatusCode: http.StatusBadRequest,
		}
		_ = ctx.Error(err)
		return
	}

	record := reflect.New(v).Interface()

	err := ctx.ShouldBindJSON(&record)
	if err != nil {
		err = &domain.Error{
			Message:    fmt.Sprintf("Bind JSON error: %s", err.Error()),
			Err:        errors.New(err.Error()),
			StatusCode: http.StatusBadRequest,
		}
		_ = ctx.Error(err)
		return
	}

	output := reflect.ValueOf(record).MethodByName("String").Call(nil)[0].String()
	rr, _ := dns.NewRR(output)
	err = r.recordUseCase.CreateRecord(ctx, rr)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, rr)
}

// GetRecordAPI ...
// @title GetRecordAPI
// @description Get dns record by name, qtype, qclass
// @tags Record
// @accept json
// @param name query string true "Domain Name"
// @param qtype query string true "Record Type"
// @param qclass query string true "Record Class"
// @success 200 {object} domain.A
// @failure 400 {object} domain.Error
// @router /record [GET]
func (r *recordHandler) GetRecordAPI(ctx *gin.Context) {
	var (
		question domain.Question
		rr       dns.RR
		err      error
	)

	err = ctx.ShouldBindQuery(&question)
	if err != nil {
		err = domain.Error{
			Message:    fmt.Sprintf("Bind query error: %s", err.Error()),
			StatusCode: http.StatusBadRequest,
		}
		_ = ctx.Error(err)
		return
	}

	rr, err = r.recordUseCase.GetRecord(ctx, question)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, rr)
}

// ListRecordsAPI ...
// @title ListRecordsAPI
// @description List all dns record
// @tags Record
// @accept json
// @success 200 {object} []domain.A
// @failure 400 {object} domain.Error
// @router /records [GET]
func (r *recordHandler) ListRecordsAPI(ctx *gin.Context) {
	var (
		rrs []dns.RR
		err error
	)

	rrs, err = r.recordUseCase.ListRecords(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, rrs)
}

// UpdateRecordAPI ...
// @title UpdateRecordAPI
// @description Update an existed dns record
// @tags Record
// @accept json
// @param body body domain.A true "The example of A record request body"
// @success 200 {object} domain.A
// @failure 400 {object} domain.Error
// @router /record/{recordType} [PUT]
func (r *recordHandler) UpdateRecordAPI(ctx *gin.Context) {
	recordType := ctx.Param("recordType")
	v, ok := domain.RecordTypeMap[recordType]
	if !ok {
		err := domain.Error{
			Message:    "This type doesn't support currently",
			StatusCode: http.StatusBadRequest,
		}
		_ = ctx.Error(err)
		return
	}

	record := reflect.New(v).Interface()
	err := ctx.ShouldBindJSON(&record)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	output := reflect.ValueOf(record).MethodByName("String").Call(nil)[0].String()
	rr, _ := dns.NewRR(output)
	err = r.recordUseCase.UpdateRecord(ctx, rr)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, rr)
}

// DeleteRecordAPI ...
// @title DeleteRecordAPI
// @description Delete dns record by name, qtype, qclass
// @tags Record
// @accept json
// @param body body dns.Question true "The example of Question request body"
// @success 200 {object} domain.A
// @failure 400 {object} domain.Error
// @router /record [DELETE]
func (r *recordHandler) DeleteRecordAPI(ctx *gin.Context) {
	var question domain.Question

	err := ctx.ShouldBindJSON(&question)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	err = r.recordUseCase.DeleteRecord(ctx, question)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func NewRecordHandler(injector *do.Injector) (domain.RecordHandler, error) {
	return &recordHandler{do.MustInvoke[domain.RecordUseCase](injector)}, nil
}
