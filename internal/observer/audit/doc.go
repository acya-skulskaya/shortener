// Package audit provides an auditing interface to track URL shortening operations.
// Audit events can be submitted to a file or to a remote endpoint
//
// # To submit audit events to a file, a FileAuditSubscriber needs to be created with a path to file
//
// Example:
//
//	fileAuditSubscriber := subscribers.NewFileAuditSubscriber("./path/to/audit/file.json")
//	auditPublisher.Subscribe(fileAuditSubscriber)
//
// # To submit audit events to an endpoint, a NewHTTPAuditSubscriber needs to be created with a URI
//
// Example:
//
//	httpAuditSubscriber := subscribers.NewHTTPAuditSubscriber("http://audit-subscriber.test/api/new/")
//	auditPublisher.Subscribe(httpAuditSubscriber)
//
// To submit an event a Notify method of Publisher interface should be called with audit event params
//
// Example:
//
//	auditPublisher.Notify(models.AuditEvent{
//		Timestamp:   time.Now().Unix(),
//		Action:      models.AuditEventActionTypeShorten,
//		UserID:      userID,
//		OriginalURL: url,
//	})

package audit
