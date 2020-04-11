#!/bin/bash

go clean ; go build ; strip -s OgameStats ; zip OgameStats.zip OgameStats
