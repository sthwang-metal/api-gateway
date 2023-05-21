#!/usr/bin/env bash

# cat ../../metadata-api/schema/*.graphql >> metadata-api.graphql
# cat ../schema/*.graphql >> location-api.graphql

rover supergraph compose --config supergraph.yml --output supergraph.graphql
