package templates

const SERVICES_IMinioServiceTemplate = `
/* Generated by egg v0.0.1 */

package services

import (
	"io"

	"github.com/google/uuid"
)

type IMinioService interface {
	Upload(uploaderID uuid.UUID, uploadName string, uploadFile io.Reader, size int64) error
	Get(uploaderID uuid.UUID, uploadName string) ([]byte, error)
	GetPresigned(uploaderID uuid.UUID, uploadName string) (string, error)
}
`
