import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./app/globals.css";
import { ConvexAppProvider } from "./providers/convex-provider.tsx";
import { Toaster } from "@/components/ui/toaster";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ConvexAppProvider>
      <App />
      <Toaster />
    </ConvexAppProvider>
  </React.StrictMode>
);
