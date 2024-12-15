import React from "react"; // Explicit React import
import { json } from "@remix-run/node";
import { useFetcher } from "@remix-run/react";
import type { ActionFunction, LoaderFunction } from "@remix-run/node";
import { Bar } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

type Dependency = {
  id: string;
  score: number;
  check_scores: Record<string, number>;
};

type FetcherData = {
  project_name: string;
  dependencies: Dependency[];
  main_scores?: Record<string, number>;
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
  const API_URL = process.env.API_URL || "http://localhost:8080";
  const response = await fetch(`${API_URL}/dependency/${encodedName}`, {
    method: "GET",
  });
  if (!response.ok) {
    return json(
      { error: "Failed to fetch data from the server: You haven't searched for this project before so it's not saved in the database.", project_name: projectName, dependencies: [] },
      { status: 500 }
    );
  }

  const data = await response.json();
  return json({ ...data, project_name: projectName });
};


const getRoundedCounts = (data: number[]) => {
  const counts: Record<number, number> = {};
  for (let i = 0; i <= 10; i++) {
    counts[i] = 0;
  }
  data.forEach((value) => {
    const rounded = Math.round(value);
    if (rounded in counts) {
      counts[rounded]++;
    }
  });
  return Object.entries(counts).map(([key, value]) => ({
    bucket: parseInt(key),
    count: value,
  }));
};


export default function Index() {
  const scoreNames: Record<string, string> = {
    BA: "Binary-Artifacts",
    BP: "Branch-Protection",
    CBP: "CII-Best-Practices",
    CR: "Code-Review",
    DW: "Dangerous-Workflow",
    F: "Fuzzing",
    L: "License",
    M: "Maintained",
    P: "Packaging",
    PD: "Pinned-Dependencies",
    S: "SAST",
    SP: "Security-Policy",
    SR: "Signed-Releases",
    TP: "Token-Permissions",
    V: "Vulnerabilities",
  };
  const fetcher = useFetcher<FetcherData>();
  const dependencies = fetcher.data?.dependencies || [];
  const projectName = fetcher.data?.project_name || "";
  const isLoading = fetcher.state === "loading" || fetcher.state === "submitting";

  const scores = dependencies.map((dep) => dep.score);
  const roundedData = getRoundedCounts(scores);
  const mainScores = fetcher.data?.main_scores || {};

  const chartData = {
    labels: roundedData.map((item) => item.bucket.toString()),
    datasets: [
      {
        label: "Count of Scores (rounded)",
        data: roundedData.map((item) => item.count),
        backgroundColor: "rgba(75, 192, 192, 0.6)",
        borderColor: "rgba(75, 192, 192, 1)",
        borderWidth: 1,
      },
    ],
  };


  return (
    <div className="container">
      <div className="input-section">
        <h1>What project would you like to get dependencies for?</h1>
        <fetcher.Form method="post">
          <input
            type="text"
            name="projectName"
            placeholder="github.com/cli/cli..."
            className="input-box"
          />
          <button type="submit" className="submit-button" disabled={isLoading}>
            {isLoading ? "Loading..." : "Send"}
          </button>
        </fetcher.Form>
      </div>

      <div className="main-section">
        <div className="list-section">
          <h1>Dependencies for: {isLoading ? "Loading..." : projectName || "No project selected"}</h1>
          <div className="scrollable-list">
            {isLoading ? (
              <p>Loading dependencies...</p>
            ) : dependencies.length > 0 ? (
              dependencies.map((item, index) => {
                let scoreClass = "";
                if (item.score > 8.0) {
                  scoreClass = "rectangle-high";
                } else if (item.score > 4.0) {
                  scoreClass = "rectangle-medium";
                } else {
                  scoreClass = "rectangle-low";
                }
              
                const cs = item.check_scores || {};
                const formatScore = (val: number) => (val === -1 ? "?" : val.toString());
              
                return (
                  <div
                    className={`rectangle ${scoreClass}`}
                    key={index}
                    onClick={() => {
                      if (item.score === -1 || dependencies.length === 0) {
                        alert("No known dependencies for this item.");
                      } else {
                        const formData = new FormData();
                        formData.append("projectName", item.id);
                        fetcher.submit(formData, { method: "post" });
                      }
                    }}
                  >
                    <p><strong>ID:</strong> {item.id}</p>
                    <p><strong>Score:</strong> {item.score === -1 ? "???" : item.score.toFixed(1)}</p>
              
                    {item.score !== -1 && item.check_scores && (
                    <div className="check-scores">
                      <div className="check-scores-list">
                      {Object.entries(cs).map(([key, value]) => {
                        // Find the abbreviation by matching the full key name
                        const abbreviation = Object.keys(scoreNames).find(
                          (abbr) => scoreNames[abbr] === key
                        );

                        if (!abbreviation) return null; // Skip if no abbreviation found

                        return (
                          <span key={key} className="score-item">
                            <span className="abbr">{abbreviation}</span>
                            <span className="full">{scoreNames[abbreviation]}</span>: {formatScore(value)}
                          </span>
                        );
                      })}
                    </div>
                    </div>
                  )}

                  </div>
                );
              })
              
            ) : (
              <p>{fetcher.data?.error || "Enter a project name to fetch dependencies"}</p>
            )}
          </div>
        </div>

        <div className="graph-section">
          <h1>Dependency Score Graph</h1>
          <div className="graph-placeholder">
            {dependencies.length > 0 ? (
              <Bar data={chartData} options={{ responsive: true, plugins: { legend: { position: "top" } } }} />
            ) : (
              <p>No graph data to display.</p>
            )}
          </div>
          <div className="main-scores-section">
            <h2>{projectName} Scores</h2>
            {Object.keys(mainScores).length > 0 ? (
              <div className="main-scores-list">
                {Object.entries(mainScores).map(([key, value]) => {
                  // Generate abbreviation from the key
                  const abbreviation = key
                    .split('-')
                    .map((word) => word.charAt(0))
                    .join('');
                  
                  // Get the full name from scoreNames mapping
                  const fullName = scoreNames[abbreviation] || key;
                  
                  return (
                    <span key={key} className="score-item">
                      <span className="abbr">{abbreviation}</span>
                      <span className="full">{fullName}</span>: {value === -1 ? '?' : value}
                    </span>
                  );
                })}
              </div>
            ) : (
              <p>No project selected yet.</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
