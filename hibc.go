package gohibc

import (
	"errors"
)

var (
	ErrNoSupplierLabelingFlag      = errors.New("missing HIBC Supplier flag")
	ErrInvalidSupplierLabelingFlag = errors.New("invalid HIBC Supplier flag")
	ErrMissingLIC                  = errors.New("missing LIC")
	ErrInvalidLIC                  = errors.New("invalid LIC")
	ErrInvalidPCN                  = errors.New("invalid PCN")
	ErrInvalidUnitOfMeasure        = errors.New("invalid Unit Of Measure")
)

var (
	supplierFlag = '+'
)

type HIBC struct {
	supplierFlag   rune
	di             deviceIdentifier
	pi             *productIdentifier
	checkCharacter rune
}

type deviceIdentifier struct {
	lic []rune
	pcn []rune
	um  rune
}

type productIdentifier struct {
}

func Parse(s string) (*HIBC, error) {
	code := new(HIBC)
	err := code.parse([]rune(s))
	if err != nil {
		return nil, err
	}
	return code, nil
}

func (hibc *HIBC) parse(rs []rune) error {
	return hibc.parseSupplierFlag(rs)
}

func (hibc *HIBC) parseSupplierFlag(rs []rune) error {
	if len(rs) <= 1 {
		return ErrNoSupplierLabelingFlag
	}
	if rs[0] != supplierFlag {
		return ErrInvalidSupplierLabelingFlag
	}
	hibc.supplierFlag = rs[0]

	if rs[1] == '$' {
		return nil
	}
	return hibc.parsePrimaryCode(rs[1:])
}

func (hibc *HIBC) parsePrimaryCode(rs []rune) error {
	hibc.di = deviceIdentifier{}
	return hibc.parseLIC(rs)
}

func (hibc *HIBC) parseLIC(rs []rune) error {
	if len(rs) == 0 || len(rs) < 4 {
		return ErrMissingLIC
	} else if alphabetic(rs[0]) {
		for i := 1; i <= 3; i++ {
			if !alphanumeric(rs[i]) {
				return ErrInvalidLIC
			}
		}
		hibc.di.lic = rs[0:4]
	} else {
		return ErrInvalidLIC
	}
	return hibc.parsePCN(rs[4:])
}

func (hibc *HIBC) parsePCN(rs []rune) error {
	if len(rs) < 3 {
		// minimum(3): PCN 1, UM: 1, CheckCharacter 1
		return ErrInvalidPCN
	}
	hibc.di.pcn = []rune{}
	for i := 0; i+2 < len(rs) && (rs[i+1] != '/'); i++ {
		// while we don't get to the last 2 characters (UM and CheckCharacter)
		if !alphanumeric(rs[i]) {
			return ErrInvalidPCN
		}
		hibc.di.pcn = append(hibc.di.pcn, rs[i])
	}
	return hibc.parseUM(rs[len(hibc.di.pcn):])
}

func (hibc *HIBC) parseUM(rs []rune) error {
	if !numeric(rs[0]) {
		return ErrInvalidUnitOfMeasure
	}
	hibc.di.um = rs[0]
	return nil
}

func alphabetic(r rune) bool {
	return 'A' <= r && r <= 'Z'
}

func numeric(r rune) bool {
	return '0' <= r && r <= '9'
}

func alphanumeric(r rune) bool {
	return alphabetic(r) || numeric(r)
}
