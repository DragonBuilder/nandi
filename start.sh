#!/bin/bash

export $(xargs < .env)
go run .