\c postgres

-- Inicialización de una base de datos con todo lo necesario para ejecutar el proyecto.
-- Los scripts de migraciones propias del proyecto se ejecutan mediante Flyway.

-- Elimina el esquema si ya existe
DROP SCHEMA IF EXISTS "dragon_ball" CASCADE;

-- Crea el esquema necesario
CREATE SCHEMA "dragon_ball";

-- Crea la tabla dentro del esquema
CREATE TABLE "dragon_ball"."characters"
(
    "ID" serial primary key,
    "name" varchar not null,
    "ki" varchar not null,
    "maxKi" varchar,
    "race" varchar,
    "gender" varchar,
    "description" varchar,
    "image" varchar(255),
    "affiliation" varchar,
    "deleteAt" timestamp without time zone
);

COMMENT ON TABLE "dragon_ball"."characters" IS 'Tabla que almacena todos los tipos de personaje de la serie animada Dragon Ball';

COMMENT ON COLUMN "dragon_ball"."characters"."ID" IS 'Identificador del tipo del personaje';
COMMENT ON COLUMN "dragon_ball"."characters"."name" IS 'Nombre del personaje';
COMMENT ON COLUMN "dragon_ball"."characters"."ki" IS 'Es el poder de nuestro personaje en el momento actual';
COMMENT ON COLUMN "dragon_ball"."characters"."maxKi" IS 'Es el máximo de poder que puede alcanzar';
COMMENT ON COLUMN "dragon_ball"."characters"."race" IS 'La raza con la cual ha nacido';
COMMENT ON COLUMN "dragon_ball"."characters"."gender" IS 'El género con el cual ha nacido';
COMMENT ON COLUMN "dragon_ball"."characters"."description" IS 'Descripción de nuestro personaje';
COMMENT ON COLUMN "dragon_ball"."characters"."image" IS 'El link de la imagen';
COMMENT ON COLUMN "dragon_ball"."characters"."affiliation" IS 'Temporada donde apareció por primera vez';
COMMENT ON COLUMN "dragon_ball"."characters"."deleteAt" IS 'Fecha en la cual fue eliminado';

CREATE INDEX characters_name ON "dragon_ball"."characters" (name);

