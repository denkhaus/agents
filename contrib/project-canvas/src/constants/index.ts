/**
 * Constants Index
 * Central export point for all application constants
 */

export {
  WORKSPACE_TYPES,
  NODE_TYPES,
  STORE_NAMES,
  WORKSPACE_NODE_MAPPING,
  NODE_STORE_MAPPING,
  WORKSPACE_STORE_MAPPING,
  type WorkspaceType,
  type NodeType,
  type StoreName,
  isValidWorkspaceType,
  isValidNodeType,
  isValidStoreName,
  getNodeTypeForWorkspace,
  getStoreNameForNodeType,
  getStoreNamesForWorkspace
} from './workspace-constants';