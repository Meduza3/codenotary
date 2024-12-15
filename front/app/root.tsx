import React from "react";
import "./styles/root.css";
import { Links, LiveReload, Meta, Outlet, Scripts } from "@remix-run/react";
import type { LinksFunction } from "@remix-run/node";

export const links: LinksFunction = () => [
  { rel: "stylesheet", href: "/styles/root.css" },
];
export default function App() {
  return (
    <html lang="en">
      <head>
        <Meta />
        <Links />
      </head>
      <body>
        <Outlet /> {}
        <Scripts />
        <LiveReload />
      </body>
    </html>
  );
}
