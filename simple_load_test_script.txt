#!/bin/bash
docker-compose up -d
sleep 3
hey -m PUT -c 1 -n 1 -d '{"type":"caesar", "input": "ab", "shift": -2}' http://localhost:8080/records/6667b924-463d-49a7-aaf0-f23febe00420
hey -m GET -c 1 -n 1000 http://localhost:8080/records/6667b924-463d-49a7-aaf0-f23febe00420
