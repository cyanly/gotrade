CREATE TABLE IF NOT EXISTS order_status
(
  order_status_id SERIAL,
  description     TEXT UNIQUE    NOT NULL,
  name            CHAR(3) UNIQUE NOT NULL,
  fix42value      CHAR(1),
  fix44value      CHAR(1),
  fix50value      CHAR(1),

  PRIMARY KEY (order_status_id)
);

--=====================
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (0, 'NEW', 'O  ', '0', '0', '0');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (1, 'PARTIALLY_FILLED', 'PF ', '1', '1', '1');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (2, 'FILLED', 'F  ', '2', '2', '2');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (3, 'DONE_FOR_DAY', 'DD ', '3', '3', '3');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (4, 'CANCELED', 'X  ', '4', '4', '4');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (5, 'REPLACED', 'RP ', '5', null, '5');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (6, 'PENDING_CANCEL', 'PC ', '6', '6', '6');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (7, 'STOPPED', 'ST ', '7', '7', '7');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (8, 'REJECTED', 'R  ', '8', '8', '8');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (9, 'SUSPENDED', 'SP ', '9', '9', '9');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (10, 'PENDING_NEW', 'PN ', 'A', 'A', 'A');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (11, 'CALCULATED', 'CD ', 'B', 'B', 'B');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (12, 'EXPIRED', 'EX ', 'C', 'C', null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (13, 'ACCEPTED_FOR_BIDDING', 'AB ', 'D', 'D', 'D');
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (14, 'PENDING_REPLACE', 'PR ', 'E', 'E', 'E');

INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (100, 'ORDER_RECEIVED', 'OR ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (101, 'ORDER_SENT', 'OS ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (102, 'MC_ACK_ORDER', 'OA ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (103, 'MC_SENT_ORDER', 'SO ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (110, 'CANCEL_RECEIVED', 'CR ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (111, 'CANCEL_SENT', 'CS ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (112, 'MC_ACK_CANCEL', 'CA ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (113, 'MC_SENT_CANCEL', 'CB ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (120, 'REPLACE_RECEIVED', 'RR ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (121, 'REPLACE_SENT', 'RS ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (122, 'MC_ACK_REPLACE', 'RA ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (123, 'MC_SENT_REPLACE', 'RB ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (125, 'EXECUTION_REPLACED', 'ER ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (126, 'EXECUTION_CANCELLED', 'EC ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (130, 'BOOKABLE', 'BA ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (131, 'TRADE_PUMP_RECEIVED', 'TPR', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (132, 'TRADE_PUMP_PROCESSED', 'TPP', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (134, 'SENT_PMS', 'TS ', null, null, null);
INSERT INTO order_status(order_status_id, description, name, fix42value, fix44value, fix50value) VALUES (135, 'BOOKED', 'BK ', null, null, null);
--=====================


CREATE TABLE IF NOT EXISTS exec_type
(
  exec_type_id SERIAL,
  name         TEXT UNIQUE NOT NULL,
  fix42value   CHAR(1),
  fix44value   CHAR(1),
  fix50value   CHAR(1),

  PRIMARY KEY (exec_type_id)
);


--=====================
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (0, 'NEW', '0', '0', '0');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (1, 'PARTIAL_FILL', '1', null, null);
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (2, 'FILL', null, null, '2');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (3, 'DONE_FOR_DAY', '3', '3', '3');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (4, 'CANCELED', '4', '4', '4');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (5, 'REPLACE', '5', '5', '5');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (6, 'PENDING_CANCEL', '6', '6', '6');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (7, 'STOPPED', '7', '7', '7');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (8, 'REJECTED', '8', '8', '8');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (9, 'SUSPENDED', '9', '9', '9');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (10, 'PENDING_NEW', 'A', 'A', 'A');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (11, 'CALCULATED', 'B', 'B', 'B');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (12, 'EXPIRED', 'C', 'C', 'C');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (13, 'RESTATED', 'D', 'D', 'D');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (14, 'PENDING_REPLACE', 'E', 'E', 'E');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (15, 'TRADE', null, 'F', 'F');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (16, 'TRADE_CORRECT', null, 'G', 'G');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (17, 'TRADE_CANCEL', null, 'H', 'H');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (18, 'ORDER_STATUS', null, 'I', 'I');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (19, 'CLEARING_HOLD', null, null, 'J');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (20, 'RELEASED_TO_CLEARING', null, null, 'K');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (21, 'SYSTEM', null, null, 'L');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (22, 'LOCKED', null, null, 'M');
INSERT INTO exec_type(exec_type_id, name, fix42value, fix44value, fix50value) VALUES (23, 'RELEASED', null, null, 'N');
--=====================


