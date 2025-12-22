package eventsService

import (
	"errors"
	"strings"
)

var (
	ErrUnauthorized   = errors.New("Пользователь должен быть авторизован")
	ErrEmptyComment   = errors.New("Комментарий не может быть пустым")
	ErrCommentTooLong = errors.New("Комментарий слишком длинный")
)

func (s *Service) validateUser(userHash string) error {
	if strings.TrimSpace(userHash) == "" {
		return ErrUnauthorized
	}
	return nil
}

func (s *Service) validateComment(userHash, text string) error {
	if err := s.validateUser(userHash); err != nil {
		return err
	}

	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return ErrEmptyComment
	}
	if len(trimmed) > 1000 {
		return ErrCommentTooLong
	}

	return nil
}
