import React from "react"; // Explicit React import
import "./styles/root.css"; // Standard CSS import
import { Links, LiveReload, Meta, Outlet, Scripts } from "@remix-run/react";
import type { LinksFunction } from "@remix-run/node";

export const links: LinksFunction = () => [
  { rel: "stylesheet", href: "/styles/root.css" }, // Link to the stylesheet
];
export default function App() {
  return (
    <html lang="en">
      <head>
        <Meta />
        <Links />
      </head>
      <body>
        <Outlet /> {/* Placeholder for routes */}
        <Scripts />
        <LiveReload />
      </body>
    </html>
  );
}
