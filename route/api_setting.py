from .tool.func import *
from .user_setting_skin_set_main import user_setting_skin_set_main_set_list

def api_setting(name = 'markup'):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        # from other
        ok_list_1 = ['markup']
        ok_list_1 += [for_a for for_a in user_setting_skin_set_main_set_list()]

        # from html_filter
        ok_list_2 = ['inter_wiki']
        
        if name in ok_list_1:
            curs.execute(db_change('select data from other where name = ?'), [name])
            rep_data = curs.fetchall()
            if rep_data:
                return flask.jsonify({ name : rep_data })
        elif name in ok_list_2:
            curs.execute(db_change("select html, plus, plus_t from html_filter where kind = ?"), [name])
            rep_data = curs.fetchall()
            if rep_data:
                return flask.jsonify({ name : rep_data })

        return flask.jsonify({})