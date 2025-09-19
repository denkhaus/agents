/**
 * Store Registry
 * Central registry for all store operations with type safety
 */

import { 
  NODE_STORE_MAPPING, 
  type StoreName, 
  type NodeType 
} from '@/constants';

// Base interface for all store operations
export interface StoreOperations {
  getCount(): number;
  updatePosition?(id: string, position: { x: number; y: number }): Promise<void>;
  getDisplayName(): string;
}

// Store Registry Class
export class StoreRegistry {
  private static stores = new Map<StoreName, StoreOperations>();
  
  static register(storeName: StoreName, operations: StoreOperations) {
    this.stores.set(storeName, operations);
  }
  
  static get(storeName: StoreName): StoreOperations | undefined {
    return this.stores.get(storeName);
  }
  
  static getStoreForNodeType(nodeType: NodeType): StoreOperations | undefined {
    const storeName = NODE_STORE_MAPPING[nodeType as keyof typeof NODE_STORE_MAPPING];
    return storeName ? this.get(storeName) : undefined;
  }
  
  static getAllStores(): Map<StoreName, StoreOperations> {
    return new Map(this.stores);
  }
  
  static isRegistered(storeName: StoreName): boolean {
    return this.stores.has(storeName);
  }
}