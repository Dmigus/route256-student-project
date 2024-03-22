// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: v1/loms.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on OrderCreateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *OrderCreateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrderCreateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// OrderCreateRequestMultiError, or nil if none found.
func (m *OrderCreateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *OrderCreateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for User

	if len(m.GetItems()) < 1 {
		err := OrderCreateRequestValidationError{
			field:  "Items",
			reason: "value must contain at least 1 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetItems() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, OrderCreateRequestValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, OrderCreateRequestValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OrderCreateRequestValidationError{
					field:  fmt.Sprintf("Items[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return OrderCreateRequestMultiError(errors)
	}

	return nil
}

// OrderCreateRequestMultiError is an error wrapping multiple validation errors
// returned by OrderCreateRequest.ValidateAll() if the designated constraints
// aren't met.
type OrderCreateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrderCreateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrderCreateRequestMultiError) AllErrors() []error { return m }

// OrderCreateRequestValidationError is the validation error returned by
// OrderCreateRequest.Validate if the designated constraints aren't met.
type OrderCreateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrderCreateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrderCreateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrderCreateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrderCreateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrderCreateRequestValidationError) ErrorName() string {
	return "OrderCreateRequestValidationError"
}

// Error satisfies the builtin error interface
func (e OrderCreateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrderCreateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrderCreateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrderCreateRequestValidationError{}

// Validate checks the field values on Item with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Item) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Item with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ItemMultiError, or nil if none found.
func (m *Item) ValidateAll() error {
	return m.validate(true)
}

func (m *Item) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Sku

	if m.GetCount() <= 0 {
		err := ItemValidationError{
			field:  "Count",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ItemMultiError(errors)
	}

	return nil
}

// ItemMultiError is an error wrapping multiple validation errors returned by
// Item.ValidateAll() if the designated constraints aren't met.
type ItemMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ItemMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ItemMultiError) AllErrors() []error { return m }

// ItemValidationError is the validation error returned by Item.Validate if the
// designated constraints aren't met.
type ItemValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ItemValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ItemValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ItemValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ItemValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ItemValidationError) ErrorName() string { return "ItemValidationError" }

// Error satisfies the builtin error interface
func (e ItemValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sItem.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ItemValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ItemValidationError{}

// Validate checks the field values on OrderId with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *OrderId) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrderId with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in OrderIdMultiError, or nil if none found.
func (m *OrderId) ValidateAll() error {
	return m.validate(true)
}

func (m *OrderId) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if len(errors) > 0 {
		return OrderIdMultiError(errors)
	}

	return nil
}

// OrderIdMultiError is an error wrapping multiple validation errors returned
// by OrderId.ValidateAll() if the designated constraints aren't met.
type OrderIdMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrderIdMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrderIdMultiError) AllErrors() []error { return m }

// OrderIdValidationError is the validation error returned by OrderId.Validate
// if the designated constraints aren't met.
type OrderIdValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrderIdValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrderIdValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrderIdValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrderIdValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrderIdValidationError) ErrorName() string { return "OrderIdValidationError" }

// Error satisfies the builtin error interface
func (e OrderIdValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrderId.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrderIdValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrderIdValidationError{}

// Validate checks the field values on OrderInfoResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *OrderInfoResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrderInfoResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// OrderInfoResponseMultiError, or nil if none found.
func (m *OrderInfoResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *OrderInfoResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Status

	// no validation rules for User

	for idx, item := range m.GetItems() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, OrderInfoResponseValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, OrderInfoResponseValidationError{
						field:  fmt.Sprintf("Items[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OrderInfoResponseValidationError{
					field:  fmt.Sprintf("Items[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return OrderInfoResponseMultiError(errors)
	}

	return nil
}

// OrderInfoResponseMultiError is an error wrapping multiple validation errors
// returned by OrderInfoResponse.ValidateAll() if the designated constraints
// aren't met.
type OrderInfoResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrderInfoResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrderInfoResponseMultiError) AllErrors() []error { return m }

// OrderInfoResponseValidationError is the validation error returned by
// OrderInfoResponse.Validate if the designated constraints aren't met.
type OrderInfoResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrderInfoResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrderInfoResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrderInfoResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrderInfoResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrderInfoResponseValidationError) ErrorName() string {
	return "OrderInfoResponseValidationError"
}

// Error satisfies the builtin error interface
func (e OrderInfoResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrderInfoResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrderInfoResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrderInfoResponseValidationError{}

// Validate checks the field values on StocksInfoRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *StocksInfoRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StocksInfoRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StocksInfoRequestMultiError, or nil if none found.
func (m *StocksInfoRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *StocksInfoRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Sku

	if len(errors) > 0 {
		return StocksInfoRequestMultiError(errors)
	}

	return nil
}

// StocksInfoRequestMultiError is an error wrapping multiple validation errors
// returned by StocksInfoRequest.ValidateAll() if the designated constraints
// aren't met.
type StocksInfoRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StocksInfoRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StocksInfoRequestMultiError) AllErrors() []error { return m }

// StocksInfoRequestValidationError is the validation error returned by
// StocksInfoRequest.Validate if the designated constraints aren't met.
type StocksInfoRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StocksInfoRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StocksInfoRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StocksInfoRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StocksInfoRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StocksInfoRequestValidationError) ErrorName() string {
	return "StocksInfoRequestValidationError"
}

// Error satisfies the builtin error interface
func (e StocksInfoRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStocksInfoRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StocksInfoRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StocksInfoRequestValidationError{}

// Validate checks the field values on StocksInfoResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *StocksInfoResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StocksInfoResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StocksInfoResponseMultiError, or nil if none found.
func (m *StocksInfoResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *StocksInfoResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Count

	if len(errors) > 0 {
		return StocksInfoResponseMultiError(errors)
	}

	return nil
}

// StocksInfoResponseMultiError is an error wrapping multiple validation errors
// returned by StocksInfoResponse.ValidateAll() if the designated constraints
// aren't met.
type StocksInfoResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StocksInfoResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StocksInfoResponseMultiError) AllErrors() []error { return m }

// StocksInfoResponseValidationError is the validation error returned by
// StocksInfoResponse.Validate if the designated constraints aren't met.
type StocksInfoResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StocksInfoResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StocksInfoResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StocksInfoResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StocksInfoResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StocksInfoResponseValidationError) ErrorName() string {
	return "StocksInfoResponseValidationError"
}

// Error satisfies the builtin error interface
func (e StocksInfoResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStocksInfoResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StocksInfoResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StocksInfoResponseValidationError{}
