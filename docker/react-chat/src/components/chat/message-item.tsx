"use client";

import React from "react";
import { Message } from "@/lib/types";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { formatDistanceToNow } from "date-fns";
import { ArrowRight, BrainCircuit, Cog, Play, Flag } from "lucide-react";
import { Markdown } from "@/components/ui/markdown";
import { AGENT_IDS, getAgentDisplayName } from "@/lib/constants/agents";

interface MessageItemProps {
  message: Message;
}

import { MessagePart } from "@/lib/types";

const StructuredPart = ({ part }: { part: MessagePart }) => {
  const partType = Object.keys(part)[0];
  const content = part[partType]?.content;
  if (!content) return null;

  const config = {
    planning: {
      Icon: Cog,
      title: "Planning",
      style:
        "bg-yellow-50 border-yellow-200 dark:bg-yellow-950 dark:border-yellow-800",
      titleStyle: "text-yellow-800 dark:text-yellow-200",
    },
    reasoning: {
      Icon: BrainCircuit,
      title: "Reasoning",
      style:
        "bg-purple-50 border-purple-200 dark:bg-purple-950 dark:border-purple-800",
      titleStyle: "text-purple-800 dark:text-purple-200",
    },
    action: {
      Icon: Play,
      title: "Action",
      style: "bg-blue-50 border-blue-200 dark:bg-blue-950 dark:border-blue-800",
      titleStyle: "text-blue-800 dark:text-blue-200",
    },
    final_answer: {
      Icon: Flag,
      title: "Final Answer",
      style:
        "bg-green-50 border-green-200 dark:bg-green-950 dark:border-green-800",
      titleStyle: "text-green-800 dark:text-green-200",
    },
  }[partType];

  if (!config) return null;

  const { Icon, title, style, titleStyle } = config;

  const renderContent = () => (
    <div className={`mt-2 p-3 border rounded-lg ${style}`}>
      <div className="flex items-center gap-2 mb-2">
        <Icon className={`h-5 w-5 ${titleStyle}`} />
        <strong className={titleStyle}>{title}</strong>
      </div>
      <div className="prose prose-sm max-w-none dark:prose-invert">
        <Markdown>{content}</Markdown>
      </div>
    </div>
  );

  if (partType === "reasoning") {
    return (
      <details className="mt-2">
        <summary className="cursor-pointer text-sm text-purple-600 hover:text-purple-800">
          Show Reasoning
        </summary>
        {renderContent()}
      </details>
    );
  }

  return renderContent();
};

