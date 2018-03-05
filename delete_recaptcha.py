import json
import sqlite3

json_data = open('set.json').read()
set_data = json.loads(json_data)

conn = sqlite3.connect(set_data['db'] + '.db', check_same_thread = False)
curs = conn.cursor()

curs.execute("delete from other where name = 'recaptcha'")
curs.execute("delete from other where name = 'sec_re'")
conn.commit()