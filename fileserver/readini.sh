#!/bin/bash

# Nama file konfigurasi
config_file="/home/5678cxz/wth/fileserver/config.ini"

# Memastikan file ada
if [ ! -f "$config_file" ]; then
  echo "File konfigurasi tidak ditemukan: $config_file"
  exit 1
fi

# Membaca nilai port dengan grep/awk
port=$(grep -A 1 '\[webconf\]' "$config_file" | grep 'port' | awk -F '=' '{print $2}' | xargs)

if [ -z "$port" ]; then
  echo "Port tidak ditemukan atau entri port kosong."
else
  echo "Port: $port"
fi
