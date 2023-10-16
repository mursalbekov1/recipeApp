ALTER TABLE recipes ADD CONSTRAINT recipes_runtime_check CHECK (runtime >= 0);

ALTER TABLE recipes ADD CONSTRAINT how_many_steps_check CHECK (array_length(steps, 1) BETWEEN 1 AND 10);
