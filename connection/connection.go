package connection

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Connection interface {
	Establish() error
	Terminate() error
}

func MustUngracefullyEstablish(connections ...Connection) {
	if err := UngracefullyEstablish(connections...); err != nil {
		panic(err)
	}
}

func UngracefullyEstablish(connections ...Connection) error {
	var group errgroup.Group

	for _, connection := range connections {
		group.Go(establishConnectionFunc(connection))
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("establish connections: %w", err)
	}

	return nil
}

func establishConnectionFunc(connection Connection) func() error {
	return func() error {
		if err := connection.Establish(); err != nil {
			return fmt.Errorf("establish: %w", err)
		}

		return nil
	}
}
