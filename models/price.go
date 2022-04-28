// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Price Price represents cost of an entity
//
// swagger:model Price
type Price struct {

	// value
	Value float32 `json:"value,omitempty"`

	// iso
	Iso CurrencyISO `json:"iso,omitempty"`
}

// Validate validates this price
func (m *Price) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateIso(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Price) validateIso(formats strfmt.Registry) error {
	if swag.IsZero(m.Iso) { // not required
		return nil
	}

	if err := m.Iso.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("iso")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("iso")
		}
		return err
	}

	return nil
}

// ContextValidate validate this price based on the context it is used
func (m *Price) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateIso(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Price) contextValidateIso(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Iso.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("iso")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("iso")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Price) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Price) UnmarshalBinary(b []byte) error {
	var res Price
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
