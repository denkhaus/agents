/**
 * Generic Singleton Registry
 * Professional, reusable singleton pattern for managing registries
 *
 * @template TKey - The type of keys used in the registry
 * @template TValue - The type of values stored in the registry
 */
export abstract class SingletonRegistry<TKey, TValue> {
  private static instances = new Map<string, SingletonRegistry<any, any>>();
  private registry = new Map<TKey, TValue>();

  /**
   * Protected constructor to prevent direct instantiation
   * Subclasses must implement their own getInstance method
   */
  protected constructor(private readonly registryName: string) {}

  /**
   * Get singleton instance for a specific registry type
   * This method should be overridden by subclasses
   */
  protected static getSingletonInstance<T extends SingletonRegistry<any, any>>(
    registryName: string,
    createInstance: () => T
  ): T {
    if (!this.instances.has(registryName)) {
      this.instances.set(registryName, createInstance());
    }
    return this.instances.get(registryName) as T;
  }

  /**
   * Register a key-value pair in the registry
   */
  public register(key: TKey, value: TValue): void {
    if (this.registry.has(key)) {
      console.warn(
        `Registry ${this.registryName}: Key already exists, overwriting:`,
        key
      );
    }
    this.registry.set(key, value);
  }

  /**
   * Get a value by key from the registry
   */
  public get(key: TKey): TValue | undefined {
    return this.registry.get(key);
  }

  /**
   * Check if a key exists in the registry
   */
  public has(key: TKey): boolean {
    return this.registry.has(key);
  }

  /**
   * Get all registered entries
   */
  public getAll(): Map<TKey, TValue> {
    return new Map(this.registry);
  }

  /**
   * Get all registered keys
   */
  public getKeys(): TKey[] {
    return Array.from(this.registry.keys());
  }

  /**
   * Get all registered values
   */
  public getValues(): TValue[] {
    return Array.from(this.registry.values());
  }

  /**
   * Get the number of registered entries
   */
  public size(): number {
    return this.registry.size;
  }

  /**
   * Clear all entries from the registry
   */
  public clear(): void {
    this.registry.clear();
  }

  /**
   * Remove a specific entry from the registry
   */
  public unregister(key: TKey): boolean {
    return this.registry.delete(key);
  }

  /**
   * Check if the registry has been initialized
   */
  public isInitialized(): boolean {
    return true; // Always considered initialized as there's no explicit initialization step anymore
  }

  /**
   * Get registry statistics
   */
  public getStats(): {
    name: string;
    size: number;
    keys: TKey[];
  } {
    return {
      name: this.registryName,
      size: this.size(),
      keys: this.getKeys(),
    };
  }
}

/**
 * Generic Singleton Manager
 * Manages multiple singleton registries
 */
export class SingletonManager {
  private static instance: SingletonManager | null = null;
  private registries = new Map<string, SingletonRegistry<any, any>>();

  private constructor() {}

  public static getInstance(): SingletonManager {
    if (!SingletonManager.instance) {
      SingletonManager.instance = new SingletonManager();
    }
    return SingletonManager.instance;
  }

  /**
   * Register a singleton registry
   */
  public registerRegistry<TKey, TValue>(
    name: string,
    registry: SingletonRegistry<TKey, TValue>
  ): void {
    this.registries.set(name, registry);
  }

  /**
   * Get a registered singleton registry
   */
  public getRegistry<TKey, TValue>(
    name: string
  ): SingletonRegistry<TKey, TValue> | undefined {
    return this.registries.get(name) as
      | SingletonRegistry<TKey, TValue>
      | undefined;
  }

  /**
   * Initialize all registered registries
   */
  public initializeAll(): void {
    // No need to initialize individual registries anymore as they manage their own state
  }

  /**
   * Get statistics for all registries
   */
  public getAllStats(): Record<string, any> {
    const stats: Record<string, any> = {};
    for (const [name, registry] of this.registries) {
      stats[name] = registry.getStats();
    }
    return stats;
  }

  /**
   * Clear all registries
   */
  public clearAll(): void {
    for (const registry of this.registries.values()) {
      registry.clear();
    }
  }
}
