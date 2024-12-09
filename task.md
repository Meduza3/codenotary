Build a Simple Application Using the deps.dev API

Features

    Data Management:
        Fetch dependencies of the repository https://github.com/cli/cli using the deps.dev API. DONE
        Extract and store all dependencies, including their OpenSSF scores, in an SQLite database.
        
    Backend API:
        Develop an API to:
            Add or update dependency details in the SQLite database, integrating data from deps.dev.
            Retrieve stored dependency details, including the OpenSSF score.
            Provide CRUD operations to manage dependencies.
            Enable querying by dependency name and OpenSSF score.

    Frontend:
        Create a user-friendly interface that:
            Displays dependencies and their OpenSSF scores dynamically.
            Uses a chart to visualize the scores of all dependencies.
            Allows users to search for specific dependencies by name and view their scores.

    Integration with deps.dev:
        Dynamically fetch and update package information using the deps.dev API where applicable.

    Documentation:
        Include a README with:
            Instructions for setup, running, and usage of the application.
            API documentation that describes available endpoints and their usage.
        Detail the SQLite schema used to store the dependency data.

    Ease of Deployment:
        Provide a docker-compose.yml file to run the application, ensuring seamless deployment of the backend, frontend, and database components.