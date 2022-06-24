package server

func (s *Server) Addpaths() {
	s.Router.POST("/bot/:token", HandleBotWebHook)
}
