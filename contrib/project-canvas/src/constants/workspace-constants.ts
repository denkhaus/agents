/**
 * Workspace Constants
 * Central definition of all workspace types, node types, and store names
 * Ensures type safety and prevents hardcoding throughout the application
 */

// 1. Workspace Types
export const WORKSPACE_TYPES = {
  PROJECTS: 'projects',
  AGENTS: 'agents', 
  SETTINGS: 'settings'
} as const;

export type WorkspaceType = typeof WORKSPACE_TYPES[keyof typeof WORKSPACE_TYPES];

// 2. Node Types
export const NODE_TYPES = {
  TASK: 'task',
  PROJECT: 'project', 
  AGENT: 'agent',
  WORKFLOW: 'workflow',
  RESOURCE: 'resource'
} as const;

export type NodeType = typeof NODE_TYPES[keyof typeof NODE_TYPES];

// 3. Store Names
export const STORE_NAMES = {
  PROJECT: 'projectStore',
  TASK: 'taskStore',
  AGENT: 'agentStore',
  AGENT_PROJECT: 'agentProjectStore',
  UI: 'uiStore',
  SETTINGS: 'settingsStore'
} as const;

export type StoreName = typeof STORE_NAMES[keyof typeof STORE_NAMES];

// 4. Workspace-zu-Node Mapping
export const WORKSPACE_NODE_MAPPING = {
  [WORKSPACE_TYPES.PROJECTS]: NODE_TYPES.TASK,
  [WORKSPACE_TYPES.AGENTS]: NODE_TYPES.AGENT,
  [WORKSPACE_TYPES.SETTINGS]: NODE_TYPES.PROJECT // fallback
} as const;

// 5. Node-zu-Store Mapping
export const NODE_STORE_MAPPING = {
  [NODE_TYPES.TASK]: STORE_NAMES.TASK,
  [NODE_TYPES.PROJECT]: STORE_NAMES.PROJECT,
  [NODE_TYPES.AGENT]: STORE_NAMES.AGENT
} as const;

// 6. Workspace-zu-Store Mapping
export const WORKSPACE_STORE_MAPPING = {
  [WORKSPACE_TYPES.PROJECTS]: [STORE_NAMES.PROJECT, STORE_NAMES.TASK],
  [WORKSPACE_TYPES.AGENTS]: [STORE_NAMES.AGENT, STORE_NAMES.AGENT_PROJECT],
  [WORKSPACE_TYPES.SETTINGS]: [STORE_NAMES.SETTINGS, STORE_NAMES.UI]
} as const;

// 7. Type Guards für Runtime Safety
export const isValidWorkspaceType = (value: string): value is WorkspaceType => {
  return Object.values(WORKSPACE_TYPES).includes(value as WorkspaceType);
};

export const isValidNodeType = (value: string): value is NodeType => {
  return Object.values(NODE_TYPES).includes(value as NodeType);
};

export const isValidStoreName = (value: string): value is StoreName => {
  return Object.values(STORE_NAMES).includes(value as StoreName);
};

// 8. Helper Functions
export const getNodeTypeForWorkspace = (workspace: WorkspaceType): NodeType => {
  return WORKSPACE_NODE_MAPPING[workspace];
};

export const getStoreNameForNodeType = (nodeType: NodeType): StoreName | undefined => {
  return NODE_STORE_MAPPING[nodeType as keyof typeof NODE_STORE_MAPPING];
};

export const getStoreNamesForWorkspace = (workspace: WorkspaceType): readonly StoreName[] => {
  return WORKSPACE_STORE_MAPPING[workspace];
};