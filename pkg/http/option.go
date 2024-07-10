package http

// WithDB configures this server to service the specified DB
func WithDB(db DB) Option {
	return func(s *Server) { s.d = db }
}

// WithBillProcessor configures the server with the specified billing
// engine.
func WithBillProcessor(bp BillProcessor) Option {
	return func(s *Server) { s.bp = bp }
}
