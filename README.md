# gigcity-site
[![Gitter](https://badges.gitter.im/Join Chat.svg)](https://gitter.im/GDG-Gigcity/gigcity-site?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Public facing site for GDG Gigcity

## Setup

### Install the Google Cloud SDK (Optional)

Download the most current version of the Cloud SDK from [Google Developers page](https://developers.google.com/cloud/sdk/).

During the install process, it will ask if you want it to install any other SDKs.  You will want to select to Go SDK here.

### Install the Go App Engine SDK

If the Go SDK has not been installed via the Google Cloud SDK, then it will have to be manually installed.

#### General install

On the [Appengine Download page](https://developers.google.com/appengine/downloads)
select "Google App Engine SDK for Go", download the correct package for your platform.
For installing see the directions for your OS.

#### OSX

For Mac OSX the SDKs can be installed via Homebrew.

## Running the dev server

From the project directory run

    goapp serve <path/to/app.yaml>

## Deploying the application

From the project directory run

    goapp deploy <path/to/app.yaml>

## License

This site is under the BSD 3-clause license
