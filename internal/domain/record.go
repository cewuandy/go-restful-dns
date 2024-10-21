package domain

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"reflect"
)

type Record struct {
	Name   string `json:"name"`
	RrType uint16 `json:"rrType"`
	Class  uint16 `json:"class"`
	Record string `json:"record"`
}

type ResponseType string

const (
	Answer ResponseType = "Answer"
	Ns     ResponseType = "Ns"
	Extra  ResponseType = "Extra"
)

var ResponseTypeMap = []ResponseType{Answer, Ns, Extra}

type RecordHandler interface {
	CreateRecordAPI(ctx *gin.Context)

	GetRecordAPI(ctx *gin.Context)

	ListRecordsAPI(ctx *gin.Context)

	UpdateRecordAPI(ctx *gin.Context)

	DeleteRecordAPI(ctx *gin.Context)
}

type RecordUseCase interface {
	CreateRecord(ctx context.Context, rr dns.RR) error

	GetRecord(ctx context.Context, question Question) (dns.RR, error)

	ListRecords(ctx context.Context) ([]dns.RR, error)

	UpdateRecord(ctx context.Context, rr dns.RR) error

	DeleteRecord(ctx context.Context, question Question) error
}

type RecordRepo interface {
	Create(ctx context.Context, record *Record) error

	Get(ctx context.Context, name string, rrType uint16, class uint16) (*Record, error)

	List(ctx context.Context) ([]*Record, error)

	Update(ctx context.Context, record *Record) error

	Delete(ctx context.Context, name string, rrType uint16, class uint16) error
}

var RecordTypeMap = map[string]reflect.Type{
	"a":    reflect.TypeOf(A{}),
	"aaaa": reflect.TypeOf(AAAA{}),
}
