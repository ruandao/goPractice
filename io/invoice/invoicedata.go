package main

import (
	"time"
	"io"
	"fmt"
	"encoding/json"
	"github.com/juju/errors"
)

type Invoice struct {
	Id				int
	CustomerId		int
	Raised			time.Time
	Due				time.Time
	Paid			bool
	Note			string
	Items			[]*Item
}

type Item struct {
	Id				string
	Price			float64
	Quantity		int
	Note			string
}

const (
	fileType 		=	"INVOICES"		// 用于纯文本格式
	magicNumber		=	0x125D			// 用于二进制格式
	fileVersion		=	100				// 用于所有格式
	dataFormat		=	"2006-01-02"	// 必须总是使用该日期
)

type InvoiceMarshaler interface {
	MarshalInvoices(writer io.Writer, invoices []*Invoice) error
}
type InvoiceUnMarshaler interface {
	UnmarshalInvoices(reader io.Reader) ([]*Invoice, error)
}

func readInvoices(reader io.Reader, suffix string) ([]*Invoice, error) {
	var unmarshaler	InvoiceUnMarshaler
	switch suffix {
	case ".gob":
		unmarshaler = GobMarshaler{}
	case ".inv":
		unmarshaler = InvMarshaler{}
	case ".jsn", ".json":
		unmarshaler = JSONMarshaler{}
	case ".txt":
		unmarshaler = TxtMarshaler{}
	case ".xml":
		unmarshaler = XMLMarshaler{}
	}
	if unmarshaler != nil {
		return unmarshaler.UnmarshalInvoices(reader)
	}
	return nil, fmt.Errorf("unrecognized input suffix: %s", suffix)
}
