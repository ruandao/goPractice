package main

import (
	"io"
	"bufio"
	"fmt"
	"strings"
	"time"
	"github.com/juju/errors"
)

type TXTMarshaler struct{}

func (TXTMarshaler)MarshalInvoice(writer io.Writer, invoices []*Invoice) error {
	bufferedWriter := bufio.NewWriter(writer)
	defer bufferedWriter.Flush()
	var write writerFunc = func(format string, args ...interface{}) error {
		_, err := fmt.Fprintf(bufferedWriter, format, args...)
		return err
	}
	if err := write("%s %d\n", fileType, fileVersion); err != nil {
		return err
	}
	for _, invoice := range invoices {
		if err := write.writeInvoice(invoice); err != nil {
			return err
		}
	}
	return nil
}

const noteSep = ":"

type writerFunc func(string, ...interface{}) error

func (write writerFunc)writeInvoice(invoice *Invoice) error {
	note := ""
	if invoice.Note != "" {
		note = noteSep + " " + invoice.Note
	}
	if err := write("INVOTE ID=%d CUSTOMER=%d RAISED=%s DUE=%s PAID=%t%s\n",
				invoice.Id, invoice.CustomerId, invoice.Raised.Format(dataFormat),
				invoice.Due.Format(dataFormat), invoice.Paid, note); err != nil {
		return err
	}
	if err := write.writeItems(invoice.Items); err != nil {
		return err
	}
	return write("\f\n")
}

func (write writerFunc)writeItems(items []*Item) error {
	for _, item := range items {
		note := ""
		if item.Note != "" {
			note = noteSep + " " + item.Note
		}
		if err := write("ITEM ID=%s PRICE=%.2f QUANTITY=%d%s\n", item.Id, item.Price, item.Quantity, note); err != nil {
			return err
		}
	}
	return nil
}

func (TXTMarshaler)UnmarshalInvoices(reader io.Reader) ([]*Invoice, error) {
	bufferedReader := bufio.NewReader(reader)
	if err := checkTxtVersion(bufferedReader); err != nil {
		return nil, err
	}
	var invoices []*Invoice
	eof := false
	for lino := 2; !eof; lino++ {
		line, err := bufferedReader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return nil, err
		}
		if invoices, err = parseTxtLine(lino, line, invoices); err != nil {
			return nil, err
		}
	}
	return invoices, nil
}

func parseTxtLine(lino int, line string, invoices []*Invoice) ([]*Invoice, error) {
	var err error
	if strings.HasPrefix(line, "INVOICE") {
		var invoice *Invoice
		invoice, err = parseTxtInvoice(lino, line)
		invoices = append(invoices, invoice)
	} else if strings.HasPrefix(line, "ITEM") {
		if len(invoices) == 0 {
			err = fmt.Errorf("item outside of an invoice line %d", lino)
		} else {
			var item *Item
			item, err = parseTxtItem(lino, line)
			items := &invoices[len(invoices) - 1].Items
			*items = append(*items, item)
		}
	}
	return invoices, err
}

func parseTxtInvoice(lino int, line string) (invoice *Invoice, err error) {
	invoice = &Invoice{}
	var raised, due string
	if _, err = fmt.Sscanf(line, "INVOICE ID=%d CUSTOMER=%d RAISED=%s DUE=%s PAID=%t",
								&invoice.Id, &invoice.CustomerId, &raised, &due, &invoice.Paid); err != nil {
		return nil, fmt.Errorf("invalid invoice %v line %d", err, lino)
	}
	if invoice.Raised, err = time.Parse(dataFormat, raised); err != nil {
		return nil, fmt.Errorf("invalid raised %v line %d", err, lino)
	}
	if invoice.Due, err = time.Parse(dataFormat, due); err != nil {
		return nil, fmt.Errorf("invalid due %v line %d", err, lino)
	}
	if i := strings.Index(line, noteSep); i > -1 {
		invoice.Note = strings.TrimSpace(line[i + len(noteSep):])
	}
	return invoice, nil
}

func parseTxtItem(lino int, line string) (item *Item, err error) {
	item = &Item{}
	if _, err = fmt.Sscanf(line, "ITEM ID=%s PRICE=%f QUANTITY=%d",
							&item.Id, &item.Price, &item.Quantity); err != nil {
		return nil, fmt.Errorf("invalid item %v line %d", err, lino)
	}
	if i := strings.Index(line, noteSep); i > -1 {
		item.Note = strings.TrimSpace(line[len(noteSep) + i:])
	}
	return item, nil
}

func checkTxtVersion(bufferReader *bufio.Reader) error {
	var version int
	if _, err := fmt.Fscanf(bufferReader, "INVOICES %d\n", &version); err != nil {
		return errors.New("cannot read non-invoices text file")
	} else if version > fileVersion {
		return fmt.Errorf("version %d is too new to read", version)
	}
	return nil
}