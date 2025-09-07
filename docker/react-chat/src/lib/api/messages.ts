// import { AgentId } from "@/lib/constants/agents";
// import { sseService } from "./sse";
// import { AgentRunRequest, Message, MessagePart, AgentEvent } from "@/lib/types";
// import { debug } from "@/lib/utils/debug";

// export const messageApi = {
// /**
//  * Sends a message from a user to an agent and handles the streaming response.
//  *
//  * @param agentId - The ID of the agent to send the message to
//  * @param content - The message content
//  * @param sessionId - The session ID for the conversation
//  * @param userId - The user ID (defaults to 'user')
//  * @param onResponse - Callback for handling agent response events (includes regular responses, tool calls, etc.)
//  * @param onError - Callback for handling errors
//  */
// async sendMessage(
//   agentID: AgentId,
//   content: string,
//   sessionId: string,
//   onResponse?: (event: AgentEvent) => void,
//   onError?: (error: Error) => void
// ): Promise<void> {
//   const request: AgentRunRequest = {
//     appName: agentId,
//     agentID: agentID,
//     sessionID: sessionId,
//     streaming: true,
//     newMessage: {
//       role: "user",
//       parts: [{ text: content }],
//     },
//   };

//   // Use SSE service for streaming responses
//   sseService.connectAgentRun(request, {
//     onMessage: (event) => {
//       onResponse?.(event);
//     },
//     onError: (error) => {
//       console.error("SSE error during agent communication:", error);
//       onError?.(error instanceof Error ? error : new Error("Unknown error"));
//     },
//     onConnectionStatusChange: (connected) => {
//       // Handle connection status changes gracefully
//       if (!connected) {
//         debug.log("Agent SSE connection closed normally");
//       }
//     },
//   });
// },

// async sendMultiAgentMessage(request: MultiChatRequest) {
//   return apiClient.sendMessage(request);
// },

//   /**
//    * Converts an AgentEvent (from SSE stream) to a Message object for display.
//    * Handles all types of events: regular agent responses, inter-agent communication, tool calls, etc.
//    *
//    * @param event - The agent event to convert
//    * @returns A Message object suitable for display in the chat interface
//    */
//   convertEventToMessage(event: AgentEvent): Message {
//     // Safe content extraction with null checks
//     let content = "";
//     let parts: MessagePart[] | undefined = undefined;

//     if (typeof event.content === "string") {
//       content = event.content;
//     } else if (event.content && typeof event.content === "object") {
//       // Check if content has parts array
//       if (
//         Array.isArray(event.content.parts) &&
//         event.content.parts.length > 0
//       ) {
//         // Extract text content from parts
//         const textParts = event.content.parts.filter((part) => part.text);
//         content = textParts.map((part) => part.text).join("");

//         // Convert parts to MessagePart format
//         parts = event.content.parts.map((part) => {
//           const messagePart: MessagePart = {};

//           if (part.text) {
//             messagePart.text = part.text;
//           }

//           if (part.functionCall) {
//             messagePart.functionCall = {
//               name: part.functionCall.name,
//               args: part.functionCall.args,
//               id: part.functionCall.id,
//             };
//           }

//           if (part.functionResponse) {
//             messagePart.functionResponse = {
//               name: part.functionResponse.name,
//               response: part.functionResponse.response,
//               id: part.functionResponse.id,
//             };
//           }

//           return messagePart;
//         });
//       } else {
//         // Fallback: try to extract text from content object
//         content = JSON.stringify(event.content);
//       }
//     }

//     // Handle timestamp - could be Unix timestamp or already a Date
//     let timestamp: Date;
//     if (typeof event.timestamp === "number") {
//       // Unix timestamp (seconds or milliseconds)
//       timestamp =
//         event.timestamp > 1000000000000
//           ? new Date(event.timestamp)
//           : new Date(event.timestamp * 1000);
//     } else {
//       timestamp = new Date();
//     }

//     // Determine message type based on event properties
//     let messageType: "user" | "agent" | "inter_agent" | "system" = "system";

//     if (event.type === "inter_agent" || event.type === "communication") {
//       messageType = "inter_agent";
//     } else if (
//       event.object === "tool_call" ||
//       event.object === "tool_response"
//     ) {
//       messageType = "system";
//     } else if (event.author && event.author !== "system") {
//       messageType = "agent";
//     }

//     return {
//       id: event.id || event.invocationId || Date.now().toString(),
//       content,
//       timestamp,
//       sender: (event.fromAgent || event.author || "system") as AgentId,
//       type: messageType,
//       metadata: {
//         fromAgent: event.fromAgent,
//         toAgent: event.toAgent,
//         eventType: event.type,
//         partial: event.partial,
//         done: event.done,
//       },
//       parts,
//       usageMetadata: event.usageMetadata,
//     };
//   },
// };
