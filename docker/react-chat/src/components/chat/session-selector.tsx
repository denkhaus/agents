"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useChatStore } from "@/lib/store";
import { formatDistanceToNow } from "date-fns";
import { Plus, MessageCircle, Clock, Trash2 } from "lucide-react";
import { toast } from "sonner";

interface SessionSelectorProps {
  agentId: string;
  onSessionSelect?: (sessionId: string) => void;
}

export function SessionSelector({
  agentId,
  onSessionSelect,
}: SessionSelectorProps) {
  const {
    availableSessions,
    currentSessionId,
    createSession,
    loadSessionMessages,
    setCurrentSession,
    deleteSession,
  } = useChatStore();

  const [isCreating, setIsCreating] = useState(false);
  const [deletingSessionId, setDeletingSessionId] = useState<string | null>(
    null
  );
  const sessions = availableSessions[agentId] || [];

  const handleCreateSession = async () => {
    setIsCreating(true);
    try {
      await createSession(agentId);
    } finally {
      setIsCreating(false);
    }
  };

  const handleSelectSession = async (sessionId: string) => {
    if (sessionId !== currentSessionId) {
      await loadSessionMessages(agentId, sessionId);
      setCurrentSession(sessionId);
      onSessionSelect?.(sessionId);
    }
  };

  const handleDeleteSession = async (
    sessionId: string,
    event: React.MouseEvent
  ) => {
    event.stopPropagation(); // Prevent session selection when clicking delete

    if (sessions.length <= 1) {
      // Don't allow deleting the last session
      return;
    }

    setDeletingSessionId(sessionId);
    try {
      console.log('Attempting to delete session:', { agentId, sessionId });
      await deleteSession(agentId, sessionId);
      console.log('Session deletion completed successfully');
      toast.success('Session deleted successfully');
    } catch (error) {
      console.error("Failed to delete session:", error);
      // Show user-friendly error message
      if (error instanceof Error) {
        if (error.message.includes('Network error')) {
          toast.warning('Unable to connect to server. Session removed locally.');
        } else if (error.message.includes('CORS')) {
          toast.warning('Session removed locally (server CORS restriction).');
        } else {
          toast.error(`Failed to delete session: ${error.message}`);
        }
      } else {
        toast.error('An unexpected error occurred while deleting the session.');
      }
    } finally {
      setDeletingSessionId(null);
    }
  };

  return (
    <Card className="h-full flex flex-col m-2">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">Chat Sessions</CardTitle>
          <Button
            onClick={handleCreateSession}
            disabled={isCreating}
            size="sm"
            className="h-8"
          >
            <Plus className="h-4 w-4 mr-1" />
            New
          </Button>
        </div>
        <p className="text-sm text-muted-foreground">Agent: {agentId}</p>
      </CardHeader>

      <CardContent className="p-4 flex-1">
        {sessions.length === 0 ? (
          <div className="p-4 text-center text-muted-foreground">
            <MessageCircle className="h-8 w-8 mx-auto mb-2 opacity-50" />
            <p className="text-sm">No sessions yet</p>
            <p className="text-xs">Create a new session to start chatting</p>
          </div>
        ) : (
          <div className="space-y-2 flex-1 overflow-y-auto">
            {sessions.map((session) => (
              <div
                key={session.id}
                className={`group p-3 rounded-lg border cursor-pointer transition-colors hover:bg-muted/50 ${
                  currentSessionId === session.id
                    ? "bg-primary/10 border-primary"
                    : "bg-background border-border"
                }`}
                onClick={() => handleSelectSession(session.id)}
              >
                <div className="flex items-start justify-between mb-2">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate">
                      Session {session.id.slice(-8)}
                    </p>
                    <div className="flex items-center gap-1 mt-1">
                      <Clock className="h-3 w-3 text-muted-foreground" />
                      <span className="text-xs text-muted-foreground">
                        {formatDistanceToNow(
                          new Date(session.lastUpdateTime * 1000),
                          {
                            addSuffix: true,
                          }
                        )}
                      </span>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {currentSessionId === session.id && (
                      <Badge variant="secondary" className="text-xs">
                        Active
                      </Badge>
                    )}
                    {sessions.length > 1 && (
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity hover:bg-destructive hover:text-destructive-foreground"
                        onClick={(e) => handleDeleteSession(session.id, e)}
                        disabled={deletingSessionId === session.id}
                        title="Delete session"
                      >
                        <Trash2 className="h-3 w-3" />
                      </Button>
                    )}
                  </div>
                </div>

                {session.events && session.events.length > 0 && (
                  <div className="text-xs text-muted-foreground">
                    <p className="truncate">
                      {session.events.length} message
                      {session.events.length !== 1 ? "s" : ""}
                    </p>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
