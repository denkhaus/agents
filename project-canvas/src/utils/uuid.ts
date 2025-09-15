/**
 * UUID utility functions
 * Type-safe UUID generation and validation
 */

import { UUID } from '../types/project.types';

/**
 * Generate a new UUID v4
 */
export function generateUUID(): UUID {
  return crypto.randomUUID();
}

/**
 * Validate UUID format
 */
export function isValidUUID(uuid: string): uuid is UUID {
  const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
  return uuidRegex.test(uuid);
}

/**
 * Type guard for UUID
 */
export function assertUUID(value: string): asserts value is UUID {
  if (!isValidUUID(value)) {
    throw new Error(`Invalid UUID format: ${value}`);
  }
}

/**
 * Convert string to UUID with validation
 */
export function toUUID(value: string): UUID {
  assertUUID(value);
  return value;
}

/**
 * Generate multiple UUIDs
 */
export function generateUUIDs(count: number): UUID[] {
  return Array.from({ length: count }, () => generateUUID());
}

/**
 * Create a UUID from a seed (for testing/dummy data)
 * Note: This is not cryptographically secure, only for development
 */
export function createTestUUID(seed: string): UUID {
  // Simple hash function for consistent test UUIDs
  let hash = 0;
  for (let i = 0; i < seed.length; i++) {
    const char = seed.charCodeAt(i);
    hash = ((hash << 5) - hash) + char;
    hash = hash & hash; // Convert to 32-bit integer
  }
  
  // Convert to hex and pad
  const hex = Math.abs(hash).toString(16).padStart(8, '0');
  
  // Format as UUID v4
  return `${hex.slice(0, 8)}-${hex.slice(0, 4)}-4${hex.slice(1, 4)}-8${hex.slice(2, 5)}-${hex.slice(0, 12)}` as UUID;
}