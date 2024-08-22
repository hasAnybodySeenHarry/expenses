package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"harry2an.com/expenses/cmd/proto/users"
	"harry2an.com/expenses/internal/data"
	"harry2an.com/expenses/internal/validator"
)

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if input.Password == "" {
		app.badRequest(w, r, errors.New("password cannot be blank"))
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Create(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Validate() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.badRequest(w, r, errors.New("user with this email already exists in the system"))
		default:
			app.serverError(w, r, err)
		}
		return
	}

	token, err := app.models.Tokens.New(user.ID, data.ScopeActivation, 15*time.Minute)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.background(func() {
		err := app.mailer.Send("user-welcome", user.Name, user.Email, token.Plaintext)
		if err != nil {
			app.logger.Println(fmt.Errorf("%s", err))
		}
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{
		"user": user,
	}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

type userServiceServer struct {
	users.UnimplementedUserServiceServer
	models *data.Models
}

func (s *userServiceServer) GetUserForToken(ctx context.Context, req *users.GetUserRequest) (*users.GetUserResponse, error) {
	v := validator.New()
	if data.ValidateToken(v, req.Token); !v.Validate() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid token: failed to validate")
	}

	user, err := s.models.Users.GetForToken(req.Token, data.ScopeAuthentication)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			return nil, status.Errorf(codes.Unauthenticated, "invalid authentication token")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
		}
	}

	return data.UserToProto(user), nil
}
