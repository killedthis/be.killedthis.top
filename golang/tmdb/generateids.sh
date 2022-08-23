#!/bin/bash

readarray -t SHOWS <shows.txt
for INDEX in "${!SHOWS[@]}"; do
	./tmdb.go -s ${SHOWS[$INDEX]}
done
