"use client";

import { MainLayout } from "@/components/layout/main-layout";
import { useAgentConnection } from "@/hooks/use-agent-connection";
import React from 'react'

export default function Home() {
  // Initialize agent connection
  useAgentConnection();

  return <MainLayout />;
}
