-- Need to specify PGPASSWORD; PGUSER; PGDATABASE env variable
DROP SCHEMA IF EXISTS agora CASCADE;

\i schemas.sql
\i functions.sql

\i seeds/seed.sql