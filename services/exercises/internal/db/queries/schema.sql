CREATE TYPE equipment_t 
AS 
ENUM(
  'Dumbbells',
  'Barbell',
  'Machine',
  'Bodyweight',
  'Medicine Ball',
  'Kettlebells',
  'Streches',
  'Cables',
  'Band',
  'Plate',
  'TRX',
  'Bosu Ball',
  'Foam roll',
  'Exercise Ball',
  'Other'
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS visuals (
  id SERIAL PRIMARY KEY,
  path VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS exercises (
  id SERIAL PRIMARY KEY,
  equipment equipment_t,
  visuals_id INT,
  FOREIGN KEY (visuals_id) REFERENCES visuals (id)
);

CREATE TABLE IF NOT EXISTS muscles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE
);

CREATE TABLE IF NOT EXISTS exercise_names (
  id SERIAL PRIMARY KEY,
  exercise_id INT NOT NULL,
  name VARCHAR(255) NOT NULL UNIQUE,
  FOREIGN KEY (exercise_id) REFERENCES exercises (id)
);

CREATE TABLE IF NOT EXISTS exercise_muscle (
  exercise_id INT NOT NULL,
  muscle_id INT,
  PRIMARY KEY (exercise_id, muscle_id),
  FOREIGN KEY (exercise_id) REFERENCES exercises (id),
  FOREIGN KEY (muscle_id) REFERENCES muscles (id)
);

CREATE TABLE IF NOT EXISTS programs (
  id UUID DEFAULT uuid_generate_v4(),
  idx INT NOT NULL,
  exercise_id INT NOT NULL,
  sets INT NOT NULL,
  reps INT NOT NULL,
  PRIMARY KEY (id, exercise_id),
  FOREIGN KEY (exercise_id) REFERENCES exercises (id)
);