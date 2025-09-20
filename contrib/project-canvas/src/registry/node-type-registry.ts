/**
 * Node Type Registry
 * Central registry for all node types with unified operations
 */

import {
  getNodeTypeForWorkspace,
  type NodeType,
  type StoreName,
  type WorkspaceType,
} from "@/constants";
// Node Type Configuration Interface
export interface NodeTypeConfig {
  storeName: StoreName;
  displayName: string;
}

// Node Type Registry Class
export class NodeTypeRegistry {
  private static registry = new Map<NodeType, NodeTypeConfig>();

  static register(type: NodeType, config: NodeTypeConfig) {
    this.registry.set(type, config);
  }

  static get(type: NodeType): NodeTypeConfig | undefined {
    return this.registry.get(type);
  }

  static getAllTypes(): NodeType[] {
    return Array.from(this.registry.keys());
  }

  static getNodeTypeForWorkspace(workspace: WorkspaceType): NodeType {
    return getNodeTypeForWorkspace(workspace);
  }

  static getConfigForWorkspace(
    workspace: WorkspaceType
  ): NodeTypeConfig | undefined {
    const nodeType = this.getNodeTypeForWorkspace(workspace);
    return this.get(nodeType);
  }

  static isRegistered(type: NodeType): boolean {
    return this.registry.has(type);
  }

  static getDisplayName(type: NodeType): string {
    const config = this.get(type);
    return config?.displayName || "Items";
  }

  static getStoreName(type: NodeType): StoreName | undefined {
    const config = this.get(type);
    return config?.storeName;
  }
}
