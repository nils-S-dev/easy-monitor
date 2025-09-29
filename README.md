<a id="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">

  <h1 align="center">Easy Monitor</h1>

  <p align="center">
    A minimalist docker-ready self-hosted monitoring tool for checking the functionality of your services.
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#configuration">Configuration</a>
    </li>
    <li>
      <a href="#deployment">Deployment</a>
    </li>
    <li>
      <a href="#usage">Usage</a>
    </li>
    <li>
    <li>
      <a href="#development">Development</a>
    </li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

This simple monitoring tool can be fully configured via a JSON file. You can either check your monitors via an endpoint, which will trigger all monitors on-demand, or additionally setup cron jobs either globally for all of your monitors or different cron patterns for inidivual ones. 

Additionally you can configure SMTP settings to get notified about which monitors are failing

<p>(<a href="#readme-top">back to top</a>)</p>

### Built With

[![Go][Go]][Go-url]

<p>(<a href="#readme-top">back to top</a>)</p>

## Configuration

There are two files where the configuration for your easy monitor takes place

1. the `config.json` where you can configure your monitors whith endpoints and expected results as well as cron patterns to trigger them regularly.
2. the env variables in docker-compose.yaml where you can optionally configure an SMTP provider to send mails if monitors fail at those points in time you have configured in `config.json`

### `config.json`

* Copy the `config.example.json` to make yourself a template for the configuration.
* Name it `config.json`

* `cron` a global cron pattern that will trigger all of your monitors
* `notify` a list of receiver email addresses to notify about failing monitors
* `monitors` your monitors
  * `name` a custom name that will help you distinguish your monitor
  * `cron` an optional cron pattern that will override the global setting and applied to this monitor exclusively
  * `endpoint` the endpoint that will be called
  * `method` the HTTP method to use for the call
  * `body` the body to send to the endpoint
  * `expect` the expected results
    * `status`the expected HTTP status
    * `body` the expected response body as plain text or JSON. when JSON is returned the object will be parsed and compared by property and value.
    * `headers` the expected response headers with values

If any of the values provided in `expected` does not match the response from the endpoint the monitor will be considered failed.

### `docker-compose.yaml`

* Copy the `docker-compose.example.yaml` to make yourself a template for the configuration.
* Name it `docker-compose.yaml`

In the `env:` section specifically you can configure your email settings:

* `SMTP_ENABLED` will enabled email sending if `"true"`
* `SMTP_HOST` your SMTP host, such as `in-v3.mailjet.com`
* `SMTP_PORT` your SMTP port, usually `587`
* `SMTP_USER` your SMTP user name or API key depending on your provider
* `SMTP_PASS` your SMTP password or key depending on your provider
* `SMTP_FROM` the email address used as sender for the notifications

## Deployment

Start your easy monitor app with docker:

```
docker compose up -d
```

## Usage
Your monitors will be tested based on the cron patterns you provided in `config.json`.

Additionally there is an endpoint `/api/v1/monitors` that allows you to check your monitors all at once on demand.

## Development

### Prerequisites

* [![Go][Go]][Go-url]

### Start the Project

Start the project using

```
go run ./cmd
```

There is also a mock server for testing configs that you can start using

```
go run ./cmd/mock
```

### Tests

Integration tests are maintainend in the `tests` package

To run tests preferably use

```
go test ./tests
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

Thanks A LOT to othneildrew for [this amazing README.md template](https://github.com/othneildrew/Best-README-Template)!

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[Go]: https://img.shields.io/badge/golang-00ADD8?logo=go&logoColor=white&style=plastic
[Go-url]: https://go.dev/