package gocsv

import (
	"encoding/csv"
	"io"
)

func ReadWithHeader(f io.Reader) (<-chan map[string]string, <-chan error) {
	datach := make(chan map[string]string, 0)
	errch := make(chan error, 0)

	var headers []string
	header := true
	go read(f, func(record []string, err error) {
		defer func() {
			close(datach)
			close(errch)
		}()

		if err != nil {
			if err == io.EOF {
				return
			}
			errch <- err
			return
		}

		if header {
			headers = make([]string, len(record))
			copy(headers, record)
		}

		data := map[string]string{}

		for idx, h := range headers {
			data[h] = record[idx]
		}
		datach <- data
	})

	return datach, errch
}

func Read(f io.Reader) (<-chan []string, <-chan error) {
	datach := make(chan []string, 0)
	errch := make(chan error, 0)

	go read(f, func(record []string, err error) {
		defer func() {
			close(datach)
			close(errch)
		}()

		if err != nil {
			if err == io.EOF {
				return
			}
			errch <- err
			return
		}

		datach <- record
	})

	return datach, errch
}

// ReadAll is a simple wrapper around csv#Reader.ReadAll
func ReadAll(f io.Reader) ([][]string, error) {
	r := csv.NewReader(f)
	r.ReuseRecord = true
	r.TrimLeadingSpace = true
	return r.ReadAll()
}

func read(f io.Reader, fn func([]string, error)) {
	r := csv.NewReader(f)
	r.ReuseRecord = true
	r.TrimLeadingSpace = true

	for {
		record, err := r.Read()
		if err != nil {
			fn(nil, err)
			return
		}
		fn(record, nil)
	}
}
