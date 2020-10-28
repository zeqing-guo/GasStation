#!/bin/bash

export $(less .env | xargs)
./bin/gas
