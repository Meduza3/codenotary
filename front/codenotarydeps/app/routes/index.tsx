import React from "react"; // Explicit React import
import { json } from "@remix-run/node";
import { useFetcher } from "@remix-run/react";
import type { ActionFunction, LoaderFunction } from "@remix-run/node";

type Dependency = {
  id: string;
  score: number;
};

type FetcherData = {
  project_name: string;
  dependencies: Dependency[];
  message?: string;
  error?: string;
};

export const loader: LoaderFunction = async () => {
  return json({
    message: "Enter a project name to fetch dependencies",
    project_name: "",
    dependencies: [],
  });
};

export const action: ActionFunction = async ({ request }) => {
  const formData = await request.formData();
  const projectName = formData.get("projectName");

  if (!projectName) {
    return json({ message: "Invalid input", project_name: "", dependencies: [] });
  }

  const encodedName = encodeURIComponent(projectName as string);
  const response = await fetch(`http://localhost:8080/dependency/${encodedName}`, {
    method: "GET",
  });

  if (!response.ok) {
    return json(
      { error: "Failed to fetch data from the server", project_name: projectName, dependencies: [] },
      { status: 500 }
    );
  }

  const data = await response.json();
  return json({ ...data, project_name: projectName });
};

export default function Index() {
  const fetcher = useFetcher<FetcherData>();
  const dependencies = fetcher.data?.dependencies || [];
  const projectName = fetcher.data?.project_name || "";
  const isLoading = fetcher.state === "loading" || fetcher.state === "submitting";

  return (
    <div className="container">
      <div className="input-section">
        <h1>What project would you like to get dependencies for?</h1>
        <fetcher.Form method="post">
          <input
            type="text"
            name="projectName"
            placeholder="Enter project name here..."
            className="input-box"
          />
          <button type="submit" className="submit-button" disabled={isLoading}>
            {isLoading ? "Loading..." : "Send"}
          </button>
        </fetcher.Form>
      </div>

      <div className="list-section">
        <h1>Dependencies for: {isLoading ? "Loading..." : projectName || "No project selected"}</h1>
        <div className="scrollable-list">
          {isLoading ? (
            <p>Loading dependencies...</p>
          ) : dependencies.length > 0 ? (
            dependencies.map((item, index) => (
              <div className="rectangle" key={index}>
                <p>
                  <strong>ID:</strong> {item.id}
                </p>
                <p>
                  <strong>Score:</strong> {item.score.toFixed(1)}
                </p>
              </div>
            ))
          ) : (
            <p>{fetcher.data?.message || "Enter a project name to fetch dependencies"}</p>
          )}
        </div>
      </div>
    </div>
  );
}