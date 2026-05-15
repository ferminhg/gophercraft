package command_test

import (
	"context"

	domainevent "github.com/fermin/gophercraft/internal/domain/event"
)

type recordingPublisher struct {
	events []domainevent.DomainEvent
	err    error
}

func (p *recordingPublisher) Publish(_ context.Context, events ...domainevent.DomainEvent) error {
	if p.err != nil {
		return p.err
	}
	p.events = append(p.events, events...)
	return nil
}
