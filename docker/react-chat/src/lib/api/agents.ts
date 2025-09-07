import { debug } from "../utils/debug";
import { apiClient } from "./client";
import { AgentInfo } from "@/lib/types";

export const agentApi = {
  async getAgents(): Promise<AgentInfo[]> {
    try {
      const agentsResult = await apiClient.getAgents();

      if (!Array.isArray(agentsResult)) {
        console.error(
          "Expected array, got:",
          typeof agentsResult,
          agentsResult
        );
        return [];
      }

      const applicationName = agentsResult.applicationName;
      const agentsResponse = agentsResult.agents || [];

      // Convert agent names to Agent objects with default properties
      const agents = agentsResponse.map<AgentInfo>((agent) => ({
        applicationName,
        id: agent.id,
        name: agent.name,
        role: agent.role,
        description: agent.description,
        status: "online",
        lastActivity: new Date(),
      }));

      return agents;
    } catch (error) {
      debug.error("Error fetching agents:", error);
      throw error;
    }
  },
};
