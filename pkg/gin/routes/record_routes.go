package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/cewuandy/go-restful-dns/internal/domain"
)

func RegisterRecordRoutes(r *gin.Engine, handler domain.RecordHandler) {
	group := r.Group(api).Group(v1)
	routes := []Route{
		{
			Name:    "Create DNS Record",
			Group:   record,
			Pattern: ":recordType",
			Method:  http.MethodPost,
			Handler: handler.CreateRecordAPI,
		},
		{
			Name:    "Get DNS Record",
			Group:   record,
			Pattern: "",
			Method:  http.MethodGet,
			Handler: handler.GetRecordAPI,
		},
		{
			Name:    "List all DNS Records",
			Group:   fmt.Sprintf("%ss", record),
			Pattern: "",
			Method:  http.MethodGet,
			Handler: handler.ListRecordsAPI,
		},
		{
			Name:    "Update DNS Record",
			Group:   record,
			Pattern: ":recordType",
			Method:  http.MethodPut,
			Handler: handler.UpdateRecordAPI,
		},
		{
			Name:    "Delete DNS Record",
			Group:   record,
			Pattern: "",
			Method:  http.MethodDelete,
			Handler: handler.DeleteRecordAPI,
		},
	}

	for i := 0; i < len(routes); i++ {
		routes[i].registerURL(group)
	}
}
