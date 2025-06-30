package validation

import (
	"context"

	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	genericvalidation "github.com/ArthurWang23/miniblog/pkg/validation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (v *Validator) ValidatePostRules() genericvalidation.Rules {
	return genericvalidation.Rules{
		"PostID": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("postID cannot be empty")
			}
			return nil
		},
		"Title": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("title cannot be empty")
			}
			return nil
		},
		"Content": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("context cannot be empty")
			}
			return nil
		},
	}
}

func (v *Validator) ValidateCreatePostRequest(ctx context.Context, rq *apiv1.CreatePostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

func (v *Validator) ValidateUpdatePostRequest(ctx context.Context, rq *apiv1.UpdatePostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

func (v *Validator) ValidateDeletePostRequest(ctx context.Context, rq *apiv1.DeletePostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

func (v *Validator) ValidateGetPostRequest(ctx context.Context, rq *apiv1.GetPostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

func (v *Validator) ValidateListPostRequest(ctx context.Context, rq *apiv1.ListPostRequest) error {
	if err := validation.Validate(rq.GetTitle(), validation.Length(5, 100), is.URL); err != nil {
		return errno.ErrInvalidArgument.WithMessage(err.Error())
	}
	return genericvalidation.ValidateSelectedFields(rq, v.ValidatePostRules(), "Offset", "Limit")

	// if rq.Title != nil && len(rq.Title) > 200 {
	// 	return errno.ErrInvalidArgument.WithMessage("title cannot be longer than 200 characters")
	// }
	// return genericvalidation.ValidateSelectedFields(rq,v.ValidatePostRules(),"Offset","Limit")
}
