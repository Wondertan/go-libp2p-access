package access

import (
	"context"
	"errors"
	"io"

	"github.com/libp2p/go-libp2p-core/peer"
)

var ErrHandshake = errors.New("access: handshake failed")

var errToken Token = "NG"

func GiveHand(ctx context.Context, s io.ReadWriter) error {
	in, err := GetToken(ctx)
	if err != nil {
		return err
	}

	_, err = WriteToken(s, in)
	if err != nil {
		return err
	}

	out, _, err := ReadToken(s)
	if err != nil {
		return err
	}

	if in != out {
		if out == errToken {
			return ErrNotGranted
		}

		return ErrHandshake
	}

	return nil
}

func TakeHand(g Granter, s io.ReadWriter, peer peer.ID) (chan<- error, error) {
	t, _, err := ReadToken(s)
	if err != nil {
		return nil, err
	}

	errs, err := g.Granted(t, peer)
	if err != nil {
		if err == ErrNotGranted {
			_, err = WriteToken(s, errToken)
			if err != nil {
				return nil, err
			}

			return nil, ErrNotGranted
		}

		return nil, err
	}

	_, err = WriteToken(s, t)
	if err != nil {
		return nil, err
	}

	return errs, nil
}
