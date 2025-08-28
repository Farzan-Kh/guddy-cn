-- Fetch Exercises with name, equipment and muscle filters (all optional)
-- name: GetExercises :many
SELECT 
    e.id,
    string_agg(DISTINCT e_names.name, ', ') AS names_grouped,
    e.equipment,
    string_agg(DISTINCT m.name, ', ') AS muscles_grouped,
    string_agg(DISTINCT v.path, ', ') AS visuals_grouped
FROM exercises e
INNER JOIN exercise_names e_names ON e_names.exercise_id = e.id
INNER JOIN exercise_muscle e_m ON e_m.exercise_id = e.id
INNER JOIN muscles m ON m.id = e_m.muscle_id
LEFT JOIN visuals v ON v.id = e.visuals_id
WHERE
  (coalesce(sqlc.narg('name')) IS NULL OR e_names.name ILIKE '%' || @name::text || '%') AND
  (coalesce(sqlc.narg('equipment')) IS NULL OR e.equipment = @equipment::equipment_t) AND
  (coalesce(sqlc.narg('muscle')) IS NULL OR m.name = @muscle::text) AND
  (coalesce(sqlc.narg('exercise_id')) IS NULL OR e.id = @exercise_id::int)
GROUP BY
    e.id
ORDER BY e.id
LIMIT coalesce(sqlc.narg('limit'), 50)
OFFSET coalesce(sqlc.narg('offset'), 0);


-- Fetch all Muscles
-- name: GetMuscles :many
SELECT
    id,
    name
FROM
    muscles;


-- Fetch equipment_t values
-- SELECT
--   e.enumlabel AS value
-- FROM
--   pg_type t
--   JOIN pg_enum e ON t.oid = e.enumtypid
-- WHERE
--   t.typname = 'equipment_t'
-- ORDER BY
--   e.enumsortorder;


-- Fetch Full Program by id
-- name: GetFullProgramById :many
SELECT
  idx,
  string_agg(DISTINCT e_names.name, ', ') AS names_grouped,
  e.equipment,
  sets,
  reps,
  string_agg(DISTINCT m.name, ', ') AS muscles_grouped,
  string_agg(DISTINCT v.path, ', ') AS visuals_grouped
FROM 
  programs p
  INNER JOIN exercises e ON e.id = p.exercise_id
  INNER JOIN exercise_names e_names ON e_names.exercise_id = e.id
  INNER JOIN exercise_muscle e_m ON e_m.exercise_id = e.id
  INNER JOIN muscles m ON m.id = e_m.muscle_id
  LEFT JOIN visuals v ON v.id = e.visuals_id
WHERE
  p.id = @program_id::uuid
GROUP BY
  e.id, idx, sets, reps
ORDER BY
  idx;

-- Fetch Program by id
-- name: GetProgramById :many
SELECT
  idx,
  e.id AS exercise_id,
  sets,
  reps
FROM 
  programs p
  INNER JOIN exercises e ON e.id = p.exercise_id
WHERE
  p.id = @program_id::uuid
GROUP BY
  e.id, idx, sets, reps
ORDER BY
  idx;


-- Insert into exercise_names
-- name: InsertToExerciseNames :one
INSERT INTO 
  exercise_names(exercise_id, name)
VALUES
  (@exercise_id::int, @name::text)
RETURNING id;

-- Insert into exercises
-- name: InsertToExercises :one
INSERT INTO
  exercises(equipment, visuals_id)
VALUES
  (@equipment::equipment_t, @visuals_id::int)
RETURNING id;

-- Insert into muscles
-- name: InsertToMuscles :one
INSERT INTO 
  muscles(name)
VALUES
  (@name::text)
RETURNING id;

-- Insert into exercise_muscles
-- name: InsertToExerciseMuscle :exec
INSERT INTO
  exercise_muscle(exercise_id, muscle_id)
VALUES
  (@exercise_id::int, @muscle_id::int);

-- Insert into visuals
-- name: InsertToVisuals :one
INSERT INTO
  visuals(path)
VALUES
  (@path::text)
RETURNING id;

-- name: InsertToProgramsById :exec
INSERT INTO
  programs(id, idx, exercise_id, sets, reps)
VALUES
  (@id::uuid, @idx::int, @exercise_id::int, @sets::int, @reps::int);
