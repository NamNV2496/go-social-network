package validator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode/utf8"

	vietnamese "github.com/go-playground/locales/vi"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/vi"
)

var (
	fieldnameTranslation = map[string]string{
		"user_id":      "ID người dùng",
		"content_text": "Nội dung bài post",
		"images":       "Hình ảnh",
	}
	validationOpts = map[string]validationOption{
		"v_user_id":      {"omitempty,min=2,max=60,printascii", "Username phải là chữ thường và có ít hơn 60 kí tự, vui lòng chỉnh sửa lại."},
		"v_content_text": {"content_text", "Nội dung là bắt buộc, có độ dài phù hợp, và bắt đầu bằng chữ cái hợp lệ."},
	}

	regex = make(map[string]*regexp.Regexp)
)

type IPostValidator interface {
	Validate(ctx context.Context, in any) error
}

type PostValidator struct {
	va    *validator.Validate
	trans ut.Translator
}

type validationOption struct {
	validation  string
	translation string
}

func matchPattern(name string, value string) bool {
	return regex[name].MatchString(value)
}

func validateUserId(ctx context.Context, fl validator.FieldLevel) bool {
	return matchPattern("p_user_id", strings.ToLower(fl.Field().String()))
}

func validateContentCharacters(ctx context.Context, fl validator.FieldLevel) bool {
	fieldVal := fl.Field().String()
	fieldLen := utf8.RuneCountInString(fieldVal)
	return fieldLen > 1 && fieldLen < 2001 && matchPattern("p_content_text", strings.ToLower(fieldVal)) && !matchPattern("p_html", strings.ToLower(fieldVal))
}

func validateLengthSlices(ctx context.Context, fl validator.FieldLevel) bool {
	fval := fl.Field().String()
	if fval == "" {
		return true
	}
	if len(fval) > 5 {
		return false
	}
	return true
}

func init() {
	patterns := map[string]string{
		"p_user_id":      `^[\x20-\x{23CD}\x{2E80}-\x{4DFF}\x{4E00}-\x{9FFF}\x{F900}-\x{FAFF}\x{FF00}-\x{FF5A}\x{1BCA0}-\x{2A6DF}\x{2A700}-\x{2B73F}\x{2F800}-\x{2FA1F}]+`,
		"p_content_text": `^[\x20-\x{2FA1F}]+`,
		"p_html":         `<.*?>`,
	}

	for name, pattern := range patterns {
		regex[name] = regexp.MustCompile(pattern)
	}
}

func NewPostValidator() (*PostValidator, error) {
	va := validator.New()
	vie := vietnamese.New()
	uni := ut.New(vie, vie)
	trans, found := uni.GetTranslator("vi")
	if !found {
		return nil, fmt.Errorf("cannot initialize vietnamese translator")
	}
	if err := vi.RegisterDefaultTranslations(va, trans); err != nil {
		return nil, fmt.Errorf("cannot register vietnamese translator, err: %w", err)
	}
	instance := &PostValidator{
		va:    va,
		trans: trans,
	}
	var customValidateFuncs = map[string]validator.FuncCtx{
		"user_id":       validateUserId,
		"content_text":  validateContentCharacters,
		"length_slices": validateLengthSlices,
	}
	for k, v := range customValidateFuncs {
		va.RegisterValidationCtx(k, v)
	}
	for alias, opt := range validationOpts {
		validation := opt.validation
		translation := opt.translation
		va.RegisterAlias(alias, validation)
		if err := va.RegisterTranslation(alias, trans, registrationFunc(alias, translation, false), translateFunc); err != nil {
			return nil, fmt.Errorf("error register translator: %w", err)
		}
	}
	return instance, nil
}

func (v *PostValidator) Validate(ctx context.Context, in any) error {
	if in == nil {
		return errors.New("nil input")
	}
	if err := v.va.StructCtx(ctx, in); err != nil {
		ve, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
		for _, fe := range ve {
			return errors.New(fe.Translate(v.trans))
		}
	}
	return nil
}
func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}
		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	fieldName := fe.Field()
	if tf, found := fieldnameTranslation[fieldName]; found {
		fieldName = tf
	}
	t, err := ut.T(fe.Tag(), fieldName)
	if err != nil {
		log.Printf("cảnh báo: lỗi chuyển ngữ FieldError: %#v", fe)
		return fe.(error).Error()
	}
	return t
}
