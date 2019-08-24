import os
import json

def init():
    db_type = {}
    try:
        open('DB_Type.json', mode='r', encoding='utf-8')
    except FileNotFoundError:
        db_type["DBMS"] = input('Select DBMS Type[mariadb, sqlite] (default: mariadb) : ')
        if db_type["DBMS"] == '':
            db_type["DBMS"] = 'mariadb'
        
        db_type["etc"] = input('Please type your systems Python Execute Command (default: python) : ')
        if db_type["etc"] == '':
            db_type["etc"] = 'python'
        
        with open('DB_Type.json', mode='w', encoding='utf-8') as fileMake:
            json.dump(db_type, fileMake, ensure_ascii=False, indent="\t")
        pass

    with open("DB_Type.json") as fileRead:
        db_type = json.load(fileRead)

    if db_type["DBMS"] == "mariadb":
        os.system(db_type["etc"] + " MariaDB.py")
    elif db_type["DBMS"] == "sqlite":
        os.system(db_type["etc"] + " SQLite.py")
    else:
        print("DBMS Type Error")
init()