// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: newsfeed_core/v1/newsfeed.proto

package newsfeedv1

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

// Validate checks the field values on NewsfeedPost with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *NewsfeedPost) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on NewsfeedPost with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in NewsfeedPostMultiError, or
// nil if none found.
func (m *NewsfeedPost) ValidateAll() error {
	return m.validate(true)
}

func (m *NewsfeedPost) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	// no validation rules for ContentText

	// no validation rules for Visible

	// no validation rules for Date

	if len(errors) > 0 {
		return NewsfeedPostMultiError(errors)
	}

	return nil
}

// NewsfeedPostMultiError is an error wrapping multiple validation errors
// returned by NewsfeedPost.ValidateAll() if the designated constraints aren't met.
type NewsfeedPostMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m NewsfeedPostMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m NewsfeedPostMultiError) AllErrors() []error { return m }

// NewsfeedPostValidationError is the validation error returned by
// NewsfeedPost.Validate if the designated constraints aren't met.
type NewsfeedPostValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e NewsfeedPostValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e NewsfeedPostValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e NewsfeedPostValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e NewsfeedPostValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e NewsfeedPostValidationError) ErrorName() string { return "NewsfeedPostValidationError" }

// Error satisfies the builtin error interface
func (e NewsfeedPostValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sNewsfeedPost.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = NewsfeedPostValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = NewsfeedPostValidationError{}

// Validate checks the field values on GetNewsfeedRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetNewsfeedRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetNewsfeedRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetNewsfeedRequestMultiError, or nil if none found.
func (m *GetNewsfeedRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetNewsfeedRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	if len(errors) > 0 {
		return GetNewsfeedRequestMultiError(errors)
	}

	return nil
}

// GetNewsfeedRequestMultiError is an error wrapping multiple validation errors
// returned by GetNewsfeedRequest.ValidateAll() if the designated constraints
// aren't met.
type GetNewsfeedRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetNewsfeedRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetNewsfeedRequestMultiError) AllErrors() []error { return m }

// GetNewsfeedRequestValidationError is the validation error returned by
// GetNewsfeedRequest.Validate if the designated constraints aren't met.
type GetNewsfeedRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetNewsfeedRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetNewsfeedRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetNewsfeedRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetNewsfeedRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetNewsfeedRequestValidationError) ErrorName() string {
	return "GetNewsfeedRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetNewsfeedRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetNewsfeedRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetNewsfeedRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetNewsfeedRequestValidationError{}

// Validate checks the field values on GetNewsfeedResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetNewsfeedResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetNewsfeedResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetNewsfeedResponseMultiError, or nil if none found.
func (m *GetNewsfeedResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetNewsfeedResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetPosts() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetNewsfeedResponseValidationError{
						field:  fmt.Sprintf("Posts[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetNewsfeedResponseValidationError{
						field:  fmt.Sprintf("Posts[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetNewsfeedResponseValidationError{
					field:  fmt.Sprintf("Posts[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetNewsfeedResponseMultiError(errors)
	}

	return nil
}

// GetNewsfeedResponseMultiError is an error wrapping multiple validation
// errors returned by GetNewsfeedResponse.ValidateAll() if the designated
// constraints aren't met.
type GetNewsfeedResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetNewsfeedResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetNewsfeedResponseMultiError) AllErrors() []error { return m }

// GetNewsfeedResponseValidationError is the validation error returned by
// GetNewsfeedResponse.Validate if the designated constraints aren't met.
type GetNewsfeedResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetNewsfeedResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetNewsfeedResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetNewsfeedResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetNewsfeedResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetNewsfeedResponseValidationError) ErrorName() string {
	return "GetNewsfeedResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetNewsfeedResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetNewsfeedResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetNewsfeedResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetNewsfeedResponseValidationError{}