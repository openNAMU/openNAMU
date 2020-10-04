from .tool.func import *
import flask
import json
import time

__api_all_title_data_0 = []
__api_all_title_data_1 = []
__api_all_title_last_update = 0



def api_all_title_2(conn):
    global __api_all_title_last_update
    if (time.time() - __api_all_title_last_update) > 120:
        __api_all_title_update(conn)
        __api_all_title_last_update = time.time()
    data = __api_all_title_data_1
    resp = flask.jsonify({"datas":data})
    resp.headers['Cache-Control'] = 'max-age='+str(60*10)
    return resp

def __api_all_title_update(conn):
    #print("__update")
    global __api_all_title_data_0
    global __api_all_title_data_1
    __api_all_title_data_0 = []
    curs = conn.cursor()
    curs.execute(db_change("SELECT title FROM data"))
    data = curs.fetchall()
    #print(data)
    if data:
        for tmp in data:
            __api_all_title_data_0.append(tmp[0])
    else:
        data = []
    #print(__api_all_title_data_0)
    __api_all_title_data_1 = __api_all_title_data_0[:]
