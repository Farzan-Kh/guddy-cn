#! /usr/bin/env python3

import json
import uuid
from sqlalchemy import create_engine, text
import os

eq_mapper = {
    "body only": "Bodyweight",
    "machine": "Machine",
    "kettlebells": "Kettlebells",
    "dumbbell": "Dumbbells",
    "cable": "Cables",
    "bands": "Band",
    "exercise ball": "Exercise Ball",
    "barbell": "Barbell",
    "e-z curl bar": "Barbell",
    "medicine ball": "Medicine Ball",
    "other": "Other",
    "foam roll": "Foam roll",
    "": "Other"
}


DB_URL = os.getenv("DATABASE_URL")

def main():
    # Load exercises from JSON
    with open("../internal/db/data/exercises.json", "r") as f:
        data = json.load(f)
        exercises = data["exercises"]

    # Create SQLAlchemy engine
    engine = create_engine(DB_URL)

    with engine.connect() as conn:
        # Check if exercises already exist
        result = conn.execute(text("SELECT COUNT(*) FROM exercises"))
        count = result.scalar()
        if count and count > 0:
            print("Data Already exists, exiting.")
            return
        for idx, ex in enumerate(exercises):
            print('Inserting row: ', idx)
            name = ex["name"]
            muscles = ex['primaryMuscles'] + ex['secondaryMuscles']
            equipment = 'Other'
            if ex['equipment'] != None:
                equipment = eq_mapper[ex['equipment']]
            instructions = ex['instructions']
            exercise_id = idx + 1  # Using integer IDs as per updated schema
            # 1. Insert into exercises table
            result = conn.execute(
                text("""
                    INSERT INTO exercises (id, instructions, equipment, visuals_id)
                    VALUES (:id, :instructions, :equipment, NULL)
                """),
                {"id": exercise_id, "instructions": instructions, "equipment": equipment}
            )
            # 2. Insert into exercise_names table
            conn.execute(
                text("""
                    INSERT INTO exercise_names (name, exercise_id)
                    VALUES (:name, :exercise_id)
                """),
                {"name": name, "exercise_id": exercise_id}
            )
            # 3. For each muscle, ensure muscle exists and link
            for muscle in muscles:
                # Try to insert muscle, ignore if exists
                conn.execute(
                    text("""
                        INSERT INTO muscles (name)
                        VALUES (:name)
                        ON CONFLICT (name) DO NOTHING
                    """),
                    {"name": muscle}
                )
                # Get muscle id
                muscle_id_result = conn.execute(
                    text("SELECT id FROM muscles WHERE name = :name"),
                    {"name": muscle}
                )
                muscle_id = muscle_id_result.scalar()
                # Link exercise and muscle
                conn.execute(
                    text("""
                        INSERT INTO exercise_muscle (exercise_id, muscle_id)
                        VALUES (:exercise_id, :muscle_id)
                        ON CONFLICT DO NOTHING
                    """),
                    {"exercise_id": exercise_id, "muscle_id": muscle_id}
                )
        conn.commit()
        print("Import complete.")

if __name__ == "__main__":
    main()