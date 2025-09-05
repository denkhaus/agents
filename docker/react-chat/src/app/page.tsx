"use client";

import { MainLayout } from "@/components/layout/main-layout";
import { useAgentSetup } from "@/hooks/use-agent-setup";
import React from 'react'

export default function Home() {
  // Initialize agent setup (replaces the old useAgentConnection)
  const { isLoading, error } = useAgentSetup();

  if (isLoading) {
    return (
      <div className="h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p>Loading agents...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="h-screen flex items-center justify-center">
        <div className="text-center text-red-500">
          <p className="text-lg font-semibold">Failed to load agents</p>
          <p className="text-sm">{error.message}</p>
        </div>
      </div>
    );
  }

  return <MainLayout />;
}
