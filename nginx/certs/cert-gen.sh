#! /bin/bash

# Generates a self signed tls cert for nginx to use

certDir=$(dirname $0)

sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout "$certDir/localhost.key" -out "$certDir/localhost.crt" -config "$certDir/localhost.conf"