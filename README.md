# iptv-filter

###docker-compose example
```bigquery
version: "3.9"

services:
  iptv:
    image: sebastianoo/iptv-filter
    restart: always
    environment:
      - TARGET=M3U URL
      - EPG_URL=EPG URL 
      - PUBLIC_URL=localhost:65341
      - INCLUDE=|EU|  ITALIA,|EU| NETHERLAND
      - CACHE_TIME=1
    volumes:
      - iptv:/out
    ports:
      - "65341:65341"
volumes:
  iptv:
    external: false
```