package domain

import (
	"github.com/miekg/dns"
)

type Question struct {
	dns.Question
	Name   string `json:"name" form:"name" binding:"required"`
	Qtype  RRType `json:"qtype" form:"qtype" binding:"required"`
	Qclass Class  `json:"qclass" form:"qclass" binding:"required"`
}

func (q *Question) String() string {
	if len(q.Name) == 0 && len(q.Qtype) == 0 && len(q.Qclass) == 0 {
		return q.Question.String()
	}
	q.Question.Name = q.Name
	q.Question.Qtype = RRTypeMap[q.Qtype]
	q.Question.Qclass = ClassMap[q.Qclass]
	return q.Question.String()
}
