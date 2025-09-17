/**
 * Property Panel Types
 * Defines interfaces for unified property display system
 */

import { ReactNode } from 'react';
import { UUID } from './project.types';

// Base interface for all nodes that can display properties
export interface PropertyPanelNode {
  id: UUID;
  type: string;
  getPropertyInfo(): PropertyInfo;
}

// Property information structure
export interface PropertyInfo {
  id: UUID;
  type: string;
  title: string;
  description?: string;
  component: ReactNode;
}

// Form data structure for editable properties
export interface EditableNodeProperties {
  title: string;
  description: string;
}

// Callback type for property updates
export type PropertyUpdateCallback<T = EditableNodeProperties> = (
  nodeId: UUID,
  updates: Partial<T>
) => void;

// Property panel state
export interface PropertyPanelState {
  selectedNodeId: UUID | null;
  isEditing: boolean;
  hasUnsavedChanges: boolean;
}