#!/usr/bin/env bash
bd_name="pgadmin4_container"
docker cp ./dist_for_pg.csv $bd_name:/
docker cp ./embeddings.csv $bd_name:/
