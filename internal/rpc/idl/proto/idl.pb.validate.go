// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: idl.proto

package proto

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

// Validate checks the field values on HelloRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *HelloRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on HelloRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in HelloRequestMultiError, or
// nil if none found.
func (m *HelloRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *HelloRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetName()); l < 5 || l > 10 {
		err := HelloRequestValidationError{
			field:  "Name",
			reason: "value length must be between 5 and 10 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return HelloRequestMultiError(errors)
	}

	return nil
}

// HelloRequestMultiError is an error wrapping multiple validation errors
// returned by HelloRequest.ValidateAll() if the designated constraints aren't met.
type HelloRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m HelloRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m HelloRequestMultiError) AllErrors() []error { return m }

// HelloRequestValidationError is the validation error returned by
// HelloRequest.Validate if the designated constraints aren't met.
type HelloRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e HelloRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e HelloRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e HelloRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e HelloRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e HelloRequestValidationError) ErrorName() string { return "HelloRequestValidationError" }

// Error satisfies the builtin error interface
func (e HelloRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHelloRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = HelloRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = HelloRequestValidationError{}

// Validate checks the field values on HelloReply with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *HelloReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on HelloReply with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in HelloReplyMultiError, or
// nil if none found.
func (m *HelloReply) ValidateAll() error {
	return m.validate(true)
}

func (m *HelloReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetMessage()); l < 5 || l > 10 {
		err := HelloReplyValidationError{
			field:  "Message",
			reason: "value length must be between 5 and 10 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return HelloReplyMultiError(errors)
	}

	return nil
}

// HelloReplyMultiError is an error wrapping multiple validation errors
// returned by HelloReply.ValidateAll() if the designated constraints aren't met.
type HelloReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m HelloReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m HelloReplyMultiError) AllErrors() []error { return m }

// HelloReplyValidationError is the validation error returned by
// HelloReply.Validate if the designated constraints aren't met.
type HelloReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e HelloReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e HelloReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e HelloReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e HelloReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e HelloReplyValidationError) ErrorName() string { return "HelloReplyValidationError" }

// Error satisfies the builtin error interface
func (e HelloReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHelloReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = HelloReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = HelloReplyValidationError{}
