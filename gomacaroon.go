// Copyright 2015 Canonical Ltd.
// Licensed under the LGPL, see LICENCE file for details.

package macarooncompat

import (
	"gopkg.in/macaroon.v1"
)

type goMacaroon struct {
	*macaroon.Macaroon
}

func (m goMacaroon) clone() goMacaroon {
	return goMacaroon{m.Macaroon.Clone()}
}

func (m goMacaroon) WithFirstPartyCaveat(caveatId string) (Macaroon, error) {
	m = m.clone()
	if err := m.Macaroon.AddFirstPartyCaveat(caveatId); err != nil {
		return nil, err
	}
	return m, nil
}

func (m goMacaroon) WithThirdPartyCaveat(rootKey []byte, caveatId string, loc string) (Macaroon, error) {
	m = m.clone()
	if err := m.Macaroon.AddThirdPartyCaveat(rootKey, caveatId, loc); err != nil {
		return nil, err
	}
	return m, nil
}

func (m goMacaroon) Bind(primary Macaroon) (Macaroon, error) {
	m = m.clone()
	m.Macaroon.Bind(primary.Signature())
	return m, nil
}

func (m goMacaroon) Verify(rootKey []byte, check Checker, discharges []Macaroon) error {
	discharges1 := make([]*macaroon.Macaroon, len(discharges))
	for i, m := range discharges {
		discharges1[i] = m.(goMacaroon).Macaroon
	}
	return m.Macaroon.Verify(rootKey, check.Check, discharges1)
}

type goMacaroonPackage struct{}

func (goMacaroonPackage) New(rootKey []byte, id, loc string) (Macaroon, error) {
	m, err := macaroon.New(rootKey, id, loc)
	if err != nil {
		return nil, err
	}
	return goMacaroon{m}, nil
}

func (goMacaroonPackage) UnmarshalJSON(data []byte) (Macaroon, error) {
	var m macaroon.Macaroon
	if err := m.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return goMacaroon{&m}, nil
}

func (goMacaroonPackage) UnmarshalBinary(data []byte) (Macaroon, error) {
	var m macaroon.Macaroon
	if err := m.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	return goMacaroon{&m}, nil
}
