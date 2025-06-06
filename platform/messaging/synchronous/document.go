package synchronous

import (
	"log/slog"
	"unterlagen/features/archive"
)

var _ archive.DocumentMessages = &DocumentMessages{}

type DocumentMessages struct {
	documentAnalyzedSubscribers []func(document archive.Document) error
	documentUploadedSubscribers []func(document archive.Document) error
	documentDeletedSubscribers  []func(document archive.Document) error
}

func (d *DocumentMessages) PublishDocumentAnalyzed(document archive.Document) error {
	for _, subscriber := range d.documentAnalyzedSubscribers {
		err := subscriber(document)
		if err != nil {
			slog.Error("failed to process document analyzed event", slog.String("error", err.Error()))
		}
	}
	return nil
}

func (d *DocumentMessages) PublishDocumentDeleted(document archive.Document) error {
	for _, subscriber := range d.documentDeletedSubscribers {
		err := subscriber(document)
		if err != nil {
			slog.Error("failed to process document deleted event", slog.String("error", err.Error()))
		}
	}
	return nil
}

func (d *DocumentMessages) PublishDocumentUploaded(document archive.Document) error {
	for _, subscriber := range d.documentUploadedSubscribers {
		err := subscriber(document)
		if err != nil {
			slog.Error("failed to process document uploaded event", slog.String("error", err.Error()))
		}
	}
	return nil
}

func (d *DocumentMessages) SubscribeDocumentAnalyzed(subscriber func(document archive.Document) error) error {
	d.documentAnalyzedSubscribers = append(d.documentAnalyzedSubscribers, subscriber)
	return nil
}

func (d *DocumentMessages) SubscribeDocumentDeleted(subscriber func(document archive.Document) error) error {
	d.documentDeletedSubscribers = append(d.documentDeletedSubscribers, subscriber)
	return nil
}

func (d *DocumentMessages) SubscribeDocumentUploaded(subscriber func(document archive.Document) error) error {
	d.documentUploadedSubscribers = append(d.documentUploadedSubscribers, subscriber)
	return nil
}

func NewDocumentMessages() *DocumentMessages {
	return &DocumentMessages{
		documentAnalyzedSubscribers: []func(document archive.Document) error{},
		documentUploadedSubscribers: []func(document archive.Document) error{},
		documentDeletedSubscribers:  []func(document archive.Document) error{},
	}
}
