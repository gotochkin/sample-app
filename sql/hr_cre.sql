drop table employees;

CREATE TABLE IF NOT EXISTS employees (
		employee_id SERIAL NOT NULL,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		hire_date DATE NOT NULL,
		manager_id BIGINT,
		PRIMARY KEY (employee_id)
	);




INSERT INTO employees VALUES 
        ( 100
        , 'Steven'
        , 'King'
        , TO_DATE('17-06-2003', 'dd-MM-yyyy')
        , 100
        );
INSERT INTO employees VALUES 
        ( 101
        , 'Neena'
        , 'Kochhar'
        , TO_DATE('21-09-2005', 'dd-MM-yyyy')
        , 100
        );
INSERT INTO employees VALUES 
        ( 102
        , 'Lex'
        , 'De Haan'
        , TO_DATE('13-01-2001', 'dd-MM-yyyy')
        , 102
        );
INSERT INTO employees VALUES 
        ( 103
        , 'Alexander'
        , 'Hunold'
        , TO_DATE('03-01-2006', 'dd-MM-yyyy')
        , 102
        );
INSERT INTO employees VALUES 
        ( 104
        , 'Bruce'
        , 'Ernst'
        , TO_DATE('21-05-2007', 'dd-MM-yyyy')
        , 103
        );
INSERT INTO employees VALUES 
        ( 105
        , 'David'
        , 'Austin'
        , TO_DATE('25-06-2005', 'dd-MM-yyyy')
        , 103
        );
INSERT INTO employees VALUES 
        ( 106
        , 'Valli'
        , 'Pataballa'
        , TO_DATE('05-02-2006', 'dd-MM-yyyy')
        , 103
        );
commit;