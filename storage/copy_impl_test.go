package storage

import (
	"context"
	"io"
	"testing"

	"github.com/lyft/flytestdlib/ioutils"
	"github.com/lyft/flytestdlib/promutils"
	"github.com/stretchr/testify/assert"
)

type notSeekerReader struct {
	bytesCount int
}

func (notSeekerReader) Close() error {
	return nil
}

func (r *notSeekerReader) Read(p []byte) (n int, err error) {
	if len(p) < 1 {
		return 0, nil
	}

	p[0] = byte(10)

	r.bytesCount--
	if r.bytesCount <= 0 {
		return 0, io.EOF
	}

	return 1, nil
}

func newNotSeekerReader(bytesCount int) *notSeekerReader {
	return &notSeekerReader{
		bytesCount: bytesCount,
	}
}

func TestCopyRaw(t *testing.T) {
	t.Run("Called", func(t *testing.T) {
		readerCalled := false
		writerCalled := false
		store := dummyStore{
			ReadRawCb: func(ctx context.Context, reference DataReference) (closer io.ReadCloser, e error) {
				readerCalled = true
				return ioutils.NewBytesReadCloser([]byte{}), nil
			},
			WriteRawCb: func(ctx context.Context, reference DataReference, size int64, opts Options, raw io.Reader) error {
				writerCalled = true
				return nil
			},
		}

		copier := newCopyImpl(&store, promutils.NewTestScope())
		assert.NoError(t, copier.CopyRaw(context.Background(), DataReference("source.pb"), DataReference("dest.pb"), Options{}))
		assert.True(t, readerCalled)
		assert.True(t, writerCalled)
	})

	t.Run("Not Seeker", func(t *testing.T) {
		readerCalled := false
		writerCalled := false
		store := dummyStore{
			ReadRawCb: func(ctx context.Context, reference DataReference) (closer io.ReadCloser, e error) {
				readerCalled = true
				return newNotSeekerReader(10), nil
			},
			WriteRawCb: func(ctx context.Context, reference DataReference, size int64, opts Options, raw io.Reader) error {
				writerCalled = true
				return nil
			},
		}

		copier := newCopyImpl(&store, promutils.NewTestScope())
		assert.NoError(t, copier.CopyRaw(context.Background(), DataReference("source.pb"), DataReference("dest.pb"), Options{}))
		assert.True(t, readerCalled)
		assert.True(t, writerCalled)
	})
}
