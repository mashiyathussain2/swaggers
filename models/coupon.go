// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Coupon coupon
//
// swagger:model Coupon
type Coupon struct {

	// code
	Code string `json:"code,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// status
	Status string `json:"status,omitempty"`

	// type
	Type string `json:"type,omitempty"`

	// valid after
	// Format: date-time
	ValidAfter strfmt.DateTime `json:"valid_after,omitempty"`

	// valid before
	// Format: date-time
	ValidBefore strfmt.DateTime `json:"valid_before,omitempty"`

	// value
	Value int64 `json:"value,omitempty"`

	// applicable on
	ApplicableOn *ApplicableON `json:"applicable_on,omitempty"`

	// id
	ID ObjectID `json:"id,omitempty"`

	// max discount
	MaxDiscount *Price `json:"max_discount,omitempty"`

	// min purchase value
	MinPurchaseValue *Price `json:"min_purchase_value,omitempty"`
}

// Validate validates this coupon
func (m *Coupon) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateValidAfter(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateValidBefore(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateApplicableOn(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMaxDiscount(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMinPurchaseValue(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Coupon) validateValidAfter(formats strfmt.Registry) error {
	if swag.IsZero(m.ValidAfter) { // not required
		return nil
	}

	if err := validate.FormatOf("valid_after", "body", "date-time", m.ValidAfter.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Coupon) validateValidBefore(formats strfmt.Registry) error {
	if swag.IsZero(m.ValidBefore) { // not required
		return nil
	}

	if err := validate.FormatOf("valid_before", "body", "date-time", m.ValidBefore.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Coupon) validateApplicableOn(formats strfmt.Registry) error {
	if swag.IsZero(m.ApplicableOn) { // not required
		return nil
	}

	if m.ApplicableOn != nil {
		if err := m.ApplicableOn.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("applicable_on")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("applicable_on")
			}
			return err
		}
	}

	return nil
}

func (m *Coupon) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := m.ID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("id")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("id")
		}
		return err
	}

	return nil
}

func (m *Coupon) validateMaxDiscount(formats strfmt.Registry) error {
	if swag.IsZero(m.MaxDiscount) { // not required
		return nil
	}

	if m.MaxDiscount != nil {
		if err := m.MaxDiscount.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("max_discount")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("max_discount")
			}
			return err
		}
	}

	return nil
}

func (m *Coupon) validateMinPurchaseValue(formats strfmt.Registry) error {
	if swag.IsZero(m.MinPurchaseValue) { // not required
		return nil
	}

	if m.MinPurchaseValue != nil {
		if err := m.MinPurchaseValue.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("min_purchase_value")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("min_purchase_value")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this coupon based on the context it is used
func (m *Coupon) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateApplicableOn(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateMaxDiscount(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateMinPurchaseValue(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Coupon) contextValidateApplicableOn(ctx context.Context, formats strfmt.Registry) error {

	if m.ApplicableOn != nil {
		if err := m.ApplicableOn.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("applicable_on")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("applicable_on")
			}
			return err
		}
	}

	return nil
}

func (m *Coupon) contextValidateID(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ID.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("id")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("id")
		}
		return err
	}

	return nil
}

func (m *Coupon) contextValidateMaxDiscount(ctx context.Context, formats strfmt.Registry) error {

	if m.MaxDiscount != nil {
		if err := m.MaxDiscount.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("max_discount")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("max_discount")
			}
			return err
		}
	}

	return nil
}

func (m *Coupon) contextValidateMinPurchaseValue(ctx context.Context, formats strfmt.Registry) error {

	if m.MinPurchaseValue != nil {
		if err := m.MinPurchaseValue.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("min_purchase_value")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("min_purchase_value")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Coupon) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Coupon) UnmarshalBinary(b []byte) error {
	var res Coupon
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
