from .tool.func import *

def api_search_2(conn, name):
    curs = conn.cursor()

    num = int(number_check(flask.request.args.get('num', '10')))
    if not num > 0:
        num = 1

    if num > 1000:
        num = 1

    curs.execute("" + \
        "select distinct title, case when title like ? then 'title' else 'data' " + \
        "end from data where title like ? or data like ? order by case " + \
        "when title like ? then 1 else 2 end limit ?",
        ['%' + name + '%', '%' + name + '%', '%' + name + '%', '%' + name + '%', str(num)]
    )
    all_list = curs.fetchall()
    if all_list:
        return flask.jsonify(all_list)
    else:
        return flask.jsonify({})