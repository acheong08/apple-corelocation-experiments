UPDATE location_changes SET change_time = DATETIME(SUBSTR(change_time, 1, INSTR(change_time, ' +') - 1));
