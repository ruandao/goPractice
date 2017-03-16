#!/usr/bin/env bash

echo Building doaminfinder...
go build -o doaminfinder

echo Building synonyms...
cd ../synonyms
go build -o ../domainfinder/lib/synonyms

echo Building available...
cd ../available
go build -o ../domainfinder/lib/available

echo Building sprinkle...
cd ../sprinkle
go build -o ../domainfinder/lib/sprinkle

echo Building domainify...
cd ../domainify
go build -o ../domainfinder/lib/domainify

echo Done.