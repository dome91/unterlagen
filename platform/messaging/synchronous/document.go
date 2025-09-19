package synchronous

import (
	"log/slog"
	"unterlagen/features/archive"
)

var _ archive.DocumentMessages = &DocumentMessages{}

type DocumentMessages struct {
	documentTextExtractedSubscribers []func(document archive.Document) error
	documentUploadedSubscribers      []func(document archive.Document) error
	documentDeletedSubscribers       []func(document archive.Document) error
}

func (d *DocumentMessages) PublishDocumentTextExtracted(document archive.Document) error {
	for _, subscriber := range d.documentTextExtractedSubscribers {
		err := subscriber(document)
		if err != nil {
			slog.Error("failed to process document text extracted event", slog.String("error", err.Error()))
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

func (d *DocumentMessages) PublishDocumentUpserted(document archive.Document) error {
	for _, subscriber := range d.documentUploadedSubscribers {
		err := subscriber(document)
		if err != nil {
			slog.Error("failed to process document uploaded event", slog.String("error", err.Error()))
		}
	}
	return nil
}

func (d *DocumentMessages) SubscribeDocumentTextExtracted(subscriber func(document archive.Document) error) error {
	d.documentTextExtractedSubscribers = append(d.documentTextExtractedSubscribers, subscriber)
	return nil
}

func (d *DocumentMessages) SubscribeDocumentDeleted(subscriber func(document archive.Document) error) error {
	d.documentDeletedSubscribers = append(d.documentDeletedSubscribers, subscriber)
	return nil
}

func (d *DocumentMessages) SubscribeDocumentUpserted(subscriber func(document archive.Document) error) error {
	d.documentUploadedSubscribers = append(d.documentUploadedSubscribers, subscriber)
	return nil
}

func NewDocumentMessages() *DocumentMessages {
	return &DocumentMessages{
		documentTextExtractedSubscribers: []func(document archive.Document) error{},
		documentUploadedSubscribers:      []func(document archive.Document) error{},
		documentDeletedSubscribers:       []func(document archive.Document) error{},
	}
}
