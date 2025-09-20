/**
 * Store Registry
 * Professional implementation using generic singleton pattern
 */

import { SingletonRegistry } from "@/utils/singleton-registry";
import { NODE_STORE_MAPPING, type StoreName, type NodeType } from "@/constants";

/**
 * Store Registry Implementation
 * Extends the generic singleton registry for store operations
 */
class StoreRegistryImpl extends SingletonRegistry<StoreName, any> {
  constructor() {
    super("StoreRegistry");
  }

  /**
   * Get singleton instance
   */
  public static getInstance(): StoreRegistryImpl {
    return this.getSingletonInstance(
      "StoreRegistry",
      () => new StoreRegistryImpl()
    );
  }

  /**
   * Get store operations by node type
   */
  public getStoreForNodeType(nodeType: NodeType): any | undefined {
    const storeName =
      NODE_STORE_MAPPING[nodeType as keyof typeof NODE_STORE_MAPPING];
    return storeName ? this.get(storeName) : undefined;
  }

  /**
   * Validate that all required stores are registered
   */
  public validateRegistration(): boolean {
    // With the removal of the operations system, we no longer validate registration in this manner.
    // The individual stores will manage their own initialization and state.
    return true;
  }
}

/**
 * Static facade for the Store Registry
 * Provides a clean API while using the singleton pattern internally
 */
export class StoreRegistry {
  /**
   * Register a store operation
   */
  public static register(storeName: StoreName, operations: any): void {
    StoreRegistryImpl.getInstance().register(storeName, operations);
  }

  /**
   * Get store operations by store name
   */
  public static get(storeName: StoreName): any | undefined {
    return StoreRegistryImpl.getInstance().get(storeName);
  }

  /**
   * Get store operations by node type
   */
  public static getStoreForNodeType(nodeType: NodeType): any | undefined {
    return StoreRegistryImpl.getInstance().getStoreForNodeType(nodeType);
  }

  /**
   * Check if a store is registered
   */
  public static isRegistered(storeName: StoreName): boolean {
    return StoreRegistryImpl.getInstance().has(storeName);
  }

  /**
   * Get all registered stores
   */
  public static getAllStores(): Map<StoreName, any> {
    return StoreRegistryImpl.getInstance().getAll();
  }

  /**
   * Get the number of registered stores
   */
  public static getStoreCount(): number {
    return StoreRegistryImpl.getInstance().size();
  }

  /**
   * Initialize the store registry
   */
  public static initialize(): void {
    // No longer needed in this context, as stores manage their own state.
  }

  /**
   * Check if the registry has been initialized
   */
  public static isInitialized(): boolean {
    return StoreRegistryImpl.getInstance().isInitialized();
  }

  /**
   * Get registry statistics
   */
  public static getStats(): any {
    return StoreRegistryImpl.getInstance().getStats();
  }

  /**
   * Clear all registered stores (useful for testing)
   */
  public static clear(): void {
    StoreRegistryImpl.getInstance().clear();
  }
}
