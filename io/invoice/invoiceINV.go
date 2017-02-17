package main

import (
	"encoding/binary"
	"io"
	"time"
	"strconv"
	"github.com/juju/errors"
	"fmt"
)

type InvMarshaler struct{}

var byteOrder = binary.LittleEndian

func (InvMarshaler)MarshalInvoices(writer io.Writer, invoices []*Invoice) error {
	var write invWriterFunc = func(x interface{}) error {
		return binary.Write(writer, byteOrder, x)
	}
	if err := write(uint32(magicNumber)); err != nil {
		return err
	}
	if err := write(uint16(fileVersion)); err != nil {
		return err
	}
	if err := write(uint32(len(invoices))); err != nil {
		return err
	}
	for _, invoice := range invoices {
		if err := write.WriteInvoice(invoice); err != nil {
			return err
		}
	}
	return nil
}

type invWriterFunc func(interface{})error

func (write invWriterFunc)WriteInvoice(invoice *Invoice) error {
	for _, i := range []int{invoice.Id, invoice.CustomerId} {
		if err := write(int32(i)); err != nil {
			return err
		}
	}
	for _, date := range []time.Time{invoice.Raised, invoice.Due} {
		if err := write.WriteDate(date); err != nil {
			return err
		}
	}
	if err := write.WriteBool(invoice.Paid); err != nil {
		return err
	}
	if err := write.WriteString(invoice.Note); err != nil {
		return err
	}
	if err := write(uint32(invoice.Items)); err != nil {
		return err
	}
	for _, item := range invoice.Items {
		if err := write.WriteItem(item); err != nil {
			return err
		}
	}
	return nil
}

const invDateFormat  = "20160102"

func (write invWriterFunc)WriteDate(date time.Time) error {
	i, err := strconv.Atoi(date.Format(invDateFormat))
	if err != nil {
		return err
	}
	return write(int32(i))
}
func (write invWriterFunc)WriteBool(b bool) error {
	var v int8
	if b {
		v = 1
	}
	return write(b)
}
func (write invWriterFunc)WriteString(note string) error {
	if err := write(uint32(len(note))); err != nil {
		return err
	}
	return write([]byte(note))
}

func (write invWriterFunc)WriteItem(item *Item) error {
	if err := write.WriteString(item.Id); err != nil {
		return err
	}
	if err := write(item.Price); err != nil {
		return err
	}
	if err := write(int16(item.Quantity)); err != nil {
		return err
	}
	if err := write(item.Note); err != nil {
		return err
	}
	return nil
}

func (InvMarshaler)UnmarshalInvoice(reader io.Reader) (invoices []*Invoice, err error) {
	if err := checkInvVersion(reader); err != nil {
		return nil, err
	}
	count, err := readIntFromInt32(reader)
	if err != nil {
		return nil, err
	}
	invoices = make([]*Invoice, 0, count)
	for i := 0; i < count; i++ {
		invoice, err := readInvInvoice(reader)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

func checkInvVersion(reader io.Reader) error {
	var magic uint32
	if err := binary.Read(reader, byteOrder, &magic); err != nil {
		return err
	}
	if magic != magicNumber {
		return  errors.New("cannot read non-invoices inv file")
	}
	var version uint16
	if err := binary.Read(reader, byteOrder, &version); err != nil {
		return err
	}
	if version > fileVersion {
		return fmt.Errorf("version %d is too new to read", version)
	}
	return nil
}

func readIntFromInt32(reader io.Reader) (int, error) {
	var count int32
	err := binary.Read(reader, byteOrder, &count)
	return int(count), err
}

func readInvInvoice(reader io.Reader) (invoice *Invoice, err error) {
	invoice = &Invoice{}
	for _, pId := range []*int{&invoice.Id, &invoice.CustomerId} {
		if *pId, err = readIntFromInt32(reader); err != nil {
			return nil, err
		}
	}
	for _, pDate := range []*time.Time{&invoice.Raised, &invoice.Due} {
		if *pDate, err = readInvDate(reader); err != nil {
			return nil, err
		}
	}
	if invoice.Paid, err = readBoolFromInt8(reader); err != nil {
		return nil, err
	}
	if invoice.Note, err = readInvString(reader); err != nil {
		return nil, err
	}
	var count int
	if count, err = readIntFromInt32(reader); err != nil {
		return nil, err
	}
	invoice.Items, err = readInvItems(reader, count)
	return invoice, err
}

func readInvDate(reader io.Reader) (time.Time, error) {
	var n int32
	if err := binary.Read(reader, byteOrder, &n); err != nil {
		return time.Time{}, err
	}
	return time.Parse(invDateFormat, fmt.Sprint(n))
}

func readBoolFromInt8(reader io.Reader) (bool, error) {
	var n int8
	err := binary.Read(reader, byteOrder, &n)
	return n == 1, err
}
func readInvString(reader io.Reader) (string, error) {
	var count int32
	if err := binary.Read(reader, byteOrder, &count); err != nil {
		return "", err
	}
	data := make([]byte, count)
	err := binary.Read(reader, byteOrder, &data)
	return string(data), err
}
func readInvItems(reader io.Reader, count int) (items []*Item, err error) {
	items = make([]*Item, 0, count)
	for i:=0; i<count; i++ {
		if item, err := readInvItem(reader); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
func readInvItem(reader io.Reader) (item *Item, err error) {
	item = &Item{}
	if item.Id, err = readInvString(reader); err != nil {
		return nil, err
	}
	if err = binary.Read(reader, byteOrder, &item.Price); err != nil {
		return nil, err
	}
	var quantity int16
	if err = binary.Read(reader, byteOrder, &quantity); err != nil {
		return nil, err
	}
	item.Quantity = int(quantity)
	item.Note, err = readInvString(reader)
	return item, err
}