export function MessageItem({ message }: MessageItemProps) {
  const isUser = message.sender === AGENT_IDS.HUMAN;
  const isInterAgent = message.type === "inter_agent";
  const isSystem = message.type === "system";
  const isToolCall = message.parts?.some((part) => part.functionCall);
  const isToolResponse = message.parts?.some((part) => part.functionResponse);
  const isToolCode =
    message.type === "system" && message.metadata?.object === "tool_code";
  const hasStructuredThoughts = message.metadata?.hasStructuredThoughts;

  const formatMessageContent = () => {
    const renderedElements: JSX.Element[] = [];

    // Priority 1: Render tool calls/responses if they exist in message.parts
    if (message.parts && message.parts.length > 0) {
      message.parts.forEach((part, index) => {
        if (part.functionCall) {
          renderedElements.push(
            <ToolCallPart key={`tool-call-${index}`} part={part.functionCall} />
          );
        }
        if (part.functionResponse) {
          renderedElements.push(
            <ToolResponsePart
              key={`tool-response-${index}`}
              part={part.functionResponse}
            />
          );
        }
      });
    }

    // Priority 2: Render structured thoughts if they exist
    if (hasStructuredThoughts && message.parts && message.parts.length > 0) {
      message.parts.forEach((part, index) => {
        const partType = Object.keys(part)[0];
        if (
          ["planning", "reasoning", "action", "final_answer"].includes(partType)
        ) {
          renderedElements.push(
            <StructuredPart key={`structured-${index}`} part={part} />
          );
        }
      });
    }

    // Priority 3: Render tool code if it's a system message with tool_code object
    if (isToolCode && message.content) {
      renderedElements.push(
        <ToolCodePart key="tool-code-main" part={{ code: message.content }} />
      );
    }

    // Priority 4: Render main message content if no specific parts were rendered
    if (renderedElements.length === 0 && message.content) {
      renderedElements.push(
        <div
          key="main-content"
          className="prose prose-sm max-w-none dark:prose-invert"
        >
          {message.sender === AGENT_IDS.HUMAN ? (
            <p className="whitespace-pre-wrap">{message.content}</p>
          ) : (
            <Markdown>{message.content}</Markdown>
          )}
        </div>
      );
    }

    return renderedElements.length > 0 ? renderedElements : null;
  };

  const getCardStyle = () => {
    if (isUser) return "bg-primary text-primary-foreground";
    if (isInterAgent)
      return "border-purple-200 dark:border-purple-800 bg-purple-50 dark:bg-purple-950";
    if (isToolCall) return "border-blue-200 dark:border-blue-800";
    if (isToolResponse) return "border-green-200 dark:border-green-800";
    if (isToolCode) return "border-slate-200 dark:border-slate-800";
    if (isSystem) return "border-orange-200 dark:border-orange-800";
    return "";
  };

  const getMessageTypeLabel = () => {
    if (isUser) return "You";
    if (
      isInterAgent &&
      message.metadata?.fromAgent &&
      message.metadata?.toAgent
    ) {
      const fromName =
        message.metadata.fromAgent === AGENT_IDS.HUMAN
          ? "You"
          : getAgentDisplayName(message.metadata.fromAgent);
      const toName =
        message.metadata.toAgent === AGENT_IDS.HUMAN
          ? "You"
          : getAgentDisplayName(message.metadata.toAgent);
      return `${fromName} -> ${toName}`;
    }

    const displayName = getAgentDisplayName(message.sender);
    if (hasStructuredThoughts) return `${displayName} (Thinking)`;
    if (isToolCall) return `${displayName} (Tool Call)`;
    if (isToolResponse) return `${displayName} (Tool Response)`;
    if (isToolCode) return `${displayName} (Tool Code)`;
    if (isSystem) return `${displayName} (System)`;
    return displayName;
  };

  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
      <Card className={`max-w-[90%] ${getCardStyle()}`}>
        <CardContent className="p-3">
          <div className="flex items-center gap-2 mb-2">
            <span className="font-medium text-sm">{getMessageTypeLabel()}</span>

            {isInterAgent &&
              message.metadata?.fromAgent &&
              message.metadata?.toAgent && (
                <div className="flex items-center gap-1 text-xs">
                  <Badge variant="secondary" className="px-1 py-0">
                    {message.metadata.fromAgent}
                  </Badge>
                  <ArrowRight className="h-3 w-3" />
                  <Badge variant="secondary" className="px-1 py-0">
                    {message.metadata.toAgent}
                  </Badge>
                </div>
              )}

            <span className="text-xs opacity-70 ml-auto">
              {formatDistanceToNow(message.timestamp, { addSuffix: true })}
            </span>
          </div>

          <div className="text-sm">{formatMessageContent()}</div>

          <div className="flex gap-2 mt-2">
            {message.metadata?.partial && (
              <Badge variant="outline" className="text-xs">
                Streaming...
              </Badge>
            )}
            {hasStructuredThoughts && (
              <Badge
                variant="secondary"
                className="text-xs bg-purple-100 text-purple-800"
              >
                Structured Thought
              </Badge>
            )}
            {isToolCall && (
              <Badge
                variant="secondary"
                className="text-xs bg-blue-100 text-blue-800"
              >
                Tool Call
              </Badge>
            )}
            {isToolResponse && (
              <Badge
                variant="secondary"
                className="text-xs bg-green-100 text-green-800"
              >
                Tool Response
              </Badge>
            )}
            {isToolCode && (
              <Badge
                variant="secondary"
                className="text-xs bg-slate-100 text-slate-800"
              >
                Tool Code
              </Badge>
            )}
            {isInterAgent && (
              <Badge
                variant="secondary"
                className="text-xs bg-purple-100 text-purple-800"
              >
                Inter-Agent
              </Badge>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

// Helper components for different tool parts to keep the main component clean

const ToolCallPart = ({
  part,
}: {
  part: { name: string; args: unknown; id?: string };
}) => (
  <div className="mt-2 p-3 bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800 rounded-lg">
    <div className="flex items-center gap-2 mb-2">
      <strong className="text-blue-800 dark:text-blue-200">
        Tool Call: {part.name}
      </strong>
      {part.id && (
        <Badge variant="outline" className="text-xs">
          {part.id}
        </Badge>
      )}
    </div>
    <details className="mt-2">
      <summary className="cursor-pointer text-sm text-blue-600 hover:text-blue-800">
        View Parameters
      </summary>
      <pre className="mt-2 text-xs overflow-x-auto bg-white dark:bg-gray-900 p-2 rounded border max-w-full">
        {JSON.stringify(part.args, null, 2)}
      </pre>
    </details>
  </div>
);

const ToolResponsePart = ({
  part,
}: {
  part: { name: string; response: unknown; id?: string };
}) => (
  <div className="mt-2 p-3 bg-green-50 dark:bg-green-950 border border-green-200 dark:border-green-800 rounded-lg">
    <div className="flex items-center gap-2 mb-2">
      <strong className="text-green-800 dark:text-green-200">
        Tool Response: {part.name}
      </strong>
      {part.id && (
        <Badge variant="outline" className="text-xs">
          {part.id}
        </Badge>
      )}
    </div>
    <details className="mt-2">
      <summary className="cursor-pointer text-sm text-green-600 hover:text-green-800">
        View Response
      </summary>
      <pre className="mt-2 text-xs overflow-x-auto bg-white dark:bg-gray-900 p-2 rounded border max-h-40 max-w-full">
        {JSON.stringify(part.response, null, 2)}
      </pre>
    </details>
  </div>
);

const ToolCodePart = ({ part }: { part: { code: string } }) => (
  <div className="mt-2 p-3 bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-lg">
    <div className="flex items-center gap-2 mb-2">
      <strong className="text-slate-800 dark:text-slate-200">Tool Code</strong>
    </div>
    <div className="mt-2 text-xs overflow-x-auto bg-white dark:bg-gray-900 p-2 rounded border max-w-full">
      <Markdown>{`${part.code}`}</Markdown>
    </div>
  </div>
);
