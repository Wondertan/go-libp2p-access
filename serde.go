package access

import (
	"io"

	"github.com/Wondertan/go-serde"
)

func WriteToken(w io.Writer, t Token) (int, error) {
	return serde.Write(w, (*msg)(&t))
}

func ReadToken(r io.Reader) (Token, int, error) {
	var m msg
	n, err := serde.Read(r, &m)
	return Token(m), n, err
}

type msg Token

func (m msg) Size() int {
	return len(m)
}

func (m msg) MarshalTo(data []byte) (int, error) {
	return copy(data, m), nil
}

func (m *msg) Unmarshal(data []byte) error {
	b := make([]byte, len(data))
	copy(b, data)
	*m = msg(b)
	return nil
}
