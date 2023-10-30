#!/bin/bash

# Obtenga el directorio actual
current_dir=$(pwd)

# Obtenga la lista de subcarpetas
subdirectories=$(find "$current_dir" -mindepth 1 -maxdepth 1 -type d)

# Copie el archivo .env en cada subcarpeta
for subdirectory in $subdirectories; do
  cp .env "$subdirectory/.env"
done