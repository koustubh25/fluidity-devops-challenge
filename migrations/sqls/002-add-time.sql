
-- add the time to the table

ALTER TABLE average_compute_units
	ADD COLUMN created_by TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
