package web

import (
	"net/http"

	"github.com/denkhaus/agents/pkg/multi/plugins/web/schema"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"trpc.group/trpc-go/trpc-agent-go/log"
	"trpc.group/trpc-go/trpc-agent-go/session"
)

func (s *Server) handleListSessions(w http.ResponseWriter, r *http.Request) {
	log.Infof("handleListSessions called: path=%s", r.URL.Path)
	vars := mux.Vars(r)
	appName := vars["appName"]
	agentIDString := vars["agentId"]

	agentID, err := uuid.Parse(agentIDString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userKey := session.UserKey{AppName: appName, UserID: agentID.String()}
	sessions, err := s.chatProcessor.ListSessions(r.Context(), userKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert internal sessions to ADK format.
	adkSessions := make([]*schema.ADKSession, 0, len(sessions))
	for _, sess := range sessions {

		adkSess, err := schema.NewADKSession(sess)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		adkSessions = append(adkSessions, adkSess)
	}

	s.writeJSON(w, adkSessions)
}

func (s *Server) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	log.Infof("handleCreateSession called: path=%s", r.URL.Path)
	vars := mux.Vars(r)
	appName := vars["appName"]
	agentIDString := vars["agentId"]

	agentID, err := uuid.Parse(agentIDString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key := session.Key{AppName: appName, UserID: agentID.String()}
	sess, err := s.chatProcessor.CreateSession(r.Context(), key, session.StateMap{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	adkSession, err := schema.NewADKSession(sess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, adkSession)
}

func (s *Server) handleGetSession(w http.ResponseWriter, r *http.Request) {
	log.Infof("handleGetSession called: path=%s", r.URL.Path)
	vars := mux.Vars(r)
	appName := vars["appName"]
	userID := vars["userId"]
	sessionID := vars["sessionId"]

	sess, err := s.chatProcessor.GetSession(r.Context(), session.Key{
		AppName:   appName,
		UserID:    userID,
		SessionID: sessionID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if sess == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	adkSession, err := schema.NewADKSession(sess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, adkSession)
}

func (s *Server) handleDeleteSession(w http.ResponseWriter, r *http.Request) {
	log.Infof("handleDeleteSession called: path=%s", r.URL.Path)
	vars := mux.Vars(r)
	appName := vars["appName"]
	agentIDString := vars["agentId"]
	sessionIDString := vars["sessionId"]

	agentID, err := uuid.Parse(agentIDString)
	if err != nil {
		s.logger.Error("failed to parse agentID", zap.String("agent_id", agentIDString))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID, err := uuid.Parse(sessionIDString)
	if err != nil {
		s.logger.Error("failed to parse sessionID", zap.String("session_id", sessionIDString))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.chatProcessor.DeleteSession(r.Context(), session.Key{
		AppName:   appName,
		UserID:    agentID.String(),
		SessionID: sessionID.String(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
