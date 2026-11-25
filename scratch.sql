\set VERBOSITY terse
--This is for testing only 
BEGIN
INSERT INTO transactions(description)
VALUES('alice pays ')
COMMIT;

SELECT * from account_balances;