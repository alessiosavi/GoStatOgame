#!/bin/bash
go clean ; go build ; strip GoStatOgame ; rm GoStatOgame.zip ;  zip GoStatOgame.zip GoStatOgame

