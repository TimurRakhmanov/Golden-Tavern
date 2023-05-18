#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings
./bookings -dbname=bookings -dbuser=rakhmanovtimur -cache=false -production=false