/// <reference types="vite/client" />

/**
 * Convex Client Configuration
 * Sets up the Convex client for real-time data
 */

import { ConvexReactClient } from "convex/react";

const convexUrl = import.meta.env?.VITE_CONVEX_URL;

if (!convexUrl) {
  throw new Error(
    "Missing VITE_CONVEX_URL environment variable. " +
      "Run `npx convex dev` to get your Convex URL."
  );
}

export const convex = new ConvexReactClient(convexUrl);
