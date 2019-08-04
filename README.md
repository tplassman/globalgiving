# GlobalGiving Project Search

Filterable listing page for GlobalGiving projects using the GlobalGiving project search API:

[https://www.globalgiving.org/api/methods/search-projects/](https://www.globalgiving.org/api/methods/search-projects/)

## Installation

### Requirements

* [Go](https://golang.org/dl/) v1.11 or higher
* [NodeJS](https://nodejs.org/en/) v8.16.0 or higher (development only)
* [npm](https://www.npmjs.com/get-npm) v6.4.1 or higher (development only)

### Server Build Process

1. Copy `.example-env` to `.env`.

    ```
    $ cp .env-example .env
    ```

2. Add GlobalGiving API key to .env file.

3. Build or run server in `cmd` directory.

    ```
    $ cd cmd
    $ go build server.go
    $ cd ..
    $ ./cmd/server
    ```

    or

    ```
    $ go run cmd/server.go
    ```

4. Visit `localhost:8080` in web browser.

### Assets Builld Process (development only)

Styles and scripts are built from `assets` directory into `static` directory.

1. Install front end dependencies.

    ```
    $ npm install
    ```

2. Build or watch assets.

    ```
    $ gulp build
    ```

    or

    ```
    $ gulp watch
    ```

## Features

This challenge of this project was to extend the limit of the GlobalGiving project search API to allow users to search for higher limits than the maxium of 10, as defined in the API. Using a user provided limit in conjunction with the API limit we could predetermine the number of requests required to the API to populate the users search. The requests could be sent concurrently using go routines and ordered by placing the results into a map using the index of the request as a key. These results were then sorted by looping through the request count and appending results matching the current index into the final results set.

## Improvements

The project works to display a filterable list of active GlobalGiving projects.

* Filtering options are limited to search input and theme. Consider extending to country and organization since those are directly supported by the GlobalGiving API.
* Add sorting functionality to filter form to allow sorting by project attributes (e.g. date, funding progress, etc.).
* Dynamically populate filter options (e.g. theme) in filter form from API to ensure options are always up to date.
* Update URL structure so category filters such as theme and organization would use path slugs instead of query parameters for better looking URLs.
* The server currently routes all requests to the projects handler. Non-supported routes should be sent to a 404 handler.
* As more routes are added to the application, the templating logic could be abstracted into a dedicated package that would also parse an option `templates/_components` directory to allow for nesting reusable component partials into route pages.
* Add a `templates/_components` directory to encapsulate component markup.
* Add a `scripts/components` directory and corresponding component utilities to encapselate component scripts in a single class.
* Pass environment variables into templates to allow for live reloading of assets in development.
