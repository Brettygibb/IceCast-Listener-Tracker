# IceCast Listener Tracker

## Description
This Golang application monitors listener statistics from an IceCast server and logs hourly data into an Excel file.

## Features
- Fetches listener count every minute
- Compiles daily Excel reports
- Lightweight and efficient

## Technologies Used
- Golang
- IceCast JSON API
- Excel handling libraries

## How to Run
1. Clone the repo.
2. Run `go run main.go` after installing dependencies.
3. Configure the endpoint inside `config.json`.
