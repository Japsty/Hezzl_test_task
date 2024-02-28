version: '3'

services:
  pq_database:
    image: postgres
    restart: always
    environment:
      POSTGRES_DB: Efficent
      POSTGRES_USER: Mobile
      POSTGRES_PASSWORD: mobile
    ports:
      - "5432:5432"
    networks:
      - service_network

  service:
    build:
      context: .
      dockerfile: Dockerfile
    depends_on:
      - pq_database
    networks:
      - service_network
    ports:
      - "8080:8080"
    environment:
      PORT: '8080'
      MIGRATE_DB: true


  service2:
    build:
      context: .
      dockerfile: Dockerfile
    depends_on:
      - pq_database
    networks:
      - service_network
    ports:
      - "8081:8081"
    environment:
      PORT: '8081'
      MIGRATE_DB: false

  service3:
    build:
      context: .
      dockerfile: Dockerfile
    depends_on:
      - pq_database
    networks:
      - service_network
    ports:
      - "8082:8082"
    environment:
      PORT: '8082'
      MIGRATE_DB: false

  nginx_reverse_proxy:
    build:
      context: configs/nginx
      dockerfile: Dockerfile
    ports:
      - "80:80"
    depends_on:
      - service
      - service2
      - service3
    networks:
      - service_network

  prometheus:
    build:
      context: configs/prometheus
      dockerfile: Dockerfile
    networks:
      - service_network
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    networks:
      - service_network

networks:
  service_network:
    driver: bridge
