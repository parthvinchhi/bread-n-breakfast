#!/bin/bash

go build -o bookings cmd/web/*.go
./bookings -dbname=bread-n-breakfast -dbuser=postgres -dbpass=123 -cache=false -production=false