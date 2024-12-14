FROM pgvector/pgvector:0.8.0-pg16

COPY init.sql /docker-entrypoint-initdb.d/ 

RUN chmod 644 /docker-entrypoint-initdb.d/init.sql
