/**
 * Convex Provider Component
 * Wraps the app with Convex real-time client
 */

import React from 'react';
import { ConvexProvider } from "convex/react";
import { convex } from "@/lib/convex";

interface ConvexAppProviderProps {
  children: React.ReactNode;
}

export const ConvexAppProvider: React.FC<ConvexAppProviderProps> = ({ 
  children 
}) => {
  return (
    <ConvexProvider client={convex}>
      {children}
    </ConvexProvider>
  );
};