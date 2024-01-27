from .tool.func import *

def api_search(name = 'Test', num = 10, page = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = 1 if num > 1000 else num
        page = (page * (num - 1)) if page * num > 0 else 0

        # 개편 예정
        curs.execute(db_change('select data from other where name = "count_all_title"'))
        if int(curs.fetchall()[0][0]) < 30000:
            curs.execute(db_change("" + \
                "select distinct title, case " + \
                "when title like ? then 'title' else 'data' end from data " + \
                "where (title like ? or data like ?) order by case " + \
                "when title like ? then 1 else 2 end limit ?, ?"),
                ['%' + name + '%', '%' + name + '%', '%' + name + '%', '%' + name + '%', page, num]
            )
        else:
            curs.execute(db_change("select title from data where title like ? order by title limit ?, ?"),
                ['%' + name + '%', page, num]
            )
            
        all_list = curs.fetchall()
        if all_list:
            return flask.jsonify(all_list)
        else:
            return flask.jsonify({})